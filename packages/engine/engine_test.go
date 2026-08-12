package engine

import (
	"testing"

	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/finance"
	"github.com/finance-boardgame/engine/profession"
	"github.com/finance-boardgame/engine/ratrace"
)

// ─── helpers ──────────────────────────────────────────────────────────────

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

func hasEvent(events []domain.Event, want domain.EventType) bool {
	for _, ev := range events {
		if ev.Type == want {
			return true
		}
	}
	return false
}

// paydayAmount คืน (amount, true) ถ้ามี EventPayday ใน events
func paydayAmount(events []domain.Event) (domain.Money, bool) {
	for _, ev := range events {
		if ev.Type == domain.EventPayday {
			if a, ok := ev.Data["amount"].(int64); ok {
				return domain.Money(a), true
			}
			return 0, true
		}
	}
	return 0, false
}

// rollUntilPending ทอย (สลับผู้เล่นตาม CurrentTurn) จนกว่าจะตก Opportunity (Pending)
func rollUntilPending(t *testing.T, e *Engine, ids []string, maxRolls int) *domain.PendingDecision {
	t.Helper()
	for i := 0; i < maxRolls; i++ {
		if p := e.State().Pending; p != nil {
			return p
		}
		pid := ids[e.State().CurrentTurn%len(ids)]
		if _, err := e.Apply(domain.Action{PlayerID: pid, Type: domain.ActionRoll}); err != nil {
			t.Fatalf("roll %d: %v", i, err)
		}
	}
	return e.State().Pending
}

// ─── guards (คงจาก Slice 1) ──────────────────────────────────────────────

func TestApply_NotYourTurn_ReturnsError(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0), samplePlayer("p2", 0))
	_, err := e.Apply(domain.Action{PlayerID: "p2", Type: domain.ActionRoll})
	if err == nil {
		t.Error("expected error when acting out of turn, got nil")
	}
}

func TestApply_GameEnded_ReturnsError(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0))
	e.state.Phase = domain.PhaseEnded
	_, err := e.Apply(domain.Action{PlayerID: "p1", Type: domain.ActionRoll})
	if err == nil {
		t.Error("expected error when game ended, got nil")
	}
}

func TestApply_RollKeepsPositionInRange(t *testing.T) {
	e := newTestEngine(7, samplePlayer("p1", 0))
	if _, err := e.Apply(domain.Action{PlayerID: "p1", Type: domain.ActionRoll}); err != nil {
		t.Fatal(err)
	}
	pos := e.State().Players[0].Position
	if pos < 0 || pos >= ratrace.BoardSize {
		t.Errorf("Position = %d, out of [0,%d)", pos, ratrace.BoardSize)
	}
}

// ─── Payday ───────────────────────────────────────────────────────────────

// ผู้เล่นที่ pos 23 → roll ใด ๆ จะ wrap → ต้องมี EventPayday amount = MonthlyCashFlow
func TestApply_Payday_EmitsAmountOnWrap(t *testing.T) {
	p := samplePlayer("p1", 23)
	want := finance.MonthlyCashFlow(p) // 2,000
	e := newTestEngine(42, p)

	events, err := e.Apply(domain.Action{PlayerID: "p1", Type: domain.ActionRoll})
	if err != nil {
		t.Fatalf("roll: %v", err)
	}
	amt, ok := paydayAmount(events)
	if !ok {
		t.Fatal("expected EventPayday when wrapping past index 0 from pos 23")
	}
	if amt != want {
		t.Errorf("payday amount = %d, want %d", amt, want)
	}
}

// ─── Opportunity → decision phase ─────────────────────────────────────────

// ตก Opportunity → ตั้ง Pending + ตอน Pending การ roll ต้อง error
func TestApply_LandOnOpportunity_SetsPendingAndBlocksRoll(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0))
	pending := rollUntilPending(t, e, []string{"p1"}, 80)
	if pending == nil {
		t.Fatal("never landed on Opportunity in 80 rolls")
	}
	// ขณะ Pending, roll ต้อง error (ต้อง resolve ก่อน)
	if _, err := e.Apply(domain.Action{PlayerID: "p1", Type: domain.ActionRoll}); err == nil {
		t.Error("expected error when rolling while a deal is pending")
	}
}

// ซื้อดีล → เพิ่ม Asset + หักเงินดาวน์ + ล้าง Pending + เปลี่ยนเทิร์น
func TestApply_BuyAsset_AddsAssetAndAdvancesTurn(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0), samplePlayer("p2", 0))
	pending := rollUntilPending(t, e, []string{"p1", "p2"}, 80)
	if pending == nil {
		t.Fatal("never landed on Opportunity")
	}

	player := e.playerByID(pending.PlayerID)
	player.Cash = 1_000_000 // ให้เงินพอซื้อทุกดีล
	card := pending.DealCard
	assetsBefore := len(player.Assets)
	turnBefore := e.State().CurrentTurn

	events, err := e.Apply(domain.Action{PlayerID: pending.PlayerID, Type: domain.ActionBuyAsset})
	if err != nil {
		t.Fatalf("buy: %v", err)
	}
	if !hasEvent(events, domain.EventAssetBought) {
		t.Error("expected EventAssetBought")
	}
	after := e.playerByID(pending.PlayerID)
	if len(after.Assets) != assetsBefore+1 {
		t.Errorf("Assets count = %d, want %d", len(after.Assets), assetsBefore+1)
	}
	if after.Cash != 1_000_000-card.DownPayment {
		t.Errorf("Cash = %d, want %d (1,000,000 − down %d)", after.Cash, 1_000_000-card.DownPayment, card.DownPayment)
	}
	if e.State().Pending != nil {
		t.Error("Pending not cleared after buy")
	}
	if e.State().CurrentTurn == turnBefore {
		t.Error("turn not advanced after buy")
	}
}

// ซื้อเมื่อเงินไม่พอ → error (และ Pending ยังคงอยู่)
func TestApply_BuyAsset_CantAffordReturnsError(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0), samplePlayer("p2", 0))
	pending := rollUntilPending(t, e, []string{"p1", "p2"}, 80)
	if pending == nil {
		t.Fatal("never landed on Opportunity")
	}
	e.playerByID(pending.PlayerID).Cash = 100 // ดาวน์ต่ำสุดของดีล = 5,000 → ซื้อไม่ได้

	_, err := e.Apply(domain.Action{PlayerID: pending.PlayerID, Type: domain.ActionBuyAsset})
	if err == nil {
		t.Error("expected error when can't afford down payment, got nil")
	}
	if e.State().Pending == nil {
		t.Error("Pending should remain after failed buy")
	}
}

// ผ่านดีล → ล้าง Pending + เปลี่ยนเทิร์น
func TestApply_Decline_ClearsPendingAndAdvancesTurn(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0), samplePlayer("p2", 0))
	pending := rollUntilPending(t, e, []string{"p1", "p2"}, 80)
	if pending == nil {
		t.Fatal("never landed on Opportunity")
	}
	turnBefore := e.State().CurrentTurn

	_, err := e.Apply(domain.Action{PlayerID: pending.PlayerID, Type: domain.ActionDecline})
	if err != nil {
		t.Fatalf("decline: %v", err)
	}
	if e.State().Pending != nil {
		t.Error("Pending not cleared after decline")
	}
	if e.State().CurrentTurn == turnBefore {
		t.Error("turn not advanced after decline")
	}
}

// ─── Shopping / Crisis ────────────────────────────────────────────────────

// ตก Shopping/Crisis → ต้องมี EventCashChanged ค่า amount ติดลบ
func TestApply_ShoppingOrCrisis_EmitsNegativeCashChange(t *testing.T) {
	e := newTestEngine(42, samplePlayer("p1", 0))
	found := false
	for i := 0; i < 100 && !found; i++ {
		// resolve pending ก่อนถ้ามี
		if e.State().Pending != nil {
			e.Apply(domain.Action{PlayerID: e.State().Pending.PlayerID, Type: domain.ActionDecline})
			continue
		}
		events, err := e.Apply(domain.Action{PlayerID: "p1", Type: domain.ActionRoll})
		if err != nil {
			t.Fatalf("roll %d: %v", i, err)
		}
		for _, ev := range events {
			if ev.Type != domain.EventCashChanged {
				continue
			}
			kind, _ := ev.Data["kind"].(string)
			if kind == "shopping" || kind == "crisis" {
				amt, _ := ev.Data["amount"].(int64)
				if amt >= 0 {
					t.Errorf("%s amount should be negative, got %d", kind, amt)
				}
				found = true
			}
		}
	}
	if !found {
		t.Fatal("never landed on Shopping/Crisis in 100 rolls")
	}
}

// ─── Determinism ──────────────────────────────────────────────────────────

// seed เดียวกัน + action เดียวกัน (decline เมื่อ pending, roll ปกติ) → state เท่ากัน
func TestApply_Deterministic_SameSeedSameActions(t *testing.T) {
	mk := func() *Engine { return newTestEngine(123, samplePlayer("p1", 0), samplePlayer("p2", 5)) }
	a, b := mk(), mk()

	for i := 0; i < 30; i++ {
		sa := a.State()
		var act domain.Action
		if sa.Pending != nil {
			act = domain.Action{PlayerID: sa.Pending.PlayerID, Type: domain.ActionDecline}
		} else {
			pid := []string{"p1", "p2"}[sa.CurrentTurn]
			act = domain.Action{PlayerID: pid, Type: domain.ActionRoll}
		}
		if _, err := a.Apply(act); err != nil {
			t.Fatalf("a step %d: %v", i, err)
		}
		if _, err := b.Apply(act); err != nil {
			t.Fatalf("b step %d: %v", i, err)
		}
	}

	sa, sb := a.State(), b.State()
	for i := range sa.Players {
		if sa.Players[i].Position != sb.Players[i].Position ||
			sa.Players[i].Cash != sb.Players[i].Cash ||
			len(sa.Players[i].Assets) != len(sb.Players[i].Assets) {
			t.Errorf("player %d diverged: pos %d/%d cash %d/%d assets %d/%d",
				i, sa.Players[i].Position, sb.Players[i].Position,
				sa.Players[i].Cash, sb.Players[i].Cash,
				len(sa.Players[i].Assets), len(sb.Players[i].Assets))
		}
	}
	if sa.CurrentTurn != sb.CurrentTurn || sa.Round != sb.Round {
		t.Errorf("turn/round diverged: %d/%d vs %d/%d", sa.CurrentTurn, sa.Round, sb.CurrentTurn, sb.Round)
	}
	if (sa.Pending == nil) != (sb.Pending == nil) {
		t.Error("pending state diverged")
	}
}

// ─── profession / statement (คงจาก Slice 2) ──────────────────────────────

func TestNewWithRandomProfessions_CreatesValidPlayers(t *testing.T) {
	e := NewWithRandomProfessions(42, 2)
	all := profession.All()
	names := map[string]bool{}
	for _, p := range all {
		names[p.Name] = true
	}
	state := e.State()
	if len(state.Players) != 2 {
		t.Fatalf("Players count = %d, want 2", len(state.Players))
	}
	for i, p := range state.Players {
		if !names[p.Name] {
			t.Errorf("player %d: %q is not a real profession", i, p.Name)
		}
		if p.Cash != p.Profession.Savings {
			t.Errorf("player %d (%s): Cash %d != Savings %d", i, p.Name, p.Cash, p.Profession.Savings)
		}
	}
}

func TestStatement_ReturnsBreakdown(t *testing.T) {
	e := NewWithRandomProfessions(42, 3)
	fs, err := e.Statement(0)
	if err != nil {
		t.Fatalf("Statement(0): %v", err)
	}
	p := e.State().Players[0]
	if fs.Tax != p.Profession.Taxes {
		t.Errorf("Tax = %d, want %d", fs.Tax, p.Profession.Taxes)
	}
	if fs.SocialSecurity != p.Profession.SocialSecurity {
		t.Errorf("SocialSecurity = %d, want %d", fs.SocialSecurity, p.Profession.SocialSecurity)
	}
	payments := p.Profession.HomeMortgage.Payment + p.Profession.CarLoan.Payment +
		p.Profession.CreditCard.Payment + p.Profession.SchoolLoan.Payment
	want := p.Profession.Salary - fs.Tax - fs.SocialSecurity - p.Profession.OtherExpenses - payments
	if fs.MonthlyCashFlow != want {
		t.Errorf("MonthlyCashFlow = %d, want %d", fs.MonthlyCashFlow, want)
	}
}

func TestStatement_OutOfRange(t *testing.T) {
	e := NewWithRandomProfessions(42, 2)
	if _, err := e.Statement(5); err == nil {
		t.Error("expected error for out-of-range player index, got nil")
	}
}

// ─── Loans (Slice 5) ──────────────────────────────────────────────────────

// loanPlayer ผู้เล่นทดสอบสินเชื่อ: เงินเดือน 30,000 + บ้าน (ค้ำได้)
func loanPlayer(id string) domain.Player {
	return domain.Player{
		ID: id, Name: id, Cash: 50_000, Position: 0,
		Profession: domain.Profession{
			Name:         "Tester",
			Salary:       30_000,
			HomeMortgage: domain.Liability{Payment: 3_000, Balance: 400_000},
		},
		Assets:      []domain.Asset{},
		Liabilities: []domain.Liability{},
		Loans:       []domain.Loan{},
	}
}

func TestApply_TakeLoanPersonal_AddsLoanAndCash(t *testing.T) {
	e := newTestEngine(42, loanPlayer("p1"))
	_, err := e.Apply(domain.Action{
		PlayerID: "p1", Type: domain.ActionTakeLoan,
		Payload:  map[string]any{"lender": "personal", "amount": float64(50_000)},
	})
	if err != nil {
		t.Fatalf("take loan: %v", err)
	}
	p := e.State().Players[0]
	if len(p.Loans) != 1 {
		t.Fatalf("Loans count = %d, want 1", len(p.Loans))
	}
	if p.Loans[0].Principal != 50_000 {
		t.Errorf("Principal = %d, want 50000", p.Loans[0].Principal)
	}
	if p.Cash != 100_000 {
		t.Errorf("Cash = %d, want 100000 (50000 + 50000 loan)", p.Cash)
	}
}

func TestApply_TakeLoan_OverLimitRejected(t *testing.T) {
	e := newTestEngine(42, loanPlayer("p1")) // salary 30,000 → personal max 150,000
	_, err := e.Apply(domain.Action{
		PlayerID: "p1", Type: domain.ActionTakeLoan,
		Payload:  map[string]any{"lender": "personal", "amount": float64(200_000)},
	})
	if err == nil {
		t.Error("expected over-limit rejection, got nil")
	}
	if len(e.State().Players[0].Loans) != 0 {
		t.Error("loan should not be added on rejection")
	}
}

func TestApply_TakeLoan_NotYourTurn(t *testing.T) {
	e := newTestEngine(42, loanPlayer("p1"), loanPlayer("p2")) // current = p1
	_, err := e.Apply(domain.Action{
		PlayerID: "p2", Type: domain.ActionTakeLoan,
		Payload:  map[string]any{"lender": "personal", "amount": float64(10_000)},
	})
	if err == nil {
		t.Error("expected not-your-turn error, got nil")
	}
}

func TestApply_PayOffLoan_ClearsLoan(t *testing.T) {
	e := newTestEngine(42, loanPlayer("p1"))
	if _, err := e.Apply(domain.Action{
		PlayerID: "p1", Type: domain.ActionTakeLoan,
		Payload:  map[string]any{"lender": "personal", "amount": float64(50_000)},
	}); err != nil {
		t.Fatal(err)
	}
	loanID := e.State().Players[0].Loans[0].ID

	_, err := e.Apply(domain.Action{
		PlayerID: "p1", Type: domain.ActionPayOffLiability,
		Payload:  map[string]any{"loanID": loanID},
	})
	if err != nil {
		t.Fatalf("pay off: %v", err)
	}
	p := e.State().Players[0]
	if len(p.Loans) != 0 {
		t.Errorf("Loans count = %d, want 0 after payoff", len(p.Loans))
	}
	// cash = 100000 (หลังกู้) − 74000 (balance) = 26000
	if p.Cash != 26_000 {
		t.Errorf("Cash = %d, want 26000", p.Cash)
	}
}
