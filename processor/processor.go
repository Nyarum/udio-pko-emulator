package processor

import (
	"encoding/binary"
	"little/types"
)

type Mode string

const (
	Write Mode = "write"
	Read  Mode = "read"
)

type ByteOrder interface {
	binary.AppendByteOrder
	binary.ByteOrder
}

type Processor struct {
	mode  Mode
	order ByteOrder
}

func NewProcessor(mode Mode, order ...ByteOrder) *Processor {
	p := &Processor{
		mode:  mode,
		order: binary.BigEndian,
	}

	if len(order) != 0 {
		p.order = order[0]
	}

	return p
}

func (p *Processor) String(buf *[]byte, s *types.HardString, order ...ByteOrder) {
	if len(order) != 0 {
		p.order = order[0]
	}

	switch p.mode {
	case Read:
		ln := p.order.Uint16((*buf)[:2]) + 2

		s.Value = string((*buf)[2 : ln-1])
		s.Buf = (*buf)[:ln]
		*buf = (*buf)[ln:]
	case Write:
		newBuf := p.writeString(buf, s.Value, order...)

		copyBuf := make([]byte, len(newBuf))
		copy(copyBuf, newBuf)

		s.Buf = copyBuf
		*buf = append(*buf, newBuf...)
	}
}

func (p *Processor) writeString(buf *[]byte, s string, order ...ByteOrder) (newBuf []byte) {
	newBuf = make([]byte, 2)
	p.order.PutUint16(newBuf, uint16(len(s)+1))
	newBuf = append(newBuf, s...)
	newBuf = append(newBuf, 0x00)

	return
}

func (p *Processor) UInt16(buf *[]byte, v *types.HardUInt16, order ...ByteOrder) {
	if len(order) != 0 {
		p.order = order[0]
	}

	switch p.mode {
	case Read:
		v.Buf = (*buf)[:2]
		v.Value = p.order.Uint16(v.Buf)
		*buf = (*buf)[2:]
	case Write:
		newBuf := make([]byte, 2)
		p.order.PutUint16(newBuf, v.Value)

		copyBuf := make([]byte, len(newBuf))
		copy(copyBuf, newBuf)

		v.Buf = copyBuf

		*buf = append(*buf, newBuf...)
	}
}

func (p *Processor) Bytes(buf *[]byte, ln uint16, v *types.HardBytes) {
	switch p.mode {
	case Read:
		v.Buf = (*buf)[:ln]
		v.Value = v.Buf
		*buf = (*buf)[ln:]
	case Write:
		copyBuf := make([]byte, len(v.Value))
		copy(copyBuf, v.Value)

		v.Buf = copyBuf
		*buf = append(*buf, v.Value...)
	}
}

func (p *Processor) Bool(buf *[]byte, v *types.HardBool) {
	switch p.mode {
	case Read:
		v.Buf = []byte{(*buf)[0]}

		boolValue := false
		if v.Buf[0] == 1 {
			boolValue = true
		}

		v.Value = boolValue
		*buf = (*buf)[1:]
	case Write:
		intValue := uint8(0)
		if v.Value {
			intValue = 1
		}

		copyBuf := make([]byte, 1)
		copy(copyBuf, []byte{intValue})

		v.Buf = copyBuf
		*buf = append(*buf, intValue)
	}
}

func (p *Processor) UInt8(buf *[]byte, v *types.HardUInt8) {
	switch p.mode {
	case Read:
		v.Buf = []byte{(*buf)[0]}
		v.Value = v.Buf[0]
		*buf = (*buf)[1:]
	case Write:
		copyBuf := make([]byte, 1)
		copy(copyBuf, []byte{v.Value})

		v.Buf = copyBuf
		*buf = append(*buf, v.Value)
	}
}

func (p *Processor) UInt32(buf *[]byte, v *types.HardUInt32, order ...ByteOrder) {
	if len(order) != 0 {
		p.order = order[0]
	}

	switch p.mode {
	case Read:
		v.Buf = (*buf)[:4]
		v.Value = p.order.Uint32(v.Buf)
		*buf = (*buf)[4:]
	case Write:
		newBuf := make([]byte, 4)
		p.order.PutUint32(newBuf, v.Value)

		copyBuf := make([]byte, len(newBuf))
		copy(copyBuf, newBuf)

		v.Buf = copyBuf
		*buf = append(*buf, newBuf...)
	}
}
