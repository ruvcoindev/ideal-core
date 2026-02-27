package neural

import "strings"

// BodyMarker — категория телесного маркера
type BodyMarker string

const (
	MarkerIntimacy  BodyMarker = "INTIMACY"   // "пися", "оргазм"
	MarkerCleansing BodyMarker = "CLEANSING"  // "клизма", "прокакалась"
	MarkerWater     BodyMarker = "WATER"      // "купалась", "душ"
	MarkerPain      BodyMarker = "PAIN"       // "болит", "слабость"
)

// BodyMarkers — словарь маркеров
var BodyMarkers = map[BodyMarker][]string{
	MarkerIntimacy:  {"пися", "писю", "кончаю", "оргазм", "куни"},
	MarkerCleansing: {"клизма", "прокакалась", "какала", "свеча"},
	MarkerWater:     {"купалась", "душ", "ванна", "соли"},
	MarkerPain:      {"болит", "слабость", "трясет", "плита"},
}

// DetectBodyMarker ищет маркеры в тексте
func DetectBodyMarker(text string) (BodyMarker, bool) {
	textLower := strings.ToLower(text)
	for marker, keywords := range BodyMarkers {
		for _, kw := range keywords {
			if strings.Contains(textLower, kw) {
				return marker, true
			}
		}
	}
	return "", false
}

// IntimacyLevel оценивает уровень доверия по маркеру
func IntimacyLevel(marker BodyMarker) int {
	switch marker {
	case MarkerIntimacy:
		return 9
	case MarkerCleansing:
		return 7
	case MarkerWater:
		return 5
	case MarkerPain:
		return 6
	default:
		return 0
	}
}
