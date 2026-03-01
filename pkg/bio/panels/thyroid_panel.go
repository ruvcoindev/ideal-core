package panels

import "ideal-core/pkg/bio"

// ThyroidPanel — панель щитовидной железы
//
// Состав:
// - ТТГ (основной маркер функции щитовидной)
// - Free T4 (тироксин, основной гормон)
// - Free T3 (трийодтиронин, активная форма)
// - Антитела к ТПО (аутоиммунный процесс)
// - Антитела к ТГ (аутоиммунный процесс)
// - Тиреоглобулин (маркер ткани щитовидной)
var ThyroidPanel = LabTestPanel{
	ID:             "thyroid_panel",
	Name:           "Щитовидная железа",
	NameEn:         "Thyroid Panel",
	Description:    "Комплексное исследование функции щитовидной железы",
	Cost:           4000,
	TurnaroundTime: 48,
	Tests: []bio.LabTest{
		{ID: "tsh", Name: "ТТГ", Cost: 500},
		{ID: "free_t4", Name: "Free T4", Cost: 500},
		{ID: "free_t3", Name: "Free T3", Cost: 500},
		{ID: "tpo_antibodies", Name: "Антитела к ТПО", Cost: 1000},
		{ID: "tg_antibodies", Name: "Антитела к ТГ", Cost: 1000},
		{ID: "thyroglobulin", Name: "Тиреоглобулин", Cost: 500},
	},
	ChakraCorrelations: []int{4},
	RecommendedFor: []string{
		"Усталость",
		"Набор веса",
		"Выпадение волос",
		"Зябкость",
		"«Ком в горле»",
		"Нарушения цикла",
	},
}

// GetThyroidPanel возвращает панель щитовидной железы
func GetThyroidPanel() LabTestPanel {
	return ThyroidPanel
}
