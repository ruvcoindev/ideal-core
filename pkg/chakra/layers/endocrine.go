// Package layers — слои модели человека
//
// Контекст:
// Человек — многомерная система. Каждый слой влияет на остальные:
// - Энергетический блок → гормональный сбой → психологическая проблема
// - Психологическая травма → энергетический зажим → физический симптом
//
// Этот файл реализует эндокринный слой.
package layers

// EndocrineLayer — эндокринный слой для одной чакры
type EndocrineLayer struct {
	ChakraIndex      int      `json:"chakra_index"`
	Gland            string   `json:"gland"`
	GlandEn          string   `json:"gland_en"`
	Hormones         []string `json:"hormones"`
	Functions        []string `json:"functions"`
	ImbalanceSigns   []string `json:"imbalance_signs"`
	LabTests         []string `json:"lab_tests"`
	SupportMethods   []string `json:"support_methods"`
}

// GetEndocrineLayers возвращает эндокринные слои для всех чакр
func GetEndocrineLayers() []EndocrineLayer {
	return []EndocrineLayer{
		{
			ChakraIndex:    0,
			Gland:          "Надпочечники",
			GlandEn:        "Adrenal Glands",
			Hormones:       []string{"Кортизол", "Адреналин", "ДГЭА-С"},
			Functions:      []string{"Реакция на стресс", "Регуляция давления", "Иммунная модуляция"},
			ImbalanceSigns: []string{"Хроническая усталость", "Тревожность", "Нарушения веса"},
			LabTests:       []string{"Кортизол (слюна, 4 точки)", "ДГЭА-С", "Адреналин/Норадреналин"},
			SupportMethods: []string{"Ашваганда", "Витамин C", "Сон до 23:00"},
		},
		{
			ChakraIndex:    1,
			Gland:          "Гонады",
			GlandEn:        "Gonads",
			Hormones:       []string{"Эстроген", "Тестостерон", "Прогестерон"},
			Functions:      []string{"Репродукция", "Либидо", "Эмоциональный баланс"},
			ImbalanceSigns: []string{"Гормональные сбои", "ПМС", "Снижение либидо"},
			LabTests:       []string{"Эстрадиол", "Тестостерон общий", "Прогестерон"},
			SupportMethods: []string{"Витекс", "Мака", "Омега-3"},
		},
		{
			ChakraIndex:    2,
			Gland:          "Поджелудочная железа",
			GlandEn:        "Pancreas",
			Hormones:       []string{"Инсулин", "Глюкагон"},
			Functions:      []string{"Регуляция сахара", "Пищеварение", "Метаболизм"},
			ImbalanceSigns: []string{"Диабет", "Нарушения пищеварения", "Проблемы с весом"},
			LabTests:       []string{"Глюкоза", "Инсулин", "Гликированный гемоглобин"},
			SupportMethods: []string{"Хром", "Корица", "Регулярное питание"},
		},
		{
			ChakraIndex:    3,
			Gland:          "Тимус",
			GlandEn:        "Thymus",
			Hormones:       []string{"Тимозин", "Т-лимфоциты"},
			Functions:      []string{"Иммунная регуляция", "Созревание Т-клеток"},
			ImbalanceSigns: []string{"Аутоиммунные заболевания", "Частые инфекции"},
			LabTests:       []string{"Иммунограмма", "Т-клетки субпопуляции"},
			SupportMethods: []string{"Витамин D", "Цинк", "Глубокое дыхание"},
		},
		{
			ChakraIndex:    4,
			Gland:          "Щитовидная железа",
			GlandEn:        "Thyroid",
			Hormones:       []string{"Т3", "Т4", "Кальцитонин"},
			Functions:      []string{"Метаболизм", "Рост и развитие", "Энергетический баланс"},
			ImbalanceSigns: []string{"Гипотиреоз/гипертиреоз", "Нарушения веса", "Усталость"},
			LabTests:       []string{"ТТГ", "Свободный Т3", "Свободный Т4", "Антитела к ТПО"},
			SupportMethods: []string{"Йод", "Селен", "Управление стрессом"},
		},
		{
			ChakraIndex:    5,
			Gland:          "Гипофиз",
			GlandEn:        "Pituitary",
			Hormones:       []string{"TSH", "ACTH", "GH", "FSH", "LH"},
			Functions:      []string{"Мастер-железа", "Регуляция других желёз", "Цикл сна"},
			ImbalanceSigns: []string{"Гормональные сбои", "Нарушения сна", "Неврологические расстройства"},
			LabTests:       []string{"Пролактин", "СТГ", "АКТГ"},
			SupportMethods: []string{"Медитация", "Достаточный сон", "Омега-3"},
		},
		{
			ChakraIndex:    6,
			Gland:          "Эпифиз",
			GlandEn:        "Pineal",
			Hormones:       []string{"Мелатонин", "Серотонин"},
			Functions:      []string{"Циркадные ритмы", "Сон и бодрствование", "Духовное восприятие"},
			ImbalanceSigns: []string{"Нарушения сна", "Депрессия", "Сезонное расстройство"},
			LabTests:       []string{"Мелатонин (слюна)", "Серотонин"],
			SupportMethods: []string{"Тёмная комната", "Солнечный свет днём", "Магний"},
		},
	}
}
