// Package symptoms — система симптомов для диагностики
//
// Контекст:
// Этот модуль — ОБЩИЙ для всех систем проекта ideal-core:
// - chakra (энергетическая)
// - bio (физиологическая)
// - psycho (психологическая)
// - neural (нервная)
//
// Философия:
// - Симптомы — это сигналы тела о дисбалансе
// - Один симптом может указывать на несколько систем
// - Система помогает найти корневую причину, а не лечить следствие
//
// Использование:
// 1. Пользователь вводит симптомы
// 2. Система сопоставляет с возможными дисбалансами
// 3. Рекомендует анализы/практики/специалистов
//
// Важно:
// Это НЕ замена врачу! Система помогает:
// - Подготовиться к визиту к врачу
// - Понять, какие вопросы задать
// - Отслеживать динамику
package symptoms

// Symptom — отдельный симптом
//
// Концепция:
// Каждый симптом имеет:
// - Уникальный ID для связки между системами
// - Категорию (физический, эмоциональный, ментальный)
// - Связь с системами организма
// - Приоритет (когда срочно к врачу)
type Symptom struct {
	// ID — уникальный идентификатор
	ID string `json:"id"`

	// Name — название на русском
	Name string `json:"name"`

	// NameEn — название на английском (для медицинских источников)
	NameEn string `json:"name_en"`

	// Category — категория симптома
	Category string `json:"category"` // physical | emotional | mental | spiritual

	// Description — подробное описание
	Description string `json:"description"`

	// RelatedSystems — связанные системы организма
	RelatedSystems []string `json:"related_systems"` // ["endocrine", "nervous", "chakra"]

	// RelatedChakras — связанные чакры (индексы 0-6)
	RelatedChakras []int `json:"related_chakras"`

	// RelatedHormones — связанные гормоны (для анализов)
	RelatedHormones []string `json:"related_hormones"`

	// Priority — приоритет внимания
	Priority int `json:"priority"` // 1-критичный, 2-важный, 3-рекомендуемый

	// WhenToSeeDoctor — когда срочно к врачу
	WhenToSeeDoctor string `json:"when_to_see_doctor"`

	// CommonCauses — распространённые причины
	CommonCauses []string `json:"common_causes"`

	// AssociatedSymptoms — сопутствующие симптомы
	AssociatedSymptoms []string `json:"associated_symptoms"`
}

// SymptomDatabase — база всех симптомов
type SymptomDatabase struct {
	Symptoms []Symptom `json:"symptoms"`
}

// SymptomMatch — результат поиска по симптомам
type SymptomMatch struct {
	Symptom      Symptom   `json:"symptom"`
	Confidence   float64   `json:"confidence"` // 0.0-1.0
	RelatedDiagnoses []string `json:"related_diagnoses"`
}

// GetSymptomDatabase возвращает базу симптомов
//
// Источники:
// - МКБ-10 (международная классификация болезней)
// - Клинические рекомендации Минздрава РФ
// - Функциональная медицина
// - Психосоматические справочники
func GetSymptomDatabase() SymptomDatabase {
	return SymptomDatabase{
		Symptoms: []Symptom{
			// === ФИЗИЧЕСКИЕ СИМПТОМЫ ===
			{
				ID:          "PHYS_001",
				Name:        "Хроническая усталость",
				NameEn:      "Chronic Fatigue",
				Category:    "physical",
				Description: "Постоянное чувство усталости, не проходящее после отдыха",
				RelatedSystems: []string{"endocrine", "nervous", "immune", "chakra"},
				RelatedChakras: []int{0, 2, 6}, // Муладхара, Манипура, Сахасрара
				RelatedHormones: []string{"Cortisol", "TSH", "Free T4", "Vitamin D", "Ferritin"},
				Priority:      2,
				WhenToSeeDoctor: "Если длится более 2 недель, сопровождается потерей веса или температурой",
				CommonCauses: []string{
					"Выгорание надпочечников",
					"Гипотиреоз",
					"Дефицит железа",
					"Дефицит витамина D",
					"Хронический стресс",
				},
				AssociatedSymptoms: []string{"бессонница", "туман в голове", "слабость в мышцах"},
			},
			{
				ID:          "PHYS_002",
				Name:        "Головные боли",
				NameEn:      "Headaches",
				Category:    "physical",
				Description: "Боль в области головы различной интенсивности",
				RelatedSystems: []string{"nervous", "vascular", "chakra"},
				RelatedChakras: []int{5, 6}, // Аджна, Сахасрара
				RelatedHormones: []string{"Cortisol", "Estrogen", "Serotonin"},
				Priority:      2,
				WhenToSeeDoctor: "При внезапной сильной боли, с нарушением зрения или речи",
				CommonCauses: []string{
					"Стресс и напряжение",
					"Гормональные колебания",
					"Обезвоживание",
					"Проблемы с шейным отделом",
				},
				AssociatedSymptoms: []string{"тошнота", "светочувствительность", "головокружение"},
			},
			{
				ID:          "PHYS_003",
				Name:        "Проблемы с пищеварением",
				NameEn:      "Digestive Issues",
				Category:    "physical",
				Description: "Вздутие, боли в животе, нарушения стула",
				RelatedSystems: []string{"digestive", "endocrine", "nervous", "chakra"},
				RelatedChakras: []int{2}, // Манипура
				RelatedHormones: []string{"Cortisol", "Insulin", "Gastrin"},
				Priority:      2,
				WhenToSeeDoctor: "При крови в стуле, потере веса, постоянной боли",
				CommonCauses: []string{
					"Стресс",
					"Неправильное питание",
					"Дисбактериоз",
					"Пищевая непереносимость",
				},
				AssociatedSymptoms: []string{"тошнота", "изжога", "вздутие"},
			},

			// === ЭМОЦИОНАЛЬНЫЕ СИМПТОМЫ ===
			{
				ID:          "EMO_001",
				Name:        "Тревожность",
				NameEn:      "Anxiety",
				Category:    "emotional",
				Description: "Постоянное чувство беспокойства, страха без явной причины",
				RelatedSystems: []string{"nervous", "endocrine", "chakra"},
				RelatedChakras: []int{0, 1, 4}, // Муладхара, Свадхистхана, Анахата
				RelatedHormones: []string{"Cortisol", "Adrenaline", "GABA", "Serotonin"},
				Priority:      2,
				WhenToSeeDoctor: "При панических атаках, нарушении сна, суицидальных мыслях",
				CommonCauses: []string{
					"Хронический стресс",
					"Дисбаланс нейромедиаторов",
					"Травма привязанности",
					"Гормональные нарушения",
				},
				AssociatedSymptoms: []string{"учащённое сердцебиение", "потливость", "бессонница"},
			},
			{
				ID:          "EMO_002",
				Name:        "Депрессивное состояние",
				NameEn:      "Depression",
				Category:    "emotional",
				Description: "Стойкое снижение настроения, потеря интереса к жизни",
				RelatedSystems: []string{"nervous", "endocrine", "chakra"},
				RelatedChakras: []int{3, 6}, // Анахата, Сахасрара
				RelatedHormones: []string{"Serotonin", "Dopamine", "TSH", "Vitamin D"},
				Priority:      1,
				WhenToSeeDoctor: "При суицидальных мыслях, невозможности функционировать",
				CommonCauses: []string{
					"Химический дисбаланс мозга",
					"Хронический стресс",
					"Травма",
					"Гормональные нарушения",
					"Дефицит витаминов",
				},
				AssociatedSymptoms: []string{"апатия", "нарушения сна", "изменения аппетита"},
			},
			{
				ID:          "EMO_003",
				Name:        "Раздражительность",
				NameEn:      "Irritability",
				Category:    "emotional",
				Description: "Повышенная реакция на раздражители, гнев",
				RelatedSystems: []string{"nervous", "endocrine", "chakra"},
				RelatedChakras: []int{2, 4}, // Манипура, Анахата
				RelatedHormones: []string{"Cortisol", "Testosterone", "Thyroid"},
				Priority:      3,
				WhenToSeeDoctor: "При агрессии, нарушении отношений, физическом насилии",
				CommonCauses: []string{
					"Стресс и выгорание",
					"Гормональные колебания",
					"Недостаток сна",
					"Низкий сахар в крови",
				},
				AssociatedSymptoms: []string{"напряжение в теле", "учащённое дыхание", "бессонница"},
			},

			// === МЕНТАЛЬНЫЕ СИМПТОМЫ ===
			{
				ID:          "MENT_001",
				Name:        "Туман в голове",
				NameEn:      "Brain Fog",
				Category:    "mental",
				Description: "Нечёткость мышления, трудности с концентрацией",
				RelatedSystems: []string{"nervous", "endocrine", "chakra"},
				RelatedChakras: []int{5, 6}, // Вишудха, Аджна
				RelatedHormones: []string{"TSH", "Cortisol", "Estrogen", "Vitamin B12"},
				Priority:      2,
				WhenToSeeDoctor: "При внезапном появлении с нарушением речи или памяти",
				CommonCauses: []string{
					"Хронический стресс",
					"Недостаток сна",
					"Гормональные нарушения",
					"Дефицит витаминов группы B",
				},
				AssociatedSymptoms: []string{"усталость", "забывчивость", "трудности с фокусом"},
			},
			{
				ID:          "MENT_002",
				Name:        "Нарушения памяти",
				NameEn:      "Memory Issues",
				Category:    "mental",
				Description: "Трудности с запоминанием или вспоминанием информации",
				RelatedSystems: []string{"nervous", "chakra"},
				RelatedChakras: []int{5, 6}, // Вишудха, Аджна
				RelatedHormones: []string{"TSH", "Estrogen", "Cortisol", "Vitamin B12"},
				Priority:      2,
				WhenToSeeDoctor: "При внезапной потере памяти, дезориентации",
				CommonCauses: []string{
					"Стресс и переутомление",
					"Недостаток сна",
					"Возрастные изменения",
					"Дефицит витаминов",
				},
				AssociatedSymptoms: []string{"туман в голове", "усталость", "трудности с фокусом"},
			},

			// === ДУХОВНЫЕ СИМПТОМЫ ===
			{
				ID:          "SPIR_001",
				Name:        "Потеря смысла",
				NameEn:      "Loss of Meaning",
				Category:    "spiritual",
				Description: "Ощущение бессмысленности существования, экзистенциальный кризис",
				RelatedSystems: []string{"chakra", "psycho"},
				RelatedChakras: []int{6}, // Сахасрара
				RelatedHormones: []string{"Serotonin", "Dopamine", "Melatonin"},
				Priority:      2,
				WhenToSeeDoctor: "При суицидальных мыслях, полной апатии",
				CommonCauses: []string{
					"Экзистенциальный кризис",
					"Выгорание",
					"Потеря близких",
					"Депрессия",
				},
				AssociatedSymptoms: []string{"апатия", "депрессия", "изоляция"},
			},
			{
				ID:          "SPIR_002",
				Name:        "Чувство изоляции",
				NameEn:      "Feeling of Isolation",
				Category:    "spiritual",
				Description: "Ощущение отделённости от мира и других людей",
				RelatedSystems: []string{"chakra", "psycho"},
				RelatedChakras: []int{3, 6}, // Анахата, Сахасрара
				RelatedHormones: []string{"Oxytocin", "Serotonin"},
				Priority:      2,
				WhenToSeeDoctor: "При полной социальной изоляции, депрессии",
				CommonCauses: []string{
					"Травма привязанности",
					"Переезд или потеря социальных связей",
					"Депрессия",
					"Тревожное расстройство",
				},
				AssociatedSymptoms: []string{"одиночество", "тревожность", "депрессия"},
			},
		},
	}
}

// FindBySymptomName ищет симптомы по названию
func (db *SymptomDatabase) FindBySymptomName(query string) []Symptom {
	var results []Symptom
	for _, symptom := range db.Symptoms {
		if containsIgnoreCase(symptom.Name, query) ||
			containsIgnoreCase(symptom.NameEn, query) {
			results = append(results, symptom)
		}
	}
	return results
}

// FindByChakra ищет симптомы, связанные с чакрой
func (db *SymptomDatabase) FindByChakra(chakraIndex int) []Symptom {
	var results []Symptom
	for _, symptom := range db.Symptoms {
		for _, index := range symptom.RelatedChakras {
			if index == chakraIndex {
				results = append(results, symptom)
				break
			}
		}
	}
	return results
}

// FindByHormone ищет симптомы, связанные с гормоном
func (db *SymptomDatabase) FindByHormone(hormone string) []Symptom {
	var results []Symptom
	for _, symptom := range db.Symptoms {
		for _, h := range symptom.RelatedHormones {
			if containsIgnoreCase(h, hormone) {
				results = append(results, symptom)
				break
			}
		}
	}
	return results
}

// containsIgnoreCase проверяет наличие подстроки без учёта регистра
func containsIgnoreCase(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			sChar := s[i+j]
			subChar := substr[j]
			// Простая конверсия для ASCII
			if sChar >= 'A' && sChar <= 'Z' {
				sChar = sChar + ('a' - 'A')
			}
			if subChar >= 'A' && subChar <= 'Z' {
				subChar = subChar + ('a' - 'A')
			}
			if sChar != subChar {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
