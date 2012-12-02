package helper

import (
    "bufio"
    "bytes"
    "fmt"
    "encoding/gob"
    "io"
    "os"
    "net"
    "strconv"
    "time"
)

// Helper Process
func Handler(helpChan chan Message, okToUse chan bool, myID int) {
    helperLC := 0
    hosts, err := readLines(os.Args[1])
    if err != nil {
        fmt.Println("Fatal Err:" + err.Error())
        os.Exit(1)
    }
    numHosts := len(hosts)
    fmt.Print("Number of Hosts:")
    fmt.Println(numHosts)

    service := os.Args[2]
    tcpAddr, err := net.ResolveTCPAddr("tcp", service)
    check(err)
    listener, err := net.ListenTCP("tcp", tcpAddr)
    check(err)

    msgQueue := NewMessageQueue(100)

    sendConn := make([]net.Conn, numHosts - 1)
    rcvConn  := make([]net.Conn, numHosts - 1)

    sendList := make([]*gob.Encoder, numHosts - 1)
    rcvList := make([]*gob.Decoder, numHosts - 1)

    seenHost := false

    for n := 0; n < len(hosts); n++ {
        fmt.Println("Host" + strconv.Itoa(n) + " " + hosts[n])
        x := n
        if hosts[n] == service {
            seenHost = true
        }

        if seenHost {
            x = n - 1
        }
        if hosts[n] != service {
            sendConn[x], err = net.Dial("tcp", hosts[n])
            fmt.Println("Connecting to Host:" + hosts[n])
            for err != nil {
                sendConn[x], err = net.Dial("tcp", hosts[n])
            }
        } 
    }

    time.Sleep(time.Second * 3)

    for n := 0; n < len(hosts) - 1; n++ {
        rcvConn[n],err = listener.Accept()
        check(err)
        rcvList[n] = gob.NewDecoder(rcvConn[n])
        check(err)
        sendList[n] = gob.NewEncoder(sendConn[n])
        check(err)
    }
    okToUse <- true
    startMsg := Message{Sender:myID, Kind:"ACK", Timestamp:-1}
    broadcast(sendList, startMsg)
    for {
            receiveMsg(rcvList, sendList, msgQueue, helperLC, myID)
            fmt.Println("Done Receiving Messages")
            msg, ok:= <-helpChan
            fmt.Println("Host:" + service + " OK:" + strconv.FormatBool(ok))
            if ok {
                if msg.Kind == "P_op" {
                    broadcast(sendList, msg)
                    time.Sleep(time.Second * 2)
                    okToUse <- true
                } else {
                    broadcast(sendList, msg)
                }
            }
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
func readLines(path string) (lines []string, err error) {
    var (
        file *os.File
        part []byte
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

func receiveMsg(rcvList []*gob.Decoder, sendList []*gob.Encoder, msgQueue *MessageQueue, helperLC int, myID int) {
    ack := Message{Sender:myID, Kind:"ACK", Timestamp:helperLC}
    helperLC++
    fmt.Println("Receiving Messages")
    var foo Message
    for ndx,_ := range rcvList {
        fmt.Println(len(rcvList))
        fmt.Print(ndx)
        fmt.Println(" is being Received")
        err := rcvList[ndx].Decode(&foo)
        if err != nil || err != io.EOF {
            fmt.Println(foo)
            msgQueue.Append(foo)
        }
        fmt.Print(ndx)
        fmt.Println(" is done Received")
        fmt.Println("Sending ACK")
        sendList[ndx].Encode(ack)
    }
}
