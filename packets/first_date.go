package packets

import (
	"fmt"
	"little/processor"
	"little/types"
	"time"
)

type Date struct {
	Time types.HardString
}

func (d Date) Opcode() uint16 {
	return 940
}

func (d *Date) Process(buf *[]byte, mode ...processor.Mode) {
	d.Time.Value = d.getCurrentTime()

	p := processor.NewProcessor(processor.Write)
	p.String(buf, &d.Time)
}

func (d *Date) getCurrentTime() string {
	timeNow := time.Now()

	return fmt.Sprintf("[%02d-%02d %02d:%02d:%02d:%03d]", timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second(), timeNow.Nanosecond()/1000000)
}

func (d *Date) Print() {
	DebugPrint(d)
}
