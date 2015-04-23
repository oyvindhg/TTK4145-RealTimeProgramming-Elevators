package fileManager

import (
	"os"
	."fmt"
	//"bufio"
    "strings"
    "strconv"
    "io/ioutil"
    ."../network"
)

func FileManager(fileOutChan chan Message, fileInChan chan Message) {
	for {
		select {
		case message := <- fileInChan:
			switch{
				case message.Type == "writeIP":
					writeIP(message.Content)
				case message.Type == "writeInside":
					writeInside(message.Floor, message.Value)
				case message.Type == "readIP":
					IPs := readIP()
					if message.Value < len(IPs) {
						message.Content = IPs[message.Value]
						fileOutChan <- message
					} else {
						message.Content = "noIP"
						fileOutChan <- message
					}
				case message.Type == "readInside":
					inside := readInside()
					if message.Floor < len(inside) {
						message.Value = inside[message.Floor]
						fileOutChan <- message
					} else {
						message.Value = -1
						fileOutChan <- message
					}
			}
		}
	}
}

func readInside() []int{

	directory, err := os.Getwd()
    if err != nil {
        Println(err)
        os.Exit(1)
    }

	content, err := ioutil.ReadFile(directory + "/insideOrders.txt")
	inside := []int{0,0,0,0,0}
	if err != nil {
    	Println("File corrupted (you might have lost the orders) or file not created. Stay tuned as we will make you a new one.")

		f, err := os.Create(directory + "/insideOrders.txt")
    	if err != nil {
		    Println("Could not open file location!")
		}
		strOrders := ""

		for i := range inside{
			strOrders = strOrders + strconv.Itoa(inside[i])
			if i < len(inside) - 1 {
				strOrders += "\t"
			}
		}

    	_, err = f.WriteString(strOrders)
    	if err != nil {
		    Println("Error while writing to file!")
		}

    	// Issue a `Sync` to flush writes to stable storage.
    	f.Sync()

    	return inside
	}
	strOrders := strings.Split(string(content), "\t")
	inside[0] = 0
	for i := range strOrders{
		inside[i],_ = strconv.Atoi(strOrders[i])
	}
	return inside
}


func writeInside(floor int, value int){

	tempInside := readInside()

	directory, err := os.Getwd()
    if err != nil {
        Println(err)
        os.Exit(1)
    }

	f, err := os.Create(directory + "/insideOrders.txt")
	if err != nil {
	    Println("Could not open file location!")
	}

	tempInside[floor] = value;

	strOrders := ""

	for i := range tempInside{
		strOrders += strconv.Itoa(tempInside[i])
		if i < len(tempInside) - 1 {
			strOrders += "\t"
		}
	}

	_, err = f.WriteString(strOrders)
	if err != nil {
	    Println("Error while writing to file!")
	}

	// Issue a `Sync` to flush writes to stable storage.
	f.Sync()
}


func readIP() []string{


	directory, err := os.Getwd()
    if err != nil {
        Println(err)
        os.Exit(1)
    }

	content, err := ioutil.ReadFile(directory + "/IP.txt")
	IPs := []string{}

	if err != nil {
	    Println("File does not exist. New file made \n")

		f, err := os.Create(directory + "/IP.txt")
    	if err != nil {
		    Println("Could not open file location!")
		}

    	// Issue a `Sync` to flush writes to stable storage.
    	f.Sync()

    	return IPs

	} else {

		IPs = append(IPs, "")
		IPs = append(IPs, strings.Split(string(content), "\n")...)

		return IPs
	}
}


func writeIP(IP string){

	tempIPs := readIP()

	directory, err := os.Getwd()
    if err != nil {
        Println(err)
        os.Exit(1)
    }

	f, err := os.Create(directory + "/IP.txt")
	if err != nil {
	    Println("Could not open file location!")
	}

	newIP := true

	for i := range tempIPs{
		if IP == tempIPs[i]{
			newIP = false
		}
	}

	if newIP {
		tempIPs = append(tempIPs, IP)
	}

	strIPs := ""

	for i := range tempIPs{
		strIPs += tempIPs[i]
		if i < len(tempIPs) - 1 && tempIPs[i] != ""{
			strIPs += "\n"
		}
		
	}

	_, err = f.WriteString(strIPs)
	if err != nil {
	    Println("Error while writing to file!")
	}

	// Issue a `Sync` to flush writes to stable storage.
	f.Sync()
}
