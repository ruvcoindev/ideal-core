package chakra

// ChakraInfo — полная информация о чакре
type ChakraInfo struct {
	Name        string   // Название
	Sanskrit    string   // Санскритское имя
	Location    string   // Расположение
	Color       string   // Цвет
	HexColor    string   // HEX код цвета
	Sound       string   // Звук (биджа-мантра)
	Element     string   // Элемент
	Theme       string   // Тема
	Psychosoma  string   // Психосоматика
	Imbalance   string   // Дисбаланс
	Affirmation string   // Аффирмация
}

// GetFullChakraSystem возвращает полную информацию о 7 чакрах
func GetFullChakraSystem() []ChakraInfo {
	return []ChakraInfo{
		{
			Name:        "Муладхара",
			Sanskrit:    "Muladhara",
			Location:    "Основание позвоночника, промежность",
			Color:       "Красный",
			HexColor:    "#FF0000",
			Sound:       "ЛАМ (LAM)",
			Element:     "Земля",
			Theme:       "Выживание, безопасность, право на существование, базовые потребности",
			Psychosoma:  "Проблемы с ногами, позвоночником, прямой кишкой, иммунитетом",
			Imbalance:   "Страх, тревога, ощущение «меня не поддерживают», финансовые проблемы",
			Affirmation: "Я в безопасности. Я имею право быть здесь. Земля поддерживает меня.",
		},
		{
			Name:        "Свадхистхана",
			Sanskrit:    "Svadhisthana",
			Location:    "Низ живота (2 пальца ниже пупка)",
			Color:       "Оранжевый",
			HexColor:    "#FFA500",
			Sound:       "ВАМ (VAM)",
			Element:     "Вода",
			Theme:       "Удовольствие, сексуальность, творчество, эмоции, желание",
			Psychosoma:  "Проблемы с почками, мочевым пузырем, репродуктивной системой",
			Imbalance:   "Чувство вины за удовольствие, подавление желаний, сексуальные проблемы",
			Affirmation: "Я имею право хотеть. Мои желания священны. Я наслаждаюсь жизнью.",
		},
		{
			Name:        "Манипура",
			Sanskrit:    "Manipura",
			Location:    "Область пупка, солнечное сплетение",
			Color:       "Желтый",
			HexColor:    "#FFD700",
			Sound:       "РАМ (RAM)",
			Element:     "Огонь",
			Theme:       "Личная сила, воля, уверенность, контроль, самооценка",
			Psychosoma:  "Проблемы с желудком, печенью, поджелудочной, пищеварением",
			Imbalance:   "Низкая самооценка, потребность контролировать, гнев, стыд",
			Affirmation: "Я силен. Я имею право на свою силу. Я управляю своей жизнью.",
		},
		{
			Name:        "Анахата",
			Sanskrit:    "Anahata",
			Location:    "Центр груди",
			Color:       "Зеленый или Розовый",
			HexColor:    "#00FF00",
			Sound:       "ЯМ (YAM)",
			Element:     "Воздух",
			Theme:       "Любовь, сострадание, прощение, принятие, связь",
			Psychosoma:  "Проблемы с сердцем, легкими, грудной клеткой, руками",
			Imbalance:   "Обиды, закрытое сердце, страх близости, ревность, одиночество",
			Affirmation: "Я люблю и любим. Мое сердце открыто. Я прощаю себя и других.",
		},
		{
			Name:        "Вишудха",
			Sanskrit:    "Vishuddha",
			Location:    "Горло, щитовидная железа",
			Color:       "Голубой",
			HexColor:    "#00BFFF",
			Sound:       "ХАМ (HAM)",
			Element:     "Эфир (пространство)",
			Theme:       "Самовыражение, правда, коммуникация, творчество через слово",
			Psychosoma:  "Проблемы с горлом, щитовидкой, шеей, зубами, ушами",
			Imbalance:   "Страх говорить, ложь, подавленный голос, «ком в горле»",
			Affirmation: "Я говорю свою правду. Мой голос важен. Я выражаю себя свободно.",
		},
		{
			Name:        "Аджна",
			Sanskrit:    "Ajna",
			Location:    "Между бровями",
			Color:       "Индиго",
			HexColor:    "#4B0082",
			Sound:       "ОМ (OM) или ШАН (SHAM)",
			Element:     "Свет",
			Theme:       "Интуиция, мудрость, видение, воображение, ясность ума",
			Psychosoma:  "Проблемы с глазами, головные боли, мигрени, синусы",
			Imbalance:   "Отсутствие видения, иллюзии, чрезмерный анализ, паралич решения",
			Affirmation: "Я вижу ясно. Я доверяю своей интуиции. Мое видение чисто.",
		},
		{
			Name:        "Сахасрара",
			Sanskrit:    "Sahasrara",
			Location:    "Макушка головы",
			Color:       "Фиолетовый или Белый, Золотой",
			HexColor:    "#8000FF",
			Sound:       "Беззвучный звук (тишина) или ОМ",
			Element:     "Космическое сознание",
			Theme:       "Духовность, единство, связь с Богом/Вселенной, просветление",
			Psychosoma:  "Депрессия, апатия, рассеянность, проблемы с мозгом",
			Imbalance:   "Отрыв от реальности, духовная гордыня, отсутствие смысла",
			Affirmation: "Я един со Вселенной. Я открыт божественному потоку. Я есть.",
		},
	}
}

// GetChakraByIndex возвращает чакру по индексу (0-6)
func GetChakraByIndex(index int) ChakraInfo {
	chakras := GetFullChakraSystem()
	if index >= 0 && index < len(chakras) {
		return chakras[index]
	}
	return ChakraInfo{}
}

// GetChakraByName возвращает чакру по названию
func GetChakraByName(name string) ChakraInfo {
	chakras := GetFullChakraSystem()
	for _, chakra := range chakras {
		if chakra.Name == name || chakra.Sanskrit == name {
			return chakra
		}
	}
	return ChakraInfo{}
}
