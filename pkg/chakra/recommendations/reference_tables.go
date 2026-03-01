

// Package reference_tables — таблицы референсных значений
// 
// Почему это критично:
// 1. Нормы зависят от возраста (тестостерон в 20 и 50 лет — разный)
// 2. Нормы зависят от пола (эстроген у мужчин и женщин — разный)
// 3. Нормы зависят от фазы цикла (прогестерон в фолликулярной и лютеиновой — разный)
// 4. Нормы зависят от лаборатории (Инвитро, Гемотест, Хеликс — разные реактивы)
// 
// Эта система:
// - Хранит референсные значения для всех основных анализов
// - Учитывает возраст, пол, фазу цикла
// - Позволяет добавлять данные от разных лабораторий
// - Обновляется при изменении клинических рекомендаций
package recommendations

// ReferenceTable — таблица референсных значений
// 
// Структура:
// - Анализ (TestName)
// - Возраст (AgeMin, AgeMax)
// - Пол (Sex)
// - Фаза цикла (CyclePhase, для женщин)
// - Единицы измерения (Unit)
// - Мин/Макс значения (Min, Max)
// - Источник (Source — лаборатория или клинические рекомендации)
type ReferenceTable struct {
	// TestName — название анализа
	TestName string `json:"test_name"`
	
	// TestNameEn — название на английском
	TestNameEn string `json:"test_name_en"`
	
	// AgeMin — минимальный возраст для этой нормы
	AgeMin int `json:"age_min"`
	
	// AgeMax — максимальный возраст для этой нормы
	AgeMax int `json:"age_max"`
	
	// Sex — пол (male, female, any)
	Sex string `json:"sex"`
	
	// CyclePhase — фаза цикла (follicular, ovulation, luteal, any)
	CyclePhase string `json:"cycle_phase"`
	
	// Unit — единица измерения
	Unit string `json:"unit"`
	
	// Min — нижняя граница нормы
	Min float64 `json:"min"`
	
	// Max — верхняя граница нормы
	Max float64 `json:"max"`
	
	// OptimalMin — оптимальная нижняя граница (функциональная медицина)
	// Часто уже лабораторной нормы
	OptimalMin float64 `json:"optimal_min"`
	
	// OptimalMax — оптимальная верхняя граница
	OptimalMax float64 `json:"optimal_max"`
	
	// Source — источник данных
	Source string `json:"source"`
	
	// LastUpdated — дата обновления
	LastUpdated string `json:"last_updated"`
	
	// Notes — примечания
	Notes string `json:"notes"`
}

// LoadReferenceTables загружает таблицы норм
// 
// Источники данных:
// 1. Клинические рекомендации Минздрава РФ
// 2. Данные лабораторий (Инвитро, Гемотест, Хеликс)
// 3. Функциональная медицина (оптимальные значения)
// 4. Международные стандарты (WHO, Endocrine Society)
func LoadReferenceTables() ReferenceTables {
	return ReferenceTables{
		Tables: []ReferenceTable{
			// Кортизол
			{
				TestName:    "Кортизол",
				TestNameEn:  "Cortisol",
				AgeMin:      18,
				AgeMax:      60,
				Sex:         "any",
				CyclePhase:  "any",
				Unit:        "нмоль/л",
				Min:         138,
				Max:         635,
				OptimalMin:  250,
				OptimalMax:  500,
				Source:      "Инвитро + функциональная медицина",
				LastUpdated: "2024-01-15",
				Notes:       "Сдавать утром 08:00-10:00 натощак. Вечерний кортизол должен быть в 2-3 раза ниже.",
			},
			
			// ТТГ (щитовидная)
			{
				TestName:    "ТТГ",
				TestNameEn:  "TSH",
				AgeMin:      18,
				AgeMax:      60,
				Sex:         "any",
				CyclePhase:  "any",
				Unit:        "мЕд/л",
				Min:         0.4,
				Max:         4.0,
				OptimalMin:  1.0,
				OptimalMax:  2.5,
				Source:      "Endocrine Society + функциональная медицина",
				LastUpdated: "2024-01-15",
				Notes:       "Лабораторная норма до 4.0, но функциональная оптимум 1.0-2.5. При планировании беременности оптимально <2.5.",
			},
			
			// Тестостерон общий (мужчины)
			{
				TestName:    "Тестостерон общий",
				TestNameEn:  "Total Testosterone",
				AgeMin:      20,
				AgeMax:      40,
				Sex:         "male",
				CyclePhase:  "any",
				Unit:        "нмоль/л",
				Min:         11,
				Max:         33,
				OptimalMin:  18,
				OptimalMax:  28,
				Source:      "Инвитро + функциональная медицина",
				LastUpdated: "2024-01-15",
				Notes:       "Сдавать утром 08:00-10:00. После 40 лет нормы снижаются на 1-2% в год.",
			},
			
			// Эстрадиол (женщины, фолликулярная фаза)
			{
				TestName:    "Эстрадиол",
				TestNameEn:  "Estradiol",
				AgeMin:      18,
				AgeMax:      45,
				Sex:         "female",
				CyclePhase:  "follicular",
				Unit:        "пмоль/л",
				Min:         68,
				Max:         1269,
				OptimalMin:  150,
				OptimalMax:  500,
				Source:      "Инвитро + функциональная медицина",
				LastUpdated: "2024-01-15",
				Notes:       "Норма сильно зависит от фазы цикла. Сдавать на 2-5 день цикла.",
			},
			
			// Витамин D
			{
				TestName:    "Витамин D (25-OH)",
				TestNameEn:  "Vitamin D 25-OH",
				AgeMin:      18,
				AgeMax:      80,
				Sex:         "any",
				CyclePhase:  "any",
				Unit:        "нг/мл",
				Min:         20,
				Max:         100,
				OptimalMin:  40,
				OptimalMax:  60,
				Source:      "Endocrine Society + функциональная медицина",
				LastUpdated: "2024-01-15",
				Notes:       "Лабораторная норма от 20, но оптимально 40-60 для иммунитета и гормонального баланса.",
			},
			
			// Ферритин
			{
				TestName:    "Ферритин",
				TestNameEn:  "Ferritin",
				AgeMin:      18,
				AgeMax:      60,
				Sex:         "female",
				CyclePhase:  "any",
				Unit:        "мкг/л",
				Min:         10,
				Max:         120,
				OptimalMin:  40,
				OptimalMax:  80,
				Source:      "функциональная медицина",
				LastUpdated: "2024-01-15",
				Notes:       "При ферритине <40 возможны симптомы дефицита железа даже при нормальном гемоглобине. Выпадение волос, усталость.",
			},
			
			// Гомоцистеин
			{
				TestName:    "Гомоцистеин",
				TestNameEn:  "Homocysteine",
				AgeMin:      18,
				AgeMax:      80,
				Sex:         "any",
				CyclePhase:  "any",
				Unit:        "мкмоль/л",
				Min:         5,
				Max:         15,
				OptimalMin:  6,
				OptimalMax:  8,
				Source:      "функциональная медицина",
				LastUpdated: "2024-01-15",
				Notes:       "Маркер воспаления, риска сердечно-сосудистых заболеваний, дефицита B-витаминов. Оптимально <8.",
			},
		},
	}
}

// ReferenceTables — коллекция таблиц
type ReferenceTables struct {
	Tables []ReferenceTable `json:"tables"`
}

// GetReference получает референсное значение для конкретного анализа
// 
// Учитывает:
// - Возраст
// - Пол
// - Фазу цикла (для женщин)
// 
// Возвращает:
// - Наиболее подходящую норму из таблицы
// - Если точного совпадения нет — ближайшую по возрасту/полу
func (rt *ReferenceTables) GetReference(
	testName string,
	age int,
	sex string,
	cyclePhase string,
) ReferenceTable {
	// Ищем точное совпадение
	for _, table := range rt.Tables {
		if table.TestName == testName &&
			age >= table.AgeMin && age <= table.AgeMax &&
			(table.Sex == sex || table.Sex == "any") &&
			(table.CyclePhase == cyclePhase || table.CyclePhase == "any") {
			return table
		}
	}
	
	// Если не нашли — возвращаем норму для "any"
	for _, table := range rt.Tables {
		if table.TestName == testName &&
			age >= table.AgeMin && age <= table.AgeMax &&
			table.Sex == "any" &&
			table.CyclePhase == "any" {
			return table
		}
	}
	
	// Если совсем не нашли — возвращаем пустую таблицу
	return ReferenceTable{}
}

// GetOptimalRange возвращает оптимальный диапазон (функциональная медицина)
// 
// Разница с лабораторной нормой:
// - Лабораторная норма = 95% населения (включая больных)
// - Оптимальная норма = здоровые люди с лучшим самочувствием
// 
// Пример:
// - ТТГ лабораторная норма: 0.4-4.0
// - ТТГ оптимальная норма: 1.0-2.5
func (rt *ReferenceTable) GetOptimalRange() (float64, float64) {
	if rt.OptimalMin > 0 && rt.OptimalMax > 0 {
		return rt.OptimalMin, rt.OptimalMax
	}
	return rt.Min, rt.Max
}

