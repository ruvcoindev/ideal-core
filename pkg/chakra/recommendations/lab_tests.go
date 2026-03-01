

// Package recommendations — рекомендательная система
// 
// Философия:
// Обычный человек не знает:
// 1. Какие анализы сдавать
// 2. Как интерпретировать результаты
// 3. Что делать с полученными данными
// 
// Эта система:
// - Сопоставляет симптомы с возможными дисбалансами
// - Рекомендует анализы по приоритету
// - Интерпретирует результаты в контексте чакр
// - Даёт персонализированные рекомендации
// 
// Важно:
// Это НЕ замена врачу! Система помогает:
// - Подготовиться к визиту к врачу
// - Понять, какие вопросы задать
// - Отслеживать динамику
package recommendations

import (
	"ideal-core/pkg/chakra/layers"
	"ideal-core/pkg/bio"
)

// LabTestEngine — движок рекомендаций по анализам
// 
// Как работает:
// 1. Пользователь вводит симптомы
// 2. Система определяет возможные дисбалансы по чакрам
// 3. Сопоставляет с эндокринными корреляциями
// 4. Формирует персонализированный список анализов
// 5. Объясняет, зачем каждый анализ
// 6. После получения результатов — интерпретирует в контексте
type LabTestEngine struct {
	// EndocrineLayers — эндокринные слои чакр
	EndocrineLayers []layers.EndocrineLayer
	
	// SymptomDatabase — база симптомов и их связей с системами
	SymptomDatabase SymptomDB
	
	// ReferenceTables — таблицы норм (возраст, пол, циклы)
	ReferenceTables ReferenceTables
}

// TestRecommendation — персонализированная рекомендация по анализу
type TestRecommendation struct {
	// Priority — приоритет (1-критичный, 2-важный, 3-рекомендуемый)
	Priority int `json:"priority"`
	
	// TestName — название анализа
	TestName string `json:"test_name"`
	
	// TestNameEn — название для лаборатории
	TestNameEn string `json:"test_name_en"`
	
	// EstimatedCost — ориентировочная стоимость (руб)
	EstimatedCost int `json:"estimated_cost"`
	
	// WhyThisTest — объяснение, зачем этот анализ
	WhyThisTest string `json:"why_this_test"`
	
	// RelatedChakras — связанные чакры
	RelatedChakras []string `json:"related_chakras"`
	
	// RelatedSystems — связанные системы (эндокринная, нервная и т.д.)
	RelatedSystems []string `json:"related_systems"`
	
	// Preparation — как подготовиться
	Preparation []string `json:"preparation"`
	
	// NextSteps — что делать после получения результата
	NextSteps []string `json:"next_steps"`
}

// InterpretationResult — интерпретация результатов анализов
// 
// Концепция:
// Лаборатория выдаёт цифры, но человек не понимает:
// - Это норма или нет?
// - Если не норма, то насколько критично?
// - С чем это связано?
// - Что делать?
// 
// Эта структура даёт контекстную интерпретацию.
type InterpretationResult struct {
	// TestName — название анализа
	TestName string `json:"test_name"`
	
	// UserValue — значение пользователя
	UserValue float64 `json:"user_value"`
	
	// Unit — единица измерения
	Unit string `json:"unit"`
	
	// Status — статус (low, normal, high, critical)
	Status string `json:"status"`
	
	// Percentile — процентиль (насколько далеко от нормы)
	Percentile float64 `json:"percentile"`
	
	// ChakraCorrelation — какая чакра связана
	ChakraCorrelation string `json:"chakra_correlation"`
	
	// PossibleCauses — возможные причины отклонения
	PossibleCauses []string `json:"possible_causes"`
	
	// Recommendations — рекомендации
	Recommendations []string `json:"recommendations"`
	
	// WhenToSeeDoctor — когда нужно к врачу
	WhenToSeeDoctor string `json:"when_to_see_doctor"`
	
	// FollowUpTests — какие анализы сдать дополнительно
	FollowUpTests []string `json:"follow_up_tests"`
}

// GenerateTestRecommendations генерирует рекомендации по анализам
// 
// Входные данные:
// - symptoms — симптомы пользователя
// - age — возраст (влияет на нормы)
// - sex — пол (влияет на нормы)
// - cyclePhase — фаза цикла для женщин (влияет на гормональные нормы)
// 
// Выходные данные:
// - Отсортированный список анализов по приоритету
// - С объяснением, зачем каждый анализ
// - С ориентировочной стоимостью
// 
// Пример:
// ```go
// engine := NewLabTestEngine()
// recommendations := engine.GenerateTestRecommendations(
//     []string{"усталость", "тревожность", "проблемы со сном"},
//     35,
//     "female",
//     "follicular",
// )
// ```
func (e *LabTestEngine) GenerateTestRecommendations(
	symptoms []string,
	age int,
	sex string,
	cyclePhase string,
) []TestRecommendation {
	// 1. Определяем возможные дисбалансы по симптомам
	possibleImbalances := e.SymptomDatabase.MatchSymptoms(symptoms)
	
	// 2. Сопоставляем с чакрами
	relatedChakras := e.mapImbalancesToChakras(possibleImbalances)
	
	// 3. Получаем эндокринные корреляции
	endocrineLayers := GetEndocrineLayers()
	
	// 4. Формируем список анализов
	var recommendations []TestRecommendation
	
	for _, chakraIndex := range relatedChakras {
		layer := endocrineLayers[chakraIndex]
		
		for _, test := range layer.LabTests {
			rec := TestRecommendation{
				Priority:       test.Priority,
				TestName:       test.TestName,
				TestNameEn:     test.TestNameEn,
				EstimatedCost:  test.Cost,
				WhyThisTest:    test.Why,
				RelatedChakras: []string{getChakraName(chakraIndex)},
				RelatedSystems: []string{"Эндокринная", layer.Gland},
				Preparation:    test.Preparation,
				NextSteps: []string{
					"Сдать анализ в лаборатории",
					"Загрузить результаты в систему",
					"Получить интерпретацию",
				},
			}
			recommendations = append(recommendations, rec)
		}
	}
	
	// 5. Сортируем по приоритету
	sortByPriority(recommendations)
	
	// 6. Удаляем дубликаты
	recommendations = deduplicateTests(recommendations)
	
	return recommendations
}

// InterpretResults интерпретирует результаты анализов
// 
// Концепция:
// Лаборатория: "Кортизол = 800 нмоль/л"
// Система: "Кортизол повышен на 26% от верхней границы нормы.
//           Это указывает на хронический стресс и возможное выгорание надпочечников.
//           Связано с дисбалансом Муладхары и Манипуры.
//           Рекомендации: ..."
func (e *LabTestEngine) InterpretResults(
	testName string,
	value float64,
	age int,
	sex string,
	cyclePhase string,
) InterpretationResult {
	// 1. Находим референсные значения
	reference := e.ReferenceTables.GetReference(testName, age, sex, cyclePhase)
	
	// 2. Определяем статус
	status := determineStatus(value, reference.Min, reference.Max)
	
	// 3. Вычисляем процентиль
	percentile := calculatePercentile(value, reference.Min, reference.Max)
	
	// 4. Находим чакральную корреляцию
	chakraCorrelation := e.mapTestToChakra(testName)
	
	// 5. Формируем возможные причины
	possibleCauses := e.getPossibleCauses(testName, status)
	
	// 6. Генерируем рекомендации
	recommendations := e.generateRecommendations(testName, status, chakraCorrelation)
	
	// 7. Определяем, когда к врачу
	whenToSeeDoctor := determineWhenToSeeDoctor(testName, status, value)
	
	// 8. Рекомендуемые дополнительные анализы
	followUpTests := e.getFollowUpTests(testName, status)
	
	return InterpretationResult{
		TestName:          testName,
		UserValue:         value,
		Unit:              reference.Unit,
		Status:            status,
		Percentile:        percentile,
		ChakraCorrelation: chakraCorrelation,
		PossibleCauses:    possibleCauses,
		Recommendations:   recommendations,
		WhenToSeeDoctor:   whenToSeeDoctor,
		FollowUpTests:     followUpTests,
	}
}

// NewLabTestEngine создаёт новый движок рекомендаций
func NewLabTestEngine() *LabTestEngine {
	return &LabTestEngine{
		EndocrineLayers: layers.GetEndocrineLayers(),
		SymptomDatabase: LoadSymptomDatabase(),
		ReferenceTables: LoadReferenceTables(),
	}
}

