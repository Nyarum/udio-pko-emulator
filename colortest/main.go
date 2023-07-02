package main

import (
	"fmt"
	"little/processor"
	"little/types"

	"github.com/fatih/structs"
	"github.com/gookit/color"
)

type Auth struct {
	Key           types.HardString
	Login         types.HardString
	Password      types.HardString
	MAC           types.HardString
	IsCheat       types.HardUInt16
	ClientVersion types.HardUInt16
}

func (a Auth) Opcode() uint16 {
	return 431
}

func (a *Auth) Process(buf *[]byte, mode ...processor.Mode) {
	p := processor.NewProcessor(processor.Read)

	p.String(buf, &a.Key)
	p.String(buf, &a.Login)
	p.String(buf, &a.Password)
	p.String(buf, &a.MAC)
	p.UInt16(buf, &a.IsCheat)
	p.UInt16(buf, &a.ClientVersion)
}

func renderHex(buf []byte, v interface{}, field string, t string) {
	color.Redp(fmt.Sprintf("%# x", buf))
	fmt.Print(" | ")
	color.Greenp(fmt.Sprintf("%v", v))
	fmt.Print(" | ")
	color.Cyanp(field)
	fmt.Print(" | ")
	color.Bluep(t)
	fmt.Print("\n")
}

func main() {
	authPacket := &Auth{}
	authPacket.Process(&[]byte{
		0x00, 0x05, 0x74, 0x65, 0x73, 0x74, 0x00,
		0x00, 0x05, 0x74, 0x65, 0x73, 0x74, 0x00,
		0x00, 0x05, 0x74, 0x65, 0x73, 0x74, 0x00,
		0x00, 0x05, 0x74, 0x65, 0x73, 0x74, 0x00,
		0x00, 0x02,
		0x00, 0x03,
	})

	f := structs.Fields(authPacket)

	for _, v := range f {
		pullContext := v.Value().(types.PullContext)

		getValue := pullContext.GetValue()

		switch typeConvValue := getValue.(type) {
		case string:
			renderHex(pullContext.GetBuf(), pullContext.GetValue(), v.Name(), pullContext.GetType()+fmt.Sprintf(" (%v)", len(typeConvValue)))
		default:
			renderHex(pullContext.GetBuf(), pullContext.GetValue(), v.Name(), pullContext.GetType())
		}
	}
}
