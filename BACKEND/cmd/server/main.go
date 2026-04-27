package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	handlers "vibe-check-backend/api/v1"
	"vibe-check-backend/internal/repository"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/vibe_check?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewMoodRepository(db)
	// Теперь env берется из пакета v1 (которому мы дали имя handlers выше)
	env := &handlers.Env{Repo: repo}

	// Маршруты остаются без изменений, так как мы сохранили имя 'handlers' в импорте
	http.HandleFunc("/api/register", corsMiddleware(env.RegisterHandler))
	http.HandleFunc("/api/login", corsMiddleware(env.LoginHandler))
	http.HandleFunc("/api/logs", corsMiddleware(env.GetMoodsHandler))
	http.HandleFunc("/api/logs/save", corsMiddleware(env.CreateMoodHandler))
	http.HandleFunc("/api/search", corsMiddleware(env.SearchHandler))
	http.HandleFunc("/api/stats", corsMiddleware(env.GetStatsHandler))
	http.HandleFunc("/api/support", corsMiddleware(env.GetSupportContent))
	http.HandleFunc("/api/team/aggregate", corsMiddleware(env.GetTeamStatsHandler))
	http.HandleFunc("/api/tags/top", corsMiddleware(env.GetTopTagsHandler))
	http.HandleFunc("/api/user/reminders", corsMiddleware(env.SetReminderHandler))
	http.HandleFunc("/api/export", corsMiddleware(env.ExportHandler))
	http.HandleFunc("/", corsMiddleware(env.NotFoundHandler))
	http.HandleFunc("/api/ping", corsMiddleware(env.PingHandler))

	// Фоновый процесс
	go func() {
		for {
			currentTime := time.Now().Format("15:04")
			fmt.Println("Проверка уведомлений... Время:", currentTime)
			time.Sleep(1 * time.Minute)
		}
	}()

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method == "OPTIONS" {
			return
		}
		next(w, r)
	}
}
