package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BirthDate string    `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
}

type Person struct {
	ID         string    `json:"person_id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	BirthDate  string    `json:"birth_date"`
	Coords     string    `json:"coords"`
	SumFreq    int       `json:"sum_freq"`
	FlowStatus string    `json:"flow_status"`
	CreatedAt  time.Time `json:"created_at"`
}

type PersonSymptom struct {
	PersonID    string    `json:"person_id"`
	SymptomKey  string    `json:"symptom_key"`
	CustomLabel string    `json:"custom_label,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Feedback struct {
	ID              int       `json:"id"`
	PersonID        string    `json:"person_id"`
	IntentionHash   string    `json:"intention_hash"`
	Action          string    `json:"action"`
	Timestamp       time.Time `json:"timestamp"`
	Reward          float64   `json:"reward"`
}

type Affirmation struct {
	ID            int      `json:"id"`
	Author        string   `json:"author"`
	OldThought    string   `json:"old_thought"`
	NewThought    string   `json:"new_thought"`
	ChakraIndex   int      `json:"chakra_index"`
	Symptoms      []string `json:"symptoms"`
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Основные таблицы
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			birth_date TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS people (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			birth_date TEXT NOT NULL,
			coords TEXT,
			sum_freq INTEGER,
			flow_status TEXT DEFAULT 'Active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS person_symptoms (
			person_id TEXT REFERENCES people(id),
			symptom_key TEXT,
			custom_label TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (person_id, symptom_key)
		)`,
		`CREATE TABLE IF NOT EXISTS feedback (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			person_id TEXT,
			intention_hash TEXT,
			action TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			reward REAL
		)`,
		`CREATE TABLE IF NOT EXISTS affirmations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT,
			old_thought TEXT,
			new_thought TEXT,
			chakra_index INTEGER,
			symptoms TEXT
		)`,
	}

	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			return nil, err
		}
	}

	log.Println("✅ DB initialized:", path)
	return &Database{db: db}, nil
}

func (d *Database) CreateUser(u User) error {
	_, err := d.db.Exec(
		"INSERT OR REPLACE INTO users (id, name, birth_date) VALUES (?, ?, ?)",
		u.ID, u.Name, u.BirthDate,
	)
	return err
}

func (d *Database) GetUser(id string) (*User, error) {
	var u User
	err := d.db.QueryRow(
		"SELECT id, name, birth_date, created_at FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Name, &u.BirthDate, &u.CreatedAt)
	return &u, err
}

func (d *Database) AddPerson(p Person) error {
	_, err := d.db.Exec(
		`INSERT OR REPLACE INTO people 
		(id, user_id, name, birth_date, coords, sum_freq, flow_status) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.UserID, p.Name, p.BirthDate, p.Coords, p.SumFreq, p.FlowStatus,
	)
	return err
}

func (d *Database) GetPeopleByUser(userID string) ([]Person, error) {
	rows, err := d.db.Query(
		"SELECT id, user_id, name, birth_date, coords, sum_freq, flow_status, created_at FROM people WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var p Person
		err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.BirthDate, &p.Coords, &p.SumFreq, &p.FlowStatus, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		people = append(people, p)
	}
	return people, nil
}

func (d *Database) AddSymptom(ps PersonSymptom) error {
	_, err := d.db.Exec(
		"INSERT OR REPLACE INTO person_symptoms (person_id, symptom_key, custom_label) VALUES (?, ?, ?)",
		ps.PersonID, ps.SymptomKey, ps.CustomLabel,
	)
	return err
}

func (d *Database) GetSymptoms(personID string) ([]PersonSymptom, error) {
	rows, err := d.db.Query(
		"SELECT person_id, symptom_key, custom_label, created_at FROM person_symptoms WHERE person_id = ?",
		personID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var symptoms []PersonSymptom
	for rows.Next() {
		var s PersonSymptom
		err := rows.Scan(&s.PersonID, &s.SymptomKey, &s.CustomLabel, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		symptoms = append(symptoms, s)
	}
	return symptoms, nil
}

func (d *Database) AddFeedback(f Feedback) error {
	_, err := d.db.Exec(
		"INSERT INTO feedback (person_id, intention_hash, action, reward) VALUES (?, ?, ?, ?)",
		f.PersonID, f.IntentionHash, f.Action, f.Reward,
	)
	return err
}

func (d *Database) AddAffirmation(a Affirmation) error {
	_, err := d.db.Exec(
		"INSERT INTO affirmations (author, old_thought, new_thought, chakra_index, symptoms) VALUES (?, ?, ?, ?, ?)",
		a.Author, a.OldThought, a.NewThought, a.ChakraIndex, a.Symptoms,
	)
	return err
}

func (d *Database) Close() {
	d.db.Close()
}
