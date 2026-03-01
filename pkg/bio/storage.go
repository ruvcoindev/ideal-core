// Package bio — интерфейс хранилища биомаркеров
package bio

import (
	"context"
	"time"
)

// BioStore — интерфейс хранилища биомаркеров
type BioStore interface {
	SaveResult(ctx context.Context, result UserLabResult) error
	GetResultsByUser(ctx context.Context, userID string) ([]UserLabResult, error)
	GetResultsByTest(ctx context.Context, userID string, testID string) ([]UserLabResult, error)
	GetLatestResult(ctx context.Context, userID string, testID string) (*UserLabResult, error)
	GetResultsByDateRange(ctx context.Context, userID string, start, end time.Time) ([]UserLabResult, error)
	GetResultsByStatus(ctx context.Context, userID string, status BiomarkerStatus) ([]UserLabResult, error)
	GetResultsByChakra(ctx context.Context, userID string, chakraIndex int) ([]UserLabResult, error)
	DeleteResult(ctx context.Context, resultID string) error
	GetTrends(ctx context.Context, userID string, testID string, period time.Duration) ([]TrendPoint, error)
	UpdateInterpretation(ctx context.Context, resultID string, interpretation ResultInterpretation) error
	AddRecommendation(ctx context.Context, resultID string, recommendation Recommendation) error
	Close() error
}

// BioStoreConfig — конфигурация хранилища
type BioStoreConfig struct {
	Driver         string `json:"driver"`
	DataSource     string `json:"data_source"`
	EncryptionKey  string `json:"encryption_key"`
	MaxConnections int    `json:"max_connections"`
	Timeout        int    `json:"timeout"`
}

// NewBioStore создаёт новое хранилище
func NewBioStore(config BioStoreConfig) (BioStore, error) {
	switch config.Driver {
	case "sqlite":
		return newSQLiteStore(config)
	case "postgres":
		return newPostgresStore(config)
	case "memory":
		return newMemoryStore(config)
	default:
		return nil, ErrUnsupportedDriver
	}
}

// ErrUnsupportedDriver — ошибка неподдерживаемого драйвера
var ErrUnsupportedDriver = &BioError{
	Code:    "UNSUPPORTED_DRIVER",
	Message: "Неподдерживаемый драйвер хранилища",
}

// BioError — ошибка операции с биомаркерами
type BioError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *BioError) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Message + " (" + e.Err.Error() + ")"
	}
	return e.Code + ": " + e.Message
}

func (e *BioError) Unwrap() error {
	return e.Err
}

func newSQLiteStore(config BioStoreConfig) (BioStore, error) {
	return nil, nil
}

func newPostgresStore(config BioStoreConfig) (BioStore, error) {
	return nil, nil
}

func newMemoryStore(config BioStoreConfig) (BioStore, error) {
	return &memoryStore{
		results: make(map[string]UserLabResult),
	}, nil
}

// memoryStore — in-memory реализация для тестов
type memoryStore struct {
	results map[string]UserLabResult
}

func (s *memoryStore) SaveResult(ctx context.Context, result UserLabResult) error {
	s.results[result.ID] = result
	return nil
}

func (s *memoryStore) GetResultsByUser(ctx context.Context, userID string) ([]UserLabResult, error) {
	var results []UserLabResult
	for _, r := range s.results {
		if r.UserID == userID {
			results = append(results, r)
		}
	}
	return results, nil
}

func (s *memoryStore) GetResultsByTest(ctx context.Context, userID string, testID string) ([]UserLabResult, error) {
	var results []UserLabResult
	for _, r := range s.results {
		if r.UserID == userID && r.TestID == testID {
			results = append(results, r)
		}
	}
	return results, nil
}

func (s *memoryStore) GetLatestResult(ctx context.Context, userID string, testID string) (*UserLabResult, error) {
	results, err := s.GetResultsByTest(ctx, userID, testID)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return &results[len(results)-1], nil
}

func (s *memoryStore) GetResultsByDateRange(ctx context.Context, userID string, start, end time.Time) ([]UserLabResult, error) {
	results, err := s.GetResultsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var filtered []UserLabResult
	for _, r := range results {
		if r.CreatedAt.After(start) && r.CreatedAt.Before(end) {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

func (s *memoryStore) GetResultsByStatus(ctx context.Context, userID string, status BiomarkerStatus) ([]UserLabResult, error) {
	results, err := s.GetResultsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return GetResultsByStatus(results, status), nil
}

func (s *memoryStore) GetResultsByChakra(ctx context.Context, userID string, chakraIndex int) ([]UserLabResult, error) {
	results, err := s.GetResultsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return GetResultsByChakra(results, chakraIndex), nil
}

func (s *memoryStore) DeleteResult(ctx context.Context, resultID string) error {
	delete(s.results, resultID)
	return nil
}

func (s *memoryStore) GetTrends(ctx context.Context, userID string, testID string, period time.Duration) ([]TrendPoint, error) {
	results, err := s.GetResultsByTest(ctx, userID, testID)
	if err != nil {
		return nil, err
	}
	var points []TrendPoint
	for _, r := range results {
		points = append(points, TrendPoint{
			Date:   r.CreatedAt,
			Value:  r.Value,
			Status: r.Status,
		})
	}
	return points, nil
}

func (s *memoryStore) UpdateInterpretation(ctx context.Context, resultID string, interpretation ResultInterpretation) error {
	result, ok := s.results[resultID]
	if !ok {
		return &BioError{Code: "NOT_FOUND", Message: "Результат не найден"}
	}
	result.Interpretation = &interpretation
	result.Interpreted = true
	s.results[resultID] = result
	return nil
}

func (s *memoryStore) AddRecommendation(ctx context.Context, resultID string, recommendation Recommendation) error {
	result, ok := s.results[resultID]
	if !ok {
		return &BioError{Code: "NOT_FOUND", Message: "Результат не найден"}
	}
	result.Recommendations = append(result.Recommendations, recommendation)
	s.results[resultID] = result
	return nil
}

func (s *memoryStore) Close() error {
	return nil
}
