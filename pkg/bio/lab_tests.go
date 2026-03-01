// Package bio — справочник лабораторных анализов
package bio

import (
	"strings"
	"time"
)

// LabTest — лабораторный анализ
type LabTest struct {
	ID                string    `json:"id" db:"id"`
	Code              string    `json:"code" db:"code"`
	Name              string    `json:"name" db:"name"`
	NameEn            string    `json:"name_en" db:"name_en"`
	BiomarkerCode     string    `json:"biomarker_code" db:"biomarker_code"`
	Method            string    `json:"method" db:"method"`
	SampleType        string    `json:"sample_type" db:"sample_type"`
	Volume            float64   `json:"volume" db:"volume"`
	TurnaroundTime    int       `json:"turnaround_time" db:"turnaround_time"`
	Cost              int       `json:"cost" db:"cost"`
	Available         bool      `json:"available" db:"available"`
	Preparation       []string  `json:"preparation" db:"preparation"`
	Contraindications []string  `json:"contraindications" db:"contraindications"`
	Interferences     []string  `json:"interferences" db:"interferences"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// LabTestDatabase — база данных лабораторных тестов
type LabTestDatabase struct {
	Tests       []LabTest        `json:"tests"`
	References  []ReferenceRange `json:"references"`
	Categories  []TestCategory   `json:"categories"`
	LastUpdated time.Time        `json:"last_updated"`
}

// TestCategory — категория тестов
type TestCategory struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// LoadLabTestDatabase загружает базу данных лабораторных тестов
func LoadLabTestDatabase() LabTestDatabase {
	return LabTestDatabase{
		Tests:       loadAllLabTests(),
		References:  loadAllReferenceRanges(),
		Categories:  loadTestCategories(),
		LastUpdated: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
}

func loadAllLabTests() []LabTest {
	return []LabTest{
		{
			ID:             "cortisol_blood",
			Code:           "cortisol",
			Name:           "Кортизол в крови",
			NameEn:         "Cortisol, Blood",
			BiomarkerCode:  "cortisol",
			Method:         "CLIA",
			SampleType:     "blood",
			Volume:         5,
			TurnaroundTime: 24,
			Cost:           800,
			Available:      true,
			Preparation: []string{
				"Натощак (8-12 часов)",
				"Избегать стресса перед сдачей",
				"Не курить за 3 часа",
				"Сдавать 08:00-10:00",
			},
			Contraindications: []string{
				"Острое заболевание",
				"Приём глюкокортикоидов",
			},
			Interferences: []string{
				"Стресс",
				"Физическая нагрузка",
				"Беременность",
			},
		},
		{
			ID:             "tsh",
			Code:           "tsh",
			Name:           "ТТГ (тиреотропный гормон)",
			NameEn:         "TSH (Thyroid Stimulating Hormone)",
			BiomarkerCode:  "tsh",
			Method:         "CLIA",
			SampleType:     "blood",
			Volume:         5,
			TurnaroundTime: 24,
			Cost:           500,
			Available:      true,
			Preparation: []string{
				"Натощак (8-12 часов)",
				"Избегать стресса",
			},
			Contraindications: []string{"Острое заболевание"},
			Interferences: []string{
				"Биотин (за 48 часов)",
				"Глюкокортикоиды",
			},
		},
		{
			ID:             "vitamin_d",
			Code:           "vitamin_d",
			Name:           "Витамин D (25-OH)",
			NameEn:         "Vitamin D 25-OH",
			BiomarkerCode:  "vitamin_d",
			Method:         "CLIA",
			SampleType:     "blood",
			Volume:         5,
			TurnaroundTime: 24,
			Cost:           2000,
			Available:      true,
			Preparation: []string{
				"Натощак (8-12 часов)",
				"Любое время суток",
			},
			Contraindications: []string{},
			Interferences:     []string{"Приём витамина D (за 24 часа)"},
		},
	}
}

func loadAllReferenceRanges() []ReferenceRange {
	return []ReferenceRange{
		{
			BiomarkerCode: "cortisol",
			RangeType:     "laboratory",
			AgeMin:        18,
			AgeMax:        60,
			Sex:           "any",
			CyclePhase:    "any",
			Min:           138,
			Max:           635,
			Unit:          "нмоль/л",
			Source:        "invitro",
			Notes:         "Сдавать утром 08:00-10:00",
			LastUpdated:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			BiomarkerCode: "cortisol",
			RangeType:     "optimal",
			AgeMin:        18,
			AgeMax:        60,
			Sex:           "any",
			CyclePhase:    "any",
			Min:           250,
			Max:           500,
			Unit:          "нмоль/л",
			Source:        "functional_medicine",
			Notes:         "Оптимальный диапазон для хорошего самочувствия",
			LastUpdated:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			BiomarkerCode: "tsh",
			RangeType:     "laboratory",
			AgeMin:        18,
			AgeMax:        60,
			Sex:           "any",
			CyclePhase:    "any",
			Min:           0.4,
			Max:           4.0,
			Unit:          "мЕд/л",
			Source:        "invitro",
			Notes:         "Лабораторная норма",
			LastUpdated:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			BiomarkerCode: "tsh",
			RangeType:     "optimal",
			AgeMin:        18,
			AgeMax:        60,
			Sex:           "any",
			CyclePhase:    "any",
			Min:           1.0,
			Max:           2.5,
			Unit:          "мЕд/л",
			Source:        "functional_medicine",
			Notes:         "Оптимально для хорошего самочувствия",
			LastUpdated:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}
}

func loadTestCategories() []TestCategory {
	return []TestCategory{
		{Code: "hormones", Name: "Гормоны", Description: "Гормональные исследования", Icon: "🧬"},
		{Code: "vitamins", Name: "Витамины", Description: "Витаминный статус", Icon: "💊"},
		{Code: "minerals", Name: "Минералы", Description: "Микро- и макроэлементы", Icon: "⚗️"},
		{Code: "metabolic", Name: "Метаболизм", Description: "Обмен веществ", Icon: "🔥"},
	}
}

// GetLabTestByCode получает тест по коду
func GetLabTestByCode(code string) (*LabTest, bool) {
	db := LoadLabTestDatabase()
	for _, test := range db.Tests {
		if test.Code == code {
			return &test, true
		}
	}
	return nil, false
}

// GetLabTestsByBiomarker получает все тесты для биомаркера
func GetLabTestsByBiomarker(biomarkerCode string) []LabTest {
	db := LoadLabTestDatabase()
	var result []LabTest
	for _, test := range db.Tests {
		if test.BiomarkerCode == biomarkerCode {
			result = append(result, test)
		}
	}
	return result
}

// GetReferenceRange получает референсный диапазон
func GetReferenceRange(biomarkerCode, rangeType, sex string, age int, cyclePhase string) *ReferenceRange {
	db := LoadLabTestDatabase()
	for _, ref := range db.References {
		if ref.BiomarkerCode == biomarkerCode &&
			ref.RangeType == rangeType &&
			(ref.Sex == sex || ref.Sex == "any") &&
			age >= ref.AgeMin && age <= ref.AgeMax &&
			(ref.CyclePhase == cyclePhase || ref.CyclePhase == "any") {
			return &ref
		}
	}
	for _, ref := range db.References {
		if ref.BiomarkerCode == biomarkerCode &&
			ref.RangeType == rangeType &&
			ref.Sex == "any" &&
			ref.CyclePhase == "any" {
			return &ref
		}
	}
	return nil
}

// BiomarkerChakraCorrelation — корреляция биомаркера с чакрой
type BiomarkerChakraCorrelation struct {
	ChakraIndex          int               `json:"chakra_index"`
	ChakraName           string            `json:"chakra_name"`
	CorrelationStrength  float64           `json:"correlation_strength"`
	Mechanism            string            `json:"mechanism"`
	Direction            map[string]string `json:"direction"`
	Symptoms             []string          `json:"symptoms"`
	RecommendedPractices []string          `json:"recommended_practices"`
}

// BiomarkerChakraMap — карта связей биомаркеров с чакрами
var BiomarkerChakraMap = map[string]BiomarkerChakraCorrelation{
	"cortisol": {
		ChakraIndex:         0,
		ChakraName:          "Муладхара",
		CorrelationStrength: 0.95,
		Mechanism:           "direct",
		Direction: map[string]string{
			"high": "Хронический стресс, тревога, гиперконтроль",
			"low":  "Истощение надпочечников, хроническая усталость",
		},
		Symptoms: []string{
			"Усталость после пробуждения",
			"Тревожность",
			"Тяга к солёному",
		},
		RecommendedPractices: []string{
			"Заземляющие медитации",
			"Ходьба босиком",
			"Дыхательные практики (4-7-8)",
		},
	},
	"tsh": {
		ChakraIndex:         4,
		ChakraName:          "Вишудха",
		CorrelationStrength: 0.95,
		Mechanism:           "direct",
		Direction: map[string]string{
			"high": "Гипотиреоз — усталость, набор веса, депрессия",
			"low":  "Гипертиреоз — тревога, потеря веса, раздражительность",
		},
		Symptoms: []string{
			"«Ком в горле»",
			"Трудности с выражением себя",
			"Проблемы с голосом",
		},
		RecommendedPractices: []string{
			"Пение, мантры",
			"Ведение дневника",
			"Йод, селен, цинк",
		},
	},
	"vitamin_d": {
		ChakraIndex:         3,
		ChakraName:          "Анахата",
		CorrelationStrength: 0.70,
		Mechanism:           "indirect",
		Direction: map[string]string{
			"low": "Депрессия, частые болезни, усталость, боли в костях",
		},
		Symptoms: []string{
			"Частые простуды",
			"Депрессия",
			"Усталость",
		},
		RecommendedPractices: []string{
			"Солнечный свет 15-30 мин/день",
			"Витамин D3 + K2",
			"Практики открытия сердца",
		},
	},
}

// GetChakraCorrelationsForBiomarker получает корреляции биомаркера с чакрами
func GetChakraCorrelationsForBiomarker(biomarkerName string) []BiomarkerChakraCorrelation {
	normalized := normalizeBiomarkerName(biomarkerName)
	if corr, ok := BiomarkerChakraMap[normalized]; ok {
		return []BiomarkerChakraCorrelation{corr}
	}
	var results []BiomarkerChakraCorrelation
	for name, corr := range BiomarkerChakraMap {
		if containsIgnoreCase(name, normalized) {
			results = append(results, corr)
		}
	}
	return results
}

// GetBiomarkersForChakra получает все биомаркеры, связанные с чакрой
func GetBiomarkersForChakra(chakraIndex int) []BiomarkerChakraCorrelation {
	var results []BiomarkerChakraCorrelation
	for _, corr := range BiomarkerChakraMap {
		if corr.ChakraIndex == chakraIndex {
			results = append(results, corr)
		}
	}
	return results
}

func normalizeBiomarkerName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
