package psycho

import "testing"

func TestDetectPsychAge(t *testing.T) {
	// Тест на детский возраст (страх, беспомощность)
	textChild := "мне страшно, я не хочу, мама не поймёт, тяжело"
	if DetectPsychAge(textChild) != AgeChild {
		t.Errorf("Expected AgeChild for fear text, got %d", DetectPsychAge(textChild))
	}

	// Тест на взрослый возраст (решимость, действие)
	textAdult := "я решу, я сделаю, у меня есть план, моя цель"
	if DetectPsychAge(textAdult) != AgeAdult {
		t.Errorf("Expected AgeAdult for decisive text, got %d", DetectPsychAge(textAdult))
	}
}
