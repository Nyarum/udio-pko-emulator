package packets

import (
	"encoding/json"
	"fmt"
	"little/processor"
	"little/types"
)

type PoseLean struct {
	State  types.HardUInt8
	Pose   types.HardUInt32
	Angle  types.HardUInt32
	PosX   types.HardUInt32
	PosY   types.HardUInt32
	Height types.HardUInt32
}

func (p *PoseLean) Process(buf *[]byte, mode ...processor.Mode) {

}

type PoseSeat struct {
	Pose  types.HardUInt16
	Angle types.HardUInt16
}

func (p *PoseSeat) Process(buf *[]byte, mode ...processor.Mode) {

}

// ViewObjectResponse хочет пакет с стейтом остановки шага в обратку
type ViewObjectResponse struct {
	SeeType             types.HardUInt8
	CharacterBase       CharacterBase
	NpcType             types.HardUInt8
	NpcState            types.HardUInt8
	StatePose           types.HardUInt16
	PoseLean            PoseLean // only if StatePose == PoseLean
	PoseSeat            PoseSeat // only if StatePose == PoseSeat
	CharacterAttribute  CharacterAttribute
	CharacterSkillState CharacterSkillState
}

func (v ViewObjectResponse) Opcode() uint16 {
	return 504
}

func (v *ViewObjectResponse) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &v.SeeType)
	(&v.CharacterBase).Process(buf, defaultMode)
	p.UInt8(buf, &v.NpcType)
	p.UInt8(buf, &v.NpcState)
	p.UInt16(buf, &v.StatePose)

	switch v.StatePose.Value {
	case types.PoseLean:
		(&v.PoseLean).Process(buf, defaultMode)
	case types.PoseSeat:
		(&v.PoseSeat).Process(buf, defaultMode)

	}

	(&v.CharacterAttribute).Process(buf, defaultMode)
	(&v.CharacterSkillState).Process(buf, defaultMode)
}

func (v *ViewObjectResponse) Print() {
	DebugPrint(v)
}

func (v *ViewObjectResponse) Basic() {
	err := json.Unmarshal([]byte(`{"SeeType":{"Value":0,"Buf":"AA=="},"CharacterBase":{"ChaID":{"Value":1,"Buf":"AAAAAQ=="},"WorldID":{"Value":10329,"Buf":"AAAoWQ=="},"CommID":{"Value":10329,"Buf":"AAAoWQ=="},"CommName":{"Value":"MrTProg","Buf":"AAhNclRQcm9nAA=="},"GmLvl":{"Value":99,"Buf":"Yw=="},"Handle":{"Value":33565846,"Buf":"AgAslg=="},"CtrlType":{"Value":1,"Buf":"AQ=="},"Name":{"Value":"MrTProg","Buf":"AAhNclRQcm9nAA=="},"MottoName":{"Value":"","Buf":"AAEA"},"Icon":{"Value":1,"Buf":"AAE="},"GuildID":{"Value":0,"Buf":"AAAAAA=="},"GuildName":{"Value":"","Buf":"AAEA"},"GuildMotto":{"Value":"","Buf":"AAEA"},"StallName":{"Value":"","Buf":"AAEA"},"State":{"Value":3,"Buf":"AAM="},"Position":{"X":{"Value":222270,"Buf":"AANkPg=="},"Y":{"Value":268288,"Buf":"AAQYAA=="},"Radius":{"Value":40,"Buf":"AAAAKA=="}},"Angle":{"Value":270,"Buf":"AQ4="},"TeamLeaderID":{"Value":0,"Buf":"AAAAAA=="},"Side":{"SideID":{"Value":0,"Buf":"AA=="}},"EntityEvent":{"EntityID":{"Value":10329,"Buf":"AAAoWQ=="},"EntityType":{"Value":1,"Buf":"AQ=="},"EventID":{"Value":0,"Buf":"AAA="},"EventName":{"Value":"","Buf":"AAEA"}},"Look":{"SynType":{"Value":0,"Buf":"AA=="},"TypeID":{"Value":1,"Buf":"AAE="},"IsBoat":{"Value":0,"Buf":"AA=="},"LookBoat":{"PosID":{"Value":0,"Buf":null},"BoatID":{"Value":0,"Buf":null},"Header":{"Value":0,"Buf":null},"Body":{"Value":0,"Buf":null},"Engine":{"Value":0,"Buf":null},"Cannon":{"Value":0,"Buf":null},"Equipment":{"Value":0,"Buf":null}},"LookHuman":{"HairID":{"Value":2000,"Buf":"B9A="},"ItemGrid":[{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":2554,"Buf":"Cfo="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":"AAA="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":0,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":289,"Buf":"ASE="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":9481,"Buf":"JQk="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":1,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":9216,"Buf":"ACQ="},"Value":{"Value":512,"Buf":"AAI="}},{"ID":{"Value":12032,"Buf":"AC8="},"Value":{"Value":256,"Buf":"AAE="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}}]},{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":641,"Buf":"AoE="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":9477,"Buf":"JQU="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":1,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":9984,"Buf":"ACc="},"Value":{"Value":256,"Buf":"AAE="}},{"ID":{"Value":9216,"Buf":"ACQ="},"Value":{"Value":512,"Buf":"AAI="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}}]},{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":9999,"Buf":"Jw8="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":"AAA="},"Endure":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":0,"Buf":"AA=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":0,"Buf":"AA=="},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":0,"Buf":"AAA="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":0,"Buf":null},"Endure":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"Energy":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"ForgeLevel":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"IsDBParams":{"Value":0,"Buf":null},"DBParams":[{"Value":0,"Buf":null},{"Value":0,"Buf":null}],"IsInstAttrs":{"Value":0,"Buf":null},"InstAttrs":[{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}},{"ID":{"Value":0,"Buf":null},"Value":{"Value":0,"Buf":null}}]},{"ID":{"Value":8,"Buf":"AAg="},"ItemSync":{"Endure":{"Value":0,"Buf":null},"Energy":{"Value":0,"Buf":null},"IsValid":{"Value":0,"Buf":null}},"ItemShow":{"Num":{"Value":1,"Buf":"AAE="},"Endure":[{"Value":9497,"Buf":"JRk="},{"Value":10000,"Buf":"JxA="}],"Energy":[{"Value":0,"Buf":"AAA="},{"Value":0,"Buf":"AAA="}],"ForgeLevel":{"Value":0,"Buf":"AA=="},"IsValid":{"Value":1,"Buf":"AQ=="}},"IsDBParams":{"Value":1,"Buf":"AQ=="},"DBParams":[{"Value":0,"Buf":"AAAAAA=="},{"Value":0,"Buf":"AAAAAA=="}],"IsInstAttrs":{"Value":1,"Buf":"AQ=="},"InstAttrs":[{"ID":{"Value":8704,"Buf":"ACI="},"Value":{"Value":3584,"Buf":"AA4="}},{"ID":{"Value":8960,"Buf":"ACM="},"Value":{"Value":4608,"Buf":"ABI="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}},{"ID":{"Value":0,"Buf":"AAA="},"Value":{"Value":0,"Buf":"AAA="}}]}]}},"PkCtrl":{"PkCtrl":{"Value":0,"Buf":"AA=="}},"LookAppend":[{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}},{"LookID":{"Value":0,"Buf":"AAA="},"IsValid":{"Value":0,"Buf":null}}]},"NpcType":{"Value":0,"Buf":"AA=="},"NpcState":{"Value":0,"Buf":"AA=="},"StatePose":{"Value":0,"Buf":"AAA="},"PoseLean":{"State":{"Value":0,"Buf":null},"Pose":{"Value":0,"Buf":null},"Angle":{"Value":0,"Buf":null},"PosX":{"Value":0,"Buf":null},"PosY":{"Value":0,"Buf":null},"Height":{"Value":0,"Buf":null}},"PoseSeat":{"Pose":{"Value":0,"Buf":null},"Angle":{"Value":0,"Buf":null}},"CharacterAttribute":{"Type":{"Value":0,"Buf":"AA=="},"Num":{"Value":74,"Buf":"AEo="},"Attributes":[{"ID":{"Value":0,"Buf":"AA=="},"Value":{"Value":98,"Buf":"AAAAYg=="}},{"ID":{"Value":1,"Buf":"AQ=="},"Value":{"Value":3055,"Buf":"AAAL7w=="}},{"ID":{"Value":2,"Buf":"Ag=="},"Value":{"Value":314,"Buf":"AAABOg=="}},{"ID":{"Value":3,"Buf":"Aw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":4,"Buf":"BA=="},"Value":{"Value":9,"Buf":"AAAACQ=="}},{"ID":{"Value":5,"Buf":"BQ=="},"Value":{"Value":87,"Buf":"AAAAVw=="}},{"ID":{"Value":6,"Buf":"Bg=="},"Value":{"Value":176,"Buf":"AAAAsA=="}},{"ID":{"Value":7,"Buf":"Bw=="},"Value":{"Value":96,"Buf":"AAAAYA=="}},{"ID":{"Value":8,"Buf":"CA=="},"Value":{"Value":586140,"Buf":"AAjxnA=="}},{"ID":{"Value":9,"Buf":"CQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":10,"Buf":"Cg=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":11,"Buf":"Cw=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":12,"Buf":"DA=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":13,"Buf":"DQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":14,"Buf":"Dg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":15,"Buf":"Dw=="},"Value":{"Value":3871833472,"Buf":"5sd9gA=="}},{"ID":{"Value":16,"Buf":"EA=="},"Value":{"Value":3924326032,"Buf":"6eh2kA=="}},{"ID":{"Value":17,"Buf":"EQ=="},"Value":{"Value":3706378269,"Buf":"3OrYHQ=="}},{"ID":{"Value":18,"Buf":"Eg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":19,"Buf":"Ew=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":20,"Buf":"FA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":21,"Buf":"FQ=="},"Value":{"Value":625,"Buf":"AAACcQ=="}},{"ID":{"Value":22,"Buf":"Fg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":23,"Buf":"Fw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":24,"Buf":"GA=="},"Value":{"Value":1500,"Buf":"AAAF3A=="}},{"ID":{"Value":25,"Buf":"GQ=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":26,"Buf":"Gg=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":27,"Buf":"Gw=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":28,"Buf":"HA=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":29,"Buf":"HQ=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":30,"Buf":"Hg=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":31,"Buf":"Hw=="},"Value":{"Value":3055,"Buf":"AAAL7w=="}},{"ID":{"Value":32,"Buf":"IA=="},"Value":{"Value":314,"Buf":"AAABOg=="}},{"ID":{"Value":33,"Buf":"IQ=="},"Value":{"Value":21,"Buf":"AAAAFQ=="}},{"ID":{"Value":34,"Buf":"Ig=="},"Value":{"Value":25,"Buf":"AAAAGQ=="}},{"ID":{"Value":35,"Buf":"Iw=="},"Value":{"Value":9,"Buf":"AAAACQ=="}},{"ID":{"Value":36,"Buf":"JA=="},"Value":{"Value":204,"Buf":"AAAAzA=="}},{"ID":{"Value":37,"Buf":"JQ=="},"Value":{"Value":205,"Buf":"AAAAzQ=="}},{"ID":{"Value":38,"Buf":"Jg=="},"Value":{"Value":105,"Buf":"AAAAaQ=="}},{"ID":{"Value":39,"Buf":"Jw=="},"Value":{"Value":15,"Buf":"AAAADw=="}},{"ID":{"Value":40,"Buf":"KA=="},"Value":{"Value":35,"Buf":"AAAAIw=="}},{"ID":{"Value":41,"Buf":"KQ=="},"Value":{"Value":2,"Buf":"AAAAAg=="}},{"ID":{"Value":42,"Buf":"Kg=="},"Value":{"Value":1428,"Buf":"AAAFlA=="}},{"ID":{"Value":43,"Buf":"Kw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":44,"Buf":"LA=="},"Value":{"Value":480,"Buf":"AAAB4A=="}},{"ID":{"Value":45,"Buf":"LQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":46,"Buf":"Lg=="},"Value":{"Value":1,"Buf":"AAAAAQ=="}},{"ID":{"Value":47,"Buf":"Lw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":48,"Buf":"MA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":49,"Buf":"MQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":50,"Buf":"Mg=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":51,"Buf":"Mw=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":52,"Buf":"NA=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":53,"Buf":"NQ=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":54,"Buf":"Ng=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":55,"Buf":"Nw=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":56,"Buf":"OA=="},"Value":{"Value":3055,"Buf":"AAAL7w=="}},{"ID":{"Value":57,"Buf":"OQ=="},"Value":{"Value":314,"Buf":"AAABOg=="}},{"ID":{"Value":58,"Buf":"Og=="},"Value":{"Value":7,"Buf":"AAAABw=="}},{"ID":{"Value":59,"Buf":"Ow=="},"Value":{"Value":7,"Buf":"AAAABw=="}},{"ID":{"Value":60,"Buf":"PA=="},"Value":{"Value":5,"Buf":"AAAABQ=="}},{"ID":{"Value":61,"Buf":"PQ=="},"Value":{"Value":204,"Buf":"AAAAzA=="}},{"ID":{"Value":62,"Buf":"Pg=="},"Value":{"Value":204,"Buf":"AAAAzA=="}},{"ID":{"Value":63,"Buf":"Pw=="},"Value":{"Value":105,"Buf":"AAAAaQ=="}},{"ID":{"Value":64,"Buf":"QA=="},"Value":{"Value":15,"Buf":"AAAADw=="}},{"ID":{"Value":65,"Buf":"QQ=="},"Value":{"Value":35,"Buf":"AAAAIw=="}},{"ID":{"Value":66,"Buf":"Qg=="},"Value":{"Value":2,"Buf":"AAAAAg=="}},{"ID":{"Value":67,"Buf":"Qw=="},"Value":{"Value":1428,"Buf":"AAAFlA=="}},{"ID":{"Value":68,"Buf":"RA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":69,"Buf":"RQ=="},"Value":{"Value":480,"Buf":"AAAB4A=="}},{"ID":{"Value":70,"Buf":"Rg=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":71,"Buf":"Rw=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":72,"Buf":"SA=="},"Value":{"Value":0,"Buf":"AAAAAA=="}},{"ID":{"Value":73,"Buf":"SQ=="},"Value":{"Value":0,"Buf":"AAAAAA=="}}]},"CharacterSkillState":{"StatesLen":{"Value":0,"Buf":"AA=="},"States":null}}`),
		v,
	)
	if err != nil {
		fmt.Println(err)
	}
}

type EndViewObject struct {
	SeeType types.HardUInt8
	ID      types.HardUInt32
}

func (e EndViewObject) Opcode() uint16 {
	return 505
}

func (e *EndViewObject) Process(buf *[]byte, mode ...processor.Mode) {
	defaultMode := processor.Write
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	p := processor.NewProcessor(defaultMode)
	p.UInt8(buf, &e.SeeType)
	p.UInt32(buf, &e.ID)
}

func (e *EndViewObject) Print() {
	DebugPrint(e)
}
