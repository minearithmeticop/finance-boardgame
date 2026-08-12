package engine

import (
	"testing"

	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/finance"
	"github.com/finance-boardgame/engine/ratrace"
)

// ─── helpers ──────────────────────────────────────────────────────────────

// newTestEngine สร้าง engine พร้อมผู้เล่นที่กำหนด (seed ตายตัว → deterministic)
func newTestEngine(seed int64, players ...domain.Player) *Engine {
	return New(Config{Seed: seed, Players: players})
}

// samplePlayer ผู้เล่นทดสอบ: เงินเดือน 2,000 ไม่มีรายจ่าย → MonthlyCashFlow = 2,000
func samplePlayer(id string, pos int) domain.Player {
	return domain.Player{
		ID:       id,
		Name:     id,
		Cash:     1000,
		Position: pos,
		Profession: domain.Profession{
			Name:   "Tester",
			Salary: domain.Money(2000),
		},
	}
}

// mustRoll ทอยเต๋าให้ผู้เล่นที่กำหนด ถ้า error ให้ fail test ทันที
func mustRoll(t *testing.T, e *Engine, playerID string) {
	t.Helper()
	if _, err := e.Apply(domain.Action{PlayerID: playerID, Type: domain.ActionRoll}); err != nil {
		t.Fatalf("Apply roll for %s: %v", playerID, err)
	}
}

// hasEvent เช็คว่ามี event type ที่ต้องการในลิสต์หรือไม่
func hasEvent(events []domain.Event, want domain.EventType) bool {
	for _, ev := range events {
		if ev.Type == want {
			return true
		}
	}
	return false
}

// ─── tests ────────────────────────────────────────────────────────────────

// TestApply_RollAdvancesToNextPlayer — หลัง roll เทิร์นต้องเลื่อนไปผู้เล่นถัดไป
func TestApply_RollAdvancesToNextPlayer(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0), samplePlayer("p2", 0))

	mustRoll(t, e, "p1")

	if got := e.State().CurrentTurn; got != 1 {
		t.Errorf("CurrentTurn = %d, want 1", got)
	}
}

// TestApply_FullLapIncrementsRound — ทุกผู้เล่น roll ครั้งละ 1 = ครบ 1 รอบ → Round = 1
func TestApply_FullLapIncrementsRound(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0), samplePlayer("p2", 0))

	mustRoll(t, e, "p1")
	mustRoll(t, e, "p2")

	if e.State().Round != 1 {
		t.Errorf("Round = %d, want 1", e.State().Round)
	}
	if e.State().CurrentTurn != 0 {
		t.Errorf("CurrentTurn = %d, want 0 (wrapped)", e.State().CurrentTurn)
	}
}

// TestApply_PassPayday_ReceivesMonthlyCashFlow
// ผู้เล่นอยู่ช่องสุดท้าย (23) — roll อะไรก็ตาม (1-6) จะ wrap ผ่าน Payday (index 0)
// → เงินสดต้องเพิ่มขึ้นเท่ากับ MonthlyCashFlow และ position ต้องอยู่ใน [0,5]
func TestApply_PassPayday_ReceivesMonthlyCashFlow(t *testing.T) {
	p := samplePlayer("p1", 23)
	expectedPay := finance.MonthlyCashFlow(p) // 2,000
	startCash := p.Cash                        // 1,000

	e := newTestEngine(42, p)
	events, err := e.Apply(domain.Action{PlayerID: "p1", Type: domain.ActionRoll})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	got := e.State().Players[0]
	if got.Cash != startCash+expectedPay {
		t.Errorf("Cash = %d, want %d (start %d + payday %d)", got.Cash, startCash+expectedPay, startCash, expectedPay)
	}
	if got.Position < 0 || got.Position > 5 {
		t.Errorf("Position = %d, want 0..5 (wrapped past Payday at index 0)", got.Position)
	}
	if !hasEvent(events, domain.EventPayday) {
		t.Error("expected EventPayday in emitted events")
	}
}

// TestApply_RollKeepsPositionInRange — หลัง roll position ต้องอยู่ใน [0, BoardSize)
func TestApply_RollKeepsPositionInRange(t *testing.T) {
	e := newTestEngine(7, samplePlayer("p1", 0))

	mustRoll(t, e, "p1")

	pos := e.State().Players[0].Position
	if pos < 0 || pos >= ratrace.BoardSize {
		t.Errorf("Position = %d, out of [0,%d)", pos, ratrace.BoardSize)
	}
}

// TestApply_Deterministic_SameSeedSameActions
// engine สองตัว seed เดียวกัน + actions เดียวกัน → state เท่ากันทุกด้าน (replay ได้)
func TestApply_Deterministic_SameSeedSameActions(t *testing.T) {
	a := newTestEngine(123, samplePlayer("p1", 0), samplePlayer("p2", 5))
	b := newTestEngine(123, samplePlayer("p1", 0), samplePlayer("p2", 5))

	for i := 0; i < 4; i++ {
		pid := []string{"p1", "p2"}[a.State().CurrentTurn]
		mustRoll(t, a, pid)
		mustRoll(t, b, pid)
	}

	sa, sb := a.State(), b.State()
	for i := range sa.Players {
		if sa.Players[i].Position != sb.Players[i].Position {
			t.Errorf("player %d position: %d vs %d", i, sa.Players[i].Position, sb.Players[i].Position)
		}
		if sa.Players[i].Cash != sb.Players[i].Cash {
			t.Errorf("player %d cash: %d vs %d", i, sa.Players[i].Cash, sb.Players[i].Cash)
		}
	}
	if sa.CurrentTurn != sb.CurrentTurn {
		t.Errorf("CurrentTurn: %d vs %d", sa.CurrentTurn, sb.CurrentTurn)
	}
	if sa.Round != sb.Round {
		t.Errorf("Round: %d vs %d", sa.Round, sb.Round)
	}
}

// TestApply_NotYourTurn_ReturnsError — ส่ง action ของคนที่ไม่ใช่เทิร์นปัจจุบัน → error
func TestApply_NotYourTurn_ReturnsError(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0), samplePlayer("p2", 0))

	_, err := e.Apply(domain.Action{PlayerID: "p2", Type: domain.ActionRoll})
	if err == nil {
		t.Error("expected error when acting out of turn, got nil")
	}
}

// TestApply_GameEnded_ReturnsError — เกมจบแล้ว → action ใด ๆ ต้อง error
func TestApply_GameEnded_ReturnsError(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0))
	e.state.Phase = domain.PhaseEnded

	_, err := e.Apply(domain.Action{PlayerID: "p1", Type: domain.ActionRoll})
	if err == nil {
		t.Error("expected error when game ended, got nil")
	}
}
