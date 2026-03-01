package chakra

// ChakraSystem — состояние системы чакр человека
type ChakraSystem struct {
	PersonID    string
	Active      []ChakraIndex
	Blocked     []ChakraIndex
	LastUpdated string
}

// CalculateSystemState рассчитывает состояние системы по симптомам и векторам
func CalculateSystemState(personID string, symptoms []string, vectors [3][3]int) ChakraSystem {
	active := DetectActiveChakras(symptoms)
	
	// Эвристика: вектор года влияет на верхние чакры
	if vectors[2][2] < -5 { // токсичный сброс
		active = append(active, Ajna, Sahasrara) // гипер-активация или блок
	}
	
	return ChakraSystem{
		PersonID:    personID,
		Active:      active,
		Blocked:     findBlocked(active),
		LastUpdated: "now",
	}
}

// findBlocked находит заблокированные чакры (упрощённо)
func findBlocked(active []ChakraIndex) []ChakraIndex {
	all := []ChakraIndex{Muladhara, Svadhisthana, Manipura, Anahata, Vishuddha, Ajna, Sahasrara}
	var blocked []ChakraIndex
	for _, idx := range all {
		if !containsChakra(active, idx) {
			blocked = append(blocked, idx)
		}
	}
	return blocked
}

func containsChakra(list []ChakraIndex, item ChakraIndex) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

// GetRecommendations возвращает практики для активных чакр
func (cs ChakraSystem) GetRecommendations() []string {
	var recs []string
	for _, idx := range cs.Active {
		info := GetChakraInfo(idx)
		recs = append(recs, info.Practice)
	}
	return recs
}
