package network

import (
	"net"
	."fmt"
	."time"
	."encoding/json"
	."../fileManager"
)

const ELEV_COUNT = 3
const FLOOR_COUNT = 4
const MASTER_INIT_IP = "129.241.187.92"
const PORT = ":20001"

type Message struct {
	Type string
	Content string
	Floor int
	Value int
	From int
	To int
}

func Network(networkReceive chan Message, networkSend chan Message) {
	recievedChannel := make(chan Message)
	go listen(recievedChannel)
	message := Message{}
	appendable := true
	destination := 0
	readableFile := false
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
		message.Type = "broadcast"
		go startBroadcast(message, IPlist)
	}

	if /* FILE PRESENT */ destination == -2 {
		readableFile = true
		//FILL IN IP-adresses
	}

	for{
		select{
			case message = <- recievedChannel:
				switch{
				case message.Type == "newElev":
					IPlist = append(IPlist, message.Content)
					WriteIP(message.Content)
					message.Type = "addElev"
					message.To = len(IPlist) - 1
					if len(IPlist) > 2 && IPlist[0] == IPlist[1] {
						for i := 1; i < len(IPlist); i++ {
							message.Content = IPlist[i]
							networkSend <- message
						}
					}
				case message.Type == "addElev":
					for i := 1; i < len(IPlist); i++ {
						if IPlist[i] == message.Content {
							appendable = false
						}
					}
					if appendable {
						IPlist = append(IPlist, message.Content)
					}
				case message.Type == "elevOffline":
					IPlist = append(IPlist[:message.From], IPlist[message.From+1:]...)
				}
				networkReceive <- message

			case message = <- networkSend:
				switch{
				case message.Type == "imAlive":
					destination = -1
				case message.Type == "newElev":
					message.Content = IPlist[0]
					if !readableFile {
						destination = -1
					} else {
						destination = 0
					}
				case message.Type == "addElev":
					destination = message.To
				case message.Type == "IPpackage":

				case message.Type == "newOrder":
					destination = 0
				case message.Type == "deleteOrder":
					destination = 0
				case message.Type == "newTarget":
					destination = message.To
				case message.Type == "stateUpdate":
					destination = 0
				case message.Type == "floorReached":
					for i := 1; i < len(IPlist); i++ {
						if IPlist[0] == IPlist[i] {
							destination = i
							break
						}
					}
				}
				if destination == 0 {
					for i := 1; i < len(IPlist); i++ {
						go send(message, i, IPlist)
					}
				} else {
					go send(message, destination, IPlist)
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

func send(message Message, destination int, IPlist[] string) {
	recipient := ""
	switch{
	case destination == -1:
		recipient = MASTER_INIT_IP
	case destination > 0:
		recipient = IPlist[destination]
	}
	connection, error := net.Dial("tcp", recipient+ PORT)
	if error != nil {
		Println("Send connection error: ", error)

		/*
		ELEVATOR OFFLINE

		message.Type = "elevOffline"
		message.From = destination
		for i := 1; i < len(IPlist); i++ {
			if i != destination {
				Println("Sending elevOffline message to rank", i)
				go send(message, i, IPlist)
			}
		}
		*/
	}
	defer connection.Close()
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

func startBroadcast(message Message, IPlist[] string) {
	Sleep(400 * Millisecond)
	go send(message, -1, IPlist)
}
