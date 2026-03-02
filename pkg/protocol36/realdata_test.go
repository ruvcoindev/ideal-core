package protocol36

import (
	"testing"
	"time"
)

func TestProtocol36_RealData(t *testing.T) {
	now := time.Date(2026, 3, 1, 12, 44, 0, 0, time.UTC)
	
	// Дина: блокировка 01.03.2026 12:44
	dinaLastContact := time.Date(2026, 3, 1, 12, 44, 0, 0, time.UTC)
	config, days := CalculatePhase(dinaLastContact, now)
	
	if config.Name != "Preparation" || days != 0 {
		t.Logf("Дина: только что заблокировала — фаза %s, день %d", config.Name, days)
	}
	
	// Елена: последний контакт 22.11.2025
	elenaLastContact := time.Date(2025, 11, 22, 0, 0, 0, 0, time.UTC)
	config2, days2 := CalculatePhase(elenaLastContact, now)
	
	if days2 < 90 {
		t.Errorf("Елена: ожидали ~100 дней, получили %d", days2)
	}
	if config2.Name != PhaseIntegration {
		t.Logf("Елена: фаза %s — возможна реинтеграция", config2.Name)
	}
	
	// Валя: гипотетический последний контакт 10.02.2026
	valyaLastContact := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
	config3, days3 := CalculatePhase(valyaLastContact, now)
	
	if days3 < 19 {
		t.Errorf("Валя: ожидали ~19 дней, получили %d", days3)
	}
	if config3.Name != PhaseRewire {
		t.Logf("Валя: фаза %s — период перепрошивки", config3.Name)
	}
}

func TestProtocol36_TasksGeneration(t *testing.T) {
	lastContact := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)
	now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	
	config, days := CalculatePhase(lastContact, now)
	
	if days != 10 {
		t.Errorf("Expected day 10, got %d", days)
	}
	if config.Name != PhaseRewire {
		t.Errorf("Expected PhaseRewire, got %s", config.Name)
	}
	if len(config.DailyTasks) == 0 {
		t.Error("Expected non-empty DailyTasks for Rewire phase")
	}
}
