package packets

import (
	"little/processor"
)

type Packet interface {
	Opcode() uint16
	Process(buf *[]byte, mode ...processor.Mode)
	Print()
}

type Packets struct {
	Packets []Packet
	List    map[uint16]Packet
}

func NewPackets() *Packets {
	packets := &Packets{
		Packets: []Packet{
			&Date{},
			&Auth{},
			&CharactersChoice{},
			&CharacterCreate{},
			&CharacterCreateReply{},
			&EnterGameRequest{},
			&KitbagSyncRequest{},
			&SayRequest{},
			&BeginActionRequest{},
			&BeginActionEndRequest{},
		},
		List: make(map[uint16]Packet),
	}

	for _, p := range packets.Packets {
		packets.List[p.Opcode()] = p
	}

	return packets
}

func (p *Packets) GetPacket(opcode uint16) Packet {
	return p.List[opcode]
}
