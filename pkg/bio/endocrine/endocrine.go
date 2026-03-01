// Package endocrine — эндокринная система человека
//
// Контекст:
// ЭТОТ ПАКЕТ — ТОЛЬКО эндокринная система (гормоны, железы).
// Не включает чакры, психологию, симптомы.
//
// Связи с другими системами:
// - pkg/chakra/ — корреляции через pkg/correlations/
// - pkg/symptoms/ — симптомы гормональных нарушений
// - pkg/bio/ — другие биомаркеры
package endocrine

// GlandInfo — информация об эндокринной железе
type GlandInfo struct {
	Index          int      `json:"index"`           // 0-6 (корреляция с чакрой)
	Name           string   `json:"name"`            // Название на русском
	NameEn         string   `json:"name_en"`         // Название на английском
	Hormones       []string `json:"hormones"`        // Производимые гормоны
	Functions      []string `json:"functions"`       // Физиологические функции
	ImbalanceSigns []string `json:"imbalance_signs"` // Признаки дисбаланса
	LabTests       []string `json:"lab_tests"`       // Рекомендуемые анализы
}

// HormoneInfo — информация о гормоне
type HormoneInfo struct {
	Name            string  `json:"name"`             // Название
	NameEn          string  `json:"name_en"`          // Для лабораторий
	Unit            string  `json:"unit"`             // Единица измерения
	ReferenceMin    float64 `json:"reference_min"`    // Нижняя граница нормы
	ReferenceMax    float64 `json:"reference_max"`    // Верхняя граница нормы
	OptimalMin      float64 `json:"optimal_min"`      // Функциональная норма (мин)
	OptimalMax      float64 `json:"optimal_max"`      // Функциональная норма (макс)
	TimeOfDay       string  `json:"time_of_day"`      // Оптимальное время сдачи
	FastingRequired bool    `json:"fasting_required"` // Натощак или нет
}

// GetAllGlands возвращает все эндокринные железы
func GetAllGlands() []GlandInfo {
	return []GlandInfo{
		{
			Index:          0,
			Name:           "Надпочечники",
			NameEn:         "Adrenal Glands",
			Hormones:       []string{"Кортизол", "Адреналин", "ДГЭА-С"},
			Functions:      []string{"Реакция на стресс", "Регуляция давления", "Иммунная модуляция"},
			ImbalanceSigns: []string{"Хроническая усталость", "Тревожность", "Нарушения веса"},
			LabTests:       []string{"Кортизол (слюна, 4 точки)", "ДГЭА-С", "Адреналин/Норадреналин"},
		},
		{
			Index:          1,
			Name:           "Гонады",
			NameEn:         "Gonads",
			Hormones:       []string{"Эстроген", "Тестостерон", "Прогестерон"},
			Functions:      []string{"Репродукция", "Либидо", "Эмоциональный баланс"},
			ImbalanceSigns: []string{"Гормональные сбои", "ПМС", "Снижение либидо"},
			LabTests:       []string{"Эстрадиол", "Тестостерон общий", "Прогестерон"},
		},
		{
			Index:          2,
			Name:           "Поджелудочная железа",
			NameEn:         "Pancreas",
			Hormones:       []string{"Инсулин", "Глюкагон"},
			Functions:      []string{"Регуляция сахара", "Пищеварение", "Метаболизм"},
			ImbalanceSigns: []string{"Диабет", "Нарушения пищеварения", "Проблемы с весом"},
			LabTests:       []string{"Глюкоза", "Инсулин", "Гликированный гемоглобин"},
		},
		{
			Index:          3,
			Name:           "Тимус",
			NameEn:         "Thymus",
			Hormones:       []string{"Тимозин", "Т-лимфоциты"},
			Functions:      []string{"Иммунная регуляция", "Созревание Т-клеток"},
			ImbalanceSigns: []string{"Аутоиммунные заболевания", "Частые инфекции"},
			LabTests:       []string{"Иммунограмма", "Т-клетки субпопуляции"},
		},
		{
			Index:          4,
			Name:           "Щитовидная железа",
			NameEn:         "Thyroid",
			Hormones:       []string{"Т3", "Т4", "Кальцитонин"},
			Functions:      []string{"Метаболизм", "Рост и развитие", "Энергетический баланс"},
			ImbalanceSigns: []string{"Гипотиреоз/гипертиреоз", "Нарушения веса", "Усталость"},
			LabTests:       []string{"ТТГ", "Свободный Т3", "Свободный Т4", "Антитела к ТПО"},
		},
		{
			Index:          5,
			Name:           "Гипофиз",
			NameEn:         "Pituitary",
			Hormones:       []string{"TSH", "ACTH", "GH", "FSH", "LH"},
			Functions:      []string{"Мастер-железа", "Регуляция других желёз", "Цикл сна"},
			ImbalanceSigns: []string{"Гормональные сбои", "Нарушения сна", "Неврологические расстройства"},
			LabTests:       []string{"Пролактин", "СТГ", "АКТГ"},
		},
		{
			Index:          6,
			Name:           "Эпифиз",
			NameEn:         "Pineal",
			Hormones:       []string{"Мелатонин", "Серотонин"},
			Functions:      []string{"Циркадные ритмы", "Сон и бодрствование", "Духовное восприятие"},
			ImbalanceSigns: []string{"Нарушения сна", "Депрессия", "Сезонное расстройство"},
			LabTests:       []string{"Мелатонин (слюна)", "Серотонин"},
		},
	}
}

// GetGlandByIndex возвращает железу по индексу (0-6)
func GetGlandByIndex(index int) GlandInfo {
	glands := GetAllGlands()
	for _, gland := range glands {
		if gland.Index == index {
			return gland
		}
	}
	return GlandInfo{}
}

// GetGlandByName возвращает железу по названию
func GetGlandByName(name string) GlandInfo {
	glands := GetAllGlands()
	for _, gland := range glands {
		if gland.Name == name || gland.NameEn == name {
			return gland
		}
	}
	return GlandInfo{}
}
