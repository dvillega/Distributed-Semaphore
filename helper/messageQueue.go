package helper

import (
    "sort"
    "strconv"
    "strings"
)

type MessageQueue struct {
    MQueue      []Message 
}

type Message struct {
        Sender      int
        Kind        string
        Timestamp   int
}

func (msg *Message) String() string {
    foo := strconv.Itoa(msg.Sender)
    bar := strconv.Itoa(msg.Timestamp)
    val := "Sender:" + foo + " " + msg.Kind + " TS: " + bar
    return val
}

func (s MessageQueue) Len() int      { return len(s.MQueue)}
func (s MessageQueue) Swap(i, j int) { s.MQueue[i],s.MQueue[j] = s.MQueue[j],s.MQueue[i]}
func (s MessageQueue) Less(i, j int) bool {
    if s.MQueue[i].Timestamp == s.MQueue[j].Timestamp {
        return s.MQueue[i].Sender < s.MQueue[j].Sender
    } 
    return s.MQueue[i].Timestamp < s.MQueue[j].Timestamp
}

// Returns a new Message Queue which is initialized to 0 with capacity qCap
func NewMessageQueue(qCap int) (MQ *MessageQueue) {
    // Size 0 Capacity int index=0
    MQ = &MessageQueue{make([]Message,0,qCap)}
    return
}

// Adds the message to the queue - sorts afterwards
// based on Timestamp with Sender as tiebreaker
func (base *MessageQueue) Append(msg Message) {
    base.MQueue = append(base.MQueue, msg)
    sort.Sort(base)
}

// Returns the message at location pos
func (base *MessageQueue) Get(pos int) (msg Message){
    val := base.MQueue[pos]
    return val
}

// Returns Size of the queue
func (base *MessageQueue) Size() (size int) {
    return len(base.MQueue)
}

// Removes Message at position pos
func (base *MessageQueue) RemovePos(pos int) {
    base.MQueue = append(base.MQueue[:pos], base.MQueue[pos+1:]...)
}

// Removes message if it is contained in the queue
func (base *MessageQueue) Remove(msg Message) {
    pos := base.findMessage(msg)
    if pos >= 0 {
        base.RemovePos(pos)
    }
}

// Util func
func (base *MessageQueue) findMessage(msg Message) (pos int) { 
    for pos, elem := range base.MQueue{ 
        if elem == msg {
            return pos
        } // found
    }
    return -1 // not found
} 

// Returns true if msg is in the queue
func (base *MessageQueue) Contains(msg Message) bool {
    return (base.findMessage(msg) >= 0)
}

// Returns Fully Acknowledged Messages 
// @param lc int - Logical Clock value from which to base the return on
// @returns []Message which are Fully Acknowledged from lc
func (base *MessageQueue) FullyAck(lc int) ([]Message) {
    end := len(base.MQueue) - 1
    foo := -1
    for pos, _ := range base.MQueue {
        if base.MQueue[end-pos].Timestamp <= lc {
            foo = end-pos
            break
        }
    }
    if foo != -1 {
        return base.MQueue[:foo]
    }
    return nil
}

// String Return
func (base *MessageQueue) String() string {
    foo := make([]string,0,len(base.MQueue))
    for _,val := range base.MQueue {
        foo = append(foo,val.String())
    }
    return strings.Join(foo,"\n")
}


/* Testing Main

func main() {
    msg := Message{Sender:0, Kind:"VOP", Timestamp:0}
    msg1 := Message{Sender:1, Kind:"VOP", Timestamp:1}
    msg2 := Message{Sender:0, Kind:"VOP", Timestamp:2}
    msg3 := Message{Sender:1, Kind:"VOP", Timestamp:3}
    msg4 := Message{Sender:0, Kind:"VOP", Timestamp:4}

    msg5 := Message{Sender:1, Kind:"POP", Timestamp:4}

    mq := NewMessageQueue(10)

    mq.Append(msg)
    fmt.Println("One")
    fmt.Println(mq)
    mq.Append(msg1)
    mq.Append(msg4)
    mq.Append(msg2)
    mq.Append(msg5)
    fmt.Println("Two")
    fmt.Println(mq)
    mq.Append(msg3)
    fmt.Println("Three")
    fmt.Println(mq)
}
*/
