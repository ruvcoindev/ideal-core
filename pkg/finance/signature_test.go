package finance

import (
	"testing"
	"time"
)

func TestAnalyzeSignature_Vitaly(t *testing.T) {
	txns := []Transaction{
		{Date: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), Amount: 50000, Type: "IN"},
		{Date: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), Amount: 30000, Type: "IN"},
	}

	sig := AnalyzeSignature("vitaly", txns, -5)

	if sig.TotalAmount != 80000 {
		t.Errorf("TotalAmount: got %.0f, want 80000", sig.TotalAmount)
	}
	if sig.FlowDirection != 1 {
		t.Error("FlowDirection: got 0, want 1")
	}
	if sig.VectorYearEnd != -5 {
		t.Errorf("VectorYearEnd: got %d, want -5", sig.VectorYearEnd)
	}
}

func TestIsToxicResource_Friend(t *testing.T) {
	// Бывший Друг: вектор года -8 → токсичный
	sig := FinancialSignature{
		VectorYearEnd: -8,
		TotalAmount:   100000,
		FlowDirection: 1,
	}

	if !IsToxicResource(sig) {
		t.Error("Expected toxic resource for vector -8")
	}
}

func TestIsToxicResource_Safe(t *testing.T) {
	// Безопасный ресурс: вектор года -5
	sig := FinancialSignature{
		VectorYearEnd: -5,
		TotalAmount:   50000,
		FlowDirection: 1,
	}

	if IsToxicResource(sig) {
		t.Error("Expected safe resource for vector -5")
	}
}
