package panels

import "ideal-core/pkg/bio"

// MetabolicPanel — метаболическая панель
//
// Состав:
// - Глюкоза — уровень сахара в крови
// - Гликированный гемоглобин — средний сахар за 3 месяца
// - Инсулин — чувствительность к инсулину
// - Липидограмма — холестерин и жиры
// - Гомоцистеин — воспаление, B-витамины
// - С-реактивный белок — воспаление
// - Ферритин — запасы железа
var MetabolicPanel = LabTestPanel{
	ID:             "metabolic_panel",
	Name:           "Метаболизм и диабет",
	NameEn:         "Metabolic Panel",
	Description:    "Исследование метаболизма и риска диабета",
	Cost:           5000,
	TurnaroundTime: 48,
	Tests: []bio.LabTest{
		{ID: "glucose", Name: "Глюкоза", Cost: 300},
		{ID: "hba1c", Name: "Гликированный гемоглобин", Cost: 600},
		{ID: "insulin", Name: "Инсулин", Cost: 700},
		{ID: "lipid_profile", Name: "Липидограмма", Cost: 1000},
		{ID: "homocysteine", Name: "Гомоцистеин", Cost: 1500},
		{ID: "crp", Name: "С-реактивный белок", Cost: 500},
		{ID: "ferritin", Name: "Ферритин", Cost: 700},
	},
	ChakraCorrelations: []int{2},
	RecommendedFor: []string{
		"Лишний вес",
		"Тяга к сладкому",
		"Усталость после еды",
		"Семейная история диабета",
		"Высокое давление",
	},
}

// GetMetabolicPanel возвращает метаболическую панель
func GetMetabolicPanel() LabTestPanel {
	return MetabolicPanel
}
