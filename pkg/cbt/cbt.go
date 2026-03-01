package cbt

import "strings"

// CognitiveDistortion — тип когнитивного искажения
type CognitiveDistortion string

const (
	DistortionAllOrNothing      CognitiveDistortion = "all_or_nothing"      // Чёрно-белое мышление
	DistortionOvergeneralization CognitiveDistortion = "overgeneralization" // Сверхобобщение
	DistortionMentalFilter      CognitiveDistortion = "mental_filter"       // Ментальный фильтр
	DistortionDisqualifying     CognitiveDistortion = "disqualifying"       // Обесценивание позитива
	DistortionJumpingConclusions CognitiveDistortion = "jumping_conclusions" // Поспешные выводы
	DistortionMagnification     CognitiveDistortion = "magnification"       // Катастрофизация
	DistortionEmotionalReasoning CognitiveDistortion = "emotional_reasoning" // Эмоциональное обоснование
	DistortionShouldStatements  CognitiveDistortion = "should_statements"   // Долженствования
	DistortionLabeling          CognitiveDistortion = "labeling"            // Навешивание ярлыков
	DistortionPersonalization   CognitiveDistortion = "personalization"     // Персонализация
)

// ThoughtRecord — запись автоматической мысли (CBT)
type ThoughtRecord struct {
	Situation        string                `json:"situation"`
	AutomaticThought string                `json:"automatic_thought"`
	Emotions         []string              `json:"emotions"` // "тревога", "гнев", "вина"
	Intensity        int                   `json:"intensity"` // 0-100
	Distortions      []CognitiveDistortion `json:"distortions"`
	RationalResponse string                `json:"rational_response"`
	NewIntensity     int                   `json:"new_intensity"`
}

// DetectDistortions ищет когнитивные искажения в тексте
func DetectDistortions(text string) []CognitiveDistortion {
	var found []CognitiveDistortion
	textLower := strings.ToLower(text)

	// All-or-nothing
	if containsAny(textLower, []string{"всегда", "никогда", "все", "никто", "полностью", "совсем"}) {
		found = append(found, DistortionAllOrNothing)
	}

	// Overgeneralization
	if containsAny(textLower, []string{"опять", "снова", "постоянно", "каждый раз", "всё время"}) {
		found = append(found, DistortionOvergeneralization)
	}

	// Catastrophizing
	if containsAny(textLower, []string{"ужас", "катастрофа", "конец", "никогда не", "всё пропало", "невыносимо"}) {
		found = append(found, DistortionMagnification)
	}

	// Should statements
	if containsAny(textLower, []string{"должен", "должна", "надо", "следует", "обязан", "обязана"}) {
		found = append(found, DistortionShouldStatements)
	}

	// Emotional reasoning
	if containsAny(textLower, []string{"чувствую, что", "мне кажется, что", "я знаю, что", "очевидно"}) {
		found = append(found, DistortionEmotionalReasoning)
	}

	// Labeling
	if containsAny(textLower, []string{"я неудачник", "я слабак", "он манипулятор", "она истеричка", "все такие"}) {
		found = append(found, DistortionLabeling)
	}

	return found
}

// GenerateRationalResponse помогает сформировать рациональный ответ
func GenerateRationalResponse(automaticThought string, distortions []CognitiveDistortion) string {
	response := "Альтернативный взгляд: "

	if contains(distortions, DistortionAllOrNothing) {
		response += "Возможно, есть оттенки между «всё» и «ничего». "
	}
	if contains(distortions, DistortionMagnification) {
		response += "Что самое реалистичное может произойти? "
	}
	if contains(distortions, DistortionShouldStatements) {
		response += "Что будет, если заменить «должен» на «хочу» или «предпочитаю»? "
	}
	if contains(distortions, DistortionLabeling) {
		response += "Это поведение или вся личность? Можно ли отделить поступок от человека? "
	}

	return response
}

// Helper functions
func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func contains(slice []CognitiveDistortion, item CognitiveDistortion) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// CBTSession — сессия работы с мыслями
type CBTSession struct {
	PersonID     string         `json:"person_id"`
	Records      []ThoughtRecord `json:"records"`
	CreatedAt    string         `json:"created_at"`
	Insights     []string       `json:"insights,omitempty"`
}

// AnalyzeSession возвращает инсайты по сессии
func AnalyzeSession(session CBTSession) []string {
	var insights []string

	// Подсчитываем частоту искажений
	distortionCount := make(map[CognitiveDistortion]int)
	for _, rec := range session.Records {
		for _, d := range rec.Distortions {
			distortionCount[d]++
		}
	}

	// Формируем инсайты
	for d, count := range distortionCount {
		if count >= 3 {
			insights = append(insights, 
				"Частое искажение: "+string(d)+" ("+string(rune(count))+") — стоит проработать")
		}
	}

	// Проверяем прогресс по интенсивности эмоций
	if len(session.Records) >= 2 {
		first := session.Records[0].Intensity
		last := session.Records[len(session.Records)-1].NewIntensity
		if last < first {
			insights = append(insights, "Прогресс: интенсивность эмоций снизилась с "+string(rune(first))+" до "+string(rune(last)))
		}
	}

	return insights
}
