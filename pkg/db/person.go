package db

import (
	"time"
	"ideal-core/pkg/cube"
	"ideal-core/pkg/identity"
)

// RelationType определяет тип связи в родословной
type RelationType string

const (
	RelSelf       RelationType = "SELF"
	RelMother     RelationType = "MOTHER"
	RelFather     RelationType = "FATHER"
	RelGrandmother RelationType = "GRANDMOTHER"
	RelGrandfather RelationType = "GRANDFATHER"
	RelDaughter   RelationType = "DAUGHTER"
	RelSon        RelationType = "SON"
	RelPartner    RelationType = "PARTNER"
	RelFriend     RelationType = "FRIEND"
	RelAcquaintance RelationType = "ACQUAINTANCE"
	RelLover      RelationType = "LOVER"
	RelOther      RelationType = "OTHER"
)

// Chakra представляет энергетический центр
type Chakra struct {
	Level     int    // 1-100: уровень активности
	Blocked   bool   // заблокирована ли
	LastCheck time.Time
}

// LunarProfile содержит лунные данные
type LunarProfile struct {
	BirthPhase      int    // Фаза луны при рождении (0-29)
	BirthSign       string // Знак зодиака
	CurrentPhase    int    // Текущая фаза (рассчитывается)
	FavorableDays   []int  // Благоприятные дни для действий
}

// PsychosomaticRecord связывает заболевание с аффирмацией
type PsychosomaticRecord struct {
	MedicalTerm     string   // "Алопеция"
	FolkTerm        string   // "Облысение"
	ChakraIndex     int      // 0-6: какая чакра затронута
	AffirmationID   string   // Ссылка на базу аффирмаций
	AddedAt         time.Time
	Source          string   // "Louise Hay", "Zhikarentsev", "User"
}

// Person — полная сущность пользователя
type Person struct {
	// Идентификация
	ID          identity.PublicKey
	PublicID    string              // Детерминированный ID из даты
	Name        string
	BirthDate   time.Time
	DeathDate   *time.Time
	Gender      string
	
	// Кубические данные
	CubeData    cube.Person
	
	// Энергетический профиль
	Chakras     [7]Chakra
	Lunar       LunarProfile
	
	// Здоровье и психосоматика
	Diagnoses   []PsychosomaticRecord
	
	// Связи (хранятся как хеши в Merkle tree)
	Relations   map[RelationType]string // ID связанного человека
	
	// Статус потока
	FlowStatus  string   // "Active", "Draining", "Resource", "Offline"
	Tags        []string // Контекст: "Family", "Business", "Pyramid"
	
	// Метаданные
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsAnonymized bool // Обезличен ли для общей базы
}

// NewPerson создаёт нового человека с расчётом всех профилей
func NewPerson(pubKey identity.PublicKey, name string, birthDate time.Time) *Person {
	cubePerson := cube.NewPerson("", name, birthDate)
	
	return &Person{
		ID:          pubKey,
		PublicID:    identity.GenerateID(birthDate, "ideal-core-v1"),
		Name:        name,
		BirthDate:   birthDate,
		Gender:      "",
		CubeData:    *cubePerson,
		Chakras:     initDefaultChakras(),
		Lunar:       CalculateLunarProfile(birthDate),
		Relations:   make(map[RelationType]string),
		FlowStatus:  "Active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func initDefaultChakras() [7]Chakra {
	var chakras [7]Chakra
	for i := range chakras {
		chakras[i] = Chakra{Level: 50, Blocked: false, LastCheck: time.Now()}
	}
	return chakras
}

// CalculateLunarProfile вычисляет лунный профиль по дате рождения
// (упрощённая версия — можно заменить на астрономическую библиотеку)
func CalculateLunarProfile(birthDate time.Time) LunarProfile {
	// Упрощённый расчёт фазы: (день + месяц) % 30
	phase := (birthDate.Day() + int(birthDate.Month())) % 30
	
	// Знак зодиака по месяцу и дню
	sign := getZodiacSign(birthDate.Month(), birthDate.Day())
	
	return LunarProfile{
		BirthPhase:    phase,
		BirthSign:     sign,
		CurrentPhase:  phase, // будет обновляться при использовании
		FavorableDays: calculateFavorableDays(phase),
	}
}

func getZodiacSign(month time.Month, day int) string {
	// Упрощённая таблица
	switch {
	case (month == time.December && day >= 22) || (month == time.January && day <= 19):
		return "Capricorn"
	case (month == time.January && day >= 20) || (month == time.February && day <= 18):
		return "Aquarius"
	// ... добавить остальные знаки
	default:
		return "Unknown"
	}
}

func calculateFavorableDays(birthPhase int) []int {
	// Пример: благоприятные дни — когда текущая фаза близка к натальной
	return []int{(birthPhase - 1) % 30 + 1, birthPhase, (birthPhase + 1) % 30 + 1}
}
