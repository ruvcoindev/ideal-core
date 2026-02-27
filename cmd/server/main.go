package main

import (
	"ideal-core/pkg/api"
	"ideal-core/pkg/auth"
	"ideal-core/pkg/db"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	database, err := db.NewDatabase("data/ideal.db")
	if err != nil {
		log.Fatal("❌ DB error:", err)
	}
	defer database.Close()

	authManager := auth.NewAuthManager(24 * time.Hour)

	handler := &api.Handler{
		DB:   database,
		Auth: authManager,
	}

	http.HandleFunc("/api/auth", handler.AuthHandler)
	http.HandleFunc("/api/people", handler.ListPeopleHandler)
	http.HandleFunc("/api/people/add", handler.AddPersonHandler)

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server running on port", port)
	log.Println("🌐 Open: http://localhost:" + port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}
