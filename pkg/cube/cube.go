package cube

import "time"

// Coordinates — координаты комнаты (X=день, Y=месяц, Z=год)
type Coordinates [3]int

// Vectors — векторы движения по трём осям
type Vectors [3][3]int

// Person — сущность человека в системе Куба
type Person struct {
	ID           string
	Name         string
	BirthDate    time.Time
	Coordinates  Coordinates
	Vectors      Vectors
	SumFrequency int
}

// SumDigits — сумма цифр числа
func SumDigits(n int) int {
	if n < 0 {
		n = -n
	}
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// CalcCoordinates — расчёт координат по дате рождения
func CalcCoordinates(date time.Time) Coordinates {
	day := SumDigits(date.Day())
	month := SumDigits(int(date.Month()))
	year := SumDigits(date.Year() % 1000)
	return Coordinates{day, month, year}
}

// CalcVectors — расчёт векторов движения (разности цифр)
func CalcVectors(date time.Time) Vectors {
	calcVec := func(n int) [3]int {
		d := [3]int{n / 100, (n / 10) % 10, n % 10}
		return [3]int{
			d[0] - d[1],
			d[1] - d[2],
			d[2] - d[0],
		}
	}
	return Vectors{
		calcVec(date.Day()),
		calcVec(int(date.Month())),
		calcVec(date.Year() % 1000),
	}
}

// GetSumFrequency — общая частота Σ
func GetSumFrequency(c Coordinates) int {
	return c[0] + c[1] + c[2]
}

// NewPerson — конструктор Person
func NewPerson(id, name string, birthDate time.Time) *Person {
	coords := CalcCoordinates(birthDate)
	vectors := CalcVectors(birthDate)
	sum := GetSumFrequency(coords)
	return &Person{
		ID:           id,
		Name:         name,
		BirthDate:    birthDate,
		Coordinates:  coords,
		Vectors:      vectors,
		SumFrequency: sum,
	}
}

// Distance — евклидово расстояние между двумя людьми (квадрат расстояния)
func Distance(p1, p2 Person) float64 {
	dx := float64(p1.Coordinates[0] - p2.Coordinates[0])
	dy := float64(p1.Coordinates[1] - p2.Coordinates[1])
	dz := float64(p1.Coordinates[2] - p2.Coordinates[2])
	return dx*dx + dy*dy + dz*dz
}

// EuclideanDistance — алиас для Distance
func EuclideanDistance(c1, c2 Coordinates) float64 {
	dx := float64(c1[0] - c2[0])
	dy := float64(c1[1] - c2[1])
	dz := float64(c1[2] - c2[2])
	return dx*dx + dy*dy + dz*dz
}
