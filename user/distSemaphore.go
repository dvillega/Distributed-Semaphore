package main

import (
	"UNO6401/helper"
	"fmt"
	"os"
        "time"
)

var helpChan chan helper.Message
var okToUse chan int
var userlc int
var myID int

func main() {
	helpChan = make(chan helper.Message, 20)
	okToUse = make(chan int)
	if len(os.Args) != 3 {
		fmt.Println("Usage: " + os.Args[0] + " fileToRead PortListen")
		os.Exit(1)
	}

	hosts, err := helper.ReadLines(os.Args[1])
	if err != nil {
		fmt.Println("Fatal Err:" + err.Error())
		os.Exit(1)
	}
	for pos, val := range hosts {
		if val == os.Args[2] {
			myID = pos
		}
	}

	check(err)

        var trash string
        fmt.Scanf("%s",&trash)

	go helper.Handler(helpChan, myID, hosts, okToUse)
        time.Sleep(time.Second * 3)

	userlc := 0

	for {
		helper.Prompt()
		var cmd string
		fmt.Scanf("%s", &cmd)
		if cmd == "v" || cmd == "V" {
			fmt.Println("V Command Issued")
			msg := helper.Message{Sender: myID, Kind: "reqV", Timestamp: userlc}
			helpChan <- msg
			userlc++
		} else if cmd == "p" || cmd == "P" {
			fmt.Println("P Command Issued")
			msg := helper.Message{Sender: myID, Kind: "reqP", Timestamp: userlc}
			userlc++
			helpChan <- msg
			fmt.Println("Waiting on P")
			ts := <-okToUse
			if ts+1 > userlc {
				userlc = ts + 1
			}
			userlc++
			fmt.Println("Done Waiting on P")
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
