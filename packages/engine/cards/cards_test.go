package cards

import (
	"testing"

	"github.com/finance-boardgame/engine/rng"
)

func TestDecks_NotEmpty(t *testing.T) {
	if len(DealCards()) == 0 {
		t.Error("DealCards() is empty")
	}
	if len(DoodadCards()) == 0 {
		t.Error("DoodadCards() is empty")
	}
	if len(CrisisCards()) == 0 {
		t.Error("CrisisCards() is empty")
	}
}

// TestDealCards_Decks_ReturnCopy — แก้ผลลัพธ์ต้องไม่กระทบชุดต้นฉบับ
func TestDealCards_Decks_ReturnCopy(t *testing.T) {
	a := DealCards()
	a[0].Title = "ถูกแก้"
	if DealCards()[0].Title == "ถูกแก้" {
		t.Error("DealCards() returned underlying slice (mutation leaked)")
	}
}

// TestDealCards_HasGoodAndBadDeals — ต้องมีทั้งดีล net บวก (ดี) และ net ลบ/ศูนย์ (แยก)
// ไม่งั้นผู้เล่นไม่ได้เรียนรู้การประเมินดีล
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
		t.Error("no positive-net deal — need good deals to buy")
	}
	if !hasNonPos {
		t.Error("no non-positive deal — need bad/neutral deals to teach evaluation")
	}
}

// TestDraw_ReturnsMemberOfDeck — Draw ต้องคืนการ์ดที่อยู่ในชุดเสมอ
func TestDraw_ReturnsMemberOfDeck(t *testing.T) {
	mk := func(titles []string) map[string]bool {
		m := map[string]bool{}
		for _, s := range titles {
			m[s] = true
		}
		return m
	}
	deals := mk(func() []string {
		out := make([]string, 0)
		for _, c := range DealCards() {
			out = append(out, c.Title)
		}
		return out
	}())
	doodads := mk(func() []string {
		out := make([]string, 0)
		for _, c := range DoodadCards() {
			out = append(out, c.Title)
		}
		return out
	}())
	crisis := mk(func() []string {
		out := make([]string, 0)
		for _, c := range CrisisCards() {
			out = append(out, c.Title)
		}
		return out
	}())

	r := rng.New(7)
	for i := 0; i < 100; i++ {
		if !deals[DrawDealCard(r).Title] {
			t.Error("DrawDealCard returned unknown card")
		}
		if !doodads[DrawDoodadCard(r).Title] {
			t.Error("DrawDoodadCard returned unknown card")
		}
		if !crisis[DrawCrisisCard(r).Title] {
			t.Error("DrawCrisisCard returned unknown card")
		}
	}
}
