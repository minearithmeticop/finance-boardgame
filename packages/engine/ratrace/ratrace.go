// Package ratrace implements Rat Race (วงจรหนูถีบจักร) game rules
package ratrace

import "github.com/finance-boardgame/engine/domain"

// BoardSize — จำนวนช่องบนกระดาน Rat Race
const BoardSize = 24

// DefaultBoard คืน layout กระดาน Rat Race (24 ช่อง, Payday index 0)
//
// ออกแบบให้สมดุลต่อรอบ: บวก/โอกาส (Opportunity 7 + Windfall 3 + SideJob 2 = 12)
// มากกว่าลบ (Shopping 3 + Crisis 1 + Family/Health/Donate = ~5) + News/Learn เป็น flavor
// **Crisis แค่ 1 ช่อง/รอบ** (ตามดีไซน์)
func DefaultBoard() []domain.Tile {
	return []domain.Tile{
		{Type: domain.TilePayday, Name: "Payday"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileWindfall, Name: "Windfall"},
		{Type: domain.TileNews, Name: "News"},
		{Type: domain.TileSideJob, Name: "SideJob"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileShopping, Name: "Shopping"},
		{Type: domain.TileWindfall, Name: "Windfall"},
		{Type: domain.TileNews, Name: "News"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileFamily, Name: "Family"},
		{Type: domain.TileSideJob, Name: "SideJob"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileShopping, Name: "Shopping"},
		{Type: domain.TileWindfall, Name: "Windfall"},
		{Type: domain.TileNews, Name: "News"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileDonate, Name: "Donate"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileShopping, Name: "Shopping"},
		{Type: domain.TileLearn, Name: "Learn"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileHealth, Name: "Health"},
		{Type: domain.TileCrisis, Name: "Crisis"},
	}
}

// MovePlayer เดินผู้เล่นตามจำนวนลูกเต๋า แล้วคืนตำแหน่งปลายทาง (wraps รอบกระดาน)
// ไม่ resolve เหตุการณ์ของช่องปลายทาง — เป็นหน้าที่ของ engine.Apply
func MovePlayer(p *domain.Player, steps int) int {
	if p == nil {
		return 0
	}
	p.Position = (p.Position + steps) % BoardSize
	if p.Position < 0 {
		p.Position += BoardSize
	}
	return p.Position
}
