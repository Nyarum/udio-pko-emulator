package packets

import (
	"encoding/binary"
	"little/processor"
	"little/types"
)

type Auth struct {
	KeyLen        types.HardUInt16
	Key           types.HardBytes
	Login         types.HardString
	PasswordLen   types.HardUInt16
	Password      types.HardBytes
	MAC           types.HardString
	IsCheat       types.HardUInt16
	ClientVersion types.HardUInt16
}

func (a Auth) Opcode() uint16 {
	return 431
}

func (a *Auth) Process(buf *[]byte, mode ...processor.Mode) {
	p := processor.NewProcessor(processor.Read)
	p.UInt16(buf, &a.KeyLen)
	p.Bytes(buf, a.KeyLen.Value, &a.Key)
	p.String(buf, &a.Login)
	p.UInt16(buf, &a.PasswordLen)
	p.Bytes(buf, a.PasswordLen.Value, &a.Password)
	p.String(buf, &a.MAC)
	p.UInt16(buf, &a.IsCheat)
	p.UInt16(buf, &a.ClientVersion)
}

func (a *Auth) Print() {
	DebugPrint(a)
}

type AuthError struct {
	ErrorCode types.HardUInt16
}

func (a AuthError) Opcode() uint16 {
	return 931
}

func (a *AuthError) Process(buf *[]byte, mode ...processor.Mode) {
	p := processor.NewProcessor(processor.Write)
	p.UInt16(buf, &a.ErrorCode)

	// 1002 - password incorrect
}

func (a *AuthError) Print() {
	DebugPrint(a)
}

type InstAttr struct {
	ID    types.HardUInt16
	Value types.HardUInt16
}

func (i *InstAttr) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode, binary.LittleEndian)
	p.UInt16(buf, &i.ID)
	p.UInt16(buf, &i.Value)
}

type ItemAttr struct {
	Attr   types.HardUInt16
	IsInit types.HardBool
}

func (i *ItemAttr) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode, binary.LittleEndian)
	p.UInt16(buf, &i.Attr)
	p.Bool(buf, &i.IsInit)
}

type ItemGrid struct {
	ID        types.HardUInt16
	Num       types.HardUInt16
	Endure    [2]types.HardUInt16
	Energy    [2]types.HardUInt16
	ForgeLv   types.HardUInt8
	DBParams  [2]types.HardUInt32
	InstAttrs [5]InstAttr
	ItemAttrs [40]ItemAttr
	IsChange  types.HardBool
}

func (i *ItemGrid) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode, binary.LittleEndian)

	p.UInt16(buf, &i.ID)
	p.UInt16(buf, &i.Num)

	for k := range i.Endure {
		p.UInt16(buf, &i.Endure[k])
	}

	for k := range i.Energy {
		p.UInt16(buf, &i.Energy[k])
	}

	p.UInt8(buf, &i.ForgeLv)

	for k := range i.DBParams {
		p.UInt32(buf, &i.DBParams[k])
	}

	for k := range i.InstAttrs {
		i.InstAttrs[k].Process(buf, mode...)
	}

	for k := range i.ItemAttrs {
		i.ItemAttrs[k].Process(buf, mode...)
	}

	p.Bool(buf, &i.IsChange)
}

type Look struct {
	Ver       types.HardUInt16
	TypeID    types.HardUInt16
	ItemGrids [10]ItemGrid
	Hair      types.HardUInt16
}

func (c *Look) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode, binary.LittleEndian)
	p.UInt16(buf, &c.Ver)
	p.UInt16(buf, &c.TypeID)

	for k := range c.ItemGrids {
		c.ItemGrids[k].Process(buf, mode...)
	}

	p.UInt16(buf, &c.Hair)
}

type Character struct {
	IsActive types.HardBool
	Name     types.HardString
	Job      types.HardString
	Map      types.HardString
	Level    types.HardUInt16
	LookSize types.HardUInt16
	Look     Look
}

func (c *Character) Process(buf *[]byte, mode ...processor.Mode) {
	p := processor.NewProcessor(processor.Write)
	p.Bool(buf, &c.IsActive)

	if c.IsActive.Value {
		p.String(buf, &c.Name)
		p.String(buf, &c.Job)
		p.UInt16(buf, &c.Level)
		p.UInt16(buf, &c.LookSize)
		c.Look.Process(buf, mode...)
	}
}

type CharactersChoice struct {
	ErrorCode    types.HardUInt16
	KeyLen       types.HardUInt16
	Key          types.HardBytes
	CharacterLen types.HardUInt8
	Characters   []Character
	Pincode      types.HardUInt8
	Encryption   types.HardUInt32
	DWFlag       types.HardUInt32
}

func (c CharactersChoice) Opcode() uint16 {
	return 931
}

func (c *CharactersChoice) Process(buf *[]byte, mode ...processor.Mode) {
	c.basic()

	p := processor.NewProcessor(processor.Write)
	p.UInt16(buf, &c.ErrorCode)

	p.UInt16(buf, &c.KeyLen)
	p.Bytes(buf, c.KeyLen.Value, &c.Key)

	c.CharacterLen.Value = uint8(len(c.Characters))

	p.UInt8(buf, &c.CharacterLen)

	for k := range c.Characters {
		c.Characters[k].Process(buf, mode...)
	}

	p.UInt8(buf, &c.Pincode)
	p.UInt32(buf, &c.Encryption)
	p.UInt32(buf, &c.DWFlag)
}

func (c *CharactersChoice) Print() {
	DebugPrint(c)
}

func (c *CharactersChoice) basic() {
	c.Key = types.HardBytes{Value: []byte{0x7C, 0x35, 0x09, 0x19, 0xB2, 0x50, 0xD3, 0x49}}
	c.KeyLen = types.HardUInt16{Value: uint16(len(c.Key.Value))}
	c.DWFlag = types.HardUInt32{Value: 12820}
	c.Pincode = types.HardUInt8{Value: 1}
}

type CharacterCreate struct {
	Name     types.HardString
	Map      types.HardString
	LookSize types.HardUInt16
	Look     Look
}

func (c CharacterCreate) Opcode() uint16 {
	return 435
}

func (c *CharacterCreate) Process(buf *[]byte, mode ...processor.Mode) {
	p := processor.NewProcessor(processor.Read)
	p.String(buf, &c.Name)
	p.String(buf, &c.Map)
	p.UInt16(buf, &c.LookSize)
	c.Look.Process(buf, processor.Read)
}

func (c *CharacterCreate) Print() {
	DebugPrint(c)
}

type CharacterCreateReply struct {
	ErrorCode types.HardUInt16
}

func (c CharacterCreateReply) Opcode() uint16 {
	return 935
}

func (c *CharacterCreateReply) Process(buf *[]byte, mode ...processor.Mode) {
	p := processor.NewProcessor(processor.Write)
	p.UInt16(buf, &c.ErrorCode)
}

func (c *CharacterCreateReply) Print() {
	DebugPrint(c)
}
