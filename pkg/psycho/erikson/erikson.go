// Package erikson — теория психосоциального развития Эрика Эриксона
//
// Контекст:
// ЭТОТ ПАКЕТ — ТОЛЬКО психологическое развитие (8 стадий).
// Не включает чакры, эндокринную систему, симптомы.
//
// Связи с другими системами:
// - pkg/chakra/ — корреляции через pkg/correlations/
// - pkg/psycho/judith/ — дополняет модель Джудит
//
// 8 стадий развития:
// 1. Доверие vs Недоверие (0-1.5 года) → Надежда
// 2. Автономия vs Стыд (1.5-3 года) → Воля
// 3. Инициатива vs Вина (3-6 лет) → Цель
// 4. Трудолюбие vs Неполноценность (6-12 лет) → Компетентность
// 5. Идентичность vs Ролевое смешение (12-18 лет) → Верность
// 6. Близость vs Изоляция (18-40 лет) → Любовь
// 7. Продуктивность vs Стагнация (40-65 лет) → Забота
// 8. Целостность vs Отчаяние (65+ лет) → Мудрость
package erikson

// Stage — стадия психосоциального развития
type Stage struct {
	Number       int    `json:"number"`         // 1-8
	AgeRange     string `json:"age_range"`      // Возрастной диапазон
	Crisis       string `json:"crisis"`         // Психосоциальный кризис
	Virtue       string `json:"virtue"`         // Формируемая добродетель
	SignsOfBlock string `json:"signs_of_block"` // Признаки непрохождения
	HealingFocus string `json:"healing_focus"`  // Фокус для исцеления
}

// Profile — профиль развития человека
type Profile struct {
	CurrentStage     int     `json:"current_stage"`      // Текущая стадия
	Age              int     `json:"age"`                // Возраст пользователя
	BlockedStages    []int   `json:"blocked_stages"`     // Заблокированные стадии
	PassedStages     []int   `json:"passed_stages"`      // Пройденные стадии
	OverallProgress  float64 `json:"overall_progress"`   // 0.0-1.0
}

// GetAllStages возвращает все 8 стадий развития
func GetAllStages() []Stage {
	return []Stage{
		{
			Number:       1,
			AgeRange:     "0-1.5 года",
			Crisis:       "Доверие vs Недоверие",
			Virtue:       "Надежда",
			SignsOfBlock: "Базовое недоверие к миру, тревожность, трудности с зависимостью от других",
			HealingFocus: "Создание чувства безопасности, заземляющие практики",
		},
		{
			Number:       2,
			AgeRange:     "1.5-3 года",
			Crisis:       "Автономия vs Стыд/Сомнение",
			Virtue:       "Воля",
			SignsOfBlock: "Стыд за свои желания, трудности с автономией, чрезмерный контроль",
			HealingFocus: "Развитие здоровой воли, принятие эмоций",
		},
		{
			Number:       3,
			AgeRange:     "3-6 лет",
			Crisis:       "Инициатива vs Вина",
			Virtue:       "Цель",
			SignsOfBlock: "Чувство вины за инициативу, страх действовать, пассивность",
			HealingFocus: "Развитие здоровой инициативы, работа с чувством вины",
		},
		{
			Number:       4,
			AgeRange:     "6-12 лет",
			Crisis:       "Трудолюбие vs Неполноценность",
			Virtue:       "Компетентность",
			SignsOfBlock: "Чувство неполноценности, трудности с признанием достижений",
			HealingFocus: "Развитие самоценности, принятие достижений",
		},
		{
			Number:       5,
			AgeRange:     "12-18 лет",
			Crisis:       "Идентичность vs Ролевое смешение",
			Virtue:       "Верность",
			SignsOfBlock: "Путаница в идентичности, трудности с самовыражением",
			HealingFocus: "Развитие аутентичности, поиск своего голоса",
		},
		{
			Number:       6,
			AgeRange:     "18-40 лет",
			Crisis:       "Близость vs Изоляция",
			Virtue:       "Любовь",
			SignsOfBlock: "Трудности с близостью, изоляция, страх уязвимости",
			HealingFocus: "Развитие способности к близости, баланс автономии и связи",
		},
		{
			Number:       7,
			AgeRange:     "40-65 лет",
			Crisis:       "Продуктивность vs Стагнация",
			Virtue:       "Забота",
			SignsOfBlock: "Стагнация, отсутствие продуктивности, эгоцентризм",
			HealingFocus: "Развитие заботы о следующих поколениях, поиск смысла",
		},
		{
			Number:       8,
			AgeRange:     "65+ лет",
			Crisis:       "Целостность vs Отчаяние",
			Virtue:       "Мудрость",
			SignsOfBlock: "Отчаяние, сожаления, страх смерти",
			HealingFocus: "Принятие жизненного пути, интеграция опыта",
		},
	}
}

// GetStageByNumber возвращает стадию по номеру (1-8)
func GetStageByNumber(number int) Stage {
	stages := GetAllStages()
	for _, stage := range stages {
		if stage.Number == number {
			return stage
		}
	}
	return Stage{}
}

// GetStageByAge возвращает стадию по возрасту
func GetStageByAge(age int) Stage {
	stages := GetAllStages()
	for _, stage := range stages {
		if ageInRange(age, stage.AgeRange) {
			return stage
		}
	}
	return stages[len(stages)-1] // Последняя стадия для 65+
}

// ageInRange проверяет, попадает ли возраст в диапазон
func ageInRange(age int, rangeStr string) bool {
	// Простая реализация (можно улучшить парсингом)
	switch rangeStr {
	case "0-1.5 года":
		return age <= 1
	case "1.5-3 года":
		return age >= 1 && age <= 3
	case "3-6 лет":
		return age >= 3 && age <= 6
	case "6-12 лет":
		return age >= 6 && age <= 12
	case "12-18 лет":
		return age >= 12 && age <= 18
	case "18-40 лет":
		return age >= 18 && age <= 40
	case "40-65 лет":
		return age >= 40 && age <= 65
	case "65+ лет":
		return age >= 65
	}
	return false
}
