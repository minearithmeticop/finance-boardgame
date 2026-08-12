package loan

import (
	"strings"
	"testing"

	"github.com/finance-boardgame/engine/domain"
)

func mkPlayer(salary, homeBal, homePay, carBal, carPay domain.Money) domain.Player {
	return domain.Player{
		ID: "p1", Cash: 50_000,
		Profession: domain.Profession{
			Name:         "Tester",
			Salary:       salary,
			HomeMortgage: domain.Liability{Payment: homePay, Balance: homeBal},
			CarLoan:      domain.Liability{Payment: carPay, Balance: carBal},
		},
	}
}

// ─── ComputeLoan ───────────────────────────────────────────────────────────

func TestComputeLoan_PersonalMath(t *testing.T) {
	// 50,000 @ 24%/ปี × 24 เดือน → interest 24,000 → balance 74,000 → 3,083/เดือน
	ln := ComputeLoan(LenderPersonal, 50_000)
	if ln.Balance != 74_000 {
		t.Errorf("Balance = %d, want 74000", ln.Balance)
	}
	if ln.MonthlyPay != 3_083 {
		t.Errorf("MonthlyPay = %d, want 3083", ln.MonthlyPay)
	}
}

func TestComputeLoan_InformalMostExpensive(t *testing.T) {
	// 50,000 informal (120%/ปี × 6 เดือน) → interest 30,000 → balance 80,000 → 13,333/เดือน
	inf := ComputeLoan(LenderInformal, 50_000)
	per := ComputeLoan(LenderPersonal, 50_000)
	if inf.MonthlyPay != 13_333 {
		t.Errorf("informal MonthlyPay = %d, want 13333", inf.MonthlyPay)
	}
	if inf.MonthlyPay <= per.MonthlyPay {
		t.Error("informal should be far more expensive than personal")
	}
}

// ─── Request: Personal ─────────────────────────────────────────────────────

func TestRequest_PersonalApproved(t *testing.T) {
	p := mkPlayer(20_000, 0, 0, 0, 0)
	ln, err := Request(p, LenderPersonal, 50_000, "", "")
	if err != nil {
		t.Fatalf("expected approval, got %v", err)
	}
	if ln.Principal != 50_000 {
		t.Errorf("Principal = %d, want 50000", ln.Principal)
	}
}

func TestRequest_PersonalOverLimit(t *testing.T) {
	p := mkPlayer(20_000, 0, 0, 0, 0) // 5× salary = 100,000
	if _, err := Request(p, LenderPersonal, 150_000, "", ""); err == nil ||
		!strings.Contains(err.Error(), "5 เท่า") {
		t.Errorf("expected over-limit rejection, got %v", err)
	}
}

func TestRequest_PersonalHighDSCRRejected(t *testing.T) {
	// หนี้เดิม 12,000/เดือน (60% ของเงินเดือน 20,000) + ค่างวดใหม่ → DSCR > 50%
	p := mkPlayer(20_000, 0, 12_000, 0, 0)
	if _, err := Request(p, LenderPersonal, 50_000, "", ""); err == nil ||
		!strings.Contains(err.Error(), "หนี้เกินเกณฑ์") {
		t.Errorf("expected DSCR rejection, got %v", err)
	}
}

// ─── Request: Secured ──────────────────────────────────────────────────────

func TestRequest_SecuredNeedsCollateral(t *testing.T) {
	p := mkPlayer(20_000, 0, 0, 0, 0) // ไม่มีบ้าน/รถ
	if _, err := Request(p, LenderSecured, 100_000, "home", ""); err == nil ||
		!strings.Contains(err.Error(), "หลักค้ำ") {
		t.Errorf("expected no-collateral rejection, got %v", err)
	}
}

func TestRequest_SecuredHomeApproved(t *testing.T) {
	p := mkPlayer(20_000, 400_000, 3_000, 0, 0) // บ้าน 400,000 → max 280,000
	ln, err := Request(p, LenderSecured, 100_000, "home", "")
	if err != nil {
		t.Fatalf("expected approval, got %v", err)
	}
	if ln.Collateral != "ค้ำบ้าน" || ln.CollatKind != "home" {
		t.Errorf("collateral = %q/%q, want ค้ำบ้าน/home", ln.Collateral, ln.CollatKind)
	}
}

func TestRequest_SecuredLTVExceeded(t *testing.T) {
	p := mkPlayer(20_000, 400_000, 3_000, 0, 0) // max 280,000
	if _, err := Request(p, LenderSecured, 300_000, "home", ""); err == nil ||
		!strings.Contains(err.Error(), "70%") {
		t.Errorf("expected LTV rejection, got %v", err)
	}
}

func TestRequest_SecuredCollateralInUse(t *testing.T) {
	p := mkPlayer(20_000, 400_000, 3_000, 0, 0)
	// กู้ครั้งแรก (ใช้บ้านค้ำ)
	first, err := Request(p, LenderSecured, 100_000, "home", "")
	if err != nil {
		t.Fatal(err)
	}
	p.Loans = append(p.Loans, first)
	// กู้ครั้งที่สองค้ำบ้านซ้ำ → ต้องปฏิเสธ
	if _, err := Request(p, LenderSecured, 50_000, "home", ""); err == nil ||
		!strings.Contains(err.Error(), "ใช้ค้ำสินเชื่ออื่น") {
		t.Errorf("expected in-use rejection, got %v", err)
	}
}

// ─── Request: Informal ─────────────────────────────────────────────────────

func TestRequest_InformalAlwaysApproved(t *testing.T) {
	// ผู้เล่นเงินเดือนต่ำ + ไม่มีหลักค้ำ + วงเงินใหญ่ → นอกระบบก็อนุมัติ!
	p := mkPlayer(9_000, 0, 0, 0, 0)
	ln, err := Request(p, LenderInformal, 200_000, "", "")
	if err != nil {
		t.Fatalf("informal should always approve, got %v", err)
	}
	if ln.RateYear != 120 {
		t.Errorf("informal RateYear = %d, want 120", ln.RateYear)
	}
}
