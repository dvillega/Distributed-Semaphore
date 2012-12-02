package helper

type WatermarkArray struct {
    LCval  []int
}

func NewWMA(wmaCap int) (WA *WatermarkArray) {
    WA = &WatermarkArray{LCval:make([]int,wmaCap)}
    return
}

func (WA *WatermarkArray) FullyAck() int {
    foo := 0
    for _,elem := range WA.LCval{
        if foo < elem {
            foo = elem
        }
    }
    return foo
}

func (WA *WatermarkArray) Update(msg Message) {
    if WA.LCval[msg.Sender] < msg.Timestamp {
        WA.LCval[msg.Sender] = msg.Timestamp
    }
}
