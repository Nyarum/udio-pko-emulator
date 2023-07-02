package packets

import (
	"encoding/json"
	"fmt"
	"little/processor"
	"little/types"
)

type CharacterBoat struct {
	CharacterBase       CharacterBase
	CharacterAttribute  CharacterAttribute
	CharacterKitbag     CharacterKitbag
	CharacterSkillState CharacterSkillState
}

func (c *CharacterBoat) Process(buf *[]byte, mode ...processor.Mode) {
	c.CharacterBase.Process(buf, mode...)
	c.CharacterAttribute.Process(buf, mode...)
	c.CharacterKitbag.Process(buf, mode...)
	c.CharacterSkillState.Process(buf, mode...)
}

type Shortcut struct {
	Type   types.HardUInt8
	GridID types.HardUInt16
}

func (c *Shortcut) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &c.Type)
	p.UInt16(buf, &c.GridID)
}

type CharacterShortcut struct {
	Shortcuts [types.SHORT_CUT_NUM]Shortcut
}

func (c *CharacterShortcut) Process(buf *[]byte, mode ...processor.Mode) {
	for k := range c.Shortcuts {
		c.Shortcuts[k].Process(buf, mode...)
	}
}

type KitbagItem struct {
	GridID       types.HardUInt16
	ID           types.HardUInt16
	Num          types.HardUInt16
	Endure       [2]types.HardUInt16
	Energy       [2]types.HardUInt16
	ForgeLevel   types.HardUInt8
	IsValid      types.HardBool
	ItemDBInstID types.HardUInt32
	ItemDBForge  types.HardUInt32
	IsParams     types.HardBool
	InstAttrs    [5]InstAttr
}

func (k *KitbagItem) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &k.GridID)

	if k.GridID.Value == 65535 {
		return
	}

	p.UInt16(buf, &k.ID)

	if k.ID.Value > 0 {
		p.UInt16(buf, &k.Num)

		for v := range k.Endure {
			p.UInt16(buf, &k.Endure[v])
		}

		for v := range k.Energy {
			p.UInt16(buf, &k.Energy[v])
		}

		p.UInt8(buf, &k.ForgeLevel)
		p.Bool(buf, &k.IsValid)

		//if "item_info.type" == "boat" {
		if k.ID.Value == 3988 {
			p.UInt32(buf, &k.ItemDBInstID)
		}

		p.UInt32(buf, &k.ItemDBForge)

		//if "item_info.type" == "boat" {
		if k.ID.Value == 3988 {
			v := types.HardUInt32{}
			p.UInt32(buf, &v)
		} else {
			p.UInt32(buf, &k.ItemDBInstID)
		}

		p.Bool(buf, &k.IsParams)

		if k.IsParams.Value {
			for ki := range k.InstAttrs {
				k.InstAttrs[ki].Process(buf, mode...)
			}
		}
	}
}

type CharacterKitbag struct {
	Type      types.HardUInt8
	KeybagNum types.HardUInt16
	Items     []KitbagItem
}

func (c *CharacterKitbag) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &c.Type)

	if c.Type.Value == types.SYN_KITBAG_INIT {
		p.UInt16(buf, &c.KeybagNum)
	}

	c.KeybagNum.Value = c.KeybagNum.Value + 1

	if len(c.Items) != int(c.KeybagNum.Value) {
		c.Items = make([]KitbagItem, c.KeybagNum.Value)
	}

	for k := range c.Items {
		c.Items[k].Process(buf, mode...)
	}
}

type Attribute struct {
	ID    types.HardUInt8
	Value types.HardUInt32
}

func (a *Attribute) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &a.ID)

	//v := &a.Value
	/*
		if a.ID.Value == 44 {
			v.Value = 1200
		}
	*/

	p.UInt32(buf, &a.Value)
}

type CharacterAttribute struct {
	Type       types.HardUInt8
	Num        types.HardUInt16
	Attributes []Attribute
}

func (c *CharacterAttribute) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &c.Type)
	p.UInt16(buf, &c.Num)

	if len(c.Attributes) != int(c.Num.Value) {
		c.Attributes = make([]Attribute, c.Num.Value)
	}

	for k := range c.Attributes {
		c.Attributes[k].Process(buf, mode...)
	}
}

type SkillState struct {
	ID    types.HardUInt8
	Level types.HardUInt8
}

func (s *SkillState) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &s.ID)
	p.UInt8(buf, &s.Level)
}

type CharacterSkillState struct {
	StatesLen types.HardUInt8
	States    []SkillState
}

func (c *CharacterSkillState) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &c.StatesLen)

	if len(c.States) != int(c.StatesLen.Value) {
		c.States = make([]SkillState, c.StatesLen.Value)
	}

	for k := range c.States {
		c.States[k].Process(buf, mode...)
	}
}

type CharacterSkill struct {
	ID         types.HardUInt16
	State      types.HardUInt8
	Level      types.HardUInt8
	UseSP      types.HardUInt16
	UseEndure  types.HardUInt16
	UseEnergy  types.HardUInt16
	ResumeTime types.HardUInt32
	RangeType  types.HardUInt16
	Params     []types.HardUInt16 // ?
}

func (c *CharacterSkill) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &c.ID)
	p.UInt8(buf, &c.State)
	p.UInt8(buf, &c.Level)
	p.UInt16(buf, &c.UseSP)
	p.UInt16(buf, &c.UseEndure)
	p.UInt16(buf, &c.UseEnergy)
	p.UInt32(buf, &c.ResumeTime)
	p.UInt16(buf, &c.RangeType)

	for k := range c.Params {
		p.UInt16(buf, &c.Params[k])
	}
}

type CharacterSkillBag struct {
	SkillID  types.HardUInt16
	Type     types.HardUInt8
	SkillNum types.HardUInt16
	Skills   []CharacterSkill
}

func (c *CharacterSkillBag) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &c.SkillID)
	p.UInt8(buf, &c.Type)
	p.UInt16(buf, &c.SkillNum)

	if len(c.Skills) != int(c.SkillNum.Value) {
		c.Skills = make([]CharacterSkill, c.SkillNum.Value)
	}

	for k := range c.Skills {
		c.Skills[k].Process(buf, mode...)
	}
}

type CharacterAppendLook struct {
	LookID  types.HardUInt16
	IsValid types.HardUInt8
}

func (c *CharacterAppendLook) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &c.LookID)

	if c.LookID.Value != 0 {
		p.UInt8(buf, &c.IsValid)
	}
}

type CharacterPK struct {
	PkCtrl types.HardUInt8
}

func (c *CharacterPK) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &c.PkCtrl)
}

type CharacterLookBoat struct {
	PosID     types.HardUInt16
	BoatID    types.HardUInt16
	Header    types.HardUInt16
	Body      types.HardUInt16
	Engine    types.HardUInt16
	Cannon    types.HardUInt16
	Equipment types.HardUInt16
}

func (c *CharacterLookBoat) Process(buf *[]byte, synType types.HardUInt8, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &c.PosID)
	p.UInt16(buf, &c.BoatID)
	p.UInt16(buf, &c.Header)
	p.UInt16(buf, &c.Body)
	p.UInt16(buf, &c.Engine)
	p.UInt16(buf, &c.Cannon)
	p.UInt16(buf, &c.Equipment)
}

type CharacterLookItemSync struct {
	Endure  types.HardUInt16
	Energy  types.HardUInt16
	IsValid types.HardUInt8
}

func (c *CharacterLookItemSync) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &c.Endure)
	p.UInt16(buf, &c.Energy)
	p.UInt8(buf, &c.IsValid)
}

type CharacterLookItemShow struct {
	Num        types.HardUInt16
	Endure     [2]types.HardUInt16
	Energy     [2]types.HardUInt16
	ForgeLevel types.HardUInt8
	IsValid    types.HardUInt8
}

func (c *CharacterLookItemShow) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &c.Num)

	for k := range c.Endure {
		p.UInt16(buf, &c.Endure[k])
	}

	for k := range c.Energy {
		p.UInt16(buf, &c.Energy[k])
	}

	p.UInt8(buf, &c.ForgeLevel)
	p.UInt8(buf, &c.IsValid)
}

type CharacterLookItem struct {
	ID          types.HardUInt16
	ItemSync    CharacterLookItemSync
	ItemShow    CharacterLookItemShow
	IsDBParams  types.HardUInt8
	DBParams    [2]types.HardUInt32
	IsInstAttrs types.HardUInt8
	InstAttrs   [5]InstAttr
}

func (c *CharacterLookItem) Process(buf *[]byte, synType types.HardUInt8, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &c.ID)

	if c.ID.Value != 0 {
		if synType.Value == types.SynLookChange {
			c.ItemSync.Process(buf, mode...)
		} else {
			c.ItemShow.Process(buf, mode...)
			p.UInt8(buf, &c.IsDBParams)

			if c.IsDBParams.Value != 0 {
				for k := range c.DBParams {
					p.UInt32(buf, &c.DBParams[k])
				}

				p.UInt8(buf, &c.IsInstAttrs)

				if c.IsInstAttrs.Value != 0 {
					for k := range c.InstAttrs {
						c.InstAttrs[k].Process(buf, mode...)
					}
				}
			}
		}
	}
}

type CharacterLookHuman struct {
	HairID   types.HardUInt16
	ItemGrid [10]CharacterLookItem
}

func (c *CharacterLookHuman) Process(buf *[]byte, synType types.HardUInt8, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &c.HairID)

	for k := range c.ItemGrid {
		c.ItemGrid[k].Process(buf, synType, mode...)
	}

}

type CharacterLook struct {
	SynType   types.HardUInt8
	TypeID    types.HardUInt16
	IsBoat    types.HardUInt8
	LookBoat  CharacterLookBoat
	LookHuman CharacterLookHuman
}

func (c *CharacterLook) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &c.SynType)
	p.UInt16(buf, &c.TypeID)
	p.UInt8(buf, &c.IsBoat)

	if c.IsBoat.Value == 1 {
		(&c.LookBoat).Process(buf, c.SynType, mode...)
	} else {
		(&c.LookHuman).Process(buf, c.SynType, mode...)
	}
}

type EntityEvent struct {
	EntityID   types.HardUInt32
	EntityType types.HardUInt8
	EventID    types.HardUInt16
	EventName  types.HardString
}

func (e *EntityEvent) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt32(buf, &e.EntityID)
	p.UInt8(buf, &e.EntityType)
	p.UInt16(buf, &e.EventID)
	p.String(buf, &e.EventName)
}

type CharacterSide struct {
	SideID types.HardUInt8
}

func (c *CharacterSide) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &c.SideID)
}

type Position struct {
	X      types.HardUInt32
	Y      types.HardUInt32
	Radius types.HardUInt32
}

func (ps *Position) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt32(buf, &ps.X)
	p.UInt32(buf, &ps.Y)
	p.UInt32(buf, &ps.Radius)
}

type CharacterBase struct {
	ChaID        types.HardUInt32
	WorldID      types.HardUInt32
	CommID       types.HardUInt32
	CommName     types.HardString
	GmLvl        types.HardUInt8
	Handle       types.HardUInt32
	CtrlType     types.HardUInt8
	Name         types.HardString
	MottoName    types.HardString
	Icon         types.HardUInt16
	GuildID      types.HardUInt32
	GuildName    types.HardString
	GuildMotto   types.HardString
	StallName    types.HardString
	State        types.HardUInt16
	Position     Position
	Angle        types.HardUInt16
	TeamLeaderID types.HardUInt32
	Side         CharacterSide
	EntityEvent  EntityEvent
	Look         CharacterLook
	PkCtrl       CharacterPK
	LookAppend   [types.ESPE_KBGRID_NUM]CharacterAppendLook
}

func (c *CharacterBase) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt32(buf, &c.ChaID)
	p.UInt32(buf, &c.WorldID)
	p.UInt32(buf, &c.CommID)
	p.String(buf, &c.CommName)
	p.UInt8(buf, &c.GmLvl)
	p.UInt32(buf, &c.Handle)
	p.UInt8(buf, &c.CtrlType)
	p.String(buf, &c.Name)
	p.String(buf, &c.MottoName)
	p.UInt16(buf, &c.Icon)
	p.UInt32(buf, &c.GuildID)
	p.String(buf, &c.GuildName)
	p.String(buf, &c.GuildMotto)
	p.String(buf, &c.StallName)
	p.UInt16(buf, &c.State)
	c.Position.Process(buf, mode...)
	p.UInt16(buf, &c.Angle)
	p.UInt32(buf, &c.TeamLeaderID)
	c.Side.Process(buf, mode...)
	c.EntityEvent.Process(buf, mode...)
	c.Look.Process(buf, mode...)
	c.PkCtrl.Process(buf, mode...)

	for k := range c.LookAppend {
		c.LookAppend[k].Process(buf, mode...)
	}
}

type EnterGame struct {
	EnterRet            types.HardUInt16
	AutoLock            types.HardUInt8
	KitbagLock          types.HardUInt8
	EnterType           types.HardUInt8
	IsNewChar           types.HardUInt8
	MapName             types.HardString
	CanTeam             types.HardUInt8
	CharacterBase       CharacterBase
	CharacterSkillBag   CharacterSkillBag
	CharacterSkillState CharacterSkillState
	CharacterAttribute  CharacterAttribute
	CharacterKitbag     CharacterKitbag
	CharacterShortcut   CharacterShortcut
	BoatLen             types.HardUInt8
	CharacterBoats      []CharacterBoat
	ChaMainID           types.HardUInt32
}

func (e EnterGame) Opcode() uint16 {
	return 516
}

func (e *EnterGame) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt16(buf, &e.EnterRet)
	p.UInt8(buf, &e.AutoLock)
	p.UInt8(buf, &e.KitbagLock)
	p.UInt8(buf, &e.EnterType)
	p.UInt8(buf, &e.IsNewChar)
	p.String(buf, &e.MapName)
	p.UInt8(buf, &e.CanTeam)
	(&e.CharacterBase).Process(buf, mode...)
	(&e.CharacterSkillBag).Process(buf, mode...)
	(&e.CharacterSkillState).Process(buf, mode...)
	(&e.CharacterAttribute).Process(buf, mode...)

	(&e.CharacterKitbag).Process(buf, mode...)

	(&e.CharacterShortcut).Process(buf, mode...)
	e.BoatLen = types.HardUInt8{Value: 0}
	p.UInt8(buf, &e.BoatLen)

	if len(e.CharacterBoats) != int(e.BoatLen.Value) {
		e.CharacterBoats = make([]CharacterBoat, e.BoatLen.Value)
	}

	for k := range e.CharacterBoats {
		e.CharacterBoats[k].Process(buf, mode...)
	}

	p.UInt32(buf, &e.ChaMainID)
}

func (e *EnterGame) Print() {
	DebugPrint(e)
}

func (e *EnterGame) Basic() {
	err := json.Unmarshal([]byte(`{"EnterRet":{"Value":0,"Buf":"AAA="},"AutoLock":{"Value":0,"Buf":"AA=="},"KitbagLock":{"Value":0,"Buf":"AA=="},"EnterType":{"Value":1,"Buf":"AQ=="},"IsNewChar":{"Value":0,"Buf":"AA=="},"MapName":{"Value":"garner","Buf":"AAdnYXJuZXIA"},"CanTeam":{"Value":1,"Buf":"AQ=="},"CharacterBase":{"ChaID":{"Value":4,"Buf":"AAAABA=="},"WorldID":{"Value":10271,"Buf":"AAAoHw=="},"CommID":{"Value":10271,"Buf":"AAAoHw=="},"CommName":{"Value":"ingrysty","Buf":"AAlpbmdyeXN0eQA="},"GmLvl":{"Value":0,"Buf":"AA=="},"Handle":{"Value":33565845,"Buf":"AgAslQ=="},"CtrlType":{"Value":1,"Buf":"AQ=="},"Name":{"Value":"ingrysty","Buf":"AAlpbmdyeXN0eQA="},"MottoName":{"Value":"","Buf":"AAEA"},"Icon":{"Value":4,"Buf":"AAQ="},"GuildID":{"Value":0,"Buf":"AAAAAA=="},"GuildName":{"Value":"","Buf":"AAEA"},"GuildMotto":{"Value":"","Buf":"AAEA"},"StallName":{"Value":"","Buf":"AAEA"},"State":{"Value":1,"Buf":"AAE="},"Position":{"X":{"Value":217475,"Buf":"AANRgw=="},"Y":{"Value":278175,"Buf":"AAQ+nw=="},"Radius":{"Value":40,"Buf":"AAAAKA=="}},"Angle":{"Value":71,"Buf":"AEc="},"TeamLeaderID":{"Value":0,"Buf":"AAAAAA=="},"Side":{"SideID":{"Value":0,"Buf":"AA=="}},"EntityEvent":{"EntityID":{"Value":10271,"Buf":"AAAoHw=="},"EntityType":{"Value":1,"Buf":"AQ=="},"EventID":{"Value":0,"Buf":"AAA="},"EventName":{"Value":"","Buf":"AAEA"}},"Look":{"SynType":{"Value":0,"Buf":"AA=="},"TypeID":{"Value":4,"Buf":"AAQ="},"IsBoat":{"Value":0,"Buf":"AA=="},"LookBoat":{"PosID":{"Value":0,"Buf":null},"BoatID":{"Value":0,"Buf":null},"Header":{"Value":0,"Buf":null},"Body":{"Value":0,"Buf":null},"Engine":{"Value":0,"Buf":null},"Cannon":{"Value":0,"Buf":null},"Equipment":{"Value":0,"Buf":null}},"LookHuman":{"HairID":{"Value":2291,"Buf":"CPM="},"ItemGrid":[{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":2554,"Buf":"Cfo="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":"AAA="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":0,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":359,"Buf":"AWc="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":9999,"Buf":"Jw8="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":3040,"Buf":"C+A="},{"Value":3040,"Buf":"C+A="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":2000000000,"Buf":"dzWUAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":1,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":9216,"Buf":"ACQ="},"Value":{"Value":2304,"Buf":"AAk="}},{"ID":{"Value":12032,"Buf":"AC8="},"Value":{"Value":1024,"Buf":"AAQ="}},{"ID":{"Value":6912,"Buf":"ABs="},"Value":{"Value":768,"Buf":"AAM="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}}]},{"ID":{"Value":4309,"Buf":"ENU="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":9996,"Buf":"Jww="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":1,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":10240,"Buf":"ACg="},"Value":{"Value":1792,"Buf":"AAc="}},{"ID":{"Value":9216,"Buf":"ACQ="},"Value":{"Value":768,"Buf":"AAM="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}}]},{"ID":{"Value":4310,"Buf":"ENY="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":9988,"Buf":"JwQ="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":1,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":9984,"Buf":"ACc="},"Value":{"Value":512,"Buf":"AAI="}},{"ID":{"Value":9216,"Buf":"ACQ="},"Value":{"Value":768,"Buf":"AAM="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}}]},{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":9999,"Buf":"Jw8="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":"AAA="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":0,"Buf":"AA=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":0,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":104,"Buf":"AGg="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":9998,"Buf":"Jw4="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":3000,"Buf":"C7g="},{"Value":3000,"Buf":"C7g="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":2000000000,"Buf":"dzWUAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":1,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":8704,"Buf":"ACI="},"Value":{"Value":19456,"Buf":"AEw="}},{"ID":{"Value":8960,"Buf":"ACM="},"Value":{"Value":24064,"Buf":"AF4="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}}]}]}},"PkCtrl":{"PkCtrl":{"Value":0,"Buf":"AA=="}},"LookAppend":[{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}}]},"CharacterSkillBag":{"SkillID":{"Value":36,"Buf":"ACQ="},"Type":{"Value":0,"Buf":"AA=="},"SkillNum":{"Value":9,"Buf":"AAk="},"Skills":[{"ID":{"Value":25,"Buf":"ABk="},"State":{"Value":0,"Buf":"AA=="},"Level":{"Value":1,"Buf":"AQ=="},"UseSP":{"Value":0,"Buf":"AAA="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":0,"Buf":"AAAAAA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null},{"ID":{"Value":28,"Buf":"ABw="},"State":{"Value":0,"Buf":"AA=="},"Level":{"Value":1,"Buf":"AQ=="},"UseSP":{"Value":0,"Buf":"AAA="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":0,"Buf":"AAAAAA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null},{"ID":{"Value":29,"Buf":"AB0="},"State":{"Value":0,"Buf":"AA=="},"Level":{"Value":1,"Buf":"AQ=="},"UseSP":{"Value":0,"Buf":"AAA="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":0,"Buf":"AAAAAA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null},{"ID":{"Value":34,"Buf":"ACI="},"State":{"Value":0,"Buf":"AA=="},"Level":{"Value":1,"Buf":"AQ=="},"UseSP":{"Value":0,"Buf":"AAA="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":0,"Buf":"AAAAAA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null},{"ID":{"Value":35,"Buf":"ACM="},"State":{"Value":0,"Buf":"AA=="},"Level":{"Value":1,"Buf":"AQ=="},"UseSP":{"Value":0,"Buf":"AAA="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":0,"Buf":"AAAAAA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null},{"ID":{"Value":36,"Buf":"ACQ="},"State":{"Value":1,"Buf":"AQ=="},"Level":{"Value":1,"Buf":"AQ=="},"UseSP":{"Value":0,"Buf":"AAA="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":0,"Buf":"AAAAAA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null},{"ID":{"Value":37,"Buf":"ACU="},"State":{"Value":0,"Buf":"AA=="},"Level":{"Value":1,"Buf":"AQ=="},"UseSP":{"Value":0,"Buf":"AAA="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":0,"Buf":"AAAAAA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null},{"ID":{"Value":97,"Buf":"AGE="},"State":{"Value":1,"Buf":"AQ=="},"Level":{"Value":1,"Buf":"AQ=="},"UseSP":{"Value":34,"Buf":"ACI="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":6700,"Buf":"AAAaLA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null},{"ID":{"Value":99,"Buf":"AGM="},"State":{"Value":1,"Buf":"AQ=="},"Level":{"Value":8,"Buf":"CA=="},"UseSP":{"Value":46,"Buf":"AC4="},"UseEndure":{"Value":0,"Buf":"AAA="},"UseEnergy":{"Value":0,"Buf":"AAA="},"ResumeTime":{"Value":3600,"Buf":"AAAOEA=="},"RangeType":{"Value":0,"Buf":"AAA="},"Params":null}]},"CharacterSkillState":{"StatesLen":{"Value":0,"Buf":"AA=="},"States":null},"CharacterAttribute":{"Type":{"Value":0,"Buf":"AA=="},"Num":{"Value":74,"Buf":"AEo="},"Attributes":[{"ID":{"Value":0,"Buf":"AA=="},"Value":{"Value":19,"Buf":"AAAAEw=="}},{"ID":{"Value":1,"Buf":"AQ=="},"Value":{"Value":640,"Buf":"AAACgA=="}},{"ID":{"Value":2,"Buf":"Ag=="},"Value":{"Value":292,"Buf":"AAABJA=="}},{"ID":{"Value":3,"Buf":"Aw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":4,"Buf":"BA=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":5,"Buf":"BQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":6,"Buf":"Bg=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":7,"Buf":"Bw=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":8,"Buf":"CA=="},"Value":{"Value":9521,"Buf":"AAAlMQ=="}},{"ID":{"Value":9,"Buf":"CQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":10,"Buf":"Cg=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":11,"Buf":"Cw=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":12,"Buf":"DA=="},"Value":{"Value":2,"Buf":"AAAAAg=="}},{"ID":{"Value":13,"Buf":"DQ=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":14,"Buf":"Dg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":15,"Buf":"Dw=="},"Value":{"Value":65385,"Buf":"AAD/aQ=="}},{"ID":{"Value":16,"Buf":"EA=="},"Value":{"Value":83656,"Buf":"AAFGyA=="}},{"ID":{"Value":17,"Buf":"EQ=="},"Value":{"Value":65306,"Buf":"AAD/Gg=="}},{"ID":{"Value":18,"Buf":"Eg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":19,"Buf":"Ew=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":20,"Buf":"FA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":21,"Buf":"FQ=="},"Value":{"Value":625,"Buf":"AAACcQ=="}},{"ID":{"Value":22,"Buf":"Fg=="},"Value":{"Value":1750,"Buf":"AAAG1g=="}},{"ID":{"Value":23,"Buf":"Fw=="},"Value":{"Value":1500,"Buf":"AAAF3A=="}},{"ID":{"Value":24,"Buf":"GA=="},"Value":{"Value":5000,"Buf":"AAATiA=="}},{"ID":{"Value":25,"Buf":"GQ=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":26,"Buf":"Gg=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":27,"Buf":"Gw=="},"Value":{"Value":8,"Buf":"AAAACA=="}},{"ID":{"Value":28,"Buf":"HA=="},"Value":{"Value":8,"Buf":"AAAACA=="}},{"ID":{"Value":29,"Buf":"HQ=="},"Value":{"Value":28,"Buf":"AAAAHA=="}},{"ID":{"Value":30,"Buf":"Hg=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":31,"Buf":"Hw=="},"Value":{"Value":640,"Buf":"AAACgA=="}},{"ID":{"Value":32,"Buf":"IA=="},"Value":{"Value":292,"Buf":"AAABJA=="}},{"ID":{"Value":33,"Buf":"IQ=="},"Value":{"Value":83,"Buf":"AAAAUw=="}},{"ID":{"Value":34,"Buf":"Ig=="},"Value":{"Value":101,"Buf":"AAAAZQ=="}},{"ID":{"Value":35,"Buf":"Iw=="},"Value":{"Value":20,"Buf":"AAAAFA=="}},{"ID":{"Value":36,"Buf":"JA=="},"Value":{"Value":53,"Buf":"AAAANQ=="}},{"ID":{"Value":37,"Buf":"JQ=="},"Value":{"Value":49,"Buf":"AAAAMQ=="}},{"ID":{"Value":38,"Buf":"Jg=="},"Value":{"Value":105,"Buf":"AAAAaQ=="}},{"ID":{"Value":39,"Buf":"Jw=="},"Value":{"Value":15,"Buf":"AAAADw=="}},{"ID":{"Value":40,"Buf":"KA=="},"Value":{"Value":10,"Buf":"AAAACg=="}},{"ID":{"Value":41,"Buf":"KQ=="},"Value":{"Value":4,"Buf":"AAAABA=="}},{"ID":{"Value":42,"Buf":"Kg=="},"Value":{"Value":1369,"Buf":"AAAFWQ=="}},{"ID":{"Value":43,"Buf":"Kw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":44,"Buf":"LA=="},"Value":{"Value":480,"Buf":"AAAB4A=="}},{"ID":{"Value":45,"Buf":"LQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":46,"Buf":"Lg=="},"Value":{"Value":4,"Buf":"AAAABA=="}},{"ID":{"Value":47,"Buf":"Lw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":48,"Buf":"MA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":49,"Buf":"MQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":50,"Buf":"Mg=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":51,"Buf":"Mw=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":52,"Buf":"NA=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":53,"Buf":"NQ=="},"Value":{"Value":8,"Buf":"AAAACA=="}},{"ID":{"Value":54,"Buf":"Ng=="},"Value":{"Value":27,"Buf":"AAAAGw=="}},{"ID":{"Value":55,"Buf":"Nw=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":56,"Buf":"OA=="},"Value":{"Value":640,"Buf":"AAACgA=="}},{"ID":{"Value":57,"Buf":"OQ=="},"Value":{"Value":292,"Buf":"AAABJA=="}},{"ID":{"Value":58,"Buf":"Og=="},"Value":{"Value":7,"Buf":"AAAABw=="}},{"ID":{"Value":59,"Buf":"Ow=="},"Value":{"Value":7,"Buf":"AAAABw=="}},{"ID":{"Value":60,"Buf":"PA=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":61,"Buf":"PQ=="},"Value":{"Value":46,"Buf":"AAAALg=="}},{"ID":{"Value":62,"Buf":"Pg=="},"Value":{"Value":47,"Buf":"AAAALw=="}},{"ID":{"Value":63,"Buf":"Pw=="},"Value":{"Value":105,"Buf":"AAAAaQ=="}},{"ID":{"Value":64,"Buf":"QA=="},"Value":{"Value":15,"Buf":"AAAADw=="}},{"ID":{"Value":65,"Buf":"QQ=="},"Value":{"Value":10,"Buf":"AAAACg=="}},{"ID":{"Value":66,"Buf":"Qg=="},"Value":{"Value":4,"Buf":"AAAABA=="}},{"ID":{"Value":67,"Buf":"Qw=="},"Value":{"Value":1369,"Buf":"AAAFWQ=="}},{"ID":{"Value":68,"Buf":"RA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":69,"Buf":"RQ=="},"Value":{"Value":480,"Buf":"AAAB4A=="}},{"ID":{"Value":70,"Buf":"Rg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":71,"Buf":"Rw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":72,"Buf":"SA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":73,"Buf":"SQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}}]},"CharacterKitbag":{"Type":{"Value":0,"Buf":"AA=="},"KeybagNum":{"Value":24,"Buf":"ABg="},"Items":[{"GridID":{"Value":0,"Buf":"AAA="},"ID":{"Value":2911,"Buf":"C18="},"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":true,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":6656,"Buf":"ABo="},"Value":{"Value":1792,"Buf":"AAc="}},{"ID":{"Value":7168,"Buf":"ABw="},"Value":{"Value":2560,"Buf":"AAo="}},{"ID":{"Value":7424,"Buf":"AB0="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":6912,"Buf":"ABs="},"Value":{"Value":14080,"Buf":"ADc="}},{"ID":{"Value":7680,"Buf":"AB4="},"Value":{"Value":768,"Buf":"AAM="}}]},{"GridID":{"Value":1,"Buf":"AAE="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":2,"Buf":"AAI="},"ID":{"Value":4118,"Buf":"EBY="},"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":3,"Buf":"AAM="},"ID":{"Value":1789,"Buf":"Bv0="},"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":4,"Buf":"AAQ="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":5,"Buf":"AAU="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":6,"Buf":"AAY="},"ID":{"Value":1848,"Buf":"Bzg="},"Num":{"Value":6,"Buf":"AAY="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":7,"Buf":"AAc="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":8,"Buf":"AAg="},"ID":{"Value":2197,"Buf":"CJU="},"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":10000,"Buf":"JxA="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":1000,"Buf":"A+g="},{"Value":1000,"Buf":"A+g="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":true,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":9216,"Buf":"ACQ="},"Value":{"Value":2816,"Buf":"AAs="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}}]},{"GridID":{"Value":9,"Buf":"AAk="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":10,"Buf":"AAo="},"ID":{"Value":3399,"Buf":"DUc="},"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":11,"Buf":"AAs="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":12,"Buf":"AAw="},"ID":{"Value":439,"Buf":"Abc="},"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":13,"Buf":"AA0="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":14,"Buf":"AA4="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":15,"Buf":"AA8="},"ID":{"Value":3963,"Buf":"D3s="},"Num":{"Value":5,"Buf":"AAU="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":16,"Buf":"ABA="},"ID":{"Value":3988,"Buf":"D5Q="},"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":2952808369,"Buf":"sABHsQ=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":17,"Buf":"ABE="},"ID":{"Value":3964,"Buf":"D3w="},"Num":{"Value":10,"Buf":"AAo="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":18,"Buf":"ABI="},"ID":{"Value":2656,"Buf":"CmA="},"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":10000,"Buf":"JxA="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":true,"Buf":"AQ=="},"ItemDBInstID":{"Value":0,"Buf":"AAAAAA=="},"ItemDBForge":{"Value":0,"Buf":"AAAAAA=="},"IsParams":{"Value":false,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":19,"Buf":"ABM="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":20,"Buf":"ABQ="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":21,"Buf":"ABU="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":22,"Buf":"ABY="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":23,"Buf":"ABc="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":65535,"Buf":"//8="},"ID":{"Value":0,"Buf":null},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]}]},"CharacterShortcut":{"Shortcuts":[{"Type":{"Value":2,"Buf":"Ag=="},"GridID":{"Value":97,"Buf":"AGE="}},{"Type":{"Value":2,"Buf":"Ag=="},"GridID":{"Value":99,"Buf":"AGM="}},{"Type":{"Value":1,"Buf":"AQ=="},"GridID":{"Value":6,"Buf":"AAY="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}},{"Type":{"Value":0,"Buf":"AA=="},"GridID":{"Value":0,"Buf":"AAA="}}]},"BoatLen":{"Value":1,"Buf":"AQ=="},"CharacterBoats":[{"CharacterBase":{"ChaID":{"Value":305,"Buf":"AAABMQ=="},"WorldID":{"Value":2952808369,"Buf":"sABHsQ=="},"CommID":{"Value":10271,"Buf":"AAAoHw=="},"CommName":{"Value":"ingrysty","Buf":"AAlpbmdyeXN0eQA="},"GmLvl":{"Value":0,"Buf":"AA=="},"Handle":{"Value":33565846,"Buf":"AgAslg=="},"CtrlType":{"Value":1,"Buf":"AQ=="},"Name":{"Value":"Lodka","Buf":"AAZMb2RrYQA="},"MottoName":{"Value":"","Buf":"AAEA"},"Icon":{"Value":4,"Buf":"AAQ="},"GuildID":{"Value":0,"Buf":"AAAAAA=="},"GuildName":{"Value":"","Buf":"AAEA"},"GuildMotto":{"Value":"","Buf":"AAEA"},"StallName":{"Value":"","Buf":"AAEA"},"State":{"Value":1,"Buf":"AAE="},"Position":{"X":{"Value":4294967295,"Buf":"/////w=="},"Y":{"Value":4294967295,"Buf":"/////w=="},"Radius":{"Value":40,"Buf":"AAAAKA=="}},"Angle":{"Value":21,"Buf":"ABU="},"TeamLeaderID":{"Value":0,"Buf":"AAAAAA=="},"Side":{"SideID":{"Value":0,"Buf":"AA=="}},"EntityEvent":{"EntityID":{"Value":2952808369,"Buf":"sABHsQ=="},"EntityType":{"Value":1,"Buf":"AQ=="},"EventID":{"Value":0,"Buf":"AAA="},"EventName":{"Value":"","Buf":"AAEA"}},"Look":{"SynType":{"Value":0,"Buf":"AA=="},"TypeID":{"Value":305,"Buf":"ATE="},"IsBoat":{"Value":1,"Buf":"AQ=="},"LookBoat":{"PosID":{"Value":405,"Buf":"AZU="},"BoatID":{"Value":1,"Buf":"AAE="},"Header":{"Value":8,"Buf":"AAg="},"Body":{"Value":43,"Buf":"ACs="},"Engine":{"Value":15,"Buf":"AA8="},"Cannon":{"Value":53,"Buf":"ADU="},"Equipment":{"Value":73,"Buf":"AEk="}},"LookHuman":{"HairID":{"Value":0,"Buf":null},"ItemGrid":[{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":null},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]}]}},"PkCtrl":{"PkCtrl":{"Value":0,"Buf":"AA=="}},"LookAppend":[{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}}]},"CharacterAttribute":{"Type":{"Value":0,"Buf":"AA=="},"Num":{"Value":74,"Buf":"AEo="},"Attributes":[{"ID":{"Value":0,"Buf":"AA=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":1,"Buf":"AQ=="},"Value":{"Value":19,"Buf":"AAAAEw=="}},{"ID":{"Value":2,"Buf":"Ag=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":3,"Buf":"Aw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":4,"Buf":"BA=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":5,"Buf":"BQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":6,"Buf":"Bg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":7,"Buf":"Bw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":8,"Buf":"CA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":9,"Buf":"CQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":10,"Buf":"Cg=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":11,"Buf":"Cw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":12,"Buf":"DA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":13,"Buf":"DQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":14,"Buf":"Dg=="},"Value":{"Value":12,"Buf":"AAAADA=="}},{"ID":{"Value":15,"Buf":"Dw=="},"Value":{"Value":1765,"Buf":"AAAG5Q=="}},{"ID":{"Value":16,"Buf":"EA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":17,"Buf":"EQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":18,"Buf":"Eg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":19,"Buf":"Ew=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":20,"Buf":"FA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":21,"Buf":"FQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":22,"Buf":"Fg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":23,"Buf":"Fw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":24,"Buf":"GA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":25,"Buf":"GQ=="},"Value":{"Value":9,"Buf":"AAAACQ=="}},{"ID":{"Value":26,"Buf":"Gg=="},"Value":{"Value":9,"Buf":"AAAACQ=="}},{"ID":{"Value":27,"Buf":"Gw=="},"Value":{"Value":9,"Buf":"AAAACQ=="}},{"ID":{"Value":28,"Buf":"HA=="},"Value":{"Value":9,"Buf":"AAAACQ=="}},{"ID":{"Value":29,"Buf":"HQ=="},"Value":{"Value":28,"Buf":"AAAAHA=="}},{"ID":{"Value":30,"Buf":"Hg=="},"Value":{"Value":9,"Buf":"AAAACQ=="}},{"ID":{"Value":31,"Buf":"Hw=="},"Value":{"Value":1919,"Buf":"AAAHfw=="}},{"ID":{"Value":32,"Buf":"IA=="},"Value":{"Value":500,"Buf":"AAAB9A=="}},{"ID":{"Value":33,"Buf":"IQ=="},"Value":{"Value":140,"Buf":"AAAAjA=="}},{"ID":{"Value":34,"Buf":"Ig=="},"Value":{"Value":210,"Buf":"AAAA0g=="}},{"ID":{"Value":35,"Buf":"Iw=="},"Value":{"Value":38,"Buf":"AAAAJg=="}},{"ID":{"Value":36,"Buf":"JA=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":37,"Buf":"JQ=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":38,"Buf":"Jg=="},"Value":{"Value":190,"Buf":"AAAAvg=="}},{"ID":{"Value":39,"Buf":"Jw=="},"Value":{"Value":10,"Buf":"AAAACg=="}},{"ID":{"Value":40,"Buf":"KA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":41,"Buf":"KQ=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":42,"Buf":"Kg=="},"Value":{"Value":1428,"Buf":"AAAFlA=="}},{"ID":{"Value":43,"Buf":"Kw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":44,"Buf":"LA=="},"Value":{"Value":406,"Buf":"AAABlg=="}},{"ID":{"Value":45,"Buf":"LQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":46,"Buf":"Lg=="},"Value":{"Value":10,"Buf":"AAAACg=="}},{"ID":{"Value":47,"Buf":"Lw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":48,"Buf":"MA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":49,"Buf":"MQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":50,"Buf":"Mg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":51,"Buf":"Mw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":52,"Buf":"NA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":53,"Buf":"NQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":54,"Buf":"Ng=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":55,"Buf":"Nw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":56,"Buf":"OA=="},"Value":{"Value":2280,"Buf":"AAAI6A=="}},{"ID":{"Value":57,"Buf":"OQ=="},"Value":{"Value":500,"Buf":"AAAB9A=="}},{"ID":{"Value":58,"Buf":"Og=="},"Value":{"Value":167,"Buf":"AAAApw=="}},{"ID":{"Value":59,"Buf":"Ow=="},"Value":{"Value":250,"Buf":"AAAA+g=="}},{"ID":{"Value":60,"Buf":"PA=="},"Value":{"Value":46,"Buf":"AAAALg=="}},{"ID":{"Value":61,"Buf":"PQ=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":62,"Buf":"Pg=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":63,"Buf":"Pw=="},"Value":{"Value":190,"Buf":"AAAAvg=="}},{"ID":{"Value":64,"Buf":"QA=="},"Value":{"Value":10,"Buf":"AAAACg=="}},{"ID":{"Value":65,"Buf":"QQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":66,"Buf":"Qg=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":67,"Buf":"Qw=="},"Value":{"Value":70,"Buf":"AAAARg=="}},{"ID":{"Value":68,"Buf":"RA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":69,"Buf":"RQ=="},"Value":{"Value":450,"Buf":"AAABwg=="}},{"ID":{"Value":70,"Buf":"Rg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":71,"Buf":"Rw=="},"Value":{"Value":10,"Buf":"AAAACg=="}},{"ID":{"Value":72,"Buf":"SA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":73,"Buf":"SQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}}]},"CharacterKitbag":{"Type":{"Value":0,"Buf":"AA=="},"KeybagNum":{"Value":24,"Buf":"ABg="},"Items":[{"GridID":{"Value":0,"Buf":"AAA="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":1,"Buf":"AAE="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":2,"Buf":"AAI="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":3,"Buf":"AAM="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":4,"Buf":"AAQ="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":5,"Buf":"AAU="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":6,"Buf":"AAY="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":7,"Buf":"AAc="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":8,"Buf":"AAg="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":9,"Buf":"AAk="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":10,"Buf":"AAo="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":11,"Buf":"AAs="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":12,"Buf":"AAw="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":13,"Buf":"AA0="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":14,"Buf":"AA4="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":15,"Buf":"AA8="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":16,"Buf":"ABA="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":17,"Buf":"ABE="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":18,"Buf":"ABI="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":19,"Buf":"ABM="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":20,"Buf":"ABQ="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":21,"Buf":"ABU="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":22,"Buf":"ABY="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":23,"Buf":"ABc="},"ID":{"Value":0,"Buf":"AAA="},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"GridID":{"Value":65535,"Buf":"//8="},"ID":{"Value":0,"Buf":null},"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":false,"Buf":null},"ItemDBInstID":{"Value":0,"Buf":null},"ItemDBForge":{"Value":0,"Buf":null},"IsParams":{"Value":false,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]}]},"CharacterSkillState":{"StatesLen":{"Value":0,"Buf":"AA=="},"States":null}}],"ChaMainID":{"Value":10271,"Buf":"AAAoHw=="}}`),
		e,
	)
	if err != nil {
		fmt.Println(err)
	}
}

type EnterGameRequest struct {
	CharacterName types.HardString
}

func (e EnterGameRequest) Opcode() uint16 {
	return 433
}

func (e *EnterGameRequest) Process(buf *[]byte, mode ...processor.Mode) {
	p := processor.NewProcessor(processor.Read)
	p.String(buf, &e.CharacterName)
}

func (e *EnterGameRequest) Print() {
	DebugPrint(e)
}
