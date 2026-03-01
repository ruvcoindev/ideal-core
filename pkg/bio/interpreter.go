// Package bio — интерпретация результатов анализов
package bio

import (
	"fmt"
	"strings"
	"time"
)

// Interpreter — движок интерпретации результатов
type Interpreter struct {
	biomarkerMap map[string]BiomarkerChakraCorrelation
}

// NewInterpreter создаёт новый движок интерпретации
func NewInterpreter() *Interpreter {
	return &Interpreter{
		biomarkerMap: BiomarkerChakraMap,
	}
}

// InterpretResult интерпретирует результат анализа
func (i *Interpreter) InterpretResult(result UserLabResult) *ResultInterpretation {
	interpretation := &ResultInterpretation{
		GeneratedAt: time.Now(),
	}

	// 1. Краткое резюме
	interpretation.Summary = i.generateSummary(result)

	// 2. Подробное объяснение
	interpretation.DetailedExplanation = i.generateDetailedExplanation(result)

	// 3. Возможные причины
	interpretation.PossibleCauses = i.getPossibleCauses(result)

	// 4. Связанные системы
	interpretation.RelatedSystems = i.getRelatedSystems(result)

	// 5. Связанные чакры
	interpretation.RelatedChakras = i.getRelatedChakras(result)

	// 6. Психологические корреляции
	interpretation.RelatedPsychological = i.getPsychologicalCorrelations(result)

	// 7. Степень серьёзности
	interpretation.Severity = i.getSeverity(result)

	// 8. Когда к врачу
	interpretation.WhenToSeeDoctor = i.getWhenToSeeDoctor(result)

	// 9. Дополнительные анализы
	interpretation.FollowUpTests = i.getFollowUpTests(result)

	// 10. Рекомендации по образу жизни
	interpretation.LifestyleRecommendations = i.getLifestyleRecommendations(result)

	// 11. Рекомендации по добавкам
	interpretation.SupplementRecommendations = i.getSupplementRecommendations(result)

	return interpretation
}

// generateSummary генерирует краткое резюме
func (i *Interpreter) generateSummary(result UserLabResult) string {
	value := result.Value
	unit := result.Unit

	switch result.Status {
	case StatusOptimal:
		return fmt.Sprintf("Ваш %s в оптимальном диапазоне (%.1f %s).", result.TestName, value, unit)
	case StatusLow:
		return fmt.Sprintf("Ваш %s ниже нормы (%.1f %s).", result.TestName, value, unit)
	case StatusHigh:
		return fmt.Sprintf("Ваш %s выше нормы (%.1f %s).", result.TestName, value, unit)
	case StatusCriticalLow:
		return fmt.Sprintf("Ваш %s критически низкий (%.1f %s). Требуется внимание!", result.TestName, value, unit)
	case StatusCriticalHigh:
		return fmt.Sprintf("Ваш %s критически высокий (%.1f %s). Требуется внимание!", result.TestName, value, unit)
	default:
		return fmt.Sprintf("Ваш %s: %.1f %s", result.TestName, value, unit)
	}
}

// generateDetailedExplanation генерирует подробное объяснение
func (i *Interpreter) generateDetailedExplanation(result UserLabResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Референсный диапазон: %.1f - %.1f %s. ",
		result.ReferenceMin, result.ReferenceMax, result.Unit))

	if result.OptimalMin > 0 && result.OptimalMax > 0 {
		sb.WriteString(fmt.Sprintf("Оптимальный диапазон: %.1f - %.1f %s. ",
			result.OptimalMin, result.OptimalMax, result.Unit))
	}

	sb.WriteString(fmt.Sprintf("Ваш результат находится на %.0f процентиле. ", result.Percentile))

	switch result.Status {
	case StatusOptimal:
		sb.WriteString("Это указывает на хорошее функционирование соответствующей системы.")
	case StatusLow, StatusCriticalLow:
		sb.WriteString("Низкие значения могут указывать на дефицит, недостаточную функцию или истощение.")
	case StatusHigh, StatusCriticalHigh:
		sb.WriteString("Высокие значения могут указывать на избыток, гиперфункцию или воспаление.")
	}

	return sb.String()
}

// getPossibleCauses возвращает возможные причины отклонения
func (i *Interpreter) getPossibleCauses(result UserLabResult) []string {
	causes := map[string][]string{
		"cortisol": {
			"Хронический стресс",
			"Нарушения сна",
			"Выгорание надпочечников",
			"Воспалительные процессы",
		},
		"tsh": {
			"Гипотиреоз",
			"Дефицит йода",
			"Аутоиммунный тиреоидит",
			"Стресс",
		},
		"vitamin_d": {
			"Недостаток солнца",
			"Нарушения всасывания",
			"Лишний вес",
			"Возраст",
		},
	}

	if c, ok := causes[result.TestID]; ok {
		return c
	}
	return []string{"Требуется консультация специалиста"}
}

// getRelatedSystems возвращает связанные системы организма
func (i *Interpreter) getRelatedSystems(result UserLabResult) []string {
	systems := map[string][]string{
		"cortisol":     {"Эндокринная", "Нервная", "Иммунная"},
		"tsh":          {"Эндокринная", "Метаболическая"},
		"vitamin_d":    {"Иммунная", "Костная", "Эндокринная"},
		"testosterone": {"Эндокринная", "Репродуктивная"},
	}

	if s, ok := systems[result.TestID]; ok {
		return s
	}
	return []string{"Требуется уточнение"}
}

// getRelatedChakras возвращает связанные чакры
func (i *Interpreter) getRelatedChakras(result UserLabResult) []string {
	correlations := i.getChakraCorrelations(result)
	var chakras []string
	for _, c := range correlations {
		chakras = append(chakras, c.ChakraName)
	}
	if len(chakras) == 0 {
		return []string{"Нет данных"}
	}
	return chakras
}

// getChakraCorrelations возвращает корреляции с чакрами
func (i *Interpreter) getChakraCorrelations(result UserLabResult) []ChakraCorrelation {
	biomarkerName := strings.ToLower(result.TestID)

	if corr, ok := i.biomarkerMap[biomarkerName]; ok {
		direction := "норма"
		if result.Status == StatusLow || result.Status == StatusCriticalLow {
			direction = "низкий"
		} else if result.Status == StatusHigh || result.Status == StatusCriticalHigh {
			direction = "высокий"
		}

		return []ChakraCorrelation{
			{
				ChakraIndex:          corr.ChakraIndex,
				ChakraName:           corr.ChakraName,
				CorrelationType:      corr.Mechanism,
				Strength:             corr.CorrelationStrength,
				Explanation:          fmt.Sprintf("%s %s связан с %s. %s", result.TestName, direction, corr.ChakraName, corr.Direction[direction]),
				RecommendedPractices: corr.RecommendedPractices,
			},
		}
	}
	return []ChakraCorrelation{}
}

// getPsychologicalCorrelations возвращает психологические корреляции
func (i *Interpreter) getPsychologicalCorrelations(result UserLabResult) []string {
	psych := map[string][]string{
		"cortisol": {
			"Тревожность",
			"Страх выживания",
			"Гиперконтроль",
			"Выгорание",
		},
		"tsh": {
			"Депрессия",
			"Замедление мышления",
			"Апатия",
		},
		"testosterone": {
			"Снижение мотивации",
			"Уверенность в себе",
			"Личная сила",
		},
	}

	if p, ok := psych[result.TestID]; ok {
		return p
	}
	return []string{}
}

// getSeverity определяет степень серьёзности
func (i *Interpreter) getSeverity(result UserLabResult) string {
	switch result.Status {
	case StatusCriticalLow, StatusCriticalHigh:
		return "critical"
	case StatusLow, StatusHigh:
		if result.Percentile < 10 || result.Percentile > 90 {
			return "important"
		}
		return "warning"
	default:
		return "info"
	}
}

// getWhenToSeeDoctor возвращает рекомендации по обращению к врачу
func (i *Interpreter) getWhenToSeeDoctor(result UserLabResult) string {
	severity := i.getSeverity(result)
	switch severity {
	case "critical":
		return "Рекомендуется обратиться к врачу в ближайшее время (1-7 дней)"
	case "important":
		return "Рекомендуется проконсультироваться с врачом в течение месяца"
	case "warning":
		return "Можно обсудить с врачом на плановом приёме"
	default:
		return "Профилактический контроль через 6-12 месяцев"
	}
}

// getFollowUpTests возвращает рекомендуемые дополнительные анализы
func (i *Interpreter) getFollowUpTests(result UserLabResult) []string {
	followUp := map[string][]string{
		"cortisol": {
			"ДГЭА-С",
			"Адреналин/Норадреналин",
			"Глюкоза натощак",
		},
		"tsh": {
			"Free T4",
			"Free T3",
			"Антитела к ТПО",
		},
		"vitamin_d": {
			"Кальций",
			"Паратгормон",
			"Магний",
		},
	}

	if f, ok := followUp[result.TestID]; ok {
		return f
	}
	return []string{}
}

// getLifestyleRecommendations возвращает рекомендации по образу жизни
func (i *Interpreter) getLifestyleRecommendations(result UserLabResult) []string {
	lifestyle := map[string][]string{
		"cortisol": {
			"Сон до 23:00",
			"Снижение стресса",
			"Дыхательные практики",
			"Прогулки на природе",
		},
		"tsh": {
			"Йодированная соль",
			"Селен (бразильский орех)",
			"Снижение стресса",
			"Регулярная физическая активность",
		},
		"vitamin_d": {
			"Солнечный свет 15-30 мин/день",
			"Жирная рыба",
			"Яичные желтки",
		},
	}

	if l, ok := lifestyle[result.TestID]; ok {
		return l
	}
	return []string{"Здоровое питание, регулярная активность, достаточный сон"}
}

// getSupplementRecommendations возвращает рекомендации по добавкам
func (i *Interpreter) getSupplementRecommendations(result UserLabResult) []string {
	supplements := map[string][]string{
		"cortisol": {
			"Ашваганда 300-500 мг",
			"Витамин C 500-1000 мг",
			"Магний 400 мг",
		},
		"tsh": {
			"Йод 150 мкг",
			"Селен 200 мкг",
			"Цинк 15-30 мг",
		},
		"vitamin_d": {
			"Витамин D3 2000-5000 МЕ",
			"Витамин K2 100 мкг",
		},
	}

	if s, ok := supplements[result.TestID]; ok {
		return s
	}
	return []string{"По назначению врача"}
}

// InterpretResults интерпретирует несколько результатов
func (i *Interpreter) InterpretResults(results []UserLabResult) []UserLabResult {
	interpreted := make([]UserLabResult, len(results))
	for idx, result := range results {
		result.Interpretation = i.InterpretResult(result)
		result.Interpreted = true
		interpreted[idx] = result
	}
	return interpreted
}
