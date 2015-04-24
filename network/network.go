package network

import (
	"net"
	."fmt"
	."time"
	."encoding/json"
)

const ELEV_COUNT = 3
const FLOOR_COUNT = 4
const MASTER_INIT_IP = "129.241.187.154"
const PORT = ":12345"

type Message struct {
	Type string
	Content string
	Floor int
	Value int
	From int
	To int
}

func NetworkInit(networkReceive chan Message, networkSend chan Message, fileOutChan chan Message, fileInChan chan Message) {

	receivedChannel := make(chan Message)
	go listen(receivedChannel)

	message := Message{}
	fileEmpty := true

	IPlist := make([]string, 1, ELEV_COUNT + 1)

	addresses, err := net.InterfaceAddrs()
	if err != nil {
		Println("Address error: ", err)
	}
	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				IPlist[0] = ipnet.IP.String()
			}
		}
	}

	if IPlist[0] == MASTER_INIT_IP {
		message.Type = "master"
		go startAliveBroadcast(message, networkSend)
	}

	for i := 1; i < ELEV_COUNT + 1; i++ {
		message.Type = "readIP"
		message.Value = i
		fileInChan <- message
		message = <- fileOutChan
		if message.Content != "noIP" {
			fileEmpty = false
			IPlist = append(IPlist, message.Content)
		}
	}

	go networkReceiver(networkReceive, receivedChannel, networkSend, fileInChan, &IPlist)
	go networkSender(networkSend, fileEmpty, &IPlist)
}

func networkReceiver(networkReceive chan Message, receivedChannel chan Message, networkSend chan Message, fileInChan chan Message, IPlist *[]string) {
	for{
		select{
		case message := <- receivedChannel:
			switch{
			case message.Type == "newElev":
				appendable := true
				for i := 1; i < len(*IPlist); i++ {
					if (*IPlist)[i] == message.Content {
						appendable = false
					}
				}
				if appendable {
					*IPlist = append(*IPlist, message.Content)
				}
				message.Type = "writeIP"
				fileInChan <- message
				message.Type = "addElev"
				message.To = len(*IPlist) - 1
				if len(*IPlist) > 2 && (*IPlist)[0] == (*IPlist)[1] {
					for i := 1; i < len(*IPlist); i++ {
						message.Content = (*IPlist)[i]
						networkSend <- message
					}
				}
			case message.Type == "addElev":
				appendable := true
				for i := 1; i < len(*IPlist); i++ {
					if (*IPlist)[i] == message.Content {
						appendable = false
					}
				}
				if appendable {
					*IPlist = append(*IPlist, message.Content)
				}
			case message.Type == "elevOffline":
				*IPlist = append((*IPlist)[:message.Value], (*IPlist)[message.Value+1:]...)
			}
			networkReceive <- message
		}
	}
}

func networkSender(networkSend chan Message, fileEmpty bool, IPlist *[]string) {
	for {
		select {
		case message := <- networkSend:
			for i := 1; i < len(*IPlist); i++ {
				if (*IPlist)[0] == (*IPlist)[i] {
					message.From = i
					break
				}
			}
			switch{								//0 = all, -1 = MASTER_INIT_IP, -2 = localhost
			case message.Type == "newElev":
				message.Content = (*IPlist)[0]
				if fileEmpty {
					message.To = -1
				} else {
					message.To = 0
				}
			case message.Type == "newOrder" || message.Type == "deleteOrder":
				if message.Content == "inside" {
					message.To = -2
				} else {
					message.To = 0
				}
			case message.Type == "stateUpdate" || message.Type == "targetUpdate" || message.Type == "floorUpdate":
				message.To = 0
			case message.Type == "floorReached":
				message.To = -2
			}
			if message.To == 0 {
				for i := 1; i < len(*IPlist); i++ {
					message.To = i
					go send(message, *IPlist, networkSend)
				}
			} else {
				go send(message, *IPlist, networkSend)
			}
		}
	}
}

func listen(receivedChannel chan Message) {
	listener, error := net.Listen("tcp", PORT)
	if error != nil {
		Println("Listen error: ", error)
	}
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if error != nil {
			Println("Listen connection error: ", err)
		}
		go receive(connection, receivedChannel)
	}
}

func send(message Message, IPlist[] string, networkSend chan Message) {
	if message.To > len(IPlist) - 1 {
		return
	}
	recipient := ""
	switch{
	case message.To == -2:
		message.To = message.From
		recipient = "localhost"
	case message.To == -1:
		recipient = MASTER_INIT_IP
	case message.To > 0:
		recipient = IPlist[message.To]
	}
	connection, error := net.Dial("tcp", recipient+ PORT)
	defer connection.Close()

	if error != nil {
		if message.From == message.To {
			connection, _ := net.Dial("tcp", "localhost"+ PORT)
			defer connection.Close()
		} else {
			Println("Send connection error: ", error)
			message.Type = "elevOffline"
			message.Value = message.To
			for i := 1; i < len(IPlist); i++ {
				if i != message.Value {
					Println("Sending elevOffline message to elev", i)				
					message.To = i
					networkSend <- message
				}
			}
			if message.To == 1 {
				Sleep(100 * Millisecond)
				message.Type = "master"
				networkSend <- message
			}
			return
		}
	}
	byteMessage, err := Marshal(message)
	if err != nil {
		Println("Send error: ", err)
	}
	connection.Write(byteMessage)
}

func receive(connection net.Conn, receivedChannel chan Message) {
	defer connection.Close()
	buffer := make([]byte, 1024)
	message := Message{}
	length, error := connection.Read(buffer)
	if error != nil {
		Println("Receive connection error: ", error)
	}
	err := Unmarshal(buffer[:length], &message)
	
	if err != nil {
		Println("Receive error: ", err)
	}
	receivedChannel <- message
}

func startAliveBroadcast(message Message, networkSend chan Message) {
	Sleep(100 * Millisecond)
	networkSend <- message
}
