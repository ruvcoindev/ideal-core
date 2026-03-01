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
	ID           string    `json:"person_id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	BirthDate    string    `json:"birth_date"`
	LastContact  time.Time `json:"last_contact,omitempty"` // 🔹 Новое поле для расчёта фаз
	Coords       string    `json:"coords"`
	SumFreq      int       `json:"sum_freq"`
	Vectors      string    `json:"vectors"`
	FlowStatus   string    `json:"flow_status"`
	PsychAge     int       `json:"psych_age,omitempty"`      // 🔹 Психовозраст (7/35/55)
	Location     string    `json:"location,omitempty"`       // 🔹 Локация: "ValyaHome", "Neutral"
	Tags         string    `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	Action          string    `json:"action"` // "copied", "printed", "ignored"
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
			last_contact DATETIME,
			coords TEXT,
			sum_freq INTEGER,
			vectors TEXT,
			flow_status TEXT DEFAULT 'Active',
			psych_age INTEGER,
			location TEXT,
			tags TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
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
		(id, user_id, name, birth_date, last_contact, coords, sum_freq, vectors, flow_status, psych_age, location, tags, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.UserID, p.Name, p.BirthDate, p.LastContact, p.Coords, p.SumFreq, p.Vectors, p.FlowStatus, p.PsychAge, p.Location, p.Tags, time.Now(),
	)
	return err
}

func (d *Database) GetPeopleByUser(userID string) ([]Person, error) {
	rows, err := d.db.Query(
		"SELECT id, user_id, name, birth_date, last_contact, coords, sum_freq, vectors, flow_status, psych_age, location, tags, created_at, updated_at FROM people WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var p Person
		var lastContact sql.NullTime
		err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.BirthDate, &lastContact, &p.Coords, &p.SumFreq, &p.Vectors, &p.FlowStatus, &p.PsychAge, &p.Location, &p.Tags, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if lastContact.Valid {
			p.LastContact = lastContact.Time
		}
		people = append(people, p)
	}
	return people, nil
}

func (d *Database) UpdateLastContact(personID string, lastContact time.Time) error {
	_, err := d.db.Exec(
		"UPDATE people SET last_contact = ?, updated_at = ? WHERE id = ?",
		lastContact, time.Now(), personID,
	)
	return err
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
