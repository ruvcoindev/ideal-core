package panels

import "ideal-core/pkg/bio"

// VitaminPanel — панель витаминов
//
// Состав:
// - Витамин D (иммунитет, кости, настроение)
// - Витамин B12 (энергия, нервы)
// - Фолиевая кислота (ДНК, беременность)
// - Витамин B6 (метаболизм, гормоны)
// - Витамин C (иммунитет, антиоксидант)
// - Витамин A (зрение, кожа)
var VitaminPanel = LabTestPanel{
	ID:             "vitamin_panel",
	Name:           "Панель витаминов",
	NameEn:         "Vitamin Panel",
	Description:    "Основные витамины и дефициты",
	Cost:           6000,
	TurnaroundTime: 48,
	Tests: []bio.LabTest{
		{ID: "vitamin_d", Name: "Витамин D (25-OH)", Cost: 2000},
		{ID: "vitamin_b12", Name: "Витамин B12", Cost: 800},
		{ID: "folate", Name: "Фолиевая кислота", Cost: 800},
		{ID: "vitamin_b6", Name: "Витамин B6", Cost: 800},
		{ID: "vitamin_c", Name: "Витамин C", Cost: 800},
		{ID: "vitamin_a", Name: "Витамин A", Cost: 800},
	},
	ChakraCorrelations: []int{3, 6},
	RecommendedFor: []string{
		"Усталость",
		"Частые простуды",
		"Выпадение волос",
		"Сухость кожи",
		"Депрессивные состояния",
	},
}

// GetVitaminPanel возвращает панель витаминов
func GetVitaminPanel() LabTestPanel {
	return VitaminPanel
}
