package psycho

import "ideal-core/pkg/db"

// Affirmation представляет утверждение для работы с подсознанием
type Affirmation struct {
	ID          string
	Text        string   // "Я люблю и одобряю себя"
	Author      string   // "Louise Hay", "Zhikarentsev"
	Language    string   // "ru", "en"
	Keywords    []string // ["self-love", "approval", "healing"]
	ChakraIndex int      // 0-6: на какую чакру направлена
}

// AffirmationDB — хранилище аффирмаций
type AffirmationDB struct {
	affirmations map[string]Affirmation
	medicalMap   map[string]string // "alopecia" -> "baldness"
}

// NewAffirmationDB создаёт базу с предустановленными данными
func NewAffirmationDB() *AffirmationDB {
	db := &AffirmationDB{
		affirmations: make(map[string]Affirmation),
		medicalMap:   make(map[string]string),
	}
	db.loadLouiseHay()
	db.loadZhikarentsev()
	db.loadMedicalMapping()
	return db
}

// loadMedicalMapping загружает соответствие мед.терминов и народных названий
func (db *AffirmationDB) loadMedicalMapping() {
	db.medicalMap["alopecia"] = "облысение"
	db.medicalMap["цефалгия"] = "головная боль"
	db.medicalMap["гастрит"] = "воспаление желудка"
	db.medicalMap["бессонница"] = "нарушение сна"
	// Добавить по мере необходимости
}

// loadLouiseHay загружает аффирмации Луизы Хей
func (db *AffirmationDB) loadLouiseHay() {
	db.affirmations["lh_001"] = Affirmation{
		ID:          "lh_001",
		Text:        "Я люблю и одобряю себя",
		Author:      "Louise Hay",
		Language:    "ru",
		Keywords:    []string{"self-love", "approval"},
		ChakraIndex: 4, // Heart
	}
	db.affirmations["lh_002"] = Affirmation{
		ID:          "lh_002",
		Text:        "Я в безопасности, я доверяю жизни",
		Author:      "Louise Hay",
		Language:    "ru",
		Keywords:    []string{"safety", "trust"},
		ChakraIndex: 1, // Root
	}
}

// loadZhikarentsev загружает аффирмации Жикаренцева
func (db *AffirmationDB) loadZhikarentsev() {
	db.affirmations["zh_001"] = Affirmation{
		ID:          "zh_001",
		Text:        "Я излучаю любовь, и мир отвечает мне любовью",
		Author:      "Zhikarentsev",
		Language:    "ru",
		Keywords:    []string{"love", "resonance"},
		ChakraIndex: 4,
	}
}

// FindByMedicalTerm находит аффирмацию по медицинскому термину
func (db *AffirmationDB) FindByMedicalTerm(term string) []Affirmation {
	// Нормализуем термин
	normalized := db.medicalMap[term]
	if normalized == "" {
		normalized = term
	}
	
	var results []Affirmation
	for _, aff := range db.affirmations {
		for _, kw := range aff.Keywords {
			if kw == normalized || contains(aff.Text, normalized) {
				results = append(results, aff)
			}
		}
	}
	return results
}

// GetPersonalAffirmation подбирает аффирмацию под профиль человека
func (db *AffirmationDB) GetPersonalAffirmation(p *db.Person) Affirmation {
	// Приоритет: заблокированная чакра → лунная фаза → кубический вектор
	for i, ch := range p.Chakras {
		if ch.Blocked {
			return db.getByChakra(i)
		}
	}
	return db.getByLunarPhase(p.Lunar.CurrentPhase)
}

func (db *AffirmationDB) getByChakra(index int) Affirmation {
	for _, aff := range db.affirmations {
		if aff.ChakraIndex == index {
			return aff
		}
	}
	return db.affirmations["lh_001"] // fallback
}

func (db *AffirmationDB) getByLunarPhase(phase int) Affirmation {
	// Упрощённая логика: фаза 0-10 → Root, 11-20 → Heart, 21-29 → Crown
	chakra := 0
	if phase > 20 {
		chakra = 6
	} else if phase > 10 {
		chakra = 4
	}
	return db.getByChakra(chakra)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)))
}
