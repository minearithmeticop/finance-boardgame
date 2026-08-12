package cards

import (
	"testing"

	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/rng"
)

// ─── DealCards ────────────────────────────────────────────────────────────

func TestDealCards_NotEmpty(t *testing.T) {
	if len(DealCards()) == 0 {
		t.Error("DealCards() is empty")
	}
}

// TestDealCards_ReturnsCopy — แก้ผลลัพธ์ต้องไม่กระทบชุดต้นฉบับ
func TestDealCards_ReturnsCopy(t *testing.T) {
	a := DealCards()
	a[0].Title = "ถูกแก้"
	if DealCards()[0].Title == "ถูกแก้" {
		t.Error("DealCards() returned underlying slice (mutation leaked)")
	}
}

// TestDealCards_HasGoodAndBadDeals — ต้องมีทั้งดีล net บวก และ net ลบ/ศูนย์
func TestDealCards_HasGoodAndBadDeals(t *testing.T) {
	hasPos, hasNonPos := false, false
	for _, c := range DealCards() {
		switch {
		case c.NetCashFlow() > 0:
			hasPos = true
		default:
			hasNonPos = true
		}
	}
	if !hasPos {
		t.Error("no positive-net deal")
	}
	if !hasNonPos {
		t.Error("no non-positive deal — need bad/neutral deals to teach evaluation")
	}
}

// ─── LifeEvents ───────────────────────────────────────────────────────────

// TestLifeEvents_AllCategoriesNotEmpty — ทุกหมวดต้องมีเนื้อหาอย่างน้อย 5 ใบ (กันซ้ำจนเบื่อ)
func TestLifeEvents_AllCategoriesNotEmpty(t *testing.T) {
	for _, cat := range LifeCategories() {
		deck := LifeEvents(cat)
		if len(deck) < 5 {
			t.Errorf("category %q has only %d events, want >= 5", cat, len(deck))
		}
	}
}

// TestLifeEvents_CategoryTagMatches — field Category ของแต่ละใบต้องตรงกับ key ของ deck
func TestLifeEvents_CategoryTagMatches(t *testing.T) {
	for _, cat := range LifeCategories() {
		for _, ev := range LifeEvents(cat) {
			if ev.Category != cat {
				t.Errorf("event in deck %q has Category=%q", cat, ev.Category)
			}
		}
	}
}

// TestLifeEvents_NewsIsFlavor — ข่าวสารทั้งหมดต้อง Amount = 0 (แค่ flavor)
func TestLifeEvents_NewsIsFlavor(t *testing.T) {
	for _, ev := range LifeEvents(domain.CatNews) {
		if ev.Amount != 0 {
			t.Errorf("news %q has Amount %d, want 0 (flavor only)", ev.Title, ev.Amount)
		}
	}
}

// TestLifeEvents_SignsCorrect — windfall/sidejob บวก, shopping/family/donate/health/crisis ลบหรือศูนย์
func TestLifeEvents_SignsCorrect(t *testing.T) {
	positives := []string{domain.CatWindfall, domain.CatSideJob}
	negatives := []string{domain.CatShopping, domain.CatFamily, domain.CatDonate, domain.CatHealth, domain.CatCrisis}
	for _, cat := range positives {
		for _, ev := range LifeEvents(cat) {
			if ev.Amount <= 0 {
				t.Errorf("%s %q Amount %d should be positive", cat, ev.Title, ev.Amount)
			}
		}
	}
	for _, cat := range negatives {
		for _, ev := range LifeEvents(cat) {
			if ev.Amount > 0 {
				t.Errorf("%s %q Amount %d should be <= 0", cat, ev.Title, ev.Amount)
			}
		}
	}
}

// ─── Draw ─────────────────────────────────────────────────────────────────

func TestDrawDealCard_ReturnsMember(t *testing.T) {
	names := map[string]bool{}
	for _, c := range DealCards() {
		names[c.Title] = true
	}
	r := rng.New(7)
	for i := 0; i < 100; i++ {
		if !names[DrawDealCard(r).Title] {
			t.Error("DrawDealCard returned unknown card")
		}
	}
}

func TestDrawLifeEvent_ReturnsMember(t *testing.T) {
	r := rng.New(7)
	for _, cat := range LifeCategories() {
		titles := map[string]bool{}
		for _, ev := range LifeEvents(cat) {
			titles[ev.Title] = true
		}
		for i := 0; i < 50; i++ {
			ev := DrawLifeEvent(r, cat)
			if !titles[ev.Title] {
				t.Errorf("DrawLifeEvent(%q) returned %q (not in deck)", cat, ev.Title)
			}
		}
	}
}

// TestDrawLifeEvent_UnknownCategorySafe — หมวดที่ไม่มี deck ต้องไม่ panic
func TestDrawLifeEvent_UnknownCategorySafe(t *testing.T) {
	r := rng.New(1)
	ev := DrawLifeEvent(r, "nope")
	if ev.Title == "" {
		t.Error("expected fallback event for unknown category")
	}
}
