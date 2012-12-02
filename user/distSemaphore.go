package main

import (
    "fmt"
    "os"
    "strconv"
    "time"
    "UNO6401/helper" 
)


var helpChan chan helper.Message
var okToUse chan bool

func main() {
    helpChan = make(chan helper.Message, 10)
    okToUse = make(chan bool)
    if len(os.Args) != 4 {
        fmt.Println("Usage: " + os.Args[0] + " fileToRead PortListen ID")
        os.Exit(1)
    }
    
    myID, err := strconv.Atoi(os.Args[3])
    check(err)

    go helper.Handler(helpChan, okToUse, myID)
    <-okToUse
    time.Sleep(time.Second * 3)

    lc := 0

    for {
        helper.Prompt()
        var cmd string
        fmt.Scanf("%s",&cmd)
        if cmd == "v" || cmd == "V" {
            fmt.Println("V Command Issued")
            msg := helper.Message{Sender:myID, Kind:"V_op", Timestamp:lc}
            helpChan <- msg
            lc = lc + 1
        } else if cmd == "p" || cmd == "P" {
            fmt.Println("P Command Issued")
            msg := helper.Message{Sender:myID, Kind:"P_op", Timestamp:lc}
            lc = lc + 1
            helpChan <- msg
            <-okToUse
        } else if cmd == "q" || cmd == "Q" {
            fmt.Println("Thank you!")
            close(helpChan)
            return
        } else {
            fmt.Println("Unrecognized Command")
        }
    }

}


// Error handling for i/o
func check(e error) {
	if e != nil {
		panic(e)
	}
}
