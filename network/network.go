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
	/*
	for i := 1; i < ELEV_COUNT + 1; i++ {
		message.Type = "readIP"
		message.Value = i
		fileInChan <- message
		message = <- fileOutChan
		if message.Content != "noIP" {
			fileEmpty = false
			IPlist[i] = message.Content
		}
	}*/
	if fileEmpty && IPlist[0] == MASTER_INIT_IP {
		message.Type = "master"
		message.To = -2
		go startAliveBroadcast(message, networkSend)
	}
	go networkReceiver(networkReceive, receivedChannel, fileInChan, fileEmpty, &IPlist)
	go networkSender(networkSend, failureChan, fileEmpty, &IPlist)
	go failureHandler(networkSend, failureChan, &IPlist)
}

func startAliveBroadcast(message Message, networkSend chan Message) {
	Sleep(100 * Millisecond)
	networkSend <- message
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
			case message.Type == "addElev":
				message.Content = (*IPlist)[message.Value]

			case message.Type == "newMaster":
				message.To = 0
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
			if message.Type == "addElev" {
				Println("\n Sending", message)
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
		if recipient == "" || Contains(recipient, "offline") {
			return
		}
	}
	connection, error := net.DialTimeout("tcp", recipient+ PORT, Duration(100)*Millisecond)
	if error != nil {
		if message.From == message.To {
			connection, _ = net.Dial("tcp", "localhost"+ PORT)
		} else {
			Println("\n", "Send connection error: ", error)
			message.Type = "connectionFailure"
			message.Content = recipient
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

func networkReceiver(networkReceive chan Message, receivedChannel chan Message, fileInChan chan Message, fileEmpty bool, IPlist *[]string) {
	for{
		select{
		case message := <- receivedChannel:
			switch{
			case message.Type == "addElev" || message.Type == "findMaster":
				if message.Value == 0 {
					for i := 1; i < ELEV_COUNT + 1; i++ {
						if (*IPlist)[i] == message.Content {
							break
						}
						if (*IPlist)[i] == "" {
							(*IPlist)[i] = message.Content
							Println("\n", *IPlist)
							break
						}
					}
				} else {
					(*IPlist)[message.Value] = message.Content
				}
				//message.Type = "writeIP"
				//fileInChan <- message
				//message.Type = "addElev"

			case message.Type == "elevOnline":
				(*IPlist)[message.Value] = TrimRight((*IPlist)[message.Value], "offline")

			case message.Type == "elevOffline":
				(*IPlist)[message.Value] += "offline"
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
				message.Value = 1
				message.To = -2
			case message.Type == "connectionFailure":
				message.Type = "elevOffline"
				message.Value = message.To
				message.To = 0
			}
			Println("\n", "FailureHandler:", message.Type, message.Content, "Value = ", message.Value, "To elev", message.To)
			networkSend <- message
		}
	}
}
