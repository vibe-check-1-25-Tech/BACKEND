package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"vibe-check-backend/models"
)

func main() {
	// 1. Маршрут для создания новой записи (POST /api/logs)
	http.HandleFunc("/api/logs", handleCreateLog)

	// 2. Маршрут для проверки, что сервер живой (просто GET /)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Backend Vibe Check работает!")
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	// Запуск сервера на порту 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}

// Функция-обработчик для создания записи
func handleCreateLog(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что к нам пришли методом POST
	if r.Method != http.MethodPost {
		http.Error(w, "Разрешен только метод POST", http.StatusMethodNotAllowed)
		return
	}

	var entry models.MoodLog

	// Декодируем JSON, который прислал фронтенд
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Ошибка в формате данных: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация (твоя работа как бэкенда)
	if entry.Score < 1 || entry.Score > 5 {
		http.Error(w, "Оценка должна быть от 1 до 5", http.StatusBadRequest)
		return
	}

	// Ставим текущее время
	entry.Timestamp = time.Now()

	// Пока базы нет, просто выводим результат в терминал
	fmt.Printf("--- Новая запись ---\n")
	fmt.Printf("Пользователь ID: %d\n", entry.UserID)
	fmt.Printf("Настроение: %d\n", entry.Score)
	fmt.Printf("Заметка: %s\n", entry.Note)
	fmt.Printf("Время: %v\n", entry.Timestamp)

	// Отвечаем фронтенду, что всё успешно
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Mood log saved (mock)",
	})
}
