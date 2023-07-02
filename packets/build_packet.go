package packets

import (
	"encoding/binary"
)

func BuildPacket(packet Packet) []byte {
	if packet == nil {
		return []byte{}
	}

	buf := []byte{}
	packet.Process(&buf)

	//packet.Print()

	newBuf := make([]byte, 0, len(buf)+8)

	newBuf = binary.BigEndian.AppendUint16(newBuf, uint16(len(buf)+8))
	newBuf = binary.BigEndian.AppendUint32(newBuf, 2147483648)
	newBuf = binary.BigEndian.AppendUint16(newBuf, packet.Opcode())

	newBuf = append(newBuf, buf...)

	return newBuf
}

type Header struct {
	Len    uint16
	ID     uint32
	Opcode uint16
}

func UnpackPacket(buf *[]byte) Header {
	h := Header{}

	h.Len = binary.BigEndian.Uint16((*buf)[0:2])
	h.ID = binary.LittleEndian.Uint32((*buf)[2:6])
	h.Opcode = binary.BigEndian.Uint16((*buf)[6:8])

	*buf = (*buf)[8:]

	return h
}
