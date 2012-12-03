package helper

import "strconv"
import "strings"
import "fmt"

type WatermarkArray struct {
	LCval []int
}

func NewWMA(wmaCap int) (WA *WatermarkArray) {
	WA = &WatermarkArray{LCval: make([]int, wmaCap)}
	return
}

func (WA *WatermarkArray) FullyAck() int {
	foo := 0
	for _, elem := range WA.LCval {
		if foo < elem {
			foo = elem
		}
	}
        fmt.Print(foo)
        fmt.Println(" Has been fully Ack'd")
	return foo
}

func (WA *WatermarkArray) Update(msg Message) {
	if WA.LCval[msg.Sender] < msg.Timestamp {
		WA.LCval[msg.Sender] = msg.Timestamp
	}
}

func (WA *WatermarkArray) String() string {
	var result string
	var foo []string
	for pos, val := range WA.LCval {
		foo[pos] = "Host " + strconv.Itoa(pos) + " LC:" + strconv.Itoa(val)
	}
	result = strings.Join(foo, "\n")
	return result
}
