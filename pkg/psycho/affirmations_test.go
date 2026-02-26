package psycho

import "testing"

func TestMedicalMapping(t *testing.T) {
	db := NewAffirmationDB()
	
	// Проверка маппинга
	if db.medicalMap["alopecia"] != "облысение" {
		t.Error("Medical mapping failed")
	}
}

func TestFindByMedicalTerm(t *testing.T) {
	db := NewAffirmationDB()
	
	// Поиск по ключевому слову
	results := db.FindByMedicalTerm("self-love")
	if len(results) == 0 {
		t.Error("Expected to find affirmations for 'self-love'")
	}
}
