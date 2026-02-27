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
			coords TEXT,
			sum_freq INTEGER,
			flow_status TEXT DEFAULT 'Active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
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

func (d *Database) Close() {
	d.db.Close()
}
