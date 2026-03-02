package journal

import (
	"os"
	"testing"
)

func TestJournal_AddGratitudeEntry(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "journal_gratitude_test")
	defer os.RemoveAll(tmpDir)
	
	j, _ := NewJournal(JournalConfig{
		DataDir:        tmpDir,
		UseOllamaEmbed: false,
	})
	
	items := []GratitudeItem{
		{
			Text:        "Я благодарен за вкусный кофе утром",
			Category:    "small_things",
			Specificity: 9, // конкретно: кофе, утро
			Emotion:     "warmth",
		},
		{
			Text:        "Я благодарен Дине за то, что она была частью моей жизни",
			Category:    "people",
			Specificity: 8,
			Emotion:     "bittersweet",
		},
	}
	
	err := j.AddGratitudeEntry(items, "День начался спокойно")
	if err != nil {
		t.Errorf("AddGratitudeEntry failed: %v", err)
	}
	
	if len(j.entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(j.entries))
	}
	
	entry := j.entries[0]
	if entry.Type != EntryTypeGratitude {
		t.Errorf("Expected type gratitude, got %s", entry.Type)
	}
	if entry.GratitudeLevel < 8 {
		t.Errorf("Expected high gratitude level, got %d", entry.GratitudeLevel)
	}
	if !containsString(entry.Tags, "gratitude_small") {
		t.Error("Expected tag 'gratitude_small'")
	}
}

func TestJournal_GetGratitudeStats(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "journal_stats_test")
	defer os.RemoveAll(tmpDir)
	
	j, _ := NewJournal(JournalConfig{DataDir: tmpDir, UseOllamaEmbed: false})
	
	// Добавляем записи благодарности
	j.AddGratitudeEntry([]GratitudeItem{
		{Text: "Благодарен за солнце", Category: "nature", Specificity: 7},
	}, "")
	j.AddGratitudeEntry([]GratitudeItem{
		{Text: "Благодарен за поддержку друга", Category: "people", Specificity: 8},
	}, "")
	
	stats := j.GetGratitudeStats()
	
	if stats.TotalEntries != 2 {
		t.Errorf("Expected 2 gratitude entries, got %d", stats.TotalEntries)
	}
	if stats.CategoryCount["nature"] != 1 {
		t.Error("Expected 1 nature category")
	}
}

func TestJournal_CombinedStats(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "journal_combined_test")
	defer os.RemoveAll(tmpDir)
	
	j, _ := NewJournal(JournalConfig{DataDir: tmpDir, UseOllamaEmbed: false})
	
	// Смешанные записи
	j.AddCBTEntry("Ситуация", "Негативная мысль", []string{"тревога"}, 70)
	j.AddGratitudeEntry([]GratitudeItem{{Text: "Благодарен", Category: "self", Specificity: 5}}, "")
	
	stats := j.GetCombinedStats()
	
	if stats.CBTEntries != 1 || stats.GratitudeEntries != 1 {
		t.Errorf("Expected 1 CBT + 1 gratitude, got %d + %d", stats.CBTEntries, stats.GratitudeEntries)
	}
	if stats.Ratio != 0.5 {
		t.Errorf("Expected ratio 0.5, got %.2f", stats.Ratio)
	}
}
