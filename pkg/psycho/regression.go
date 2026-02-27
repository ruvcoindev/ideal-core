package psycho

import "strings"

type PsychAge int

const (
	AgeChild PsychAge = 7
	AgeTeen  PsychAge = 14
	AgeAdult PsychAge = 35
	AgeElder PsychAge = 55
)

func DetectPsychAge(text string) PsychAge {
	textLower := strings.ToLower(text)
	
	// Детские маркеры
	childMarkers := []string{"страшно", "не хочу", "мама", "тяжело", "боюсь", "помогите"}
	for _, marker := range childMarkers {
		if strings.Contains(textLower, marker) {
			return AgeChild
		}
	}
	
	// Взрослые маркеры
	adultMarkers := []string{"решу", "сделаю", "план", "цель", "ответственность", "выберу"}
	for _, marker := range adultMarkers {
		if strings.Contains(textLower, marker) {
			return AgeAdult
		}
	}
	
	return AgeElder
}

func CompatibilityWithRegression(baseScore float64, age1, age2 PsychAge) float64 {
	if (age1 == AgeAdult && age2 == AgeChild) || (age1 == AgeChild && age2 == AgeAdult) {
		return baseScore * 0.5
	}
	if age1 == AgeChild && age2 == AgeChild {
		return baseScore * 0.3
	}
	return baseScore
}
