// Package bio — биомаркеры и лабораторные анализы
package bio

import (
	"time"
)

// Biomarker — базовая структура биомаркера
type Biomarker struct {
	ID              string    `json:"id" db:"id"`
	Code            string    `json:"code" db:"code"`
	Name            string    `json:"name" db:"name"`
	NameEn          string    `json:"name_en" db:"name_en"`
	Category        string    `json:"category" db:"category"`
	Description     string    `json:"description" db:"description"`
	Function        string    `json:"function" db:"function"`
	Unit            string    `json:"unit" db:"unit"`
	SampleType      string    `json:"sample_type" db:"sample_type"`
	FastingRequired bool      `json:"fasting_required" db:"fasting_required"`
	TimeOfDay       string    `json:"time_of_day" db:"time_of_day"`
	CycleDependent  bool      `json:"cycle_dependent" db:"cycle_dependent"`
	SexSpecific     bool      `json:"sex_specific" db:"sex_specific"`
	AgeSpecific     bool      `json:"age_specific" db:"age_specific"`
	HalfLife        float64   `json:"half_life" db:"half_life"`
	DiurnalRhythm   bool      `json:"diurnal_rhythm" db:"diurnal_rhythm"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ReferenceRange — референсный диапазон для биомаркера
type ReferenceRange struct {
	BiomarkerCode string    `json:"biomarker_code" db:"biomarker_code"`
	RangeType     string    `json:"range_type" db:"range_type"`
	AgeMin        int       `json:"age_min" db:"age_min"`
	AgeMax        int       `json:"age_max" db:"age_max"`
	Sex           string    `json:"sex" db:"sex"`
	CyclePhase    string    `json:"cycle_phase" db:"cycle_phase"`
	Min           float64   `json:"min" db:"min"`
	Max           float64   `json:"max" db:"max"`
	Unit          string    `json:"unit" db:"unit"`
	Source        string    `json:"source" db:"source"`
	Notes         string    `json:"notes" db:"notes"`
	LastUpdated   time.Time `json:"last_updated" db:"last_updated"`
}

// BiomarkerStatus — статус результата биомаркера
type BiomarkerStatus string

const (
	StatusCriticalLow  BiomarkerStatus = "critical_low"
	StatusLow          BiomarkerStatus = "low"
	StatusOptimal      BiomarkerStatus = "optimal"
	StatusHigh         BiomarkerStatus = "high"
	StatusCriticalHigh BiomarkerStatus = "critical_high"
)

// GetBiomarkerStatus определяет статус результата
func GetBiomarkerStatus(value, min, max, optMin, optMax float64) BiomarkerStatus {
	if value < min*0.5 {
		return StatusCriticalLow
	}
	if value > max*2 {
		return StatusCriticalHigh
	}
	if value < min {
		return StatusLow
	}
	if value > max {
		return StatusHigh
	}
	if optMin > 0 && optMax > 0 {
		if value >= optMin && value <= optMax {
			return StatusOptimal
		}
		if value < optMin {
			return StatusLow
		}
		if value > optMax {
			return StatusHigh
		}
	}
	return StatusOptimal
}

// CalculatePercentile вычисляет процентиль результата
func CalculatePercentile(value, min, max float64) float64 {
	if max == min {
		return 50.0
	}
	percentile := ((value - min) / (max - min)) * 100
	if percentile < 0 {
		return 0
	}
	if percentile > 100 {
		return 100
	}
	return percentile
}

// GetBiomarkersByCategory возвращает биомаркеры по категории
func GetBiomarkersByCategory(category string) []Biomarker {
	allBiomarkers := LoadAllBiomarkers()
	var result []Biomarker
	for _, bm := range allBiomarkers {
		if bm.Category == category {
			result = append(result, bm)
		}
	}
	return result
}

// LoadAllBiomarkers загружает все биомаркеры из справочника
func LoadAllBiomarkers() []Biomarker {
	return loadBiomarkersFromLabTests()
}

func loadBiomarkersFromLabTests() []Biomarker {
	return []Biomarker{}
}
