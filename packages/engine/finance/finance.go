// Package finance คำนวณงบการเงินประจำตัวผู้เล่น (Income Statement & Balance Sheet)
// เป็น pure functions รับ Player คืน FinancialStatement — ไม่แก้ state
package finance

import "github.com/finance-boardgame/engine/domain"

// Statement คำนวณงบการเงินเต็มรูปแบบของผู้เล่นจาก assets/liabilities/profession ปัจจุบัน
func Statement(p domain.Player) domain.FinancialStatement {
	var fs domain.FinancialStatement

	// --- Income ---
	fs.EarnedIncome = p.Profession.Salary

	for _, a := range p.Assets {
		if a.CashFlow <= 0 {
			continue
		}
		// แยก passive (อสังหา/ธุรกิจ) vs portfolio (หุ้น) ตามธรรมเนียม Cashflow
		if a.Type == domain.AssetStock {
			fs.PortfolioIncome += a.CashFlow
		} else {
			fs.PassiveIncome += a.CashFlow
		}
	}
	fs.TotalIncome = fs.EarnedIncome + fs.PassiveIncome + fs.PortfolioIncome

	// --- Expenses (หัก ณ ที่จ่าย: ภาษี + ประกันสังคม + อื่นๆ + ผ่อนหนี้) ---
	fs.Tax = p.Profession.Taxes
	fs.SocialSecurity = p.Profession.SocialSecurity
	fs.TotalExpenses = fs.Tax + fs.SocialSecurity + p.Profession.OtherExpenses
	fs.TotalExpenses += p.Profession.HomeMortgage.Payment +
		p.Profession.SchoolLoan.Payment +
		p.Profession.CarLoan.Payment +
		p.Profession.CreditCard.Payment
	for _, l := range p.Liabilities {
		fs.TotalExpenses += l.Payment
	}
	for _, ln := range p.Loans {
		fs.TotalExpenses += ln.MonthlyPay
	}

	fs.MonthlyCashFlow = fs.TotalIncome - fs.TotalExpenses

	// --- Balance Sheet ---
	fs.TotalAssets = p.Cash
	for _, a := range p.Assets {
		fs.TotalAssets += a.Cost
	}

	fs.TotalLiabilities = p.Profession.HomeMortgage.Balance +
		p.Profession.SchoolLoan.Balance +
		p.Profession.CarLoan.Balance +
		p.Profession.CreditCard.Balance
	for _, l := range p.Liabilities {
		fs.TotalLiabilities += l.Balance
	}
	for _, ln := range p.Loans {
		fs.TotalLiabilities += ln.Balance
	}

	return fs
}

// MonthlyCashFlow คืนเงินสดสุทธิต่อเดือน (shortcut)
func MonthlyCashFlow(p domain.Player) domain.Money {
	return Statement(p).MonthlyCashFlow
}

// PassiveIncome คืนรายได้มิใช่เงินเดือนรวม (passive + portfolio)
// หมายเหตุ: เกม Cashflow ใช้ "Passive Income" (รวม portfolio) เป็นเกณฑ์ออก Rat Race
func PassiveIncome(p domain.Player) domain.Money {
	fs := Statement(p)
	return fs.PassiveIncome + fs.PortfolioIncome
}

// CanEscapeRatRace เช็คเงื่อนไขออกจากวง Rat Race: PassiveIncome >= TotalExpenses
func CanEscapeRatRace(p domain.Player) bool {
	fs := Statement(p)
	return (fs.PassiveIncome + fs.PortfolioIncome) >= fs.TotalExpenses
}

// NetWorth คืนมูลค่าสุทธิ = Assets - Liabilities
func NetWorth(p domain.Player) domain.Money {
	fs := Statement(p)
	return fs.TotalAssets - fs.TotalLiabilities
}

// DebtService คืนค่าผ่อนหนี้รวมต่อเดือน (ผ่อน profession + Liabilities + ค่างวด Loans)
// ใช้คำนวณ DSCR สำหรับคุณสมบัติสินเชื่อ
func DebtService(p domain.Player) domain.Money {
	ds := p.Profession.HomeMortgage.Payment +
		p.Profession.SchoolLoan.Payment +
		p.Profession.CarLoan.Payment +
		p.Profession.CreditCard.Payment
	for _, l := range p.Liabilities {
		ds += l.Payment
	}
	for _, ln := range p.Loans {
		ds += ln.MonthlyPay
	}
	return ds
}

// TotalIncome คืนรายได้รวมต่อเดือน (เงินเดือน + รายได้จากสินทรัพย์)
func TotalIncome(p domain.Player) domain.Money {
	inc := p.Profession.Salary
	for _, a := range p.Assets {
		if a.CashFlow > 0 {
			inc += a.CashFlow
		}
	}
	return inc
}
