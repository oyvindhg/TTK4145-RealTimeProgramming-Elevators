package network

import (
	."encoding/json"
	."time"
	."fmt"
	"net"
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
	destination := 0
	readableFile := false
	IP := make([]string, 1, ELEV_COUNT + 1)
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		Println("Address error: ", err)
	}
	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				IP[0] = ipnet.IP.String()
			}
		}
	}

	if IP[0] == MASTER_INIT_IP {
		message.Type = "broadcast"
		go startBroadcast(message, IP)
	}

	if /* FILE PRESENT */ destination == -2 {
		readableFile = true
		//FILL IN IP-adresses
	}

	for{
		select{
			case message = <- recievedChannel:
				switch{
				case message.Type == "newID":
					IP = append(IP, message.Content)
				}
				networkReceive <- message


			case message = <- networkSend:
				switch{
				case message.Type == "imAlive":
					destination = -1
				case message.Type == "newID":
					if !readableFile {
						destination = -1
					} else {
						destination = 0
					}
					message.Content = IP[0]
				}
				if destination == 0 {
					for i := 1; i < len(IP); i++ {
						go send(message, i, IP)
					}
				} else {
					go send(message, destination, IP)
				}
				
		}
	}
}

func listen(recievedChannel chan Message) {
	listener, error := net.Listen("tcp", PORT)
	if error != nil {
		Println(error)
	}
	defer listener.Close()
	for {
		connection, error := listener.Accept()
		if error != nil {
			Println(error)
		}
		go receive(connection, recievedChannel)
	}
}

func send(message Message, destination int, IP[] string) {
	recipient := ""
	switch{
	case destination == -1:
		recipient = MASTER_INIT_IP
	case destination > 0:
		recipient = IP[destination]
	}
	connection, error := net.Dial("tcp", recipient+ PORT)
	if error != nil {
		Println(error)
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
		Println(error)
	}
	err := Unmarshal(buffer[:length], &message)
	
	if err != nil {
		Println("Receive error: ", err)
	}
	recievedChannel <- message
}

func startBroadcast(message Message, IP[] string) {
	Sleep(400 * Millisecond)
	go send(message, -1, IP)
}
