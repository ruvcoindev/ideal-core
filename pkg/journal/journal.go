package journal

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"ideal-core/pkg/cbt"
	"ideal-core/pkg/vector"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// EntryType — тип записи в дневнике
type EntryType string

const (
	EntryTypeCBT       EntryType = "cbt"           // Когнитивно-поведенческая терапия
	EntryTypeGratitude EntryType = "gratitude"     // Дневник благодарности
	EntryTypeReflection EntryType = "reflection"   // Свободное размышление
)

// ThoughtEntry — универсальная запись в дневнике
type ThoughtEntry struct {
	ID               string                      `json:"id"`
	Type             EntryType                   `json:"type"` // cbt | gratitude | reflection
	Timestamp        time.Time                   `json:"timestamp"`
	
	// Общие поля
	Situation        string                      `json:"situation,omitempty"`
	Notes            string                      `json:"notes,omitempty"`
	Emotions         []string                    `json:"emotions"`
	Intensity        int                         `json:"intensity,omitempty"` // 0-100
	Tags             []string                    `json:"tags"`
	Phase            string                      `json:"phase,omitempty"`     // Protocol36 phase
	PersonID         string                      `json:"person_id,omitempty"` // Связь с человеком
	Chakras          []int                       `json:"chakras,omitempty"`
	
	// Поля для режима КПТ
	AutomaticThought string                      `json:"automatic_thought,omitempty"`
	Distortions      []cbt.CognitiveDistortion   `json:"distortions,omitempty"`
	RationalResponse string                      `json:"rational_response,omitempty"`
	NewIntensity     int                         `json:"new_intensity,omitempty"`
	
	// Поля для режима благодарности
	GratitudeItems   []GratitudeItem             `json:"gratitude_items,omitempty"`
	GratitudeLevel   int                         `json:"gratitude_level,omitempty"` // 1-10, насколько глубоко прочувствовали
	
	// Вектор для семантического поиска (не сериализуется в JSON)
	Embedding        vector.Embedding            `json:"-"`
}

// GratitudeItem — один пункт в дневнике благодарности
type GratitudeItem struct {
	Text        string    `json:"text"`          // "Я благодарен за..."
	Category    string    `json:"category"`      // "people", "nature", "self", "small_things", "growth"
	Specificity int       `json:"specificity"`   // 1-10, насколько конкретно описано
	Emotion     string    `json:"emotion"`       // Какая эмоция при этом возникла
}

// JournalConfig — конфигурация дневника
type JournalConfig struct {
	DataDir         string
	OllamaHost      string
	OllamaModel     string
	UseOllamaEmbed  bool
	DefaultMode     EntryType // cbt | gratitude
}

// Journal — дневник с поддержкой нескольких режимов
type Journal struct {
	entries      []ThoughtEntry
	filePath     string
	vectorStore  vector.VectorStore
	ollamaClient *vector.OllamaEmbeddingClient
	useOllama    bool
	defaultMode  EntryType
}

// NewJournal создаёт новый дневник
func NewJournal(cfg JournalConfig) (*Journal, error) {
	if err := os.MkdirAll(cfg.DataDir, 0700); err != nil {
		return nil, err
	}
	
	j := &Journal{
		entries:     make([]ThoughtEntry, 0),
		filePath:    filepath.Join(cfg.DataDir, "thoughts.json"),
		vectorStore: vector.NewMockVectorStore(),
		useOllama:   cfg.UseOllamaEmbed,
		defaultMode: cfg.DefaultMode,
	}
	
	if cfg.UseOllamaEmbed {
		j.ollamaClient = vector.NewOllamaEmbeddingClient(cfg.OllamaHost, cfg.OllamaModel)
		if !j.ollamaClient.IsAvailable() {
			fmt.Printf("⚠️  Ollama not available, using stub embeddings\n")
			j.useOllama = false
		}
	}
	
	if err := j.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	
	return j, nil
}

// AddCBTEntry добавляет запись в режиме КПТ
func (j *Journal) AddCBTEntry(situation, automaticThought string, emotions []string, intensity int) error {
	entry := ThoughtEntry{
		Type:             EntryTypeCBT,
		Timestamp:        time.Now(),
		Situation:        situation,
		AutomaticThought: automaticThought,
		Emotions:         emotions,
		Intensity:        intensity,
	}
	return j.addEntryWithProcessing(entry)
}

// AddGratitudeEntry добавляет запись в режиме благодарности
func (j *Journal) AddGratitudeEntry(items []GratitudeItem, notes string) error {
	// Авто-расчёт уровня благодарности по конкретике
	level := 0
	for _, item := range items {
		level += item.Specificity
	}
	if len(items) > 0 {
		level = level / len(items)
	}
	
	entry := ThoughtEntry{
		Type:           EntryTypeGratitude,
		Timestamp:      time.Now(),
		GratitudeItems: items,
		GratitudeLevel: level,
		Notes:          notes,
		Emotions:       []string{"gratitude", "warmth", "peace"}, // авто-эмоции для благодарности
	}
	return j.addEntryWithProcessing(entry)
}

// addEntryWithProcessing — общая логика обработки записи
func (j *Journal) addEntryWithProcessing(entry ThoughtEntry) error {
	// Генерация ID
	hash := sha256.Sum256([]byte(entry.Timestamp.String() + entry.Notes + entry.AutomaticThought))
	entry.ID = hex.EncodeToString(hash[:8])
	
	// Обработка в зависимости от типа
	switch entry.Type {
	case EntryTypeCBT:
		entry.Distortions = cbt.DetectDistortions(entry.AutomaticThought)
		if entry.RationalResponse == "" {
			entry.RationalResponse = cbt.GenerateRationalResponse(entry.AutomaticThought, entry.Distortions)
		}
	case EntryTypeGratitude:
		// Авто-тегирование для благодарности
		entry.Tags = append(entry.Tags, j.autoTagGratitude(entry.GratitudeItems)...)
	}
	
	// Общие теги
	entry.Tags = append(entry.Tags, j.autoTagCommon(entry)...)
	entry.Tags = unique(entry.Tags)
	
	// Векторизация
	entry.Embedding = j.generateEmbedding(entry.toSearchText())
	j.vectorStore.Upsert(entry.ID, entry.Embedding, map[string]interface{}{
		"type":     string(entry.Type),
		"emotions": entry.Emotions,
		"phase":    entry.Phase,
		"person":   entry.PersonID,
		"tags":     entry.Tags,
	})
	
	j.entries = append(j.entries, entry)
	return j.Save()
}

// toSearchText возвращает текст для векторизации (объединяет все поля)
func (e *ThoughtEntry) toSearchText() string {
	parts := []string{e.Situation, e.Notes, e.AutomaticThought, e.RationalResponse}
	for _, item := range e.GratitudeItems {
		parts = append(parts, item.Text)
	}
	parts = append(parts, e.Emotions...)
	parts = append(parts, e.Tags...)
	return strings.Join(parts, " ")
}

// generateEmbedding генерирует вектор (Ollama или заглушка)
func (j *Journal) generateEmbedding(text string) vector.Embedding {
	if j.useOllama && j.ollamaClient != nil {
		emb, err := j.ollamaClient.GenerateEmbedding(text)
		if err == nil {
			return emb
		}
		fmt.Printf("⚠️  Ollama embedding failed: %v\n", err)
	}
	// Заглушка: bge-m3 = 1024 dimensions
	return make(vector.Embedding, 1024)
}

// autoTagGratitude проставляет теги для записей благодарности
func (j *Journal) autoTagGratitude(items []GratitudeItem) []string {
	var tags []string
	for _, item := range items {
		switch item.Category {
		case "people":
			tags = append(tags, "gratitude_people")
		case "nature":
			tags = append(tags, "gratitude_nature")
		case "self":
			tags = append(tags, "gratitude_self")
		case "small_things":
			tags = append(tags, "gratitude_small")
		case "growth":
			tags = append(tags, "gratitude_growth")
		}
		if item.Specificity >= 8 {
			tags = append(tags, "gratitude_specific")
		}
	}
	return tags
}

// autoTagCommon — общие теги для всех типов записей
func (j *Journal) autoTagCommon(entry ThoughtEntry) []string {
	var tags []string
	text := entry.toSearchText()
	
	// Персоналии
	if containsAny(text, []string{"Дина", "бывшая", "партнёр"}) {
		tags = append(tags, "relationship_dina")
	}
	if containsAny(text, []string{"Валя", "мама", "мать"}) {
		tags = append(tags, "family_valya")
	}
	
	// Темы
	if containsAny(text, []string{"деньги", "ресурс", "финансы", "долг"}) {
		tags = append(tags, "resource")
	}
	if containsAny(text, []string{"границы", "нет", "стоп"}) {
		tags = append(tags, "boundaries")
	}
	
	// Эмоции
	for _, em := range entry.Emotions {
		switch strings.ToLower(em) {
		case "страх", "тревога":
			tags = append(tags, "fear")
		case "гнев", "злость":
			tags = append(tags, "anger")
		case "благодарность", "gratitude":
			tags = append(tags, "gratitude")
		}
	}
	
	return tags
}

// GetEntries возвращает записи с фильтрами
func (j *Journal) GetEntries(filters EntryFilters) []ThoughtEntry {
	var result []ThoughtEntry
	for _, e := range j.entries {
		if filters.Type != "" && string(e.Type) != filters.Type {
			continue
		}
		if filters.PersonID != "" && e.PersonID != filters.PersonID {
			continue
		}
		if filters.Phase != "" && e.Phase != filters.Phase {
			continue
		}
		if filters.Tag != "" && !containsString(e.Tags, filters.Tag) {
			continue
		}
		result = append(result, e)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})
	return result
}

// EntryFilters — фильтры для поиска
type EntryFilters struct {
	Type     string // "cbt", "gratitude", "" for all
	PersonID string
	Phase    string
	Tag      string
	FromDate *time.Time
	ToDate   *time.Time
}

// SearchByMeaning — семантический поиск по всем записям
func (j *Journal) SearchByMeaning(query string, limit int) []ThoughtEntry {
	queryEmbedding := j.generateEmbedding(query)
	results := j.vectorStore.Search(queryEmbedding, limit)
	
	var entries []ThoughtEntry
	for _, r := range results {
		for _, e := range j.entries {
			if e.ID == r.ID {
				entries = append(entries, e)
				break
			}
		}
	}
	return entries
}

// GetGratitudeStats — статистика только по записям благодарности
func (j *Journal) GetGratitudeStats() GratitudeStats {
	stats := GratitudeStats{TotalEntries: 0, AvgLevel: 0, CategoryCount: make(map[string]int)}
	var totalLevel int
	
	for _, e := range j.entries {
		if e.Type != EntryTypeGratitude {
			continue
		}
		stats.TotalEntries++
		totalLevel += e.GratitudeLevel
		for _, item := range e.GratitudeItems {
			stats.CategoryCount[item.Category]++
		}
	}
	
	if stats.TotalEntries > 0 {
		stats.AvgLevel = totalLevel / stats.TotalEntries
		stats.TopCategories = topNStrings(stats.CategoryCount, 5)
	}
	return stats
}

// GratitudeStats — статистика благодарности
type GratitudeStats struct {
	TotalEntries  int            `json:"total_entries"`
	AvgLevel      int            `json:"avg_level"`
	CategoryCount map[string]int `json:"category_count"`
	TopCategories []string       `json:"top_categories"`
}

// GetCombinedStats — общая статистика по всем режимам
func (j *Journal) GetCombinedStats() CombinedStats {
	cbtCount, gratitudeCount := 0, 0
	for _, e := range j.entries {
		if e.Type == EntryTypeCBT {
			cbtCount++
		} else if e.Type == EntryTypeGratitude {
			gratitudeCount++
		}
	}
	return CombinedStats{
		TotalEntries:    len(j.entries),
		CBTEntries:      cbtCount,
		GratitudeEntries: gratitudeCount,
		Ratio:            float64(gratitudeCount) / max(1, float64(cbtCount+gratitudeCount)),
	}
}

// CombinedStats — сводная статистика
type CombinedStats struct {
	TotalEntries     int     `json:"total_entries"`
	CBTEntries       int     `json:"cbt_entries"`
	GratitudeEntries int     `json:"gratitude_entries"`
	Ratio            float64 `json:"gratitude_ratio"` // доля благодарности
}

// Save/Load/Delete/Export — методы сохранения (аналогично предыдущей версии)
// ... (код Save/Load аналогичен, с учётом новых полей)

func (j *Journal) Save() error {
	data, err := json.MarshalIndent(j.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(j.filePath, data, 0600)
}

func (j *Journal) Load() error {
	data, err := os.ReadFile(j.filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &j.entries); err != nil {
		return err
	}
	// Пересчитываем эмбеддинги при загрузке
	for i := range j.entries {
		j.entries[i].Embedding = j.generateEmbedding(j.entries[i].toSearchText())
		j.vectorStore.Upsert(j.entries[i].ID, j.entries[i].Embedding, map[string]interface{}{
			"type": string(j.entries[i].Type),
		})
	}
	return nil
}

func (j *Journal) DeleteEntry(id string) error {
	for i, e := range j.entries {
		if e.ID == id {
			j.entries = append(j.entries[:i], j.entries[i+1:]...)
			j.vectorStore.Delete(id)
			return j.Save()
		}
	}
	return os.ErrNotExist
}

// Helpers
func containsAny(s string, substrs []string) bool {
	sLower := strings.ToLower(s)
	for _, sub := range substrs {
		if strings.Contains(sLower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func unique(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func topNStrings(m map[string]int, n int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range m {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	var result []string
	for i := 0; i < n && i < len(sorted); i++ {
		result = append(result, sorted[i].Key)
	}
	return result
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// ============================================================================
// BACKWARD COMPATIBILITY WRAPPERS (для старых тестов и простого использования)
// ============================================================================

// AddEntry — универсальный метод добавления (для обратной совместимости)
// Автоматически определяет тип по наличию полей
func (j *Journal) AddEntry(entry ThoughtEntry) error {
	if entry.Type == "" {
		// Авто-определение типа
		if entry.AutomaticThought != "" {
			entry.Type = EntryTypeCBT
		} else if len(entry.GratitudeItems) > 0 {
			entry.Type = EntryTypeGratitude
		} else {
			entry.Type = EntryTypeReflection
		}
	}
	return j.addEntryWithProcessing(entry)
}

// GetStats — возвращает комбинированную статистику (для обратной совместимости)
func (j *Journal) GetStats() CombinedStats {
	return j.GetCombinedStats()
}

// ExportToMarkdown экспортирует дневник в Markdown для печати
func (j *Journal) ExportToMarkdown(outputPath string) error {
	md := "# 📓 Дневник мыслей\n\n"
	md += fmt.Sprintf("Всего записей: %d\n\n", len(j.entries))
	
	entries := j.GetEntries(EntryFilters{})
	for _, e := range entries {
		md += fmt.Sprintf("## %s\n", e.Timestamp.Format("02.01.2006 15:04"))
		md += fmt.Sprintf("**Тип:** %s", e.Type)
		if e.Phase != "" {
			md += fmt.Sprintf(" | **Фаза:** %s", e.Phase)
		}
		if e.PersonID != "" {
			md += fmt.Sprintf(" | **Человек:** %s", e.PersonID)
		}
		md += fmt.Sprintf(" | **Интенсивность:** %d/100\n\n", e.Intensity)
		
		if e.Type == EntryTypeGratitude {
			md += "### 💛 Благодарность\n"
			for _, item := range e.GratitudeItems {
				md += fmt.Sprintf("- %s [%s, конкретика: %d/10]\n", item.Text, item.Category, item.Specificity)
			}
		} else {
			md += fmt.Sprintf("### Ситуация\n%s\n\n", e.Situation)
			md += fmt.Sprintf("### Автоматическая мысль\n%s\n\n", e.AutomaticThought)
			if len(e.Distortions) > 0 {
				md += "### Искажения\n"
				for _, d := range e.Distortions {
					md += fmt.Sprintf("- %s\n", d)
				}
				md += "\n"
			}
			if e.RationalResponse != "" {
				md += fmt.Sprintf("### Рациональный ответ\n%s\n\n", e.RationalResponse)
			}
		}
		
		if len(e.Emotions) > 0 {
			md += fmt.Sprintf("**Эмоции:** %s\n\n", strings.Join(e.Emotions, ", "))
		}
		if len(e.Tags) > 0 {
			md += fmt.Sprintf("**Теги:** %s\n\n", strings.Join(e.Tags, ", "))
		}
		md += "---\n\n"
	}
	
	return os.WriteFile(outputPath, []byte(md), 0644)
}
