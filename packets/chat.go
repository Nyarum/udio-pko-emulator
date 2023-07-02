package packets

import (
	"little/processor"
	"little/types"
)

type SayRequest struct {
	Msg types.HardString
}

func (s SayRequest) Opcode() uint16 {
	return 1
}

func (s *SayRequest) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Read
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.String(buf, &s.Msg)
}

func (s *SayRequest) Print() {
	DebugPrint(s)
}

type SayResponse struct {
	ID  types.HardUInt32
	Msg types.HardString
}

func (s SayResponse) Opcode() uint16 {
	return 501
}

func (s *SayResponse) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt32(buf, &s.ID)
	p.String(buf, &s.Msg)
}

func (s *SayResponse) Print() {
	DebugPrint(s)
}
