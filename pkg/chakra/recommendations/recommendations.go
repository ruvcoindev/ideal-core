package recommendations

import (
	"ideal-core/pkg/chakra"
	"ideal-core/pkg/chakra/layers"
)

// GenerateFullRecommendation генерирует полные рекомендации для человека
func GenerateFullRecommendation(personID string, symptoms []string, activeChakras []chakra.ChakraIndex) Recommendation {
	rec := Recommendation{
		PersonID:     personID,
		Symptoms:     symptoms,
		ActiveChakras: activeChakras,
	}
	
	// 1. Практики по чакрам
	for _, idx := range activeChakras {
		info := chakra.GetChakraInfo(idx)
		rec.Practices = append(rec.Practices, Practice{
			Chakra:      info.Name,
			Title:       "Практика для " + info.Name,
			Description: info.Practice,
			Frequency:   "daily",
		})
	}
	
	// 2. Лабораторные тесты
	engine := NewLabTestEngine()
	rec.LabTests = engine.RecommendTests(symptoms, activeChakras)
	
	// 3. Эндокринные слои
	rec.EndocrineLayers = layers.GetEndocrineLayers(activeChakras)
	
	// 4. Аффирмации
	for _, idx := range activeChakras {
		info := chakra.GetChakraInfo(idx)
		rec.Affirmations = append(rec.Affirmations, info.Affirmation)
	}
	
	return rec
}

// Recommendation — полные рекомендации
type Recommendation struct {
	PersonID        string
	Symptoms        []string
	ActiveChakras   []chakra.ChakraIndex
	Practices       []Practice
	LabTests        []LabTest
	EndocrineLayers []layers.EndocrineLayer
	Affirmations    []string
}

// Practice — рекомендация по практике
type Practice struct {
	Chakra      string
	Title       string
	Description string
	Frequency   string // "daily", "weekly"
}
