package bio

import (
	"context"
	"testing"
	"time"
)

func TestGetBiomarkerStatus(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		optMin   float64
		optMax   float64
		expected BiomarkerStatus
	}{
		{"optimal", 400, 138, 635, 250, 500, StatusOptimal},
		{"low", 100, 138, 635, 250, 500, StatusLow},
		{"high", 700, 138, 635, 250, 500, StatusHigh},
		{"critical_low", 50, 138, 635, 250, 500, StatusCriticalLow},
		{"critical_high", 1300, 138, 635, 250, 500, StatusCriticalHigh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := GetBiomarkerStatus(tt.value, tt.min, tt.max, tt.optMin, tt.optMax)
			if status != tt.expected {
				t.Errorf("Ожидалось %s, получено %s", tt.expected, status)
			}
		})
	}
}

func TestCalculatePercentile(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected float64
	}{
		{"middle", 386.5, 138, 635, 50},
		{"low", 138, 138, 635, 0},
		{"high", 635, 138, 635, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentile := CalculatePercentile(tt.value, tt.min, tt.max)
			if percentile < tt.expected-1 || percentile > tt.expected+1 {
				t.Errorf("Ожидалось ~%.0f, получено %.2f", tt.expected, percentile)
			}
		})
	}
}

func TestCreateUserLabResult(t *testing.T) {
	userID := "test-user-123"

	result := CreateUserLabResult(
		userID,
		"cortisol",
		"Кортизол",
		"Cortisol",
		450,
		"нмоль/л",
		LabInfo{Name: "Инвитро"},
		[]string{},
	)

	if result.UserID != userID {
		t.Error("UserID не совпадает")
	}
	if result.Value != 450 {
		t.Error("Value не совпадает")
	}
	if result.Status != StatusOptimal {
		t.Errorf("Ожидалось StatusOptimal, получено %s", result.Status)
	}
}

func TestMemoryStore(t *testing.T) {
	store, err := newMemoryStore(BioStoreConfig{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	userID := "test-user-456"

	result := CreateUserLabResult(
		userID,
		"cortisol",
		"Кортизол",
		"Cortisol",
		450,
		"нмоль/л",
		LabInfo{Name: "Инвитро"},
		[]string{},
	)

	err = store.SaveResult(ctx, result)
	if err != nil {
		t.Fatal(err)
	}

	results, err := store.GetResultsByUser(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Errorf("Ожидалось 1 результат, получено %d", len(results))
	}

	latest, err := store.GetLatestResult(ctx, userID, "cortisol")
	if err != nil {
		t.Fatal(err)
	}
	if latest == nil {
		t.Error("Ожидался результат, получено nil")
	}

	err = store.DeleteResult(ctx, result.ID)
	if err != nil {
		t.Fatal(err)
	}

	results, err = store.GetResultsByUser(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("Ожидалось 0 результатов, получено %d", len(results))
	}
}

func TestInterpreter(t *testing.T) {
	interpreter := NewInterpreter()

	result := CreateUserLabResult(
		"test-user-789",
		"cortisol",
		"Кортизол",
		"Cortisol",
		100,
		"нмоль/л",
		LabInfo{Name: "Инвитро"},
		[]string{},
	)

	interpretation := interpreter.InterpretResult(result)

	if interpretation.Summary == "" {
		t.Error("Summary пустой")
	}
	if interpretation.DetailedExplanation == "" {
		t.Error("DetailedExplanation пустой")
	}
	if len(interpretation.RelatedChakras) == 0 {
		t.Error("RelatedChakras пустой")
	}
}

// ✅ ИСПРАВЛЕННЫЙ ТЕСТ: добавлены подтесты для всех сценариев тренда
func TestCalculateTrend(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected string
	}{
		{
			name:     "stable (изменение < 10%)",
			values:   []float64{400, 420}, // 5% рост
			expected: "stable",
		},
		{
			name:     "increasing (изменение > 10%)",
			values:   []float64{400, 450}, // 12.5% рост
			expected: "increasing",
		},
		{
			name:     "decreasing (изменение < -10%)",
			values:   []float64{400, 350}, // -12.5% падение
			expected: "decreasing",
		},
		{
			name:     "insufficient_data (менее 2 точек)",
			values:   []float64{400},
			expected: "insufficient_data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []UserLabResult
			baseTime := time.Now().AddDate(0, -1, 0)

			for i, val := range tt.values {
				results = append(results, UserLabResult{
					Value:     val,
					CreatedAt: baseTime.AddDate(0, 0, i*7), // неделя между точками
				})
			}

			trend := CalculateTrend(results)
			if trend != tt.expected {
				t.Errorf("Ожидалось %s, получено %s", tt.expected, trend)
			}
		})
	}
}

// ✅ Новый тест: проверка фильтрации по чакре
func TestGetResultsByChakra(t *testing.T) {
	results := []UserLabResult{
		{
			TestID: "cortisol",
			ChakraCorrelations: []ChakraCorrelation{
				{ChakraIndex: 0, ChakraName: "Муладхара"},
			},
		},
		{
			TestID: "tsh",
			ChakraCorrelations: []ChakraCorrelation{
				{ChakraIndex: 4, ChakraName: "Вишудха"},
			},
		},
		{
			TestID: "vitamin_d",
			ChakraCorrelations: []ChakraCorrelation{
				{ChakraIndex: 3, ChakraName: "Анахата"},
			},
		},
	}

	// Фильтруем по чакре 0 (Муладхара)
	filtered := GetResultsByChakra(results, 0)
	if len(filtered) != 1 {
		t.Errorf("Ожидалось 1 результат для чакры 0, получено %d", len(filtered))
	}
	if filtered[0].TestID != "cortisol" {
		t.Errorf("Ожидался cortisol, получено %s", filtered[0].TestID)
	}
}

// ✅ Новый тест: проверка фильтрации по статусу
func TestGetResultsByStatus(t *testing.T) {
	results := []UserLabResult{
		{TestID: "cortisol", Status: StatusOptimal},
		{TestID: "tsh", Status: StatusLow},
		{TestID: "vitamin_d", Status: StatusOptimal},
	}

	filtered := GetResultsByStatus(results, StatusOptimal)
	if len(filtered) != 2 {
		t.Errorf("Ожидалось 2 результата со статусом optimal, получено %d", len(filtered))
	}
}

func BenchmarkCreateUserLabResult(b *testing.B) {
	userID := "bench-user"

	for i := 0; i < b.N; i++ {
		CreateUserLabResult(
			userID,
			"cortisol",
			"Кортизол",
			"Cortisol",
			450,
			"нмоль/л",
			LabInfo{Name: "Инвитро"},
			[]string{},
		)
	}
}

func BenchmarkInterpretResult(b *testing.B) {
	interpreter := NewInterpreter()
	result := CreateUserLabResult(
		"bench-user",
		"cortisol",
		"Кортизол",
		"Cortisol",
		450,
		"нмоль/л",
		LabInfo{Name: "Инвитро"},
		[]string{},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		interpreter.InterpretResult(result)
	}
}
