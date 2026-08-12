// Package ratrace implements Rat Race (วงจรหนูถีบจักร) game rules
// เป็นวงเล็กที่ผู้เล่นเริ่มเกม ต้องสะสม passive income จนเท่ากับหรือมากกว่า expenses
//
// หมายเหตุ: TileType/Tile ถูกย้ายไปอยู่ที่ package domain เพื่อให้ engine และ fasttrack
// ใช้ร่วมกันได้โดยไม่เกิด import cycle
package ratrace

import "github.com/finance-boardgame/engine/domain"

// BoardSize — จำนวนช่องบนกระดาน Rat Race
const BoardSize = 24

// DefaultBoard คืน layout เริ่มต้นของกระดาน Rat Race (24 ช่อง, Payday ที่ index 0)
//
// layout นี้เป็นแบบ "Cashflow-inspired" ที่เสนอ (Payday 1 ช่อง + Opportunity/Deal
// เยอะที่สุด + Doodad/Market/Charity/Baby/Downsizing กระจาย) — ปรับเปลี่ยนได้
func DefaultBoard() []domain.Tile {
	return []domain.Tile{
		{Type: domain.TilePayday, Name: "Payday"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileDoodad, Name: "Doodad"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileMarket, Name: "Market"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileDoodad, Name: "Doodad"},
		{Type: domain.TileCharity, Name: "Charity"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileBaby, Name: "Baby"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileDoodad, Name: "Doodad"},
		{Type: domain.TileMarket, Name: "Market"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileDownsizing, Name: "Downsizing"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileDoodad, Name: "Doodad"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileMarket, Name: "Market"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileDoodad, Name: "Doodad"},
		{Type: domain.TileCharity, Name: "Charity"},
		{Type: domain.TileOpportunity, Name: "Opportunity"},
		{Type: domain.TileBlank, Name: "Blank"},
	}
}

// MovePlayer เดินผู้เล่นตามจำนวนลูกเต๋า แล้วคืนตำแหน่งปลายทาง (wraps รอบกระดาน)
//
// ไม่ resolve เหตุการณ์ของช่องปลายทาง — การ resolve (Payday/Opportunity/Doodad ฯลฯ)
// เป็นหน้าที่ของ engine.Apply เพื่อรวม logic และ emit events ที่เดียว
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
