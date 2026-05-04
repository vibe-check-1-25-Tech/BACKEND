package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	// У тебя папка api в корне, а внутри v1
	handlers "vibe-check-backend/api/v1"

	// У тебя папка repository внутри internal
	"vibe-check-backend/internal/repository"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 1. Подключение к базе данных mood_tracker на порту 3307
	dsn := "root:rootroot@tcp(127.0.0.1:3307)/mood_tracker?parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatal("База данных недоступна:", err)
	}

	// 2. Инициализация репозитория и окружения хендлеров
	repo := repository.NewMoodRepository(db)
	env := &handlers.Env{Repo: repo}

	// 3. Настройка маршрутов
	mux := http.NewServeMux()

	// Регистрация и логин
	mux.HandleFunc("/api/register", corsMiddleware(env.RegisterHandler))
	mux.HandleFunc("/api/login", corsMiddleware(env.LoginHandler))

	// Работа с настроением и мемами
	mux.HandleFunc("/api/logs", corsMiddleware(env.GetMoodsHandler))
	mux.HandleFunc("/api/logs/save", corsMiddleware(env.CreateMoodHandler))
	mux.HandleFunc("/api/search", corsMiddleware(env.SearchHandler))
	mux.HandleFunc("/api/support", corsMiddleware(env.GetSupportContent))

	// Статистика
	mux.HandleFunc("/api/stats", corsMiddleware(env.GetStatsHandler))
	mux.HandleFunc("/api/tags/top", corsMiddleware(env.GetTopTagsHandler))
	mux.HandleFunc("/api/team/aggregate", corsMiddleware(env.GetTeamStatsHandler))

	// Напоминания и экспорт
	mux.HandleFunc("/api/user/reminders", corsMiddleware(env.SetReminderHandler))
	mux.HandleFunc("/api/export", corsMiddleware(env.ExportHandler))
	mux.HandleFunc("/api/export/csv", corsMiddleware(env.ExportHandler))

	// Обработка приглашений в команду
	mux.HandleFunc("/api/team/join", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		teamID := r.URL.Query().Get("team")
		if teamID == "" {
			http.Error(w, "Team ID is missing", http.StatusBadRequest)
			return
		}
		fmt.Printf("Пользователь вступает в команду: %s\n", teamID)
		http.Redirect(w, r, "http://127.0.0.1:5500/pages/team.html?joined="+teamID, http.StatusSeeOther)
	}))

	// Раздача статических файлов (твоих мемов)
	// Создай папку assets в корне проекта BACKEND и положи мемы туда
	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/api/ping", corsMiddleware(env.PingHandler))
	mux.HandleFunc("/", corsMiddleware(env.NotFoundHandler))

	// Фоновый процесс (заглушка для уведомлений)
	go func() {
		for {
			_ = time.Now().Format("15:04")
			time.Sleep(1 * time.Minute)
		}
	}()

	log.Println("✅ Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// corsMiddleware для обработки кросс-доменных запросов
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("[%s] %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
