package handlers

import (
	"fmt"
	"little/packets"
	"log"
	"sync"
)

type Players struct {
	players map[string]*Handler
	sync.RWMutex
}

func NewPlayers() *Players {
	return &Players{
		players: make(map[string]*Handler),
		RWMutex: sync.RWMutex{},
	}
}

func (p *Players) Broadcast(ev packets.Packet) {
	p.RLock()
	defer p.RUnlock()

	for _, p := range p.players {
		p.SendPacketMyself(ev)
	}
}

func (p *Players) SendPacketToPlayer(id string, ev packets.Packet) {
	p.RLock()
	defer p.RUnlock()

	log.Printf("Send packet to player (%v) with packet opcode: %v\n", id, ev.Opcode())

	if v, ok := p.players[id]; !ok {
		fmt.Printf("Can't find a player (%v) to send packet", id)
	} else {
		v.SendPacketMyself(ev)
	}
}

func (p *Players) AddPlayer(id string, h *Handler) {
	p.Lock()
	defer p.Unlock()

	log.Printf("Added a new player with id: %v\n", id)

	p.players[id] = h
}

func (p *Players) GetPlayers() map[string]*Handler {
	p.RLock()
	defer p.RUnlock()

	return p.players
}

func (p *Players) RemovePlayer(id string, h *Handler) {
	p.Lock()
	defer p.Unlock()

	log.Printf("Removed a player with id: %v\n", id)

	delete(p.players, id)
}
