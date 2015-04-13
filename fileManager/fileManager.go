package fileManager

import (
	"os"
	."fmt"
	"bufio"
    "strings"
    "strconv"
    "io/ioutil"    
)

func fileManager() {
reader := bufio.NewReader(os.Stdin)
Print("Enter text: ")
text,_ := reader.ReadString('\n')
	switch (text){
		case "Read\n":
			inside := ReadInside()
			Println(inside)
		case "Write\n":
			WriteInside(3, 1)
		case "Readip\n":
			IP := ReadIP()
			Println(IP)
		case "WriteIP\n":
			WriteIP("2324")
	}
}

func ReadInside() []int{

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
					strOrders = strOrders + strconv.Itoa(inside[i + 1])
					if i < len(inside) - 1 {
						strOrders += "\t"
					}
				}

		    	writing, err := f.WriteString(strOrders)
		    	if err != nil {
				    Println("Error while writing to file!")
				}
		    	Printf("wrote %d bytes\n", writing)

		    	// Issue a `Sync` to flush writes to stable storage.
		    	f.Sync()

		    	return inside
		}
		strOrders := strings.Split(string(content), "\t")
		inside[0] = 0
		for i := range strOrders{
			inside[i + 1],_ = strconv.Atoi(strOrders[i])
		}
		return inside
}


func WriteInside(floor int, value int){

	tempInside := ReadInside()

	directory, err := os.Getwd()
	    if err != nil {
	        Println(err)
	        os.Exit(1)
	    }

	f, err := os.Create(directory + "/insideOrders.txt")
    	if err != nil {
		    Println("Could not open file location!")
		}

		tempInside[floor - 1] = value;

		strOrders := ""

		for i := range tempInside{
			strOrders += strconv.Itoa(tempInside[i])
			if i < len(tempInside) - 1 {
				strOrders += "\t"
			}
		}

    	writing, err := f.WriteString(strOrders)
    	if err != nil {
		    Println("Error while writing to file!")
		}
    	Printf("wrote %d bytes\n", writing)

    	// Issue a `Sync` to flush writes to stable storage.
    	f.Sync()
}


func ReadIP() []string{


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

		}	else {

			IPs = append(IPs, "")
			IPs = append(IPs, strings.Split(string(content), "\n")...)

			return IPs
		}
}


func WriteIP(IP string){

	tempIPs := ReadIP()

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

    	writing, err := f.WriteString(strIPs)
    	if err != nil {
		    Println("Error while writing to file!")
		}
    	Printf("wrote %d bytes\n", writing)

    	// Issue a `Sync` to flush writes to stable storage.
    	f.Sync()
}
