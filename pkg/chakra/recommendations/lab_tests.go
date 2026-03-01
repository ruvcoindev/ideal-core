package recommendations

import (
	"ideal-core/pkg/chakra"
	"ideal-core/pkg/chakra/layers"
	"sort"
	"strings"
)

// SymptomDB — база симптомов для сопоставления с чакрами
type SymptomDB struct {
	entries map[string][]chakra.ChakraIndex
}

// NewSymptomDB создаёт базу симптомов
func NewSymptomDB() *SymptomDB {
	db := &SymptomDB{entries: make(map[string][]chakra.ChakraIndex)}
	db.loadDefaults()
	return db
}

func (db *SymptomDB) loadDefaults() {
	// Муладхара
	db.entries["страх"] = []chakra.ChakraIndex{chakra.Muladhara}
	db.entries["тревога"] = []chakra.ChakraIndex{chakra.Muladhara}
	db.entries["усталость"] = []chakra.ChakraIndex{chakra.Muladhara, chakra.Manipura}
	
	// Свадхистана
	db.entries["вина"] = []chakra.ChakraIndex{chakra.Svadhisthana}
	db.entries["желания"] = []chakra.ChakraIndex{chakra.Svadhisthana}
	db.entries["цикл"] = []chakra.ChakraIndex{chakra.Svadhisthana}
	
	// Манипура
	db.entries["гнев"] = []chakra.ChakraIndex{chakra.Manipura}
	db.entries["самооценка"] = []chakra.ChakraIndex{chakra.Manipura}
	db.entries["сахар"] = []chakra.ChakraIndex{chakra.Manipura}
	
	// Анахата
	db.entries["обида"] = []chakra.ChakraIndex{chakra.Anahata}
	db.entries["сердце"] = []chakra.ChakraIndex{chakra.Anahata}
	db.entries["одиночество"] = []chakra.ChakraIndex{chakra.Anahata}
	
	// Вишудха
	db.entries["горло"] = []chakra.ChakraIndex{chakra.Vishuddha}
	db.entries["голос"] = []chakra.ChakraIndex{chakra.Vishuddha}
	db.entries["щитовидка"] = []chakra.ChakraIndex{chakra.Vishuddha}
	
	// Аджна
	db.entries["интуиция"] = []chakra.ChakraIndex{chakra.Ajna}
	db.entries["головная боль"] = []chakra.ChakraIndex{chakra.Ajna}
	db.entries["сны"] = []chakra.ChakraIndex{chakra.Ajna}
	
	// Сахасрара
	db.entries["смысл"] = []chakra.ChakraIndex{chakra.Sahasrara}
	db.entries["бессонница"] = []chakra.ChakraIndex{chakra.Sahasrara}
	db.entries["депрессия"] = []chakra.ChakraIndex{chakra.Sahasrara}
}

// FindChakrasBySymptom находит чакры по симптому
func (db *SymptomDB) FindChakrasBySymptom(symptom string) []chakra.ChakraIndex {
	symptomLower := strings.ToLower(symptom)
	for key, chakras := range db.entries {
		if strings.Contains(symptomLower, key) {
			return chakras
		}
	}
	return nil
}

// LabTest — лабораторный тест
type LabTest struct {
	Name        string
	Priority    int    // 1-10, чем выше — тем важнее
	Category    string // "Hormone", "Vitamin", "Mineral", "Immune"
	Chakra      chakra.ChakraIndex
	Description string
}

// LabTestEngine — движок рекомендаций анализов
type LabTestEngine struct {
	symptomDB *SymptomDB
}

// NewLabTestEngine создаёт движок
func NewLabTestEngine() *LabTestEngine {
	return &LabTestEngine{
		symptomDB: NewSymptomDB(),
	}
}

// RecommendTests генерирует рекомендации по анализам
func (e *LabTestEngine) RecommendTests(symptoms []string, activeChakras []chakra.ChakraIndex) []LabTest {
	var tests []LabTest
	
	// 1. Добавляем тесты из эндокринных слоёв
	endocrineLayers := layers.GetEndocrineLayers(activeChakras)
	for _, layer := range endocrineLayers {
		for _, testName := range layer.RecommendedTests {
			tests = append(tests, LabTest{
				Name:        testName,
				Priority:    7,
				Category:    "Hormone",
				Chakra:      layer.Chakra,
				Description: "Рекомендовано при дисбалансе " + layer.Gland,
			})
		}
	}
	
	// 2. Добавляем тесты по симптомам
	for _, symptom := range symptoms {
		chakras := e.symptomDB.FindChakrasBySymptom(symptom)
		for _, c := range chakras {
			if layer, ok := layers.EndocrineMap[c]; ok {
				for _, t := range layer.RecommendedTests {
					tests = append(tests, LabTest{
						Name:        t,
						Priority:    8,
						Category:    "SymptomBased",
						Chakra:      c,
						Description: "По симптому: " + symptom,
					})
				}
			}
		}
	}
	
	// 3. Удаляем дубликаты и сортируем
	tests = deduplicateTests(tests)
	sortByPriority(tests)
	
	return tests
}

// mapImbalancesToChakras сопоставляет дисбалансы чакрам (вспомогательная)
func (e *LabTestEngine) mapImbalancesToChakras(imbalances []string) []chakra.ChakraIndex {
	var chakras []chakra.ChakraIndex
	for _, imb := range imbalances {
		for idx, layer := range layers.EndocrineMap {
			for _, d := range layer.Dysregulation {
				if strings.Contains(strings.ToLower(imb), strings.ToLower(d)) {
					chakras = append(chakras, idx)
				}
			}
		}
	}
	return chakras
}

// mapTestToChakra возвращает чакру для теста
func (e *LabTestEngine) mapTestToChakra(testName string) chakra.ChakraIndex {
	for idx, layer := range layers.EndocrineMap {
		for _, t := range layer.RecommendedTests {
			if strings.Contains(testName, t) {
				return idx
			}
		}
	}
	return -1
}

// deduplicateTests удаляет дубликаты по имени теста
func deduplicateTests(tests []LabTest) []LabTest {
	seen := make(map[string]bool)
	var result []LabTest
	for _, t := range tests {
		if !seen[t.Name] {
			seen[t.Name] = true
			result = append(result, t)
		}
	}
	return result
}

// sortByPriority сортирует тесты по приоритету (по убыванию)
func sortByPriority(tests []LabTest) {
	sort.Slice(tests, func(i, j int) bool {
		return tests[i].Priority > tests[j].Priority
	})
}

// getChakraName возвращает название чакры
func getChakraName(idx chakra.ChakraIndex) string {
	info := chakra.GetChakraInfo(idx)
	return info.Name
}

// determineStatus определяет статус теста (упрощённо)
func determineStatus(value, refMin, refMax float64) string {
	if value < refMin {
		return "LOW"
	}
	if value > refMax {
		return "HIGH"
	}
	return "NORMAL"
}

// calculatePercentile рассчитывает процентиль (заглушка)
func calculatePercentile(value, mean, stdDev float64) float64 {
	// В продакшене: использовать статистическую библиотеку
	z := (value - mean) / stdDev
	// Упрощённая аппроксимация нормального распределения
	if z < -3 {
		return 0.1
	}
	if z > 3 {
		return 99.9
	}
	return 50 + z*17 // грубая линейная аппроксимация
}
