// Package chakra — энергетическая модель человека
//
// Контекст:
// Этот пакет реализует чакральную систему как ОДИН ИЗ СЛОЁВ
// целостной модели человека в проекте ideal-core.
//
// Философия:
// - Человек = многомерная система (энергия + физиология + психология + дух)
// - Чакры = энергетические центры, коррелирующие с физическими системами
// - Дисбаланс на одном уровне проявляется на других
//
// Архитектура:
// - chakra.go — базовые структуры (ЭТОТ ФАЙЛ)
// - chakra_system.go — системные функции (СУЩЕСТВУЮЩИЙ, НЕ УДАЛЯТЬ)
// - layers/ — слои (эндокринный, нервный, психологический)
// - correlations/ — связи между системами
// - diagnostics/ — диагностика по симптомам
// - recommendations/ — рекомендательная система
// - practices/ — практики для балансировки
// - temporal/ — временные соответствия
// - visual/ — визуализация (SVG, ASCII)
//
// Важно для разработчиков:
// - НЕ удалять chakra_system.go — обратная совместимость
// - Новые функции добавлять в соответствующие модули
// - Все изменения документировать в README.md
package chakra

// ChakraInfo — полная информация о чакре
//
// Концепция:
// Эта структура объединяет данные из всех слоёв:
// - Базовые (название, цвет, звук)
// - Физиологические (органы, железы, гормоны)
// - Психологические (Эриксон, Джудит, травмы)
// - Энергетические (лепестки, биджи, янтры)
type ChakraInfo struct {
	// === БАЗОВЫЕ ДАННЫЕ ===
	Name       string `json:"name"`        // Название на русском
	Sanskrit   string `json:"sanskrit"`    // Санскритское имя
	Location   string `json:"location"`    // Расположение в теле
	Color      string `json:"color"`       // Цвет
	HexColor   string `json:"hex_color"`   // HEX код цвета
	Sound      string `json:"sound"`       // Звук (биджа-мантра)
	Element    string `json:"element"`     // Элемент (таттва)

	// === ПСИХОЛОГИЧЕСКИЙ СЛОЙ ===
	Theme          string         `json:"theme"`           // Тема чакры
	Psychosoma     string         `json:"psychosoma"`      // Психосоматика
	Imbalance      ImbalanceModel `json:"imbalance"`       // Дисбаланс
	Affirmations   AffirmationSet `json:"affirmations"`    // Аффирмации

	// === РАЗВИТИЕ ЛИЧНОСТИ ===
	EriksonStage   EriksonStage   `json:"erikson_stage"`   // Стадия по Эриксону
	JudithInsights JudithModel    `json:"judith_insights"` // Модель Анодеи Джудит

	// === ФИЗИОЛОГИЧЕСКИЙ СЛОЙ ===
	Endocrine EndocrineCorrelation `json:"endocrine"` // Эндокринная система

	// === ТРАДИЦИОННАЯ СТРУКТУРА ===
	Traditional TraditionalStructure `json:"traditional"` // Лепестки, биджи, янтры

	// === ВИЗУАЛИЗАЦИЯ ===
	Visual VisualRepresentation `json:"visual"` // SVG, ASCII

	// === МЕТА-ДАННЫЕ ===
	Sources     []string `json:"sources"`      // Источники данных
	LastUpdated string   `json:"last_updated"` // Дата обновления
}

// ImbalanceModel — модель дисбаланса чакры
//
// Концепция:
// Дисбаланс рассматривается не как бинарное состояние,
// а как спектр с контекстом и динамикой.
type ImbalanceModel struct {
	Level       string   `json:"level"`        // "deficit" | "balanced" | "excess"
	Intensity   float64  `json:"intensity"`    // 0.0 - 1.0
	Context     []string `json:"context"`      // Сферы жизни
	Duration    string   `json:"duration"`     // "acute" | "chronic" | "cyclical"
	Triggers    []string `json:"triggers"`     // Триггеры
	Deficit     string   `json:"deficit"`      // Проявления дефицита
	Excess      string   `json:"excess"`       // Проявления избытка
}

// AffirmationSet — 7-уровневый набор аффирмаций
//
// Уровни:
// 1. Being (Я есть) — идентичность
// 2. Feeling (Я чувствую) — эмоции
// 3. Thinking (Я думаю) — убеждения
// 4. Doing (Я делаю) — действия
// 5. Allowing (Я позволяю) — принятие
// 6. Releasing (Я отпускаю) — освобождение
// 7. Trusting (Я доверяю) — интеграция
type AffirmationSet struct {
	Being     string `json:"being"`
	Feeling   string `json:"feeling"`
	Thinking  string `json:"thinking"`
	Doing     string `json:"doing"`
	Allowing  string `json:"allowing"`
	Releasing string `json:"releasing"`
	Trusting  string `json:"trusting"`
}

// EriksonStage — стадия психосоциального развития по Эриксону
type EriksonStage struct {
	Stage        int    `json:"stage"`         // 1-8
	AgeRange     string `json:"age_range"`     // Возрастной диапазон
	Crisis       string `json:"crisis"`        // Психосоциальный кризис
	Virtue       string `json:"virtue"`        // Формируемая добродетель
	ChakraIndex  int    `json:"chakra_index"`  // Индекс чакры (0-6)
	SignsOfBlock string `json:"signs_of_block"`
	HealingFocus string `json:"healing_focus"`
}

// JudithModel — психологическая модель по Анодее Джудит
type JudithModel struct {
	ChildhoodPattern  string   `json:"childhood_pattern"`
	CoreWound         string   `json:"core_wound"`
	DefenseMechanism  string   `json:"defense_mechanism"`
	AttachmentStyle   string   `json:"attachment_style"`
	HealingPath       []string `json:"healing_path"`
}

// EndocrineCorrelation — корреляция с эндокринной системой
type EndocrineCorrelation struct {
	Gland          string   `json:"gland"`
	Hormones       []string `json:"hormones"`
	Functions      []string `json:"functions"`
	ImbalanceSigns []string `json:"imbalance_signs"`
	SupportMethods []string `json:"support_methods"`
}

// PetalInfo — информация о лепестке лотоса чакры
type PetalInfo struct {
	SanskritLetter string `json:"sanskrit_letter"`
	BijaMantra     string `json:"bija_mantra"`
	Quality        string `json:"quality"`
	Color          string `json:"color"`
	ActivationNote string `json:"activation_note"`
}

// SubChakraInfo — микро-центр внутри чакры
type SubChakraInfo struct {
	Name         string      `json:"name"`
	Location     string      `json:"location"`
	Function     string      `json:"function"`
	Petals       []PetalInfo `json:"petals"`
	RelatedDeity string      `json:"related_deity"`
}

// TraditionalStructure — традиционная тантрическая структура
type TraditionalStructure struct {
	TotalPetals   int              `json:"total_petals"`
	Petals        []PetalInfo      `json:"petals"`
	SubChakras    []SubChakraInfo  `json:"sub_chakras"`
	CentralBija   string           `json:"central_bija"`
	Yantra        string           `json:"yantra"`
	CauseDeity    string           `json:"cause_deity"`
	ElementBija   string           `json:"element_bija"`
	ActivationSeq []string         `json:"activation_seq"`
}

// VisualRepresentation — визуальное отображение
type VisualRepresentation struct {
	ASCII   string                 `json:"ascii"`
	SVG     []byte                 `json:"svg"`
	Config  map[string]interface{} `json:"config"`
	Meta    VisualMeta             `json:"meta"`
}

// VisualMeta — метаданные визуализации
type VisualMeta struct {
	RecommendedMode    string `json:"recommended_mode"`
	MinRequirements    string `json:"min_requirements"`
	SVGWidth           int    `json:"svg_width"`
	SVGHeight          int    `json:"svg_height"`
	AnimationSupported bool   `json:"animation_supported"`
}
