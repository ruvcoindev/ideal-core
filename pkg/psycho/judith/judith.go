// Package judith — психологическая модель Анодеи Джудит
//
// Контекст:
// ЭТОТ ПАКЕТ — ТОЛЬКО психологические травмы и паттерны.
// Не включает чакры, Эриксона, эндокринную систему.
//
// Источники:
// - Eastern Body Western Mind (2004)
// - Wheels of Life (1987)
//
// Связи с другими системами:
// - pkg/chakra/ — корреляции через pkg/correlations/
// - pkg/psycho/erikson/ — дополняет стадии развития
package judith

// ChakraPattern — психологический паттерн для чакры
type ChakraPattern struct {
	ChakraIndex      int      `json:"chakra_index"`      // 0-6
	ChildhoodPattern string   `json:"childhood_pattern"` // Паттерн из детства
	CoreWound        string   `json:"core_wound"`        // Базовая рана
	DefenseMechanism string   `json:"defense_mechanism"` // Защитный механизм
	AttachmentStyle  string   `json:"attachment_style"`  // Стиль привязанности
	HealingPath      []string `json:"healing_path"`      // Шаги исцеления
}

// Profile — психологический профиль человека
type Profile struct {
	Patterns        []ChakraPattern `json:"patterns"`
	PrimaryWound    string          `json:"primary_wound"`
	HealingPriority []int           `json:"healing_priority"` // Индексы чакр
}

// GetAllPatterns возвращает паттерны для всех 7 чакр
func GetAllPatterns() []ChakraPattern {
	return []ChakraPattern{
		{
			ChakraIndex:      0,
			ChildhoodPattern: "Непредсказуемая забота, отсутствие безопасности в младенчестве",
			CoreWound:        "Я не в безопасности, мир опасен",
			DefenseMechanism: "Отрицание, диссоциация от тела, гиперконтроль",
			AttachmentStyle:  "Тревожный или избегающий тип",
			HealingPath: []string{
				"Заземляющие упражнения",
				"Работа с финансовыми страхами",
				"Создание ритуалов безопасности",
			},
		},
		{
			ChakraIndex:      1,
			ChildhoodPattern: "Подавление эмоций, стыд за естественные потребности",
			CoreWound:        "Мои чувства неправильны, я не должен хотеть",
			DefenseMechanism: "Подавление, отрицание, рационализация",
			AttachmentStyle:  "Тревожно-амбивалентный тип",
			HealingPath: []string{
				"Работа со стыдом и виной",
				"Творческие практики",
				"Здоровое выражение эмоций",
			},
		},
		{
			ChakraIndex:      2,
			ChildhoodPattern: "Подавление воли, критика за инициативу",
			CoreWound:        "Я не имею права на силу, мои действия неправильны",
			DefenseMechanism: "Проекция, рационализация, пассивная агрессия",
			AttachmentStyle:  "Избегающий тип",
			HealingPath: []string{
				"Работа с перфекционизмом",
				"Здоровое выражение гнева",
				"Установление личных границ",
			},
		},
		{
			ChakraIndex:      3,
			ChildhoodPattern: "Условная любовь, критика, эмоциональное пренебрежение",
			CoreWound:        "Я не достоин любви, я должен заслужить принятие",
			DefenseMechanism: "Реактивное образование, подавление, компульсивная забота",
			AttachmentStyle:  "Тревожный тип",
			HealingPath: []string{
				"Работа с обидами и прощением",
				"Развитие самопринятия",
				"Здоровые границы в отношениях",
			},
		},
		{
			ChakraIndex:      4,
			ChildhoodPattern: "Подавление голоса, критика за выражение мнения",
			CoreWound:        "Мой голос не важен, я не должен говорить правду",
			DefenseMechanism: "Подавление, интеллектуализация, пассивная агрессия",
			AttachmentStyle:  "Избегающий тип",
			HealingPath: []string{
				"Практики голоса (пение, декламация)",
				"Честное самовыражение",
				"Творческое письмо",
			},
		},
		{
			ChakraIndex:      5,
			ChildhoodPattern: "Интеллектуальное превосходство, эмоциональное пренебрежение",
			CoreWound:        "Я должен всё понять умом, интуиция ненадёжна",
			DefenseMechanism: "Интеллектуализация, рационализация, диссоциация",
			AttachmentStyle:  "Избегающий тип",
			HealingPath: []string{
				"Медитация и визуализация",
				"Развитие интуиции",
				"Баланс ума и сердца",
			},
		},
		{
			ChakraIndex:      6,
			ChildhoodPattern: "Отсутствие духовного воспитания, материалистическая среда",
			CoreWound:        "Я отделён от источника, жизнь бессмысленна",
			DefenseMechanism: "Отрицание, цинизм, духовное избегание",
			AttachmentStyle:  "Дезорганизованный тип",
			HealingPath: []string{
				"Медитация и молитва",
				"Изучение духовных текстов",
				"Служение другим",
			},
		},
	}
}

// GetPatternByChakra возвращает паттерн для конкретной чакры
func GetPatternByChakra(chakraIndex int) ChakraPattern {
	patterns := GetAllPatterns()
	for _, pattern := range patterns {
		if pattern.ChakraIndex == chakraIndex {
			return pattern
		}
	}
	return ChakraPattern{}
}
