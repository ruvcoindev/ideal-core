package hypercube

import (
	"math"
	"testing"
	"time"
)

func TestDistance4D(t *testing.T) {
	c1 := HyperCoord{X: 6, Y: 1, Z: 20, W: 51} // Виталий
	c2 := HyperCoord{X: 3, Y: 1, Z: 16, W: 7}  // Дина (регресс)

	dist := Distance4D(c1, c2)
	expected := math.Sqrt(3*3 + 0*0 + 4*4 + 44*44) // ~44.6

	if math.Abs(dist-expected) > 0.01 {
		t.Errorf("Distance4D: got %.2f, want %.2f", dist, expected)
	}
}

func TestCompatibility4D(t *testing.T) {
	p1 := Person4D{
		ID:        "V",
		BaseCoord: HyperCoord{X: 6, Y: 1, Z: 20},
		CurrentW:  51,
	}
	p2 := Person4D{
		ID:        "D",
		BaseCoord: HyperCoord{X: 3, Y: 1, Z: 16},
		CurrentW:  7, // регресс
	}

	compat := Compatibility4D(p1, p2)
	
	// W-разрыв 44 → низкая совместимость
	if compat > 0.1 {
		t.Errorf("Compatibility4D: expected low score for W-gap=44, got %.3f", compat)
	}
}

func TestPhaseShift(t *testing.T) {
	tests := []struct {
		w        int
		expected string
	}{
		{7, "Child/Trauma"},
		{25, "Adolescent/Searching"},
		{40, "Adult/Responsibility"},
		{60, "Elder/Wisdom"},
	}
	for _, tc := range tests {
		if got := PhaseShift(tc.w); got != tc.expected {
			t.Errorf("PhaseShift(%d): got %q, want %q", tc.w, got, tc.expected)
		}
	}
}

func TestWGap_Interpretation(t *testing.T) {
	if WGap_Interpretation(5) != "Близкие состояния: возможен диалог" {
		t.Error("WGap_Interpretation failed for small gap")
	}
	if WGap_Interpretation(50) != "Разные временные слои: контакт затруднён" {
		t.Error("WGap_Interpretation failed for large gap")
	}
}
