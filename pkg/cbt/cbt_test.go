package cbt

import "testing"

func TestDetectDistortions(t *testing.T) {
	tests := []struct {
		text     string
		expected []CognitiveDistortion
	}{
		{
			text:     "Он всегда меня бросает, никогда не остаётся",
			expected: []CognitiveDistortion{DistortionAllOrNothing, DistortionOvergeneralization},
		},
		{
			text:     "Это катастрофа, всё пропало, я не вынесу",
			expected: []CognitiveDistortion{DistortionMagnification},
		},
		{
			text:     "Я должен быть идеальным, иначе я неудачник",
			expected: []CognitiveDistortion{DistortionShouldStatements, DistortionLabeling},
		},
		{
			text:     "Мне кажется, что он меня не любит, значит так и есть",
			expected: []CognitiveDistortion{DistortionEmotionalReasoning},
		},
	}

	for _, tc := range tests {
		found := DetectDistortions(tc.text)
		if len(found) != len(tc.expected) {
			t.Errorf("Text %q: expected %d distortions, got %d", tc.text, len(tc.expected), len(found))
		}
	}
}

func TestGenerateRationalResponse(t *testing.T) {
	distortions := []CognitiveDistortion{DistortionAllOrNothing, DistortionMagnification}
	response := GenerateRationalResponse("Всё пропало, он меня никогда не поймёт", distortions)
	if response == "" {
		t.Error("Expected non-empty rational response")
	}
}
