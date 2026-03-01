// Package bio — результаты пользователей
package bio

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// UserLabResult — результат лабораторного анализа пользователя
type UserLabResult struct {
	ID                 string                `json:"id" db:"id"`
	UserID             string                `json:"user_id" db:"user_id"`
	TestID             string                `json:"test_id" db:"test_id"`
	TestName           string                `json:"test_name" db:"test_name"`
	TestNameEn         string                `json:"test_name_en" db:"test_name_en"`
	Value              float64               `json:"value" db:"value"`
	Unit               string                `json:"unit" db:"unit"`
	ReferenceMin       float64               `json:"reference_min" db:"reference_min"`
	ReferenceMax       float64               `json:"reference_max" db:"reference_max"`
	OptimalMin         float64               `json:"optimal_min" db:"optimal_min"`
	OptimalMax         float64               `json:"optimal_max" db:"optimal_max"`
	Status             BiomarkerStatus       `json:"status" db:"status"`
	Percentile         float64               `json:"percentile" db:"percentile"`
	LabInfo            LabInfo               `json:"lab_info" db:"lab_info"`
	SampleType         string                `json:"sample_type" db:"sample_type"`
	Fasting            bool                  `json:"fasting" db:"fasting"`
	TimeOfDay          string                `json:"time_of_day" db:"time_of_day"`
	CycleDay           int                   `json:"cycle_day" db:"cycle_day"`
	CyclePhase         string                `json:"cycle_phase" db:"cycle_phase"`
	Context            []string              `json:"context" db:"context"`
	Notes              string                `json:"notes" db:"notes"`
	CreatedAt          time.Time             `json:"created_at" db:"created_at"`
	UploadedAt         time.Time             `json:"uploaded_at" db:"uploaded_at"`
	Interpreted        bool                  `json:"interpreted" db:"interpreted"`
	Interpretation     *ResultInterpretation `json:"interpretation" db:"interpretation"`
	ChakraCorrelations []ChakraCorrelation   `json:"chakra_correlations" db:"chakra_correlations"`
	Recommendations    []Recommendation      `json:"recommendations" db:"recommendations"`
}

// LabInfo — информация о лаборатории
type LabInfo struct {
	Name            string `json:"name"`
	Code            string `json:"code"`
	Location        string `json:"location"`
	ReferenceSource string `json:"reference_source"`
}

// ResultInterpretation — автоматическая интерпретация результата
type ResultInterpretation struct {
	Summary                   string     `json:"summary"`
	DetailedExplanation       string     `json:"detailed_explanation"`
	PossibleCauses            []string   `json:"possible_causes"`
	RelatedSystems            []string   `json:"related_systems"`
	RelatedChakras            []string   `json:"related_chakras"`
	RelatedPsychological      []string   `json:"related_psychological"`
	Severity                  string     `json:"severity"`
	WhenToSeeDoctor           string     `json:"when_to_see_doctor"`
	FollowUpTests             []string   `json:"follow_up_tests"`
	LifestyleRecommendations  []string   `json:"lifestyle_recommendations"`
	SupplementRecommendations []string   `json:"supplement_recommendations"`
	GeneratedAt               time.Time  `json:"generated_at"`
	ReviewedBy                string     `json:"reviewed_by"`
	ReviewedAt                *time.Time `json:"reviewed_at"`
}

// ChakraCorrelation — связь результата с чакрами
type ChakraCorrelation struct {
	ChakraIndex          int      `json:"chakra_index"`
	ChakraName           string   `json:"chakra_name"`
	CorrelationType      string   `json:"correlation_type"`
	Strength             float64  `json:"strength"`
	Explanation          string   `json:"explanation"`
	RecommendedPractices []string `json:"recommended_practices"`
}

// Recommendation — рекомендация на основе результата
type Recommendation struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"`
	Priority      int        `json:"priority"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Action        string     `json:"action"`
	Timeline      string     `json:"timeline"`
	EstimatedCost int        `json:"estimated_cost"`
	Evidence      string     `json:"evidence"`
	Completed     bool       `json:"completed"`
	CompletedAt   *time.Time `json:"completed_at"`
	Notes         string     `json:"notes"`
}

// TrendPoint — точка для графика динамики
type TrendPoint struct {
	Date   time.Time       `json:"date"`
	Value  float64         `json:"value"`
	Status BiomarkerStatus `json:"status"`
}

// CreateUserLabResult создаёт новый результат анализа
func CreateUserLabResult(
	userID string,
	testID string,
	testName string,
	testNameEn string,
	value float64,
	unit string,
	labInfo LabInfo,
	context []string,
) UserLabResult {
	ref := GetReferenceRange(testID, "laboratory", "any", 30, "any")
	optRef := GetReferenceRange(testID, "optimal", "any", 30, "any")

	var refMin, refMax, optMin, optMax float64
	if ref != nil {
		refMin = ref.Min
		refMax = ref.Max
	}
	if optRef != nil {
		optMin = optRef.Min
		optMax = optRef.Max
	}

	status := GetBiomarkerStatus(value, refMin, refMax, optMin, optMax)
	percentile := CalculatePercentile(value, refMin, refMax)
	now := time.Now()

	return UserLabResult{
		ID:           generateResultID(),
		UserID:       userID,
		TestID:       testID,
		TestName:     testName,
		TestNameEn:   testNameEn,
		Value:        value,
		Unit:         unit,
		ReferenceMin: refMin,
		ReferenceMax: refMax,
		OptimalMin:   optMin,
		OptimalMax:   optMax,
		Status:       status,
		Percentile:   percentile,
		LabInfo:      labInfo,
		Context:      context,
		CreatedAt:    now,
		UploadedAt:   now,
		Interpreted:  false,
	}
}

func generateResultID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(bytes)
}

// CalculateTrend вычисляет тренд по результатам
func CalculateTrend(results []UserLabResult) string {
	if len(results) < 2 {
		return "insufficient_data"
	}
	first := results[0].Value
	last := results[len(results)-1].Value
	change := ((last - first) / first) * 100
	if change > 10 {
		return "increasing"
	}
	if change < -10 {
		return "decreasing"
	}
	return "stable"
}

// GetResultsByChakra фильтрует результаты по связанной чакре
func GetResultsByChakra(results []UserLabResult, chakraIndex int) []UserLabResult {
	var filtered []UserLabResult
	for _, result := range results {
		for _, corr := range result.ChakraCorrelations {
			if corr.ChakraIndex == chakraIndex {
				filtered = append(filtered, result)
				break
			}
		}
	}
	return filtered
}

// GetResultsByStatus фильтрует результаты по статусу
func GetResultsByStatus(results []UserLabResult, status BiomarkerStatus) []UserLabResult {
	var filtered []UserLabResult
	for _, result := range results {
		if result.Status == status {
			filtered = append(filtered, result)
		}
	}
	return filtered
}
