// Package payroll คำนวณการหัก ณ ที่จ่ายจากเงินเดือน: ภาษีเงินได้ (progressive brackets)
// และประกันสังคม — เป็น "Thai-style game-simplified"
//
// หัวใจที่คงไว้เพื่อสมจริง + สอน:
//   - ภาษีแบบขั้นบันได (progressive) — รายได้สูงขึ้น อัตราสูงขึ้น
//   - ค่าลดหย่อนพื้นฐาน (ส่วนบุคคล 60k, ค่าใช้จ่าย 50% สูงสุด 100k, ประกันสังคม)
//
// ⚠️ ลดทอนมาจากระบบภาษีไทยจริง (ซึ่งมีเครื่องลดหย่อนเยอะกว่านี้มาก) — เหมาะเกม
package payroll

import "github.com/finance-boardgame/engine/domain"

// ค่าคงที่การหัก/ลดหย่อน
const (
	personalAllowance    domain.Money = 60_000  // ค่าลดหย่อนส่วนบุคคล (ต่อปี)
	expenseAllowanceCap  domain.Money = 100_000 // ค่าใช้จ่าย 50% สูงสุด (ต่อปี)
	socialSecurityRate   domain.Money = 5       // ประกันสังคม 5% (ส่วนพนักงาน)
	socialSecurityWageCap domain.Money = 15_000 // เพดานฐานเงินเดือนที่คำนวณ SS
)

// taxBracket — หนึ่งขั้นของบันไดภาษี (เป็น "ขีดบนสะสม" ของช่วงนี้)
type taxBracket struct {
	upTo domain.Money // ขีดบนสะสมของรายได้ในช่วงนี้
	rate domain.Money // เปอร์เซ็นต์ของช่วงนี้
}

// brackets — อัตราภาษีเงินได้บุคคลธรรมดาไทย (PIT) แบบ game-simplified ต่อปี
var brackets = []taxBracket{
	{150_000, 0},
	{300_000, 5},
	{500_000, 10},
	{750_000, 15},
	{1_000_000, 20},
	{2_000_000, 25},
	{5_000_000, 30},
	{1 << 62, 35}, // ขั้นสูงสุด (รายได้เกิน 5,000,000)
}

// SocialSecurity คำนวณเงินสมทบประกันสังคมรายเดือนของพนักงาน
// = min(salary × 5%, 750)
func SocialSecurity(salary domain.Money) domain.Money {
	if salary <= 0 {
		return 0
	}
	base := salary
	if base > socialSecurityWageCap {
		base = socialSecurityWageCap
	}
	return base * socialSecurityRate / 100
}

// TaxableIncome คำนวณ "รายได้ต้องเสียภาษี" ต่อปี จากเงินเดือนรายเดือน
//
//	= max(0, salary×12 − ลดหย่อนส่วนบุคคล 60k − ค่าใช้จ่าย 50%(สูงสุด 100k) − ประกันสังคม×12)
func TaxableIncome(salary domain.Money) domain.Money {
	annual := salary * 12
	expense := annual / 2
	if expense > expenseAllowanceCap {
		expense = expenseAllowanceCap
	}
	ssAnnual := SocialSecurity(salary) * 12
	taxable := annual - personalAllowance - expense - ssAnnual
	if taxable < 0 {
		return 0
	}
	return taxable
}

// AnnualTax คำนวณภาษีเงินได้ทั้งปี จาก "รายได้ต้องเสียภาษี" ต่อปี (progressive brackets)
func AnnualTax(taxable domain.Money) domain.Money {
	if taxable <= 0 {
		return 0
	}
	var tax, prev domain.Money
	for _, b := range brackets {
		if taxable <= prev {
			break
		}
		amount := taxable
		if amount > b.upTo {
			amount = b.upTo
		}
		tax += (amount - prev) * b.rate / 100
		prev = b.upTo
	}
	return tax
}

// MonthlyTax คำนวณภาษีเงินได้ต่อเดือน จากเงินเดือนรายเดือน
// = AnnualTax(TaxableIncome(salary)) / 12
func MonthlyTax(salary domain.Money) domain.Money {
	return AnnualTax(TaxableIncome(salary)) / 12
}
