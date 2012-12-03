package helper

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var msgQP *MessageQueue
var msgQV *MessageQueue
var wa *WatermarkArray
var setup bool

// Helper Process
func Handler(helpChan chan Message, myID int, hosts []string, okToUse chan int) {
	// Logical Clock
	LC, s := new(int), new(int)
	*LC, *s = 0, 0
	numHosts := len(hosts)
	fmt.Print("Number of Hosts:")
	fmt.Println(numHosts)

	service := os.Args[2]
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	check(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	check(err)

	msgQP = NewMessageQueue(100)
	msgQV = NewMessageQueue(100)
	wa = NewWMA(numHosts)

	sendConn := make([]*net.Conn, numHosts)
	rcvConn := make([]*net.Conn, numHosts)

	sendList := make([]*gob.Encoder, numHosts)
	rcvList := make([]*gob.Decoder, numHosts)

	go initRcv(rcvConn, rcvList, hosts, listener)
	go initSend(sendConn, sendList, hosts)

	for !setup {
		time.Sleep(time.Second)
	}
	go receive(rcvList, helpChan)

	for {
		receiveMsg(sendList, LC, myID, helpChan, okToUse, s)
		fmt.Print(myID)
		fmt.Print(" is done rcv: LC=")
		fmt.Println(*LC)
	}
}

// User Prompt
func Prompt() {
	fmt.Println("Welcome to DisSem Technologies")
	fmt.Println("Press 'p' for P()")
	fmt.Println("Press 'v' for V()")
	fmt.Println("Press 'q' to exit")
}

// Read all info from config 
func ReadLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// Error handling for i/o
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func broadcast(sendList []*gob.Encoder, msg Message) {
	fmt.Println("Broadcasting")
	for n := range sendList {
		sendList[n].Encode(msg)
	}
}

func receiveMsg(sendList []*gob.Encoder, LC *int, myID int, helpChan chan Message, okToUse chan int, s *int) {
	foo := <-helpChan
	fmt.Print(foo)
	fmt.Println(" is being Received")
	if foo.Timestamp+1 > *LC {
		*LC = foo.Timestamp + 1
	}
	*LC++
	bar := foo.Kind
	switch bar {
	case "reqP":
		reqP := Message{Sender: myID, Kind: "POP", Timestamp: *LC}
		broadcast(sendList, reqP)
		*LC++
	case "reqV":
		reqV := Message{Sender: myID, Kind: "VOP", Timestamp: *LC}
		broadcast(sendList, reqV)
		*LC++
	case "POP":
		ack := Message{Sender: myID, Kind: "ACK", Timestamp: *LC}
		msgQP.Append(foo)
		broadcast(sendList, ack)
		*LC++
	case "VOP":
		ack := Message{Sender: myID, Kind: "ACK", Timestamp: *LC}
		msgQV.Append(foo)
		broadcast(sendList, ack)
		*LC++
	case "ACK":
		handleACK(foo, okToUse, myID, LC, s)
	}
	fmt.Print("msgQV ")
	fmt.Println(msgQV)
	fmt.Print("msgQP ")
	fmt.Println(msgQP)
}

func receive(rcvList []*gob.Decoder, helpChan chan Message) {
	for ndx, _ := range rcvList {
		go receivePerConn(rcvList[ndx], helpChan)
        }
}

func receivePerConn(receiver *gob.Decoder, helpChan chan Message) {
    for {
        var foo Message
        err := receiver.Decode(&foo)
        if err == nil {
            fmt.Print("Go Func Received:")
            fmt.Println(foo)
            helpChan <- foo
        }
    }
}

func handleACK(msg Message, okToUse chan int, myID int, LC *int, s *int) {
	wa.Update(msg)
        fmt.Print("Handling ACK of")
        fmt.Println(msg)
        fmt.Print("S val is :")
        fmt.Println(*s)
	FAVmsg := msgQV.FullyAck(wa.FullyAck())
	for _, val := range FAVmsg {
                fmt.Print("Removing ")
                fmt.Println(val)
		msgQV.Remove(val)
		*s = *s + 1
                fmt.Print("S is now:")
                fmt.Println(*s)
	}
	FAPmsg := msgQP.FullyAck(wa.FullyAck())
	for _, val := range FAPmsg {
                fmt.Print("Fully Ack'd P:")
                fmt.Println(val)
                fmt.Print("S is now:")
                fmt.Println(*s)
		if *s > 0 {
                        fmt.Print("Removing ")
                        fmt.Println(val)
			msgQP.Remove(val)
			*s = *s - 1
                        fmt.Print("S is now:")
                        fmt.Println(*s)
			if val.Sender == myID {
                                fmt.Print(myID)
                                fmt.Println(" is now okToUse")
				okToUse <- *LC
				*LC++
			}
		}
	}
}

func initRcv(rcvConn []*net.Conn, rcvList []*gob.Decoder, hosts []string, listener net.Listener) {
	for n, _ := range hosts {
		var err error
		val, err := listener.Accept()
		check(err)
		rcvConn[n] = &val
		rcvList[n] = gob.NewDecoder(*rcvConn[n])
		check(err)
	}
}

func initSend(sendConn []*net.Conn, sendList []*gob.Encoder, hosts []string) {
	for n, _ := range hosts {
		var err error
		var val net.Conn
		val, err = net.Dial("tcp", hosts[n])
		for err != nil {
			val, err = net.Dial("tcp", hosts[n])
		}
		sendConn[n] = &val
		fmt.Println("Connected to Host:" + hosts[n])
		sendList[n] = gob.NewEncoder(*sendConn[n])
		check(err)
	}
	setup = true
}
