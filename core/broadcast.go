package core

import (
	"sync"
)

const (
	DefaultBroadcastRoomName = "defaultBroadcast"
)

type BroadcastAdaptor interface {
	Join(room string, c Client) error
	Leave(room string, c Client) error
	Send(ignore []Client, room, event string, msg interface{}) error
	Len(room string) uint
}

type rooms map[string]Room

type broadcast struct {
	rooms
	sync.RWMutex
}

func newDefaultBroadcast() BroadcastAdaptor {
	b := &broadcast{
		rooms: make(rooms, 0),
	}
	b.rooms[DefaultBroadcastRoomName] = getRoom()
	return b
}

func (b *broadcast) Join(room string, c Client) error {
	b.RLock()
	r, ok := b.rooms[room]
	b.RUnlock()
	if !ok {
		r = getRoom()
		b.Lock()
		b.rooms[room] = r.SetClient(c)
		b.Unlock()
	}
	return nil
}

func (b *broadcast) Leave(room string, c Client) error {
	b.RLock()
	var r, ok = b.rooms[room]
	b.RUnlock()
	if !ok {
		return nil
	}
	var err = r.RemoveClient(c)
	if err != nil {
		return err
	}
	if r.Len() <= 0 {
		b.Lock()
		delete(b.rooms, room)
		b.Unlock()
	}
	return nil
}

func (b *broadcast) Send(ignore []Client, room, event string, msg interface{}) error {
	b.Lock()
	r, ok := b.rooms[room]
	b.Unlock()
	if !ok {
		return nil
	}
	return r.Send(ignore, event, msg)
}

func (b *broadcast) Len(room string) uint {
	b.Lock()
	r, ok := b.rooms[room]
	b.Unlock()
	if !ok {
		return -1
	}
	return r.Len()
}