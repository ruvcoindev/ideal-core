package cube

import (
	"time"
)

// Person представляет узел в системе Куба
type Person struct {
	ID           string        // Уникальный идентификатор (протокол Иггдрасиль)
	Name         string        // Имя или псевдоним
	BirthDate    time.Time     // Дата рождения
	DeathDate    *time.Time    // Опционально: дата ухода
	Gender       string        // М/Ж для лунных циклов
	
	// Кубические координаты (сумма цифр)
	Coordinates  [3]int        // [День, Месяц, Год]
	SumFrequency int           // Σ - общая частота
	
	// Векторы движения (разности соседних цифр)
	Vectors      [3][3]int     // [день][3], [месяц][3], [год][3]
	
	// Статус потока
	FlowStatus   string        // "Active", "Draining", "Resource", "Offline"
	Tags         []string      // Контекст: "Семья", "Бизнес", "Пирамида"
}

// SumDigits вычисляет сумму цифр числа
func SumDigits(n int) int {
	if n < 0 { n = -n }
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// CalcCoordinates вычисляет координаты по дате рождения
// Формат: DD.MM.YYYY → [Σ(DD), Σ(MM), Σ(YYYY%1000)]
func CalcCoordinates(date time.Time) [3]int {
	day := SumDigits(date.Day())
	month := SumDigits(int(date.Month()))
	year := SumDigits(date.Year() % 1000) // последние 3 цифры года
	return [3]int{day, month, year}
}

// CalcVectors вычисляет векторы движения по методике фильма "Куб"
// Вектор = разность соседних цифр с циклическим замыканием: [a-b, b-c, c-a]
func CalcVectors(n int) [3]int {
	// Извлекаем цифры (падим до 3 знаков)
	digits := [3]int{0, 0, 0}
	temp := n
	for i := 2; i >= 0 && temp > 0; i-- {
		digits[i] = temp % 10
		temp /= 10
	}
	// Вычисляем векторы с циклическим замыканием
	return [3]int{
		digits[0] - digits[1], // a-b
		digits[1] - digits[2], // b-c
		digits[2] - digits[0], // c-a
	}
}

// NewPerson создаёт нового человека с расчётом координат и векторов
func NewPerson(id, name string, birthDate time.Time) *Person {
	coords := CalcCoordinates(birthDate)
	p := &Person{
		ID:           id,
		Name:         name,
		BirthDate:    birthDate,
		Coordinates:  coords,
		SumFrequency: coords[0] + coords[1] + coords[2],
		Vectors: [3][3]int{
			CalcVectors(birthDate.Day()),
			CalcVectors(int(birthDate.Month())),
			CalcVectors(birthDate.Year() % 1000),
		},
		FlowStatus: "Active",
		Tags:       []string{},
	}
	return p
}

// Distance вычисляет евклидово расстояние между двумя комнатами
func Distance(p1, p2 Person) float64 {
	dx := float64(p1.Coordinates[0] - p2.Coordinates[0])
	dy := float64(p1.Coordinates[1] - p2.Coordinates[1])
	dz := float64(p1.Coordinates[2] - p2.Coordinates[2])
	return sqrt(dx*dx + dy*dy + dz*dz)
}

// Вспомогательная функция (в Go нет встроенной math.Sqrt для float64 без импорта)
func sqrt(x float64) float64 {
	if x < 0 { return 0 }
	z := x
	for i := 0; i < 20; i++ {
		z = (z + x/z) / 2
	}
	return z
}
