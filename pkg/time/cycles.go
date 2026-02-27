package time

import "time"

type CyclePhase string

const (
	PhaseThaw       CyclePhase = "THAW"
	PhasePeak       CyclePhase = "PEAK"
	PhaseFriction   CyclePhase = "FRICTION"
	PhaseBreak      CyclePhase = "BREAK"
	PhaseQuarantine CyclePhase = "QUARANTINE"
)

func DetectCyclePhase(lastContact time.Time, breakupDate time.Time) CyclePhase {
	now := time.Now()
	
	if !breakupDate.IsZero() {
		days := int(now.Sub(breakupDate).Hours() / 24)
		if days <= 36 {
			return PhaseQuarantine
		}
	}
	
	days := int(now.Sub(lastContact).Hours() / 24)
	cycleDay := days % 12
	
	if cycleDay <= 3 {
		return PhaseThaw
	} else if cycleDay <= 7 {
		return PhasePeak
	} else if cycleDay <= 10 {
		return PhaseFriction
	}
	return PhaseBreak
}

func GetRecommendation(phase CyclePhase) string {
	switch phase {
	case PhaseThaw:
		return "Ответить тепло, но кратко"
	case PhasePeak:
		return "Быть в моменте, не давить"
	case PhaseFriction:
		return "Заземлить, коротко про дела"
	case PhaseBreak, PhaseQuarantine:
		return "Переждать, не реагировать на провокации"
	default:
		return "Наблюдать"
	}
}
