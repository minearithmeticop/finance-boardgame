// Package profession เก็บข้อมูลอาชีพเริ่มต้น (Thai-flavored) และตัวสุ่มอาชีพ
// แต่ละอาชีพกำหนด เงินเดือน/รายจ่าย/หนี้สินเริ่มต้น และคำนวณ ภาษี + ประกันสังคม
// (หัก ณ ที่จ่าย) ผ่าน package payroll ตอนสร้าง — เพื่อให้ตัวเลขสอดคล้องกับ formula เสมอ
package profession

import (
	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/payroll"
	"github.com/finance-boardgame/engine/rng"
)

// allProfessions — ชุดอาชีพเริ่มต้น (ระดับรายได้ต่างกัน เพื่อสอนว่า "รายได้สูง ≠ อิสรภาพ ถ้ารายจ่ายสูงตาม")
var allProfessions = []domain.Profession{
	mk("พนักงานเสิร์ฟ", 15_000, 4_000, 500,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถจักรยานยนต์", Payment: 800, Balance: 20_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 500, Balance: 8_000},
		domain.Liability{}),
	mk("พนักงานออฟฟิศ", 20_000, 6_000, 1_000,
		domain.Liability{},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 4_000, Balance: 250_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 1_000, Balance: 20_000},
		domain.Liability{Name: "กู้กศน.", Payment: 1_000, Balance: 30_000}),
	mk("ครู", 25_000, 7_000, 2_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 3_000, Balance: 400_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 5_000, Balance: 300_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 1_500, Balance: 30_000},
		domain.Liability{}),
	mk("พยาบาล", 35_000, 9_000, 3_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 5_000, Balance: 700_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 7_000, Balance: 450_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 2_000, Balance: 50_000},
		domain.Liability{Name: "กู้กศน.", Payment: 1_000, Balance: 20_000}),
	mk("วิศวกร", 45_000, 12_000, 4_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 8_000, Balance: 1_200_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 10_000, Balance: 700_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 3_000, Balance: 80_000},
		domain.Liability{Name: "กู้กยศ.", Payment: 2_500, Balance: 200_000}),
	mk("ผู้จัดการ", 80_000, 20_000, 10_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 18_000, Balance: 3_000_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 15_000, Balance: 1_200_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 5_000, Balance: 150_000},
		domain.Liability{}),
	mk("แพทย์", 120_000, 30_000, 20_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 30_000, Balance: 5_000_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 20_000, Balance: 1_800_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 8_000, Balance: 250_000},
		domain.Liability{Name: "กู้กยศ.", Payment: 5_000, Balance: 400_000}),
	mk("ทนายความอาวุโส", 200_000, 50_000, 40_000,
		domain.Liability{Name: "ผ่อนบ้าน", Payment: 60_000, Balance: 10_000_000},
		domain.Liability{Name: "ผ่อนรถยนต์", Payment: 30_000, Balance: 3_000_000},
		domain.Liability{Name: "บัตรเครดิต", Payment: 15_000, Balance: 500_000},
		domain.Liability{}),
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
