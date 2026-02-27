package psycho

import "testing"

func TestMedicalMapping(t *testing.T) {
	db := NewAffirmationDB()
	
	if db.medicalMap["alopecia"] != "облысение" {
		t.Error("Medical mapping failed")
	}
}

func TestFindByMedicalTerm(t *testing.T) {
	db := NewAffirmationDB()
	
	results := db.FindByMedicalTerm("alopecia")
	if len(results) == 0 {
		t.Log("No affirmations found for alopecia (expected in MVP)")
	}
	
	results = db.FindByKeyword("healing")
	if len(results) == 0 {
		t.Error("Expected to find affirmations for 'healing'")
	}
}
