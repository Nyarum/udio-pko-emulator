package storage

import (
	"little/packets"
)

type Character struct {
	packets.Character

	ID                  uint32
	AccountID           string
	Base                packets.CharacterBase
	CharacterAttribute  packets.CharacterAttribute
	CharacterSkillState packets.CharacterSkillState
}
