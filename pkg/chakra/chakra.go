package chakra

// ChakraIndex — индекс чакры (0-6)
type ChakraIndex int

const (
	Muladhara ChakraIndex = iota // 0: Корень
	Svadhisthana                 // 1: Сакральный
	Manipura                     // 2: Солнечное сплетение
	Anahata                      // 3: Сердечный
	Vishuddha                    // 4: Горловой
	Ajna                         // 5: Третий глаз
	Sahasrara                    // 6: Коронный
)

// ImbalanceModel — модель дисбаланса
type ImbalanceModel struct {
	Emotions    []string `json:"emotions"`
	Behaviors   []string `json:"behaviors"`
	Physical    []string `json:"physical"`
}

// ChakraInfo — полная информация о чакре
type ChakraInfo struct {
	Index       ChakraIndex    `json:"index"`
	Name        string         `json:"name"`
	Element     string         `json:"element"`
	Color       string         `json:"color"`
	Location    string         `json:"location"`
	Imbalance   ImbalanceModel `json:"imbalance"`
	Affirmation string         `json:"affirmation"`
	Practice    string         `json:"practice"`
}

// GetChakraInfo возвращает информацию о чакре по индексу
func GetChakraInfo(idx ChakraIndex) ChakraInfo {
	chakras := map[ChakraIndex]ChakraInfo{
		Muladhara: {
			Index:     Muladhara,
			Name:      "Муладхара",
			Element:   "Земля",
			Color:     "Красный",
			Location:  "Основание позвоночника",
			Imbalance: ImbalanceModel{
				Emotions:  []string{"Страх", "Тревога", "Небезопасность"},
				Behaviors: []string{"Контроль", "Жадность", "Изоляция"},
				Physical:  []string{"Проблемы с ногами, кишечником, иммунитетом"},
			},
			Affirmation: "Я в безопасности. Земля поддерживает меня.",
			Practice:    "Гвоздестояние, заземление, работа с корнями",
		},
		Svadhisthana: {
			Index:     Svadhisthana,
			Name:      "Свадхистана",
			Element:   "Вода",
			Color:     "Оранжевый",
			Location:  "Низ живота",
			Imbalance: ImbalanceModel{
				Emotions:  []string{"Вина", "Стыд", "Подавленные желания"},
				Behaviors: []string{"Зависимости", "Эмоциональные качели"},
				Physical:  []string{"Проблемы с мочеполовой системой, поясницей"},
			},
			Affirmation: "Я позволяю себе чувствовать. Мои желания важны.",
			Practice:    "Танец, вода, творчество, принятие эмоций",
		},
		Manipura: {
			Index:     Manipura,
			Name:      "Манипура",
			Element:   "Огонь",
			Color:     "Жёлтый",
			Location:  "Солнечное сплетение",
			Imbalance: ImbalanceModel{
				Emotions:  []string{"Гнев", "Стыд", "Низкая самооценка"},
				Behaviors: []string{"Контроль других", "Перфекционизм"},
				Physical:  []string{"Проблемы с ЖКТ, диабет, усталость"},
			},
			Affirmation: "Я достоин. Моя сила — в моём выборе.",
			Practice:    "Дыхание огня, солнечные ванны, утверждение границ",
		},
		Anahata: {
			Index:     Anahata,
			Name:      "Анахата",
			Element:   "Воздух",
			Color:     "Зелёный",
			Location:  "Центр груди",
			Imbalance: ImbalanceModel{
				Emotions:  []string{"Обида", "Ревность", "Одиночество"},
				Behaviors: []string{"Закрытость", "Жертвенность"},
				Physical:  []string{"Проблемы с сердцем, лёгкими, грудным отделом"},
			},
			Affirmation: "Я люблю и принимаю себя. Моё сердце открыто.",
			Practice:    "Практика благодарности, объятия, работа с прощением",
		},
		Vishuddha: {
			Index:     Vishuddha,
			Name:      "Вишудха",
			Element:   "Эфир",
			Color:     "Голубой",
			Location:  "Горло",
			Imbalance: ImbalanceModel{
				Emotions:  []string{"Страх говорить", "Ложь", "Подавленный голос"},
				Behaviors: []string{"Молчание", "Сплетни", "Неумение сказать «нет»"},
				Physical:  []string{"Проблемы с горлом, щитовидкой, шеей"},
			},
			Affirmation: "Мой голос важен. Я говорю свою правду с любовью.",
			Practice:    "Пение, мантры, ведение дневника, честный диалог",
		},
		Ajna: {
			Index:     Ajna,
			Name:      "Аджна",
			Element:   "Свет",
			Color:     "Индиго",
			Location:  "Между бровями",
			Imbalance: ImbalanceModel{
				Emotions:  []string{"Спутанность", "Иллюзии", "Отрицание интуиции"},
				Behaviors: []string{"Гиперрациональность", "Игнорирование знаков"},
				Physical:  []string{"Головные боли, проблемы со зрением, бессонница"},
			},
			Affirmation: "Я доверяю своей интуиции. Вижу ясно.",
			Practice:    "Медитация на третий глаз, работа со снами, визуализация",
		},
		Sahasrara: {
			Index:     Sahasrara,
			Name:      "Сахасрара",
			Element:   "Сознание",
			Color:     "Фиолетовый/Белый",
			Location:  "Макушка",
			Imbalance: ImbalanceModel{
				Emotions:  []string{"Отчуждение", "Цинизм", "Потеря смысла"},
				Behaviors: []string{"Духовный материализм", "Отказ от практики"},
				Physical:  []string{"Неврологические проблемы, депрессия"},
			},
			Affirmation: "Я един со Вселенной. Мой путь имеет смысл.",
			Practice:    "Молчание, созерцание, служение, интеграция опыта",
		},
	}
	return chakras[idx]
}

// DetectActiveChakras определяет активные чакры по симптомам
func DetectActiveChakras(symptoms []string) []ChakraIndex {
	active := make(map[ChakraIndex]bool)
	
	for _, s := range symptoms {
		switch {
		case containsAny(s, []string{"страх", "тревога", "ноги", "кишечник"}):
			active[Muladhara] = true
		case containsAny(s, []string{"вина", "желания", "мочеполовая", "поясница"}):
			active[Svadhisthana] = true
		case containsAny(s, []string{"гнев", "самооценка", "жкт", "контроль"}):
			active[Manipura] = true
		case containsAny(s, []string{"обида", "сердце", "лёгкие", "одиночество"}):
			active[Anahata] = true
		case containsAny(s, []string{"горло", "голос", "щитовидка", "ложь"}):
			active[Vishuddha] = true
		case containsAny(s, []string{"интуиция", "головная боль", "зрение", "сны"}):
			active[Ajna] = true
		case containsAny(s, []string{"смысл", "отчуждение", "депрессия", "сознание"}):
			active[Sahasrara] = true
		}
	}
	
	var result []ChakraIndex
	for idx := range active {
		result = append(result, idx)
	}
	return result
}

// Helper
func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if contains(s, sub) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)
}
