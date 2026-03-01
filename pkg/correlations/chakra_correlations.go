// Package correlations — связи между системами организма
//
// Контекст:
// ЭТОТ ПАКЕТ — ТОЛЬКО связи между системами.
// Не хранит данные о системах, только корреляции.
//
// Архитектура:
// - pkg/chakra/ — чакральная система (независима)
// - pkg/psycho/erikson/ — Эриксон (независим)
// - pkg/psycho/judith/ — Джудит (независима)
// - pkg/bio/endocrine/ — эндокринная (независима)
// - pkg/symptoms/ — симптомы (независимы)
// - pkg/correlations/ — СВЯЗИ между ними (ЭТОТ ПАКЕТ)
//
// Принцип:
// Каждая система самостоятельна. Связи — через этот модуль.
package correlations

// ChakraCorrelation — связи одной чакры с другими системами
type ChakraCorrelation struct {
	ChakraIndex   int      `json:"chakra_index"`    // 0-6
	ChakraName    string   `json:"chakra_name"`     // Название чакры
	EriksonStage  int      `json:"erikson_stage"`   // Стадия Эриксона (1-8)
	JudithPattern int      `json:"judith_pattern"`  // Индекс паттерна Джудит
	EndocrineGland string  `json:"endocrine_gland"` // Железа
	Hormones      []string `json:"hormones"`        // Гормоны
	Symptoms      []string `json:"symptoms"`        // Симптомы
}

// GetAllCorrelations возвращает все корреляции
func GetAllCorrelations() []ChakraCorrelation {
	return []ChakraCorrelation{
		{
			ChakraIndex:    0,
			ChakraName:     "Муладхара",
			EriksonStage:   1,
			JudithPattern:  0,
			EndocrineGland: "Надпочечники",
			Hormones:       []string{"Кортизол", "Адреналин", "ДГЭА-С"},
			Symptoms:       []string{"Хроническая усталость", "Тревожность", "Проблемы с ногами"},
		},
		{
			ChakraIndex:    1,
			ChakraName:     "Свадхистхана",
			EriksonStage:   2,
			JudithPattern:  1,
			EndocrineGland: "Гонады",
			Hormones:       []string{"Эстроген", "Тестостерон", "Прогестерон"},
			Symptoms:       []string{"Гормональные сбои", "ПМС", "Проблемы с тазом"},
		},
		{
			ChakraIndex:    2,
			ChakraName:     "Манипура",
			EriksonStage:   3,
			JudithPattern:  2,
			EndocrineGland: "Поджелудочная железа",
			Hormones:       []string{"Инсулин", "Глюкагон"},
			Symptoms:       []string{"Диабет", "Нарушения пищеварения", "Проблемы с весом"},
		},
		{
			ChakraIndex:    3,
			ChakraName:     "Анахата",
			EriksonStage:   4,
			JudithPattern:  3,
			EndocrineGland: "Тимус",
			Hormones:       []string{"Тимозин", "Т-лимфоциты"},
			Symptoms:       []string{"Аутоиммунные заболевания", "Частые простуды", "Проблемы с сердцем"},
		},
		{
			ChakraIndex:    4,
			ChakraName:     "Вишудха",
			EriksonStage:   5,
			JudithPattern:  4,
			EndocrineGland: "Щитовидная железа",
			Hormones:       []string{"Т3", "Т4", "Кальцитонин"},
			Symptoms:       []string{"Гипотиреоз", "Проблемы с горлом", "Нарушения веса"},
		},
		{
			ChakraIndex:    5,
			ChakraName:     "Аджна",
			EriksonStage:   6,
			JudithPattern:  5,
			EndocrineGland: "Гипофиз",
			Hormones:       []string{"TSH", "ACTH", "GH", "FSH", "LH"},
			Symptoms:       []string{"Головные боли", "Нарушения сна", "Проблемы с глазами"},
		},
		{
			ChakraIndex:    6,
			ChakraName:     "Сахасрара",
			EriksonStage:   7,
			JudithPattern:  6,
			EndocrineGland: "Эпифиз",
			Hormones:       []string{"Мелатонин", "Серотонин"},
			Symptoms:       []string{"Депрессия", "Нарушения сна", "Потеря смысла"},
		},
	}
}

// GetCorrelationByChakra возвращает корреляции для конкретной чакры
func GetCorrelationByChakra(chakraIndex int) ChakraCorrelation {
	correlations := GetAllCorrelations()
	for _, corr := range correlations {
		if corr.ChakraIndex == chakraIndex {
			return corr
		}
	}
	return ChakraCorrelation{}
}

// GetCorrelationByEriksonStage возвращает корреляции по стадии Эриксона
func GetCorrelationByEriksonStage(stage int) ChakraCorrelation {
	correlations := GetAllCorrelations()
	for _, corr := range correlations {
		if corr.EriksonStage == stage {
			return corr
		}
	}
	return ChakraCorrelation{}
}

// GetCorrelationByGland возвращает корреляции по железе
func GetCorrelationByGland(glandName string) ChakraCorrelation {
	correlations := GetAllCorrelations()
	for _, corr := range correlations {
		if corr.EndocrineGland == glandName {
			return corr
		}
	}
	return ChakraCorrelation{}
}
