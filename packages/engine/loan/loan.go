// Package loan คำนวณสินเชื่อและประเมินคุณสมบัติผู้กู้แบบสมจริง
//   - สินเชื่อส่วนบุคคล (unsecured) — ดอกเบี้ยสูง วงเงินจำกัดโดยเงินเดือน + DSCR
//   - สินเชื่อค้ำหลักทรัพย์ (secured) — ดอกเบี้ยต่ำ ต้องมีหลักค้ำ + LTV 70%
//   - กู้นอกระบบ (informal) — ไม่มีเงื่อนไข แต่ดอกเบี้ยสูงลิ่ว (10%/เดือน) สอนอันตราย
package loan

import (
	"errors"
	"fmt"

	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/finance"
)

// ประเภทสินเชื่อ
const (
	LenderPersonal = "personal"
	LenderSecured  = "secured"
	LenderInformal = "informal"
)

// เกณฑ์/พารามิเตอร์สินเชื่อ
const (
	MaxDSCRPercent    = 50 // DSCR ≤ 50%
	PersonalLimitMult = 5  // ส่วนบุคคล ≤ 5× เงินเดือน
	SecuredLTVPercent = 70 // ค้ำหลัก ≤ 70% ของมูลค่าหลักค้ำ
)

type terms struct{ rateYear, termMonths int }

var termsByLender = map[string]terms{
	LenderPersonal: {24, 24},  // 24%/ปี, 24 เดือน
	LenderSecured:  {7, 60},   // 7%/ปี, 60 เดือน
	LenderInformal: {120, 6},  // 120%/ปี (10%/เดือน!), 6 เดือน
}

// ComputeLoan คำนวณค่างวด/ยอดคงค้างจาก Principal + lender
// (flat-amortized แบบเกม: totalInterest = Principal × rate × term/12)
func ComputeLoan(lender string, principal domain.Money) domain.Loan {
	t := termsByLender[lender]
	totalInterest := principal * domain.Money(t.rateYear) * domain.Money(t.termMonths) / 100 / 12
	balance := principal + totalInterest
	monthlyPay := balance / domain.Money(t.termMonths)
	return domain.Loan{
		Lender:     lender,
		Principal:  principal,
		RateYear:   t.rateYear,
		TermMonths: t.termMonths,
		MonthlyPay: monthlyPay,
		Balance:    balance,
	}
}

// Collateral — หลักค้ำประกันที่ผู้เล่นเลือก
type Collateral struct {
	Kind  string
	RefID string
	Desc  string
	Value domain.Money
}

// CollateralValue คำนวณมูลค่าหลักค้ำที่ผู้เล่นมีตาม kind (home/car/asset)
func CollateralValue(p domain.Player, kind, refID string) (Collateral, bool) {
	switch kind {
	case "home":
		if v := p.Profession.HomeMortgage.Balance; v > 0 {
			return Collateral{Kind: "home", Desc: "ค้ำบ้าน", Value: v}, true
		}
	case "car":
		if v := p.Profession.CarLoan.Balance; v > 0 {
			return Collateral{Kind: "car", Desc: "ค้ำรถ", Value: v}, true
		}
	case "asset":
		for _, a := range p.Assets {
			if a.ID == refID && a.Cost > 0 {
				return Collateral{Kind: "asset", RefID: refID, Desc: "ค้ำทรัพย์สิน (" + a.Name + ")", Value: a.Cost}, true
			}
		}
	}
	return Collateral{}, false
}

// CollateralInUse เช็คว่าหลักค้ำนี้ถูกใช้ค้ำสินเชื่ออื่นที่ยังไม่ปิดอยู่ไหม
func CollateralInUse(p domain.Player, kind, refID string) bool {
	for _, ln := range p.Loans {
		if ln.CollatKind != kind {
			continue
		}
		if kind == "asset" && ln.CollatRef != refID {
			continue
		}
		return true
	}
	return false
}

// Request ประเมินการขอสินเชื่อ → คืน Loan พร้อมใช้ หรือ error (ปฏิเสธพร้อมเหตุผลสมจริง)
func Request(p domain.Player, lender string, amount domain.Money, collatKind, collatRef string) (domain.Loan, error) {
	if amount <= 0 {
		return domain.Loan{}, errors.New("วงเงินต้องมากกว่า 0")
	}
	if _, ok := termsByLender[lender]; !ok {
		return domain.Loan{}, errors.New("ประเภทสินเชื่อไม่ถูกต้อง")
	}
	ln := ComputeLoan(lender, amount)

	switch lender {
	case LenderPersonal:
		limit := p.Profession.Salary * PersonalLimitMult
		if amount > limit {
			return domain.Loan{}, fmt.Errorf("วงเงินเกิน 5 เท่าของเงินเดือน (สูงสุด %d)", limit)
		}
		if err := checkDSCR(p, ln.MonthlyPay); err != nil {
			return domain.Loan{}, err
		}
	case LenderSecured:
		c, ok := CollateralValue(p, collatKind, collatRef)
		if !ok {
			return domain.Loan{}, errors.New("ไม่มีหลักค้ำประกันที่ใช้ได้ (ต้องมีบ้าน/รถ/สินทรัพย์)")
		}
		if CollateralInUse(p, collatKind, collatRef) {
			return domain.Loan{}, errors.New("หลักค้ำนี้ถูกใช้ค้ำสินเชื่ออื่นอยู่")
		}
		maxAmount := c.Value * SecuredLTVPercent / 100
		if amount > maxAmount {
			return domain.Loan{}, fmt.Errorf("วงเงินเกิน 70%% ของมูลค่าหลักค้ำ (สูงสุด %d)", maxAmount)
		}
		ln.Collateral = c.Desc
		ln.CollatKind = c.Kind
		ln.CollatRef = c.RefID
		if err := checkDSCR(p, ln.MonthlyPay); err != nil {
			return domain.Loan{}, err
		}
	case LenderInformal:
		// ไม่มีเงื่อนไข — อนุมัติเสมอ (ดอกเบี้ย 10%/เดือน เป็นกับดัก!)
	}
	return ln, nil
}

// checkDSCR ตรวจอัตราส่วนหนี้ต่อรายได้ ≤ 50% (รวมค่างวดใหม่) — เงื่อนไขธนาคารจริง
func checkDSCR(p domain.Player, newPay domain.Money) error {
	inc := finance.TotalIncome(p)
	if inc <= 0 {
		return errors.New("ไม่มีรายได้ — ไม่ผ่านเกณฑ์หนี้ต่อรายได้")
	}
	ds := finance.DebtService(p) + newPay
	if ds*100 > inc*MaxDSCRPercent {
		return fmt.Errorf("หนี้เกินเกณฑ์ (DSCR %d%% > 50%%)", ds*100/inc)
	}
	return nil
}
