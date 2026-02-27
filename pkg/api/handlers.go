package api

import (
	"encoding/json"
	"ideal-core/pkg/auth"
	"ideal-core/pkg/cube"
	"ideal-core/pkg/db"
	"ideal-core/pkg/identity"
	"net/http"
	"time"
)

type Handler struct {
	DB   *db.Database
	Auth *auth.AuthManager
}

func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name      string `json:"name"`
		BirthDate string `json:"birth_date"`
		Action    string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	publicID := identity.GenerateID(birthDate, "ideal-core-v1")

	if req.Action == "register" {
		if err := h.DB.CreateUser(db.User{
			ID: publicID, Name: req.Name, BirthDate: req.BirthDate,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	token, err := h.Auth.Login(publicID)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": publicID,
		"token":   token,
		"name":    req.Name,
	})
}

func (h *Handler) AddPersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	session, ok := h.Auth.Validate(token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name      string `json:"name"`
		BirthDate string `json:"birth_date"`
		Relation  string `json:"relation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	coords := cube.CalcCoordinates(birthDate)
	sumFreq := cube.GetSumFrequency(coords)

	person := db.Person{
		ID:         identity.GenerateID(birthDate, "person-"+session.UserID),
		UserID:     session.UserID,
		Name:       req.Name,
		BirthDate:  req.BirthDate,
		Coords:     formatCoords(coords),
		SumFreq:    sumFreq,
		FlowStatus: "Active",
	}

	if err := h.DB.AddPerson(person); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func (h *Handler) ListPeopleHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	session, ok := h.Auth.Validate(token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	people, err := h.DB.GetPeopleByUser(session.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func formatCoords(c [3]int) string {
	return string(rune(c[0])) + "," + string(rune(c[1])) + "," + string(rune(c[2]))
}
