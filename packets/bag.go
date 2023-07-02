package packets

import (
	"little/processor"
	"little/types"
)

type KitbagSyncRequest struct {
}

func (k KitbagSyncRequest) Opcode() uint16 {
	return 35
}

func (k *KitbagSyncRequest) Process(buf *[]byte, mode ...processor.Mode) {
	//processor.NewProcessor(processor.Read)
}

func (k *KitbagSyncRequest) Print() {
	DebugPrint(k)
}

type KitbagSyncResponse struct {
	Kitbag CharacterKitbag
}

func (k KitbagSyncResponse) Opcode() uint16 {
	return 554
}

func (k *KitbagSyncResponse) Process(buf *[]byte, mode ...processor.Mode) {
	//p := processor.NewProcessor(processor.Write)
	k.Basic()

	k.Kitbag.Process(buf, mode...)
}

func (k *KitbagSyncResponse) Basic() {
	k.Kitbag.KeybagNum = types.HardUInt16{Value: 24}

	for i := 0; i < 24; i++ {
		gridID := types.HardUInt16{Value: uint16(i)}

		k.Kitbag.Items = append(k.Kitbag.Items, KitbagItem{
			GridID: gridID,
		})
	}

	k.Kitbag.Items = append(k.Kitbag.Items, KitbagItem{
		GridID: types.HardUInt16{Value: 65535},
	})
}

func (k *KitbagSyncResponse) Print() {
	DebugPrint(k)
}
