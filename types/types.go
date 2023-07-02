package types

type PullContext interface {
	GetBuf() []byte
	GetValue() interface{}
	GetType() string
}

type Context struct {
	Buf []byte
}

func (c Context) GetBuf() []byte {
	return c.Buf
}

type HardUInt8 struct {
	Value uint8
	Context
}

func (h HardUInt8) GetValue() interface{} {
	return h.Value
}

func (h HardUInt8) GetType() string {
	return "uint8"
}

type HardUInt16 struct {
	Value uint16
	Context
}

func (h HardUInt16) GetValue() interface{} {
	return h.Value
}

func (h HardUInt16) GetType() string {
	return "uint16"
}

type HardUInt32 struct {
	Value uint32
	Context
}

func (h HardUInt32) GetValue() interface{} {
	return h.Value
}

func (h HardUInt32) GetType() string {
	return "uint32"
}

type HardString struct {
	Value string
	Context
}

func (h HardString) GetValue() interface{} {
	return h.Value
}

func (h HardString) GetType() string {
	return "string"
}

type HardBool struct {
	Value bool
	Context
}

func (h HardBool) GetValue() interface{} {
	return h.Value
}

func (h HardBool) GetType() string {
	return "bool"
}

type HardBytes struct {
	Value []byte
	Context
}

func (h HardBytes) GetValue() interface{} {
	return h.Value
}

func (h HardBytes) GetType() string {
	return "[]byte"
}
