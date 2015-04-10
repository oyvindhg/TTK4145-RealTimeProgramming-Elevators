package network

import (
	."encoding/json"
	."fmt"
	"net"
)

const PORT = ":20015"

type Message struct {
	RecipientID string
	SenderID string	
	Content string
	Command string
	ElevNumber int
	Online bool
	Rank int
	FloorNumber int
	ButtonType string
	State string
}

func Network(networkReceive chan Message, networkSend chan Message) {
	receivedChannel := make(chan Message)
	go listen(receivedChannel)
	for{
		select{
			case receivedMessage := <- receivedChannel:
				networkReceive <- receivedMessage

			case sendMessage := <- networkSend:
				go send(sendMessage)
		}
	}
}

func listen(receivedChannel chan Message) {
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
		go receive(connection, receivedChannel)
	}
}

func send(message Message) {

	connection, error := net.Dial("tcp", message.RecipientID)

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

func receive(connection net.Conn, receivedChannel chan Message) {
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
	receivedChannel <- message
}