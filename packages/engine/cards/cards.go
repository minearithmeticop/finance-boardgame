// Package cards เก็บข้อมูลการ์ดของเกม (DealCard/DoodadCard/CrisisCard) + ตัวจั่วการ์ด
//
// ตัวเลขเป็น "ค่าประมาณแบบหยาบ (game-simplified)" ปรับเปลี่ยนได้
// ดีลถูกออกแบบให้มีทั้งดีลดี (net cashflow บวก) และดีลแยก (net ลบ/ศูนย์)
// เพื่อสอนผู้เล่น "ประเมินดีล" ไม่ใช่ซื้อทุกอย่าง
package cards

import (
	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/rng"
)

// DealCards คืนชุดการ์ดดีล (จากช่อง Opportunity)
func DealCards() []domain.DealCard {
	return []domain.DealCard{
		{Title: "หุ้นกองทุนรวม", AssetType: domain.AssetStock, DownPayment: 5_000, Cost: 5_000, CashFlow: 50, LoanPayment: 0},
		{Title: "หุ้นปันผล", AssetType: domain.AssetStock, DownPayment: 10_000, Cost: 10_000, CashFlow: 150, LoanPayment: 0},
		{Title: "คอนโด 1 ห้องให้เช่า", AssetType: domain.AssetRealEstate, DownPayment: 30_000, Cost: 300_000, CashFlow: 2_500, LoanPayment: 1_800},
		{Title: "ทาวน์เฮ้าส์ให้เช่า", AssetType: domain.AssetRealEstate, DownPayment: 50_000, Cost: 600_000, CashFlow: 4_000, LoanPayment: 3_500},
		{Title: "ร้านกาแฟ", AssetType: domain.AssetBusiness, DownPayment: 40_000, Cost: 200_000, CashFlow: 3_500, LoanPayment: 1_500},
		{Title: "ซักผ้าหยอดเหรียญ", AssetType: domain.AssetBusiness, DownPayment: 20_000, Cost: 80_000, CashFlow: 1_500, LoanPayment: 800},
		{Title: "ที่ดินนอกเมือง", AssetType: domain.AssetRealEstate, DownPayment: 15_000, Cost: 15_000, CashFlow: 0, LoanPayment: 0},
		{Title: "บ้านเก่าเช่าถูก", AssetType: domain.AssetRealEstate, DownPayment: 10_000, Cost: 150_000, CashFlow: 500, LoanPayment: 1_200},
	}
}

// DoodadCards คืนชุดการ์ดช้อป (จากช่อง Shopping) — รายจ่ายฟุ่มเฟือย
func DoodadCards() []domain.DoodadCard {
	return []domain.DoodadCard{
		{Title: "อาหารเย็นนอกบ้าน", Cost: 800},
		{Title: "ดูหนัง/สังสรรค์", Cost: 1_500},
		{Title: "ช้อปเสื้อผ้า", Cost: 3_000},
		{Title: "ของขวัญ", Cost: 4_000},
		{Title: "รองเท้าแบรนด์", Cost: 5_000},
		{Title: "เครื่องใช้ไฟฟ้า", Cost: 12_000},
		{Title: "แกดเจ็ตอัปเกรด", Cost: 15_000},
		{Title: "มือถือรุ่นใหม่", Cost: 25_000},
	}
}

// CrisisCards คืนชุดการ์ดวิกฤต (จากช่อง Crisis) — เหตุการณ์ร้าย
func CrisisCards() []domain.CrisisCard {
	return []domain.CrisisCard{
		{Title: "ค่าปรับ", Amount: 5_000},
		{Title: "ทำฟัน", Amount: 8_000},
		{Title: "เครื่องใช้ในบ้านพัง", Amount: 10_000},
		{Title: "รถเสีย", Amount: 15_000},
		{Title: "ตกงาน 3 เดือน", Amount: 20_000},
		{Title: "อุบัติเหตุ", Amount: 30_000},
		{Title: "ป่วยหนักเข้า รพ.", Amount: 50_000},
	}
}

// DrawDealCard จั่วการ์ดดีลแบบสุ่ม (deterministic ผ่าน rng)
func DrawDealCard(r *rng.RNG) domain.DealCard {
	deck := DealCards()
	if r == nil {
		return deck[0]
	}
	return deck[r.Intn(len(deck))]
}

// DrawDoodadCard จั่วการ์ดช้อปแบบสุ่ม
func DrawDoodadCard(r *rng.RNG) domain.DoodadCard {
	deck := DoodadCards()
	if r == nil {
		return deck[0]
	}
	return deck[r.Intn(len(deck))]
}

// DrawCrisisCard จั่วการ์ดวิกฤตแบบสุ่ม
func DrawCrisisCard(r *rng.RNG) domain.CrisisCard {
	deck := CrisisCards()
	if r == nil {
		return deck[0]
	}
	return deck[r.Intn(len(deck))]
}
