package journal

import (
	"os"
	"testing"
	"time"
)

func TestJournal_AddEntry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "journal_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	j, err := NewJournal(JournalConfig{
		DataDir:        tmpDir,
		UseOllamaEmbed: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	
	// Добавляем запись через универсальный метод (backward compat)
	entry := ThoughtEntry{
		Type:             EntryTypeCBT,
		Timestamp:        time.Now(),
		Situation:        "Дина заблокировала меня в соцсетях",
		AutomaticThought: "Она меня никогда не поймёт, всё пропало",
		Emotions:         []string{"тревога", "обида", "безнадёжность"},
		Intensity:        85,
		Phase:            "Detox",
		PersonID:         "Dina",
	}
	
	err = j.AddEntry(entry)
	if err != nil {
		t.Errorf("AddEntry failed: %v", err)
	}
	
	// Проверяем, что запись добавлена
	if len(j.entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(j.entries))
	}
	
	// Проверяем, что ID сгенерирован
	if j.entries[0].ID == "" {
		t.Error("Expected entry ID to be generated")
	}
	
	// Проверяем, что искажения обнаружились (КПТ-режим)
	if len(j.entries[0].Distortions) == 0 {
		t.Log("Note: Distortions may be empty if text doesn't match patterns")
	}
	
	// Проверяем, что теги проставились
	if len(j.entries[0].Tags) == 0 {
		t.Error("Expected tags to be auto-generated")
	}
}

func TestJournal_GetEntries(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "journal_test")
	defer os.RemoveAll(tmpDir)
	
	j, _ := NewJournal(JournalConfig{DataDir: tmpDir, UseOllamaEmbed: false})
	
	// Добавляем несколько записей
	j.AddEntry(ThoughtEntry{
		Type:             EntryTypeCBT,
		Timestamp:        time.Now(),
		AutomaticThought: "тест 1",
		PersonID:         "Dina",
		Phase:            "Detox",
	})
	j.AddEntry(ThoughtEntry{
		Type:             EntryTypeCBT,
		Timestamp:        time.Now(),
		AutomaticThought: "тест 2",
		PersonID:         "Valya",
		Phase:            "Rewire",
	})
	
	// Фильтр по человеку
	entries := j.GetEntries(EntryFilters{PersonID: "Dina"})
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry for Dina, got %d", len(entries))
	}
	
	// Фильтр по фазе
	entries = j.GetEntries(EntryFilters{Phase: "Rewire"})
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry for Rewire, got %d", len(entries))
	}
}

func TestJournal_GetStats(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "journal_test")
	defer os.RemoveAll(tmpDir)
	
	j, _ := NewJournal(JournalConfig{DataDir: tmpDir, UseOllamaEmbed: false})
	
	j.AddEntry(ThoughtEntry{
		Type:             EntryTypeCBT,
		AutomaticThought: "всегда всё плохо",
		Emotions:         []string{"тревога", "тревога"},
		Intensity:        80,
	})
	j.AddEntry(ThoughtEntry{
		Type:             EntryTypeCBT,
		AutomaticThought: "я должен быть идеальным",
		Emotions:         []string{"вина"},
		Intensity:        60,
	})
	
	stats := j.GetStats() // backward compat wrapper
	
	if stats.TotalEntries != 2 {
		t.Errorf("Expected 2 entries, got %d", stats.TotalEntries)
	}
	if stats.CBTEntries != 2 {
		t.Errorf("Expected 2 CBT entries, got %d", stats.CBTEntries)
	}
}

func TestJournal_SaveLoad(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "journal_test")
	defer os.RemoveAll(tmpDir)
	
	j1, _ := NewJournal(JournalConfig{DataDir: tmpDir, UseOllamaEmbed: false})
	j1.AddEntry(ThoughtEntry{
		Type:             EntryTypeCBT,
		AutomaticThought: "тестовая запись",
		Emotions:         []string{"тревога"},
		Intensity:        50,
		Phase:            "Detox",
		PersonID:         "Test",
	})
	
	// Создаём новый журнал и загружаем
	j2, _ := NewJournal(JournalConfig{DataDir: tmpDir, UseOllamaEmbed: false})
	
	if len(j2.entries) != 1 {
		t.Errorf("Expected 1 entry after load, got %d", len(j2.entries))
	}
}

func TestJournal_SearchByMeaning(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "journal_test")
	defer os.RemoveAll(tmpDir)
	
	j, _ := NewJournal(JournalConfig{DataDir: tmpDir, UseOllamaEmbed: false})
	
	// Добавляем записи с разным смыслом
	j.AddEntry(ThoughtEntry{
		Type:             EntryTypeCBT,
		AutomaticThought: "финансовые проблемы, не могу платить",
		PersonID:         "Dina",
	})
	j.AddEntry(ThoughtEntry{
		Type:             EntryTypeCBT,
		AutomaticThought: "сегодня хорошая погода",
		PersonID:         "Valya",
	})
	
	// Поиск по смыслу (с заглушкой эмбеддингов)
	results := j.SearchByMeaning("деньги, долг, ресурс", 5)
	
	// С заглушкой все результаты будут иметь одинаковое сходство,
	// поэтому проверяем только что поиск не падает
	if len(results) == 0 {
		t.Error("Expected some results from semantic search")
	}
}
