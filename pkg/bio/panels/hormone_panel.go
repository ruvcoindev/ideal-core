// Package panels — готовые панели анализов
package panels

import "ideal-core/pkg/bio"

// LabTestPanel — панель лабораторных тестов
type LabTestPanel struct {
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	NameEn             string        `json:"name_en"`
	Description        string        `json:"description"`
	Cost               int           `json:"cost"`
	TurnaroundTime     int           `json:"turnaround_time"`
	Tests              []bio.LabTest `json:"tests"`
	ChakraCorrelations []int         `json:"chakra_correlations"`
	RecommendedFor     []string      `json:"recommended_for"`
}

// HormonePanel — гормональная панель
var HormonePanel = LabTestPanel{
	ID:             "hormone_panel",
	Name:           "Гормональная панель",
	NameEn:         "Hormone Panel",
	Description:    "Комплексное исследование гормонального фона",
	Cost:           8000,
	TurnaroundTime: 48,
	Tests: []bio.LabTest{
		{ID: "cortisol", Name: "Кортизол", Cost: 800},
		{ID: "tsh", Name: "ТТГ", Cost: 500},
		{ID: "free_t4", Name: "Free T4", Cost: 500},
		{ID: "free_t3", Name: "Free T3", Cost: 500},
		{ID: "testosterone_total", Name: "Тестостерон общий", Cost: 700},
		{ID: "estradiol", Name: "Эстрадиол", Cost: 700},
		{ID: "progesterone", Name: "Прогестерон", Cost: 700},
		{ID: "lh", Name: "ЛГ", Cost: 600},
		{ID: "fsh", Name: "ФСГ", Cost: 600},
		{ID: "prolactin", Name: "Пролактин", Cost: 600},
		{ID: "dhea_s", Name: "ДГЭА-С", Cost: 600},
		{ID: "shbg", Name: "ГСПГ", Cost: 700},
	},
	ChakraCorrelations: []int{0, 1, 2, 4},
	RecommendedFor: []string{
		"Усталость",
		"Нарушения цикла",
		"Снижение либидо",
		"Набор веса",
		"Перепады настроения",
	},
}

// GetHormonePanel возвращает гормональную панель
func GetHormonePanel() LabTestPanel {
	return HormonePanel
}
