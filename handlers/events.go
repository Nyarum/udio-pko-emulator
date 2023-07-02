package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"little/cryptopassword"
	"little/packets"
	"little/packetslogic"
	"little/storage"
	"little/types"
	"log"
	"math"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	defaultDistanceView = float64(24)
	defaultPassword     = "testtest"
)

type Storage interface {
	SaveAccount(sa *storage.StorageAccount) error
	GetAccount(id string) (storage.StorageAccount, bool)
	GetAccountByLogin(login string) (storage.StorageAccount, bool)
	UpdateAccount(sa *storage.StorageAccount) error
}

type Handler struct {
	storage Storage
	players *Players
	world   *World
	conn    net.Conn

	acceptEvCh chan packets.Packet

	account              storage.StorageAccount
	activeCharacter      storage.Character
	activeCharacterIndex int
	firstDate            string
}

func NewHandler(firstDate string, players *Players, world *World, storage Storage, conn net.Conn) *Handler {
	h := &Handler{
		firstDate:  firstDate,
		acceptEvCh: make(chan packets.Packet, 1),
		players:    players,
		storage:    storage,
		world:      world,
		conn:       conn,
	}

	go h.listenEvents()

	return h
}

func (h *Handler) SendPacketMyself(ev packets.Packet) {
	h.acceptEvCh <- ev
}

func (h *Handler) listenEvents() {
	for ev := range h.acceptEvCh {
		log.Printf("Got a new packet for character (%v) and packet opcode: %v ", h.activeCharacter.Name, ev.Opcode())

		//ev.Print()

		resBuf := packets.BuildPacket(ev)

		h.conn.Write(resBuf)
	}
}

func getRadians(x1, x2, y1, y2 uint32) float64 {
	cleanFromX := x1 / 100
	cleanToX := x2 / 100
	cleanFromY := y1 / 100
	cleanToY := y2 / 100

	rad2Deg := 360 / (math.Pi * 2)
	radians := math.Atan2(float64(cleanToX)-float64(cleanFromX), float64(cleanToY)-float64(cleanFromY)) * rad2Deg

	radians = 180 - radians

	return radians
}

func getMD5Password24First(pass []byte) string {
	md5Init := md5.Sum(pass)
	passMd5 := hex.EncodeToString(md5Init[:])

	return string([]byte(passMd5)[:24])
}

func (h *Handler) Ev(ctx context.Context, ev packets.Packet) (res []packets.Packet) {
	sg := h.storage

	//ev.Print()

	switch v := ev.(type) {
	case *packets.Auth:
		defaultPasswordMD5 := md5.Sum([]byte(defaultPassword))

		isNew := false
		sa, ok := sg.GetAccountByLogin(v.Login.Value)
		if !ok {

			fmt.Println("Login:", v.Login.Value)

			sa = storage.StorageAccount{
				Login:      v.Login.Value,
				Password:   strings.ToUpper(hex.EncodeToString(defaultPasswordMD5[:])),
				Created:    time.Now(),
				Characters: []storage.Character{},
			}

			err := sg.SaveAccount(&sa)
			if err != nil {
				fmt.Println(err)
			}

			isNew = true
		}

		encryptPassword, err := cryptopassword.EncryptPassword(sa.Password[:24], h.firstDate)
		if err != nil {
			fmt.Println("Can't encrypt password")
		}

		if hex.EncodeToString([]byte(encryptPassword)) != hex.EncodeToString(v.Password.Value) && !isNew {
			res = append(res, &packets.AuthError{
				ErrorCode: types.HardUInt16{Value: types.ErrAPINVALIDPWD},
			})
			return
		}

		h.account = sa

		var packetCharacters []packets.Character
		for _, ch := range sa.Characters {
			packetCharacters = append(packetCharacters, ch.Character)
		}

		h.players.AddPlayer(sa.ID.Hex(), h)

		res = append(res, &packets.CharactersChoice{
			Characters: packetCharacters,
		})
		return
	case *packets.CharacterCreate:
		fmt.Println("Create character:", v)

		h.account.Characters = append(h.account.Characters, storage.Character{
			ID: rand.Uint32(),
			Character: packets.Character{
				IsActive: types.HardBool{Value: true},
				Name:     v.Name,
				Map:      v.Map,
				Job:      types.HardString{Value: "Newbie"},
				Level:    types.HardUInt16{Value: 1},
				LookSize: v.LookSize,
				Look:     v.Look,
			},
		})

		err := sg.UpdateAccount(&h.account)
		if err != nil {
			fmt.Println("can't save account:", err)

			// TODO: rewrite to use client errors
			res = append(res, &packets.CharacterCreateReply{})
			return
		}

		res = append(res, &packets.CharacterCreateReply{})
		return
	case *packets.EnterGameRequest:
		fmt.Println("Entering by character name:", v.CharacterName)

		for index, ch := range h.account.Characters {
			if ch.Name.Value == v.CharacterName.Value {
				h.activeCharacter = ch
				h.activeCharacterIndex = index
			}
		}

		if h.activeCharacter.ID == 0 {
			fmt.Println("Character by selected name isn't found")
			return nil
		}

		enterGame := &packets.EnterGame{}
		enterGame.Basic()

		fmt.Println("Loggin with character id:", h.activeCharacter.ID)

		chaID := h.activeCharacter.ID

		if h.activeCharacter.Base.ChaID.Value != 0 {
			enterGame.CharacterBase = h.activeCharacter.Base
		} else {
			enterGame.CharacterBase.ChaID = types.HardUInt32{Value: 4}       // 3d object
			enterGame.CharacterBase.Look.TypeID = types.HardUInt16{Value: 4} // look on object
			enterGame.CharacterBase.Icon = types.HardUInt16{Value: 4}        // icon on left and up side
			enterGame.CharacterBase.Name = v.CharacterName
			enterGame.CharacterBase.CommName = v.CharacterName
			enterGame.CharacterBase.Handle = types.HardUInt32{Value: chaID}
			enterGame.CharacterBase.CommID = types.HardUInt32{Value: chaID}
			enterGame.CharacterBase.WorldID = types.HardUInt32{Value: chaID}
			enterGame.CharacterBase.EntityEvent.EntityID = types.HardUInt32{Value: chaID}
			enterGame.CharacterBase.Angle = types.HardUInt16{Value: 0}
		}

		enterGame.ChaMainID = types.HardUInt32{Value: chaID}

		h.activeCharacter.Base = enterGame.CharacterBase
		h.activeCharacter.ID = chaID
		h.activeCharacter.AccountID = h.account.ID.Hex()
		h.activeCharacter.CharacterAttribute = enterGame.CharacterAttribute
		h.activeCharacter.CharacterSkillState = enterGame.CharacterSkillState
		h.account.Characters[h.activeCharacterIndex] = h.activeCharacter

		h.world.CharacterEnterWorld(h.activeCharacter)

		position := h.activeCharacter.Base.Position

		go func() {
			time.Sleep(2000 * time.Millisecond)
			h.world.SendCharacterStep(
				h.account.ID.Hex(), h.activeCharacter.Name.Value,
				position.X.Value, position.Y.Value, float64(h.activeCharacter.Base.Angle.Value),
				0, 1,
			)
		}()

		err := sg.UpdateAccount(&h.account)
		if err != nil {
			log.Printf("can't save account: %v", err)

			// TODO: rewrite to use client errors
			res = append(res, &packets.CharacterCreateReply{})
			return
		}

		res = append(res, enterGame)
		return
	case *packets.KitbagSyncRequest:
		fmt.Println("Request kitbag sync")

		res = append(res, &packets.KitbagSyncResponse{})
		return
	case *packets.SayRequest:
		res = append(res, &packets.SayResponse{
			ID:  types.HardUInt32{Value: h.activeCharacter.ID},
			Msg: v.Msg,
		})

		h.players.Broadcast(&packets.SayResponse{
			ID:  types.HardUInt32{Value: h.activeCharacter.ID},
			Msg: v.Msg,
		})

		return
	case *packets.BeginActionRequest:
		fmt.Println("Request begin action:", v)

		cleanFromX := v.From.X.Value / 100
		cleanFromY := v.From.Y.Value / 100
		cleanToX := v.To.X.Value / 100
		cleanToY := v.To.Y.Value / 100

		radians := getRadians(v.From.X.Value, v.To.X.Value, v.From.Y.Value, v.To.Y.Value)
		steps := packetslogic.SplitFromToOnSteps(cleanFromX, cleanToX, cleanFromY, cleanToY)

		fmt.Println("Angle:", radians)

		lastFrom := packets.ActionPosition{
			X: v.From.X,
			Y: v.From.Y,
		}
		for k, step := range steps {
			h.world.SendCharacterStep(
				h.account.ID.Hex(), h.activeCharacter.Name.Value,
				step.X.Value, step.Y.Value, radians, v.MoveCount.Value,
				v.ActionType.Value,
			)

			if k == len(steps)-1 {
				res = append(res, &packets.BeginActionResponse{
					WorldID:    v.WorldID,
					MoveCount:  v.MoveCount,
					ActionType: v.ActionType,
					State:      types.HardUInt16{Value: types.MSTATE_ARRIVE},
					StopState:  types.HardUInt16{Value: types.MSTATE_ARRIVE},
					PointSize:  types.HardUInt16{Value: 16},
					From:       lastFrom,
					To:         step,
				})
			} else {
				res = append(res, &packets.BeginActionResponse{
					WorldID:    v.WorldID,
					MoveCount:  v.MoveCount,
					ActionType: v.ActionType,
					State:      types.HardUInt16{Value: types.MSTATE_ON},
					PointSize:  types.HardUInt16{Value: 16},
					From:       lastFrom,
					To:         step,
				})
			}

			lastFrom = step
		}

		h.activeCharacter.Base.Position.X.Value = v.To.X.Value
		h.activeCharacter.Base.Position.Y.Value = v.To.Y.Value
		h.activeCharacter.Base.Angle.Value = uint16(radians)

		return
	case *packets.BeginActionEndRequest:
		fmt.Println("Request stop state of action")

		res = append(res, &packets.BeginActionResponse{})
		return
	}

	return nil
}

func storageCharacterContains(s []storage.Character, chIndex int) bool {
	for k := range s {
		if k == chIndex {
			return true
		}
	}

	return false
}

func (h *Handler) Destroy() {
	if storageCharacterContains(h.account.Characters, h.activeCharacterIndex) {
		h.account.Characters[h.activeCharacterIndex] = h.activeCharacter

		err := h.storage.UpdateAccount(&h.account)
		if err != nil {
			log.Printf("can't save account: %v", err)
			return
		}

		h.world.CharacterExitWorld(h.activeCharacter)
	}

	h.players.RemovePlayer(h.account.ID.Hex(), h)
}
