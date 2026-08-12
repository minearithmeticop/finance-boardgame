// Package engine คือ orchestrator หลักของเกม — รับ Config, ถือ GameState,
// และ apply Action แบบ deterministic เพื่อผลิต Event
//
// Engine เป็น pure logic (ไม่มี I/O) จึง compile และรันได้ทั้งฝั่ง server (apps/backend)
// และฝั่ง browser (ผ่าน cmd/wasm → WebAssembly)
package engine

import (
	"errors"
	"fmt"

	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/finance"
	"github.com/finance-boardgame/engine/profession"
	"github.com/finance-boardgame/engine/ratrace"
	"github.com/finance-boardgame/engine/rng"
)

// Version ของ engine (ใช้ใน health check / WASM bridge)
const Version = "0.1.0"

// Config — การตั้งค่าเริ่มต้นของเกม
type Config struct {
	Seed    int64
	Players []domain.Player
}

// Engine ถือสถานะเกม, RNG และกระดาน ใช้งานผ่าน New() เท่านั้น
type Engine struct {
	state domain.GameState
	rng   *rng.RNG
	board []domain.Tile // กระดานของเฟสปัจจุบัน (เริ่มที่ Rat Race)
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
		rng:   rng.New(cfg.Seed),
		board: ratrace.DefaultBoard(),
	}
}

// NewWithRandomProfessions สร้างเกมใหม่โดยสุ่ม "อาชีพจริง" ให้ผู้เล่น count คน
// (สำหรับ "จบใหม่สุ่มอาชีพ") — เงินสดเริ่มต้นของแต่ละคน = เงินออมของอาชีพนั้น
func NewWithRandomProfessions(seed int64, count int) *Engine {
	if count < 1 {
		count = 1
	}
	r := rng.New(seed)
	players := make([]domain.Player, count)
	for i := range players {
		prof := profession.Random(r)
		players[i] = domain.Player{
			ID:         fmt.Sprintf("p%d", i+1),
			Name:       prof.Name,
			Cash:       prof.Savings,
			Profession: prof,
			Position:   0,
		}
	}
	return New(Config{Seed: seed, Players: players})
}

// State คืนสำเนาสถานะเกมปัจจุบัน (snapshot)
func (e *Engine) State() domain.GameState { return e.state }

// Statement คืนงบการเงินเต็มรูปแบบ (Income Statement + Balance Sheet) ของผู้เล่น index ที่กำหนด
func (e *Engine) Statement(idx int) (domain.FinancialStatement, error) {
	if idx < 0 || idx >= len(e.state.Players) {
		return domain.FinancialStatement{}, fmt.Errorf("engine: player index %d out of range", idx)
	}
	return finance.Statement(e.state.Players[idx]), nil
}

// Apply นำ action มาประมวลผลแล้วคืน events ที่เกิดขึ้น
func (e *Engine) Apply(action domain.Action) ([]domain.Event, error) {
	if e.state.Phase == domain.PhaseEnded {
		return nil, errors.New("engine: game already ended")
	}

	switch action.Type {
	case domain.ActionRoll:
		return e.applyRoll(action)
	default:
		// TODO(Slice 3+): ActionBuyAsset, ActionSellAsset, ActionPayOffLiability, ...
		return nil, fmt.Errorf("engine: action type %d not yet supported", action.Type)
	}
}

// applyRoll จัดการเทิร์นแบบ roll:
// ทอยเต๋า → เดิน → resolve ช่องปลายทาง (Payday) → emit events → เปลี่ยนเทิร์น
//
// Slice 1: auto-advance หลัง roll (ยังไม่มี decision phase — จะเกิดตอน Opportunity ใน Slice 3)
func (e *Engine) applyRoll(action domain.Action) ([]domain.Event, error) {
	current := e.state.CurrentTurn
	player := &e.state.Players[current]

	// ตรวจว่าเป็นตาของผู้เล่นที่ส่ง action จริงหรือไม่ (คืน error ก่อนทอย → ไม่ทำลาย determinism)
	if action.PlayerID != player.ID {
		return nil, fmt.Errorf("engine: not %s's turn (current player is %s)", action.PlayerID, player.ID)
	}

	startPos := player.Position
	die := e.rng.RollDice(6)
	ratrace.MovePlayer(player, die)

	var events []domain.Event
	events = append(events,
		domain.Event{Type: domain.EventMoved, PlayerID: player.ID},
		domain.Event{Type: domain.EventLanded, PlayerID: player.ID},
	)

	// Payday rule: ผ่าน "หรือ" ตกช่อง Payday (index 0) → ได้ Monthly Cash Flow
	// เนื่องจาก Payday อยู่ index 0, การ wrap กระดาน (startPos+die >= len) = ผ่าน index 0
	if startPos+die >= len(e.board) {
		pay := finance.MonthlyCashFlow(*player)
		player.Cash += pay
		events = append(events,
			domain.Event{Type: domain.EventPayday, PlayerID: player.ID},
			domain.Event{Type: domain.EventCashChanged, PlayerID: player.ID},
		)
	}

	// เปลี่ยนเทิร์น (Round++ เมื่อครบรอบผู้เล่น)
	next := current + 1
	if next >= len(e.state.Players) {
		next = 0
		e.state.Round++
	}
	e.state.CurrentTurn = next

	return events, nil
}
