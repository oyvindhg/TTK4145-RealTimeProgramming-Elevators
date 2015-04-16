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
const PORT = ":20007"

type Message struct {
	Type string
	Content string
	Floor int
	Value int
	From int
	To int
}

func Network(networkReceive chan Message, networkSend chan Message, fileInChan chan Message, fileOutChan chan Message) {

	recievedChannel := make(chan Message)
	go listen(recievedChannel)

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
		go startBroadcast(message, -1, IPlist, networkSend)
		//go send(message, IPlist, networkSend)
	}

	for i := 1; i < ELEV_COUNT + 1; i++ {
		message.Type = "readIP"
		message.Value = i
		fileOutChan <- message
		message = <- fileInChan
		if message.Content != "noIP" {
			fileEmpty = false
			IPlist = append(IPlist, message.Content)
		}
	}

	for{
		select{
			case message = <- recievedChannel:
				switch{
				case message.Type == "newElev":
					appendable := true
					for i := 0; i < len(IPlist); i++ {
						if IPlist[i] == message.Content {
							appendable = false
						}
					}
					if appendable {
						IPlist = append(IPlist, message.Content)
					}
					message.Type = "writeIP"
					fileOutChan <- message
					message.Type = "addElev"
					message.To = len(IPlist) - 1
					if len(IPlist) > 2 && IPlist[0] == IPlist[1] {
						for i := 1; i < len(IPlist); i++ {
							message.Content = IPlist[i]
							networkSend <- message
						}
					}
				case message.Type == "addElev":
					appendable := true
					for i := 1; i < len(IPlist); i++ {
						if IPlist[i] == message.Content {
							appendable = false
						}
					}
					if appendable {
						IPlist = append(IPlist, message.Content)
					}
				case message.Type == "elevOffline":
					IPlist = append(IPlist[:message.Value], IPlist[message.Value+1:]...)
				}
				networkReceive <- message

			case message = <- networkSend:
				for i := 0; i < len(IPlist); i++ {
					if IPlist[0] == IPlist[i] {
						message.From = i
						break
					}
				}
				switch{							//0 = all, -1 = MASTER_INIT_IP, -2 = localhost
				case message.Type == "newElev":
					message.Content = IPlist[0]
					if fileEmpty {
						message.To = -1
					} else {
						message.To = 0
					}
				case message.Type == "newOrder":
					message.To = 0
				case message.Type == "deleteOrder":
					message.To = 0
				case message.Type == "stateUpdate":
					message.To = 0
				case message.Type == "floorReached":
					message.To = -2
				}
				if message.To == 0 {
					for i := 1; i < len(IPlist); i++ {
						message.To = i
						go send(message, IPlist, networkSend)
					}
				} else {
					go send(message, IPlist, networkSend)
				}
				
		}
	}
}

func listen(recievedChannel chan Message) {
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
		go receive(connection, recievedChannel)
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
		Println("Send connection error: ", error)
		
		//ELEVATOR OFFLINE
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
	byteMessage, err := Marshal(message)
	if err != nil {
		Println("Send error: ", err)
	}
	connection.Write(byteMessage)
}

func receive(connection net.Conn, recievedChannel chan Message) {
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
	recievedChannel <- message
}

func startBroadcast(message Message, i int, IPlist[] string, networkSend chan Message) {
	Sleep(400 * Millisecond)
	go send(message, IPlist, networkSend)
}
