// Package profession เก็บข้อมูลอาชีพเริ่มต้น (Thai-flavored) และตัวสุ่มอาชีพ
// แต่ละอาชีพกำหนด เงินเดือน/รายจ่าย/หนี้สินเริ่มต้น และคำนวณ ภาษี + ประกันสังคม
// (หัก ณ ที่จ่าย) ผ่าน package payroll ตอนสร้าง — เพื่อให้ตัวเลขสอดคล้องกับ formula เสมอ
//
// 📊 ที่มาของเงินเดือน: ปรับตามข้อมูลอ้างอิงจริง (ค่าจ้างเฉลี่ยไทย ~15,737/เดือน, Q4/2024)
//    ดูรายละเอียด + แหล่งอ้างอิงทั้งหมดใน INFO.md § "ข้อมูลรายได้อ้างอิง"
//    ตัวเลขเป็น "ค่าประมาณแบบอ้างอิง" (ไม่ใช่ตัวเลขเที่ยงตรงถ้วน) เพราะรายได้จริงแตกต่าง
//    ตาม ภูมิภาค/ประสบการณ์/ขนาดบริษัท
package profession

import (
	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/payroll"
	"github.com/finance-boardgame/engine/rng"
)

// allProfessions — ชุดอาชีพเริ่มต้น ครอบคลุมช่วงรายได้จริงของคนไทย
// (ตั้งแต่ค่าจ้างต่ำ/นอกระบบ → มืออาชีพระดับสูง) เพื่อสอนว่า "รายได้สูง ≠ อิสรภาพ
// ถ้ารายจ่าย/หนี้สูงตาม"
var allProfessions = []domain.Profession{
	// ── ค่าจ้างต่ำ / นอกระบบ (ค่าจ้างขั้นต่ำ ~337–400/วัน ≈ 10,000–12,000/เดือน) ──
	mk("พนักงานทำความสะอาด", 9_000, 3_500, 300,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถจักรยานยนต์", Payment: 500, Balance: 20_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 300, Balance: 4_000},
		domain.Liability{}),
	mk("พนักงานเสิร์ฟ", 9_500, 3_800, 400,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถจักรยานยนต์", Payment: 600, Balance: 10_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 400, Balance: 5_000},
		domain.Liability{}),
	mk("แม่ค้าตลาด", 11_000, 4_500, 1_000,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถบรรทุก", Payment: 800, Balance: 15_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 500, Balance: 6_000},
		domain.Liability{}),
	mk("คนงานโรงงาน", 12_000, 4_000, 800,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถจักรยานยนต์", Payment: 700, Balance: 12_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 500, Balance: 5_000},
		domain.Liability{}),
	mk("พนักงานขายหน้าร้าน", 13_000, 5_000, 700,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถจักรยานยนต์", Payment: 800, Balance: 12_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 800, Balance: 10_000},
		domain.Liability{}),

	// ── พนักงานทั่วไป / จบใหม่ (ค่าจ้างเฉลี่ยรวม ~15,737/เดือน) ──
	mk("ไรเดอร์รับจ้าง", 14_000, 4_500, 1_000,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถยนต์มือสอง", Payment: 3_000, Balance: 120_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 500, Balance: 8_000},
		domain.Liability{}),
	mk("พนักงานออฟฟิศ", 16_000, 6_000, 1_500,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 3_500, Balance: 180_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 1_000, Balance: 15_000},
		domain.Liability{Name: "กู้กยศ.", Payment: 1_000, Balance: 100_000}),
	mk("ช่างซ่อมรถ", 18_000, 6_000, 1_500,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 4_000, Balance: 200_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 1_000, Balance: 15_000},
		domain.Liability{}),

	// ── ข้าราชการ / มืออาชีพต้น (ฐานเงินเดือน ก.พ. เริ่ม ~15,050) ──
	mk("ครู", 22_000, 7_000, 2_500,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 3_000, Balance: 400_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 5_000, Balance: 300_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 1_500, Balance: 30_000},
		domain.Liability{}),
	mk("พยาบาล", 25_000, 8_000, 3_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 4_000, Balance: 500_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 6_000, Balance: 350_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 2_000, Balance: 40_000},
		domain.Liability{Name: "กู้กยศ.", Payment: 1_000, Balance: 150_000}),
	mk("นักบัญชี", 28_000, 9_000, 3_500,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 5_000, Balance: 700_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 7_000, Balance: 400_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 2_000, Balance: 50_000},
		domain.Liability{Name: "กู้กยศ.", Payment: 1_500, Balance: 180_000}),

	// ── มืออาชีพ ──
	mk("วิศวกร", 32_000, 10_000, 4_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 6_000, Balance: 900_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 8_000, Balance: 500_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 2_500, Balance: 60_000},
		domain.Liability{Name: "กู้กยศ.", Payment: 2_000, Balance: 200_000}),
	mk("โปรแกรมเมอร์", 35_000, 11_000, 5_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 7_000, Balance: 1_000_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 9_000, Balance: 550_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 3_000, Balance: 80_000},
		domain.Liability{Name: "กู้กยศ.", Payment: 2_000, Balance: 250_000}),
	mk("แพทย์", 55_000, 15_000, 8_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 12_000, Balance: 2_200_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 10_000, Balance: 800_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 4_000, Balance: 120_000},
		domain.Liability{Name: "กู้กยศ.", Payment: 4_000, Balance: 400_000}),
}

// All คืนรายการอาชีพเริ่มต้นทั้งหมด (สำเนา ป้องกันการแก้ชุดต้นฉบับ)
func All() []domain.Profession {
	out := make([]domain.Profession, len(allProfessions))
	copy(out, allProfessions)
	return out
}

// Random สุ่มอาชีพหนึ่งอาชีพ (สำหรับ "จบใหม่สุ่มอาชีพ")
func Random(r *rng.RNG) domain.Profession {
	if r == nil || len(allProfessions) == 0 {
		return domain.Profession{}
	}
	return allProfessions[r.Intn(len(allProfessions))]
}

// mk สร้างอาชีพจากข้อมูลดิบ พร้อมคำนวณ ภาษี + ประกันสังคม ผ่าน payroll
func mk(name string, salary, other, savings domain.Money, mortgage, car, credit, school domain.Liability) domain.Profession {
	return domain.Profession{
		Name:           name,
		Salary:         salary,
		OtherExpenses:  other,
		Savings:        savings,
		HomeMortgage:   mortgage,
		CarLoan:        car,
		CreditCard:     credit,
		SchoolLoan:     school,
		Taxes:          payroll.MonthlyTax(salary),
		SocialSecurity: payroll.SocialSecurity(salary),
	}
}
