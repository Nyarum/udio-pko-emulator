package handlers

import (
	"fmt"
	"little/packets"
	"little/storage"
	"little/types"
	"log"
	"math"
	"sync"
)

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

type ViewObjects map[string][]string

func (v ViewObjects) IsView(p1, p2 string) bool {
	if p1Obj, ok := v[p1]; ok {
		for _, obj := range p1Obj {
			if p2 == obj {
				return true
			}
		}
	}

	return false
}

func (v ViewObjects) RemoveView(p1, p2 string) bool {
	if p1Obj, ok := v[p1]; ok {
		for k, obj := range p1Obj {
			if p2 == obj {
				v[p1] = remove(v[p1], k)
				return true
			}
		}
	}

	return false
}

type CharacterStep struct {
	AccountID  string
	Name       string
	X, Y       uint32
	Angle      float64
	MoveCount  uint32
	ActionType uint8
}

type World struct {
	players           *Players
	characters        map[string]storage.Character
	viewObjects       ViewObjects
	UpdateCharacterCh chan Handler
	charactersStepCh  chan CharacterStep
	sync.RWMutex
}

func NewWorld(players *Players) *World {
	w := &World{
		players:           players,
		characters:        make(map[string]storage.Character),
		viewObjects:       make(ViewObjects),
		UpdateCharacterCh: make(chan Handler),
		charactersStepCh:  make(chan CharacterStep, 1),
		RWMutex:           sync.RWMutex{},
	}

	go w.listenCharactersStep()

	return w
}

func (w *World) CharacterEnterWorld(ch storage.Character) {
	w.Lock()
	defer w.Unlock()

	log.Printf("Character (%v) entered to world\n", ch.Name.Value)

	w.characters[ch.Name.Value] = ch
}

func (w *World) CharacterExitWorld(ch storage.Character) {
	w.Lock()
	defer w.Unlock()

	log.Printf("Character (%v) out of a world\n", ch.Name.Value)

	delete(w.viewObjects, ch.AccountID)
	delete(w.characters, ch.Name.Value)
}

func (w *World) listenCharactersStep() {
	for cs := range w.charactersStepCh {
		fmt.Printf("Try to add a new step to existing character: %v\n", cs)

		if v, ok := w.characters[cs.Name]; ok {
			oldXPos := v.Base.Position.X
			oldYPos := v.Base.Position.Y

			v.Base.Position.X.Value = cs.X
			v.Base.Position.Y.Value = cs.Y
			v.Base.Angle.Value = uint16(cs.Angle)
			v.AccountID = cs.AccountID

			w.characters[cs.Name] = v

			chX := float64(cs.X / 100)
			chY := float64(cs.Y / 100)

			go func() {
				w.RLock()
				defer w.RUnlock()

				for _, char := range w.characters {
					if char.Name.Value == v.Name.Value {
						continue
					}

					otherChX := float64(char.Base.Position.X.Value / 100)
					otherChY := float64(char.Base.Position.Y.Value / 100)

					distance := w.calculateDistance(chX, chY, otherChX, otherChY)

					log.Printf("Check distance between (%v) and (%v), with distance: %v", v.Name.Value, char.Name.Value, distance)
					log.Printf("Coordinates of character (%v): %v, %v", v.Name.Value, chX, chY)
					log.Printf("Coordinates of character (%v): %v, %v", char.Name.Value, otherChX, otherChY)

					if distance < defaultDistanceView && !w.viewObjects.IsView(v.AccountID, char.AccountID) {
						log.Printf("Added a player (%v) to view of player (%v)", char.Name.Value, v.Name.Value)

						if _, ok := w.viewObjects[v.AccountID]; !ok {
							w.viewObjects[v.AccountID] = []string{char.AccountID}
						} else {
							w.viewObjects[v.AccountID] = append(w.viewObjects[v.AccountID], char.AccountID)
						}

						w.players.SendPacketToPlayer(v.AccountID, &packets.ViewObjectResponse{
							SeeType:             types.HardUInt8{Value: 0},
							CharacterBase:       char.Base,
							NpcType:             types.HardUInt8{Value: 0},
							NpcState:            types.HardUInt8{Value: 0},
							StatePose:           types.HardUInt16{Value: types.PoseStand},
							CharacterAttribute:  char.CharacterAttribute,
							CharacterSkillState: char.CharacterSkillState,
						})

						if _, ok := w.viewObjects[char.AccountID]; !ok {
							w.viewObjects[char.AccountID] = []string{v.AccountID}
						} else {
							w.viewObjects[char.AccountID] = append(w.viewObjects[char.AccountID], v.AccountID)
						}

						w.players.SendPacketToPlayer(char.AccountID, &packets.ViewObjectResponse{
							SeeType:             types.HardUInt8{Value: 0},
							CharacterBase:       v.Base,
							NpcType:             types.HardUInt8{Value: 0},
							NpcState:            types.HardUInt8{Value: 0},
							StatePose:           types.HardUInt16{Value: types.PoseStand},
							CharacterAttribute:  v.CharacterAttribute,
							CharacterSkillState: v.CharacterSkillState,
						})
					} else if distance > defaultDistanceView {
						log.Printf("Removed a player (%v) to view of player (%v)", char.Name.Value, v.Name.Value)

						w.viewObjects.RemoveView(v.AccountID, char.AccountID)

						w.players.SendPacketToPlayer(v.AccountID, &packets.EndViewObject{
							ID:      types.HardUInt32{Value: char.Base.WorldID.Value},
							SeeType: types.HardUInt8{Value: 0},
						})

						w.viewObjects.RemoveView(char.AccountID, v.AccountID)

						w.players.SendPacketToPlayer(char.AccountID, &packets.EndViewObject{
							ID:      types.HardUInt32{Value: v.Base.WorldID.Value},
							SeeType: types.HardUInt8{Value: 0},
						})
					}

					if w.viewObjects.IsView(char.AccountID, v.AccountID) {
						fmt.Println("From:", packets.ActionPosition{X: oldXPos, Y: oldYPos})
						fmt.Println("To:", packets.ActionPosition{X: types.HardUInt32{Value: cs.X}, Y: types.HardUInt32{Value: cs.Y}})

						w.players.SendPacketToPlayer(char.AccountID, &packets.BeginActionResponse{
							WorldID:    v.Base.WorldID,
							MoveCount:  types.HardUInt32{Value: cs.MoveCount},
							ActionType: types.HardUInt8{Value: cs.ActionType},
							State:      types.HardUInt16{Value: types.MSTATE_ARRIVE},
							StopState:  types.HardUInt16{Value: types.MSTATE_ARRIVE},
							PointSize:  types.HardUInt16{Value: 16},
							From:       packets.ActionPosition{X: oldXPos, Y: oldYPos},
							To:         packets.ActionPosition{X: types.HardUInt32{Value: cs.X}, Y: types.HardUInt32{Value: cs.Y}},
						})
					}
				}
			}()

			fmt.Printf("got a new character step with struct: %v\n", cs)
		}
	}
}

func (w *World) SendCharacterStep(
	accountID string, name string, x, y uint32,
	angle float64, moveCount uint32, actionType uint8,
) {
	w.charactersStepCh <- CharacterStep{accountID, name, x, y, angle, moveCount, actionType}
}

func (v *World) calculateDistance(x1, y1, x2, y2 float64) float64 {
	xPow := math.Pow(x2-x1, 2)
	yPow := math.Pow(y2-y1, 2)

	return math.Sqrt(xPow + yPow)
}

/*
func (v *ViewBroadcast) Listen() {
	go func() {
		for handler := range v.UpdateCharacterCh {
			ch := handler.characterBase

			v.Characters[ch.ChaID.Value] = ch

			chX := float64(ch.Position.X.Value / 100)
			chY := float64(ch.Position.Y.Value / 100)

			for _, otherCh := range v.Characters {
				otherChX := float64(otherCh.Position.X.Value / 100)
				otherChY := float64(otherCh.Position.Y.Value / 100)

				distance := v.calculateDistance(chX, otherChX, chY, otherChY)
				if distance < defaultDistanceView {

				}
			}
		}
	}()
}
*/
