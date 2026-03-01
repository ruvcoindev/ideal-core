package layers

import "ideal-core/pkg/chakra"

// EndocrineLayer связывает чакру с эндокринной железой
type EndocrineLayer struct {
	Chakra          chakra.ChakraIndex
	Gland           string
	Hormones        []string
	Dysregulation   []string
	RecommendedTests []string // Лабораторные тесты для проверки
}

// EndocrineMap — карта связей чакра-железа
var EndocrineMap = map[chakra.ChakraIndex]EndocrineLayer{
	chakra.Muladhara: {
		Chakra:          chakra.Muladhara,
		Gland:           "Надпочечники",
		Hormones:        []string{"Кортизол", "Адреналин", "Норадреналин"},
		Dysregulation:   []string{"Хронический стресс", "Усталость надпочечников", "Тревожность"},
		RecommendedTests: []string{"Кортизол слюна (4 точки)", "DHEA-S", "Адреналин/Норадреналин моча"},
	},
	chakra.Svadhisthana: {
		Chakra:          chakra.Svadhisthana,
		Gland:           "Гонады (яичники/яички)",
		Hormones:        []string{"Эстрадиол", "Прогестерон", "Тестостерон", "ЛГ", "ФСГ"},
		Dysregulation:   []string{"ПМС", "Нарушения цикла", "Снижение либидо", "Эмоциональные качели"},
		RecommendedTests: []string{"Половые гормоны панель (день 3-5 цикла)", "ГСПГ", "Пролактин"},
	},
	chakra.Manipura: {
		Chakra:          chakra.Manipura,
		Gland:           "Поджелудочная железа",
		Hormones:        []string{"Инсулин", "Глюкагон", "С-пептид"},
		Dysregulation:   []string{"Инсулинорезистентность", "Скачки сахара", "Эмоциональный голод"},
		RecommendedTests: []string{"Глюкоза натощак", "Инсулин натощак", "HOMA-IR", "Гликированный гемоглобин"},
	},
	chakra.Anahata: {
		Chakra:          chakra.Anahata,
		Gland:           "Тимус (вилочковая железа)",
		Hormones:        []string{"Тимулин", "Тимозин", "Цитокины"},
		Dysregulation:   []string{"Снижение иммунитета", "Частые простуды", "Аутоиммунные реакции"},
		RecommendedTests: []string{"Иммунограмма", "Витамин D", "Цинк", "Селен"},
	},
	chakra.Vishuddha: {
		Chakra:          chakra.Vishuddha,
		Gland:           "Щитовидная железа",
		Hormones:        []string{"ТТГ", "Т3 свободный", "Т4 свободный", "Антитела к ТПО"},
		Dysregulation:   []string{"Гипотиреоз", "Гипертиреоз", "Узелки", "Проблемы с голосом"},
		RecommendedTests: []string{"ТТГ", "Т3 св.", "Т4 св.", "АТ-ТПО", "АТ-ТГ", "УЗИ щитовидной железы"},
	},
	chakra.Ajna: {
		Chakra:          chakra.Ajna,
		Gland:           "Гипофиз",
		Hormones:        []string{"СТГ", "Пролактин", "АКТГ", "ТТГ", "ЛГ", "ФСГ"},
		Dysregulation:   []string{"Дисбаланс оси HPA", "Нарушения цикла", "Проблемы со сном"},
		RecommendedTests: []string{"Пролактин", "АКТГ", "Кортизол", "ИФР-1"},
	},
	chakra.Sahasrara: {
		Chakra:          chakra.Sahasrara,
		Gland:           "Эпифиз (шишковидная железа)",
		Hormones:        []string{"Мелатонин", "Серотонин"},
		Dysregulation:   []string{"Бессонница", "Сезонная депрессия", "Нарушения циркадных ритмов"},
		RecommendedTests: []string{"Мелатонин слюна (вечер)", "Серотонин кровь", "Витамин B6", "Магний"},
	},
}

// GetEndocrineLayers возвращает слои для активных чакр
func GetEndocrineLayers(activeChakras []chakra.ChakraIndex) []EndocrineLayer {
	var layers []EndocrineLayer
	for _, idx := range activeChakras {
		if layer, ok := EndocrineMap[idx]; ok {
			layers = append(layers, layer)
		}
	}
	return layers
}

// GetRecommendedTests возвращает список рекомендованных анализов
func GetRecommendedTests(activeChakras []chakra.ChakraIndex) []string {
	tests := make(map[string]bool)
	for _, idx := range activeChakras {
		if layer, ok := EndocrineMap[idx]; ok {
			for _, t := range layer.RecommendedTests {
				tests[t] = true
			}
		}
	}
	var result []string
	for t := range tests {
		result = append(result, t)
	}
	return result
}
