package packets

import (
	"encoding/binary"
	"little/processor"
	"little/types"
)

/*
0000   00 00 28 1f 00 00 00 3b
0010   01 00 10 b1 5a 03 00 89 18 04 00 b1 5a 03 00 51
0020   19 04 00
*/

type ActionPosition struct {
	X types.HardUInt32
	Y types.HardUInt32
}

func (po *ActionPosition) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Read
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode, binary.LittleEndian)
	p.UInt32(buf, &po.X)
	p.UInt32(buf, &po.Y)
}

type BeginActionRequest struct {
	WorldID    types.HardUInt32
	MoveCount  types.HardUInt32
	ActionType types.HardUInt8
	PointSize  types.HardUInt16
	From       ActionPosition
	To         ActionPosition
}

func (b BeginActionRequest) Opcode() uint16 {
	return 6
}

func (b *BeginActionRequest) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Read
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt32(buf, &b.WorldID)
	p.UInt32(buf, &b.MoveCount)
	p.UInt8(buf, &b.ActionType)
	p.UInt16(buf, &b.PointSize)

	(&b.From).Process(buf, defaultMode)
	(&b.To).Process(buf, defaultMode)
}

func (b *BeginActionRequest) Print() {
	DebugPrint(b)
}

type BeginActionResponse struct {
	WorldID    types.HardUInt32
	MoveCount  types.HardUInt32
	ActionType types.HardUInt8
	State      types.HardUInt16
	StopState  types.HardUInt16 // Only if State != ON
	PointSize  types.HardUInt16
	From       ActionPosition
	To         ActionPosition
}

func (b BeginActionResponse) Opcode() uint16 {
	return 508
}

func (b *BeginActionResponse) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt32(buf, &b.WorldID)
	p.UInt32(buf, &b.MoveCount)
	p.UInt8(buf, &b.ActionType)
	p.UInt16(buf, &b.State)

	if b.State.Value != types.MSTATE_ON {
		p.UInt16(buf, &b.StopState)
	}

	p.UInt16(buf, &b.PointSize)

	(&b.From).Process(buf, defaultMode)
	(&b.To).Process(buf, defaultMode)
}

func (b *BeginActionResponse) Print() {
	DebugPrint(b)
}

// BeginActionEndRequest хочет пакет с стейтом остановки шага в обратку
type BeginActionEndRequest struct {
}

func (b BeginActionEndRequest) Opcode() uint16 {
	return 7
}

func (b *BeginActionEndRequest) Process(buf *[]byte, mode ...processor.Mode) {

}

func (b *BeginActionEndRequest) Print() {
	DebugPrint(b)
}
