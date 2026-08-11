// Package room จัดการ lobby/room สำหรับ Online Multiplayer (WebSocket)
//
// TODO(Session#3): implement Hub + Client + room lifecycle:
//   - Hub: registry ของทุกห้อง, รับ/ส่ง message ระหว่าง clients
//   - Room: ถือ engine instance หนึ่งตัว, broadcast events ให้สมาชิก
//   - Client: wrapper ของ WebSocket connection หนึ่งเส้น
package room

import "sync"

// Hub คือศูนย์กลางจัดการห้องทั้งหมดใน backend instance
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

// Room คือห้องเกมหนึ่งห้อง (ผูกกับ engine instance หนึ่งตัว)
type Room struct {
	ID      string
	// TODO: players []*Client, engine *engine.Engine, broadcast chan events
}

// NewHub สร้าง hub ใหม่
func NewHub() *Hub {
	return &Hub{rooms: make(map[string]*Room)}
}

// CreateRoom สร้างห้องใหม่แล้วคืน pointer
// TODO(Session#3): implement จริง (gen ID, init engine)
func (h *Hub) CreateRoom(id string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	r := &Room{ID: id}
	h.rooms[id] = r
	return r
}

// GetRoom คืนห้องตาม id (nil ถ้าไม่มี)
func (h *Hub) GetRoom(id string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[id]
}
