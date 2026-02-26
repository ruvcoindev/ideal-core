package cube

import (
	"testing"
	"time"
)

func TestCalcCoordinates(t *testing.T) {
	// Виталий: 15.10.1974 → 015, 010, 974 → (6, 1, 20)
	date := time.Date(1974, 10, 15, 0, 0, 0, 0, time.UTC)
	coords := CalcCoordinates(date)
	
	if coords[0] != 6 { t.Errorf("Day: got %d, want 6", coords[0]) }
	if coords[1] != 1 { t.Errorf("Month: got %d, want 1", coords[1]) }
	if coords[2] != 20 { t.Errorf("Year: got %d, want 20", coords[2]) }
}

func TestNewPerson_Vitaly(t *testing.T) {
	date := time.Date(1974, 10, 15, 0, 0, 0, 0, time.UTC)
	p := NewPerson("vitaly", "Vitaly", date)
	
	if p.SumFrequency != 27 {
		t.Errorf("SumFrequency: got %d, want 27", p.SumFrequency)
	}
	// Вектор дня для 015: [0-1, 1-5, 5-0] = [-1, -4, 5]
	expectedDayVector := [3]int{-1, -4, 5}
	if p.Vectors[0] != expectedDayVector {
		t.Errorf("DayVector: got %v, want %v", p.Vectors[0], expectedDayVector)
	}
}

func TestDistance_Vitaly_Dina(t *testing.T) {
	vitaly := NewPerson("v", "Vitaly", time.Date(1974, 10, 15, 0, 0, 0, 0, time.UTC))
	dina := NewPerson("d", "Dina", time.Date(1970, 10, 3, 0, 0, 0, 0, time.UTC))
	
	dist := Distance(*vitaly, *dina)
	// Ожидаемое расстояние ~1.00 при сближении
	t.Logf("Distance Vitaly-Dina: %.2f", dist)
}
