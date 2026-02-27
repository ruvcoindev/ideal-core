package psycho

import "strings"

// Affirmation — утверждение для работы с подсознанием
type Affirmation struct {
	ID          string
	Text        string
	Author      string
	Language    string
	Keywords    []string
	ChakraIndex int
}

// AffirmationDB — хранилище
type AffirmationDB struct {
	affirmations map[string]Affirmation
	medicalMap   map[string]string
}

// NewAffirmationDB создаёт базу
func NewAffirmationDB() *AffirmationDB {
	db := &AffirmationDB{
		affirmations: make(map[string]Affirmation),
		medicalMap:   make(map[string]string),
	}
	db.loadLouiseHay()
	db.loadMedicalMapping()
	return db
}

func (db *AffirmationDB) loadMedicalMapping() {
	db.medicalMap["alopecia"] = "облысение"
	db.medicalMap["цефалгия"] = "головная боль"
	db.medicalMap["гастрит"] = "воспаление желудка"
	db.medicalMap["бессонница"] = "нарушение сна"
}

func (db *AffirmationDB) loadLouiseHay() {
	db.affirmations["lh_001"] = Affirmation{
		ID:          "lh_001",
		Text:        "Я люблю и одобряю себя",
		Author:      "Louise Hay",
		Language:    "ru",
		Keywords:    []string{"self-love", "approval", "healing"},
		ChakraIndex: 4,
	}
}

// FindByKeyword ищет аффирмации по ключевому слову
func (db *AffirmationDB) FindByKeyword(keyword string) []Affirmation {
	var results []Affirmation
	for _, aff := range db.affirmations {
		for _, kw := range aff.Keywords {
			if strings.EqualFold(kw, keyword) {
				results = append(results, aff)
			}
		}
	}
	return results
}

// FindByMedicalTerm ищет аффирмацию по медицинскому термину
func (db *AffirmationDB) FindByMedicalTerm(term string) []Affirmation {
	normalized := strings.ToLower(term)
	if val, ok := db.medicalMap[normalized]; ok {
		normalized = val
	}
	return db.FindByKeyword(normalized)
}

// GetMedicalTerm переводит медицинский термин в народное название
func (db *AffirmationDB) GetMedicalTerm(term string) string {
	if val, ok := db.medicalMap[strings.ToLower(term)]; ok {
		return val
	}
	return term
}
