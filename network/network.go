package network

import (
	"net"
	."fmt"
	."time"
	."strings"
	."encoding/json"
)

const ELEV_COUNT = 3
const FLOOR_COUNT = 4
const MASTER_INIT_IP = "129.241.187.161"
const PORT = ":12345"

type Message struct {
	Type string
	Content string
	Floor int
	Value int
	From int
	To int
}

func NetworkInit(networkReceive chan Message, networkSend chan Message, fileOutChan chan Message, fileInChan chan Message, failureChan chan Message) {

	receivedChannel := make(chan Message)
	go listen(receivedChannel)
	message := Message{}
	fileEmpty := true
	IPlist := make([]string, ELEV_COUNT + 1)
	for i := 1; i < ELEV_COUNT + 1; i++ {
		IPlist[i] = "Uninitialized"
	}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		Println("\n", "Address error: ", err)
	}
	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				IPlist[0] = ipnet.IP.String()
			}
		}
	}
	j := 1
	for i := 1; i < ELEV_COUNT + 1; i++ {

		message.Type = "readIP"
		message.Value = i
		fileInChan <- message
		message = <- fileOutChan
		if message.Content != "noIP" && message.Content != "" {
			fileEmpty = false
			IPlist[j] = message.Content
			Println("Added IP to IPlist", IPlist[j], i)
			j++
		}
	}
	go networkReceiver(networkSend, networkReceive, receivedChannel, fileInChan, fileEmpty, &IPlist)
	go networkSender(networkSend, failureChan, fileEmpty, &IPlist)
	go failureHandler(networkSend, failureChan, &IPlist)
}

func listen(receivedChannel chan Message) {
	listener, error := net.Listen("tcp", PORT)
	if error != nil {
		Println("\n", "Listen error: ", error)
	}
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if error != nil {
			Println("\n", "Listen connection error: ", err)
		}
		go receive(connection, receivedChannel)
	}
}

func receive(connection net.Conn, receivedChannel chan Message) {
	defer connection.Close()
	buffer := make([]byte, 1024)
	message := Message{}
	length, error := connection.Read(buffer)
	if error != nil {
		Println("\n", "Receive connection error: ", error)
	}
	err := Unmarshal(buffer[:length], &message)
	
	if err != nil {
		Println("\n", "Receive error: ", err)
	}
	receivedChannel <- message
}


func send(message Message, IPlist[] string, networkSend chan Message, failureChan chan Message) {
	recipient := ""
	switch{
	case message.To == -2:
		message.To = message.From
		recipient = "localhost"
	case message.To == -1:
		message.To = 1
		recipient = MASTER_INIT_IP
	case message.To > 0:
		recipient = IPlist[message.To]
		if recipient == "Uninitialized" || (message.Type != "broadcast" && Contains(recipient, "offline")) {
			return
		}
	}
	dialAddress := recipient
	connection, error := net.DialTimeout("tcp", TrimRight(dialAddress, "offline")+ PORT, Duration(100)*Millisecond)
	if error != nil {
		if message.From == message.To {
			connection, _ = net.Dial("tcp", "localhost"+ PORT)
		} else if message.Type == "broadcast" && Contains(recipient, "offline"){
			return
		} else {
			Println("\n", "Send connection error: ", error)
			message.Type = "connectionFailure"
			message.Content = recipient
			message.Value = message.To
			failureChan <- message
			return
		}
	}
	byteMessage, err := Marshal(message)
	if err != nil {
		Println("\n", "Send error: ", err)
	}
	connection.Write(byteMessage)
}

func networkSender(networkSend chan Message, failureChan chan Message, fileEmpty bool, IPlist *[]string) {
	for {
		select {
		case message := <- networkSend:
			for i := 1; i < ELEV_COUNT + 1; i++ {
				if (*IPlist)[0] == (*IPlist)[i] {
					message.From = i
					break
				} else {
					message.From = 0
				}
			}
			switch{					//0 = all, -1 = MASTER_INIT_IP, -2 = localhost
			case message.Type == "findMaster":
				message.Content = (*IPlist)[0]
				if fileEmpty {
					message.To = -1
					if message.Content == MASTER_INIT_IP {
						message.Value = 1
					}
				} else {
					message.To = 0
				}
			case message.Type == "masterOffline":
				message.To = -2

			case message.Type == "addElev":
				message.Content = (*IPlist)[message.Value]

			case message.Type == "newOrder" || message.Type == "deleteOrder":
				if message.Content == "inside" {
					message.To = -2
				} else {
					message.To = 0
				}
			case message.Type == "stateUpdate" || message.Type == "targetUpdate":
				message.To = 0
			case message.Type == "floorReached":
				message.To = -2
			}
			if message.To == 0 {
				for i := 1; i < ELEV_COUNT + 1; i++ {
					if !(message.Type == "elevOffline" && message.Value == i) {
						message.To = i
						go send(message, *IPlist, networkSend, failureChan)
					}
				}
			} else {
				go send(message, *IPlist, networkSend, failureChan)
			}
		}
	}
}

func networkReceiver(networkSend chan Message, networkReceive chan Message, receivedChannel chan Message, fileInChan chan Message, fileEmpty bool, IPlist *[]string) {
	for{
		select{
		case message := <- receivedChannel:
			switch{
			case message.Type == "broadcast":
				if Contains((*IPlist)[message.From], "offline") {
					message.Type = "elevOnline"
					message.To = 0
					message.Value = message.From
					networkSend <- message
					break
				}

			case message.Type == "addElev":
				alreadyAddedIndex := 0
				for i := 1; i < ELEV_COUNT + 1; i++ {
					if Contains((*IPlist)[i], message.Content) {
						alreadyAddedIndex = i
						break
					}
				}
				if alreadyAddedIndex == 0 {
					if message.Value == 0 {
						for i := 1; i < ELEV_COUNT + 1; i++ {
							if (*IPlist)[i] == "Uninitialized" {
								(*IPlist)[i] = message.Content
								message.Value = i
								Println("\n", *IPlist)
								break
							}
						}
					} else {
						(*IPlist)[message.Value] = message.Content
						Println("\n", *IPlist)
					}
					if message.Content != "Uninitialized" {
						message.Type = "writeIP"
						fileInChan <- message
						message.Type = "addElev"
					}
				}

			case message.Type == "findMaster":
				alreadyAddedIndex := 0
				for i := 1; i < ELEV_COUNT + 1; i++ {
					if Contains((*IPlist)[i], message.Content) {
						alreadyAddedIndex = i
						break
					}
				}
				if alreadyAddedIndex == 0 {
					if message.Value == 0 {
						for i := 1; i < ELEV_COUNT + 1; i++ {
							if (*IPlist)[i] == "Uninitialized" {
								(*IPlist)[i] = message.Content
								message.Value = i
								Println("\n", *IPlist)
								break
							}
						}
					} else {
						(*IPlist)[message.Value] = message.Content
						Println("\n", *IPlist)
					}
					message.Type = "writeIP"
					fileInChan <- message
					message.Type = "findMaster"
				} else if alreadyAddedIndex > 0 {
					message.Type = "elevOnline"
					message.Value = alreadyAddedIndex
					message.To = 0
					networkSend <- message
					break	
				}
				
			case message.Type == "elevOnline":
				(*IPlist)[message.Value] = TrimRight((*IPlist)[message.Value], "offline")

			case message.Type == "elevOffline":
				if !Contains((*IPlist)[message.Value], "offline") {
					(*IPlist)[message.Value] += "offline"
				}
			}
			networkReceive <- message
		}
	}
}

func failureHandler(networkSend chan Message, failureChan chan Message, IPlist *[]string) {
	for {
		select {
		case message := <- failureChan:
			switch {
			case message.Type == "masterOffline":
				message.Type = "elevOffline"
				message.Content = (*IPlist)[1]
				message.To = 0
			case message.Type == "connectionFailure":
				message.Type = "elevOffline"
				message.To = 0
			}
			Println("\n", "FailureHandler:", message.Type, message.Content, "number = ", message.Value)
			networkSend <- message
		}
	}
}