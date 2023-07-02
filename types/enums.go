package types

const (
	SynLookSwitch uint8 = iota
	SynLookChange
)

const (
	SYN_KITBAG_INIT uint8 = iota
	SYN_KITBAG_EQUIP
	SYN_KITBAG_UNFIX
	SYN_KITBAG_PICK
	SYN_KITBAG_THROW
	SYN_KITBAG_SWITCH
	SYN_KITBAG_TRADE
	SYN_KITBAG_FROM_NPC
	SYN_KITBAG_TO_NPC
	SYN_KITBAG_SYSTEM
	SYN_KITBAG_FORGES
	SYN_KITBAG_FORGEF
	SYN_KITBAG_BANK
	SYN_KITBAG_ATTR
)

const (
	MAX_FAST_ROW    = 3
	MAX_FAST_COL    = 12
	SHORT_CUT_NUM   = MAX_FAST_ROW * MAX_FAST_COL
	ESPE_KBGRID_NUM = 4
)

const (
	ACTION_NONE uint8 = iota
	ACTION_MOVE
	ACTION_SKILL
	ACTION_SKILL_SRC
	ACTION_SKILL_TAR
	ACTION_LOOK
	ACTION_KITBAG
	ACTION_SKILLBAG
	ACTION_ITEM_PICK
	ACTION_ITEM_THROW
	ACTION_ITEM_UNFIX
	ACTION_ITEM_USE
	ACTION_ITEM_POS
	ACTION_ITEM_DELETE
	ACTION_ITEM_INFO
	ACTION_ITEM_FAILED
	ACTION_LEAN
	ACTION_CHANGE_CHA
	ACTION_EVENT
	ACTION_FACE
	ACTION_STOP_STATE
	ACTION_SKILL_POSE
	ACTION_PK_CTRL
	ACTION_LOOK_ENERGY
	ACTION_TEMP
	ACTION_SHORTCUT
	ACTION_BANK
	ACTION_CLOSE_BANK
	ACTION_KITBAGTMP
	ACTION_KITBAGTMP_DRAG
	ACTION_GUILDBANK
	ACTION_REQUESTGUILDBANK
	ACTION_REQUESTGUILDATTR
	ACTION_REQUESTGUILDQUEST
	ACTION_ENTERGUILDHOUSE
	ACTION_GOLDSTORE
	ACTION_BAGOFHOLDING
	ACTION_BAGOFHOLDINGNAME
	MAX_ACTION_NUM
)

const (
	MSTATE_ON       uint16 = 0x00
	MSTATE_ARRIVE   uint16 = 0x01
	MSTATE_BLOCK    uint16 = 0x02
	MSTATE_CANCEL   uint16 = 0x04
	MSTATE_INRANGE  uint16 = 0x08
	MSTATE_NOTARGET uint16 = 0x10
	MSTATE_CANTMOVE uint16 = 0x20
)

const (
	PoseStand uint16 = iota
	PoseLean
	PoseSeat
)

const (
	ENTITY_SEEN_NEW = iota
	ENTITY_SEEN_SWITCH
)

const (
	ErrMCDefault uint16 = iota
	ERR_MC_NETEXCP
	ERR_MC_NOTSELCHA
	ERR_MC_NOTPLAY
	ERR_MC_NOTARRIVE
	ERR_MC_TOOMANYPLY
	ERR_MC_NOTLOGIN
	ERR_MC_VER_ERROR
	ERR_MC_ENTER_ERROR
	ERR_MC_ENTER_POS
	ERR_MC_BANUSER
	ERR_MC_PBANUSER
)

const (
	ErrAPDefault uint16 = iota + 1000
	ErrAPINVALIDUSER
	ErrAPINVALIDPWD
	ErrAPACTIVEUSER
	ErrAPLOGGED
	ErrAPDISABLELOGIN
)