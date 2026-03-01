package hypercube

import (
	"math"
	"time"
)

// HyperCoord — 4D координаты (X, Y, Z, W)
type HyperCoord struct {
	X int // День (сумма цифр)
	Y int // Месяц (сумма цифр)
	Z int // Год (сумма цифр последних 3)
	W int // 4-е измерение: Психовозраст или Состояние
}

// Person4D — человек в 4D-пространстве
type Person4D struct {
	ID           string
	Name         string
	BirthDate    time.Time
	BaseCoord    HyperCoord // Статические координаты (3D)
	CurrentW     int        // Динамическое W (психовозраст, состояние)
	Vectors4D    [4][3]int  // Векторы включая W
}

// CalcW_PsychAge рассчитывает W как психовозраст
func CalcW_PsychAge(birthDate time.Time, psychAge int) int {
	// W = паспортный возраст + (психовозраст - паспортный) * коэффициент
	// Упрощённо: W = psychAge (7, 35, 55)
	return psychAge
}

// CalcW_Cycle рассчитывает W как текущий цикл (лунный/эмоциональный)
func CalcW_Cycle(lastContact time.Time, now time.Time) int {
	days := int(now.Sub(lastContact).Hours() / 24)
	// Цикл 28 дней → W = 0..27
	return days % 28
}

// CalcW_State рассчитывает W как субъективное состояние (1-9)
func CalcW_State(moodScore int) int {
	// moodScore: -10..+10 → W: 1..9
	if moodScore < -5 {
		return 1 // низкое состояние
	}
	if moodScore > 5 {
		return 9 // высокое состояние
	}
	return 5 + moodScore/2 // линейная интерполяция
}

// Distance4D — расстояние в 4D-пространстве
func Distance4D(c1, c2 HyperCoord) float64 {
	dx := float64(c1.X - c2.X)
	dy := float64(c1.Y - c2.Y)
	dz := float64(c1.Z - c2.Z)
	dw := float64(c1.W - c2.W)
	return math.Sqrt(dx*dx + dy*dy + dz*dz + dw*dw)
}

// Compatibility4D — совместимость в 4D
func Compatibility4D(p1, p2 Person4D) float64 {
	// 3D расстояние
	dist3D := Distance3D(p1.BaseCoord, p2.BaseCoord)
	// W-разрыв
	wGap := math.Abs(float64(p1.CurrentW - p2.CurrentW))

	// Если W-разрыв > 44 — разные "временные слои" (как вы упоминали)
	if wGap > 44 {
		return 0.0 // Практически несовместимы в текущем состоянии
	}

	// Базовая совместимость: чем меньше расстояние — тем лучше
	score := 1.0 / (dist3D + wGap*0.5)

	// Бонус за совпадение по Y (эмоциональный резонанс)
	if p1.BaseCoord.Y == p2.BaseCoord.Y {
		score *= 1.3
	}

	// Штраф за противоположные векторы года
	if p1.Vectors4D[3][2]*p2.Vectors4D[3][2] < 0 {
		score *= 0.7
	}

	return score
}

// Distance3D — вспомогательная: 3D расстояние
func Distance3D(c1, c2 HyperCoord) float64 {
	dx := float64(c1.X - c2.X)
	dy := float64(c1.Y - c2.Y)
	dz := float64(c1.Z - c2.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// PhaseShift — определяет, в каком "временном слое" находится человек
func PhaseShift(w int) string {
	switch {
	case w < 18:
		return "Child/Trauma" // W < 18: регресс в травму
	case w < 35:
		return "Adolescent/Searching" // 18-34: поиск
	case w < 55:
		return "Adult/Responsibility" // 35-54: ответственность
	default:
		return "Elder/Wisdom" // 55+: мудрость или застой
	}
}

// WGap_Interpretation интерпретирует разрыв по W
func WGap_Interpretation(gap float64) string {
	switch {
	case gap < 10:
		return "Близкие состояния: возможен диалог"
	case gap < 25:
		return "Разные этапы: нужен мост понимания"
	case gap < 44:
		return "Значительный разрыв: риск непонимания"
	default:
		return "Разные временные слои: контакт затруднён"
	}
}
