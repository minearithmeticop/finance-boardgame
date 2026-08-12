// Package engine คือ orchestrator หลักของเกม — รับ Config, ถือ GameState,
// และ apply Action แบบ deterministic เพื่อผลิต Event
//
// Engine เป็น pure logic (ไม่มี I/O) จึง compile และรันได้ทั้งฝั่ง server (apps/backend)
// และฝั่ง browser (ผ่าน cmd/wasm → WebAssembly)
package engine

import (
	"errors"
	"fmt"

	"github.com/finance-boardgame/engine/cards"
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
			ID:          fmt.Sprintf("p%d", i+1),
			Name:        prof.Name,
			Cash:        prof.Savings,
			Profession:  prof,
			Position:    0,
			Assets:      []domain.Asset{},
			Liabilities: []domain.Liability{},
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
// Apply นำ action มาประมวลผลแล้วคืน events ที่เกิดขึ้น
//
// State machine: ถ้ากำลังรอ decision (Pending) → รับเฉพาะ BuyAsset/Decline
// มิฉะนั้น → รับ Roll (อย่างอื่น = Slice 4+)
func (e *Engine) Apply(action domain.Action) ([]domain.Event, error) {
	if e.state.Phase == domain.PhaseEnded {
		return nil, errors.New("engine: game already ended")
	}

	if e.state.Pending != nil {
		switch action.Type {
		case domain.ActionBuyAsset:
			return e.applyBuyAsset(action)
		case domain.ActionDecline:
			return e.applyDecline(action)
		default:
			return nil, errors.New("engine: resolve pending deal first (buy or decline)")
		}
	}

	switch action.Type {
	case domain.ActionRoll:
		return e.applyRoll(action)
	default:
		// TODO(Slice 4+): ActionSellAsset, ActionPayOffLiability, ...
		return nil, fmt.Errorf("engine: action type %d not yet supported", action.Type)
	}
}

// applyRoll จัดการเทิร์นแบบ roll: ทอย → เดิน → Payday(ถ้าผ่าน) → resolve ช่องปลายทาง
//
// Opportunity → ตั้ง Pending และ **ยังไม่เปลี่ยนเทิร์น** (รอ decision)
// Shopping/Crisis → หักเงินอัตโนมัติ แล้วเปลี่ยนเทิร์น
func (e *Engine) applyRoll(action domain.Action) ([]domain.Event, error) {
	current := e.state.CurrentTurn
	player := &e.state.Players[current]

	if action.PlayerID != player.ID {
		return nil, fmt.Errorf("engine: not %s's turn (current player is %s)", action.PlayerID, player.ID)
	}

	startPos := player.Position
	die := e.rng.RollDice(6)
	ratrace.MovePlayer(player, die)

	events := []domain.Event{eventWith(domain.EventMoved, player.ID, map[string]any{
		"steps": die, "position": player.Position,
	})}

	// Payday: ผ่าน/ตก index 0 (Payday ที่ index 0 → wrap กระดาน = ผ่าน index 0)
	if startPos+die >= len(e.board) {
		pay := finance.MonthlyCashFlow(*player)
		player.Cash += pay
		events = append(events, eventWith(domain.EventPayday, player.ID, map[string]any{
			"amount": int64(pay),
		}))
	}

	// Resolve tile ปลายทาง
	switch e.board[player.Position].Type {
	case domain.TileOpportunity:
		card := cards.DrawDealCard(e.rng)
		e.state.Pending = &domain.PendingDecision{PlayerID: player.ID, DealCard: card}
		events = append(events, eventWith(domain.EventLanded, player.ID, map[string]any{
			"kind": "opportunity", "title": card.Title,
		}))
		return events, nil // รอ decision — ห้าม advance turn
	case domain.TileShopping:
		dc := cards.DrawDoodadCard(e.rng)
		player.Cash -= dc.Cost
		events = append(events, eventWith(domain.EventCashChanged, player.ID, map[string]any{
			"kind": "shopping", "title": dc.Title, "amount": int64(-dc.Cost),
		}))
	case domain.TileCrisis:
		cc := cards.DrawCrisisCard(e.rng)
		player.Cash -= cc.Amount
		events = append(events, eventWith(domain.EventCashChanged, player.ID, map[string]any{
			"kind": "crisis", "title": cc.Title, "amount": int64(-cc.Amount),
		}))
	default:
		// Market/Baby/Charity/Downsizing/Blank — Slice 4
		events = append(events, eventWith(domain.EventLanded, player.ID, map[string]any{
			"kind": "noop", "tile": e.board[player.Position].Name,
		}))
	}

	e.advanceTurn()
	return events, nil
}

// applyBuyAsset ซื้อดีลที่กำลังตัดสินใจ (Pending) → เพิ่ม Asset + Liability + หักเงินดาวน์
func (e *Engine) applyBuyAsset(action domain.Action) ([]domain.Event, error) {
	pending := e.state.Pending
	if pending == nil {
		return nil, errors.New("engine: no pending deal to buy")
	}
	if action.PlayerID != pending.PlayerID {
		return nil, fmt.Errorf("engine: not %s's decision (pending is %s)", action.PlayerID, pending.PlayerID)
	}
	player := e.playerByID(pending.PlayerID)
	if player == nil {
		return nil, errors.New("engine: pending player not found")
	}

	card := pending.DealCard
	if player.Cash < card.DownPayment {
		return nil, fmt.Errorf("engine: can't afford down payment %d (cash %d)", card.DownPayment, player.Cash)
	}

	loan := card.Cost - card.DownPayment
	player.Assets = append(player.Assets, domain.Asset{
		ID:            fmt.Sprintf("a%d-%d", e.state.Round, len(player.Assets)),
		Type:          card.AssetType,
		Name:          card.Title,
		CashFlow:      card.CashFlow,
		Cost:          card.Cost,
		DownPayment:   card.DownPayment,
		LoanRemaining: loan,
	})
	if loan > 0 {
		player.Liabilities = append(player.Liabilities, domain.Liability{
			ID:      fmt.Sprintf("l%d-%d", e.state.Round, len(player.Liabilities)),
			Name:    card.Title + " (เงินกู้)",
			Payment: card.LoanPayment,
			Balance: loan,
		})
	}
	player.Cash -= card.DownPayment
	e.state.Pending = nil
	e.advanceTurn()

	return []domain.Event{eventWith(domain.EventAssetBought, player.ID, map[string]any{
		"title": card.Title, "down": int64(card.DownPayment),
	})}, nil
}

// applyDecline ผ่านดีลที่กำลังตัดสินใจ (Pending)
func (e *Engine) applyDecline(action domain.Action) ([]domain.Event, error) {
	pending := e.state.Pending
	if pending == nil {
		return nil, errors.New("engine: no pending deal to decline")
	}
	if action.PlayerID != pending.PlayerID {
		return nil, fmt.Errorf("engine: not %s's decision (pending is %s)", action.PlayerID, pending.PlayerID)
	}
	pid := pending.PlayerID
	title := pending.DealCard.Title
	e.state.Pending = nil
	e.advanceTurn()
	return []domain.Event{eventWith(domain.EventLanded, pid, map[string]any{
		"kind": "declined", "title": title,
	})}, nil
}

// advanceTurn เลื่อนเทิร์นไปผู้เล่นถัดไป (Round++ เมื่อครบรอบ)
func (e *Engine) advanceTurn() {
	next := e.state.CurrentTurn + 1
	if next >= len(e.state.Players) {
		next = 0
		e.state.Round++
	}
	e.state.CurrentTurn = next
}

// playerByID คืน pointer ผู้เล่นตาม ID (nil ถ้าไม่พบ)
func (e *Engine) playerByID(id string) *domain.Player {
	for i := range e.state.Players {
		if e.state.Players[i].ID == id {
			return &e.state.Players[i]
		}
	}
	return nil
}

// eventWith สร้าง Event พร้อม payload (Data)
func eventWith(t domain.EventType, playerID string, data map[string]any) domain.Event {
	return domain.Event{Type: t, PlayerID: playerID, Data: data}
}
