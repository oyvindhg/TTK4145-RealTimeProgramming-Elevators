package network

import (
	."encoding/json"
	."fmt"
	"net"
)

const ELEV_COUNT = 3
const FLOOR_COUNT = 4
const MASTER_INIT_IP = "129.241.187.148"
const PORT = ":20015"

type Message struct {
	Type string
	Content string
	Floor int
	Value int
	From int
	To int
}

func Network(networkReceive chan Message, networkSend chan Message) {
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
		networkReceive <- message
	}

	recievedChannel := make(chan Message)
	go listen(recievedChannel)

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
	switch{
	case destination == -1:
		connection, error := net.Dial("tcp", MASTER_INIT_IP+ PORT)
	case destination > 0:
		connection, error := net.Dial("tcp", IP[destination]+ PORT)
	}

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
	message := Message{"nil", "nil", "nil", "nil", 0, false, 0, 0, "nil", "nil"}
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