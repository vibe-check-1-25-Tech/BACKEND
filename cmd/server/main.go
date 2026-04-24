package main

import (
	"database/sql"
	"log"
	"net/http"
	"vibe-check-backend/internal/handlers"
	"vibe-check-backend/internal/repository"

	_ "github.com/go-sql-driver/mysql"
)

// Наша функция-обертка для CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любых адресов
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Разрешаем нужные методы
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		// Разрешаем отправку JSON в заголовках
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Если это предзапрос OPTIONS от браузера — сразу отвечаем 200 OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаем управление дальше в основной хендлер
		next.ServeHTTP(w, r)
	}
}

func main() {
	// Подключение для XAMPP
	connStr := "root:@tcp(127.0.0.1:3306)/vibe_check?parseTime=true"

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Ошибка настроек базы: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("База не отвечает! Проверь XAMPP: ", err)
	}

	repo := repository.NewMoodRepository(db)
	env := &handlers.Env{Repo: repo}

	// Оборачиваем каждый маршрут в corsMiddleware
	http.HandleFunc("/api/ping", corsMiddleware(handlers.PingHandler))
	http.HandleFunc("/api/mood", corsMiddleware(env.CreateLogHandler))
	http.HandleFunc("/api/history", corsMiddleware(env.GetHistoryHandler))

	log.Println("🚀 Сервер взлетает на MySQL с поддержкой CORS! http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
