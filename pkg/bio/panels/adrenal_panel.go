package panels

import "ideal-core/pkg/bio"

// AdrenalPanel — панель надпочечников
//
// Состав:
// - Кортизол в слюне (4 точки) — суточный ритм
// - ДГЭА-С — резерв надпочечников
// - Адреналин — острый стресс
// - Норадреналин — хронический стресс
var AdrenalPanel = LabTestPanel{
	ID:             "adrenal_panel",
	Name:           "Надпочечники и стресс",
	NameEn:         "Adrenal Panel",
	Description:    "Исследование функции надпочечников и стресс-ответа",
	Cost:           3500,
	TurnaroundTime: 48,
	Tests: []bio.LabTest{
		{ID: "cortisol_saliva", Name: "Кортизол в слюне (4 точки)", Cost: 3500},
		{ID: "dhea_s", Name: "ДГЭА-С", Cost: 600},
		{ID: "adrenaline", Name: "Адреналин", Cost: 800},
		{ID: "noradrenaline", Name: "Норадреналин", Cost: 800},
	},
	ChakraCorrelations: []int{0, 2},
	RecommendedFor: []string{
		"Хроническая усталость",
		"Тревожность",
		"Выгорание",
		"Нарушения сна",
		"Тяга к солёному",
		"Панические атаки",
	},
}

// GetAdrenalPanel возвращает панель надпочечников
func GetAdrenalPanel() LabTestPanel {
	return AdrenalPanel
}
