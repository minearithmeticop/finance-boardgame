// Package engine คือ orchestrator หลักของเกม — รับ Config, ถือ GameState,
// และ apply Action แบบ deterministic เพื่อผลิต Event
//
// Engine เป็น pure logic (ไม่มี I/O) จึง compile และรันได้ทั้งฝั่ง server (apps/backend)
// และฝั่ง browser (ผ่าน cmd/wasm → WebAssembly)
package engine

import (
	"errors"

	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/rng"
)

// Version ของ engine (ใช้ใน health check / WASM bridge)
const Version = "0.1.0"

// Config — การตั้งค่าเริ่มต้นของเกม
type Config struct {
	Seed    int64
	Players []domain.Player
}

// Engine ถือสถานะเกมและ RNG ใช้งานผ่าน New() เท่านั้น
type Engine struct {
	state domain.GameState
	rng   *rng.RNG
}

// New สร้าง engine ใหม่จาก config (เริ่มเสมอที่ Rat Race, รอบที่ 0)
func New(cfg Config) *Engine {
	players := cfg.Players
	if players == nil {
		players = []domain.Player{}
	}
	return &Engine{
		state: domain.GameState{
			Phase:       domain.PhaseRatRace,
			Players:     players,
			CurrentTurn: 0,
			Round:       0,
			Seed:        cfg.Seed,
		},
		rng: rng.New(cfg.Seed),
	}
}

// State คืนสำเนาสถานะเกมปัจจุบัน (snapshot)
func (e *Engine) State() domain.GameState { return e.state }

// Apply นำ action มาประมวลผลแล้วคืน events ที่เกิดขึ้น
//
// TODO(Session#3): dispatch ไปยัง ratrace/fasttrack ตาม phase ของผู้เล่นที่ส่ง action
func (e *Engine) Apply(action domain.Action) ([]domain.Event, error) {
	if e.state.Phase == domain.PhaseEnded {
		return nil, errors.New("engine: game already ended")
	}
	_ = action
	return nil, nil
}
