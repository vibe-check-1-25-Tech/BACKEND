package main

import (
	"database/sql"
	"log"
	"net/http"
	"vibe-check-backend/internal/handlers"
	"vibe-check-backend/internal/repository"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "user=postgres password=pass dbname=vibe_check sslmode=disable"
	db, _ := sql.Open("postgres", connStr)

	repo := repository.NewMoodRepository(db)
	env := &handlers.Env{Repo: repo}

	http.HandleFunc("/api/ping", handlers.PingHandler)
	http.HandleFunc("/api/mood", env.CreateLogHandler)
	http.HandleFunc("/api/history", env.GetHistoryHandler)

	log.Println("🚀 Сервер взлетает на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска: ", err)
	}
}
