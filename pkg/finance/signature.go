package finance

import (
	"time"
)

// FinancialSignature хранит паттерн денежного потока
type FinancialSignature struct {
	PersonID      string
	TotalAmount   float64
	FlowDirection int       // +1 рост, 0 стазис, -1 убыль
	VectorYearEnd int       // последняя цифра вектора года
	LastFlowDate  time.Time
	Tags          []string  // "Pyramid", "Business", "Family"
}

// ToxicThreshold порог токсичности вектора года
const ToxicThreshold = -7

// AnalyzeSignature оценивает человека как источник ресурса
func AnalyzeSignature(personID string, transactions []Transaction, vectorYearEnd int) FinancialSignature {
	sig := FinancialSignature{
		PersonID:      personID,
		TotalAmount:   0,
		FlowDirection: 0,
		VectorYearEnd: vectorYearEnd,
		Tags:          []string{},
	}

	// Считаем сумму и направление потока
	for _, t := range transactions {
		sig.TotalAmount += t.Amount
		if t.Date.After(sig.LastFlowDate) {
			sig.LastFlowDate = t.Date
		}
	}

	// Определяем направление
	if sig.TotalAmount > 0 {
		sig.FlowDirection = 1
	}

	return sig
}

// IsToxicResource проверяет, является ли человек "токсичным спонсором"
// (как Друг с вектором -8: деньги есть, но негатив зашкаливает)
func IsToxicResource(sig FinancialSignature) bool {
	return sig.VectorYearEnd <= ToxicThreshold
}

// Transaction модель транзакции
type Transaction struct {
	Date   time.Time
	Amount float64
	Type   string // "IN", "OUT"
}
