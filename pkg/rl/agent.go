package rl

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// State описывает состояние пользователя для RL
type State struct {
	Vectors    [3][3]int // Векторы дня/месяца/года
	Chakras    []int     // Активные чакры (0-6)
	Symptoms   []string  // Выбранные симптомы
	History    []Action  // История предыдущих действий
}

// Action — действие агента (какую аффирмацию показать)
type Action struct {
	ID          string
	Text        string
	Author      string
	ChakraIndex int
	Reason      string // Объяснение выбора (для доверия)
}

// Reward — типы наград для обучения
type Reward float64

const (
	RewardCopied   Reward = 1.0 // Клиент скопировал текст
	RewardPrinted  Reward = 2.0 // Клиент распечатал
	RewardReturned Reward = 5.0 // Клиент вернулся через неделю
	RewardIgnored  Reward = -1.0// Клиент проигнорировал
)

// Agent — Q-learning агент
type Agent struct {
	qTable       map[string]map[string]float64 // stateKey -> actionID -> Q-value
	mu           sync.RWMutex
	learningRate float64
	discount     float64
	exploration  float64 // epsilon-greedy
}

// NewAgent создаёт нового агента
func NewAgent(lr, discount, exploration float64) *Agent {
	return &Agent{
		qTable:       make(map[string]map[string]float64),
		learningRate: lr,
		discount:     discount,
		exploration:  exploration,
	}
}

// stateKey генерирует строковый ключ из State (для таблицы Q)
func stateKey(s State) string {
	// Упрощённо: хэшируем векторы + чакры
	return fmt.Sprintf("%v_%v", s.Vectors, s.Chakras)
}

// ChooseAction выбирает действие (epsilon-greedy + объяснение)
func (a *Agent) ChooseAction(s State, available []Action) Action {
	a.mu.RLock()
	defer a.mu.RUnlock()

	key := stateKey(s)
	
	// С вероятностью exploration — случайное действие (исследование)
	if rand.Float64() < a.exploration {
		act := available[rand.Intn(len(available))]
		act.Reason = "🎲 Исследование нового паттерна"
		return act
	}

	// Иначе — лучшее известное действие (эксплуатация)
	bestQ := -math.MaxFloat64
	var best Action
	for _, act := range available {
		q := a.qTable[key][act.ID]
		if q > bestQ {
			bestQ = q
			best = act
		}
	}
	
	if bestQ > 0 {
		best.Reason = fmt.Sprintf("✅ Успешно в %d похожих случаях", int(bestQ*10))
	} else {
		best.Reason = "🆕 Новый паттерн (нет истории)"
	}
	return best
}

// Learn обновляет Q-таблицу на основе награды
func (a *Agent) Learn(s State, action Action, reward Reward, nextS State) {
	a.mu.Lock()
	defer a.mu.Unlock()

	key := stateKey(s)
	nextKey := stateKey(nextS)

	// Инициализация, если нужно
	if a.qTable[key] == nil {
		a.qTable[key] = make(map[string]float64)
	}
	if a.qTable[nextKey] == nil {
		a.qTable[nextKey] = make(map[string]float64)
	}

	// Q-learning update: Q(s,a) += lr * (r + γ*max_a' Q(s',a') - Q(s,a))
	currentQ := a.qTable[key][action.ID]
	
	// Max Q for next state
	maxNextQ := -math.MaxFloat64
	for _, q := range a.qTable[nextKey] {
		if q > maxNextQ {
			maxNextQ = q
		}
	}
	if maxNextQ == -math.MaxFloat64 {
		maxNextQ = 0
	}

	newQ := currentQ + a.learningRate*(float64(reward)+a.discount*maxNextQ-currentQ)
	a.qTable[key][action.ID] = newQ
}

// ExportExperience экспортирует опыт для gossip-обмена (только агрегированные паттерны)
func (a *Agent) ExportExperience() []Experience {
	// В реальности: отправлять только статистики, не сырые данные
	return []Experience{}
}

// Experience — агрегированный опыт для обмена между узлами
type Experience struct {
	StateHash   string  // Хэш состояния (не сами данные)
	ActionID    string
	AvgReward   float64
	Count       int
	Timestamp   time.Time
}
