package handlers

import (
	"encoding/json"
	"net/http"
	"vibe-check-backend/internal/models"
	"vibe-check-backend/internal/repository"
)

type Env struct {
	Repo *repository.MoodRepository
}

// Вспомогательная функция для CORS
func setupCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Если это предзапрос (OPTIONS), отвечаем 200 и возвращаем true
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok", "message": "Сервер Аэлиты готов!"}`))
}

func (e *Env) CreateLogHandler(w http.ResponseWriter, r *http.Request) {
	// Сначала проверяем CORS
	if setupCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var entry models.MoodLog
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Ошибка в формате JSON", http.StatusBadRequest)
		return
	}

	if entry.Score < 1 || entry.Score > 5 {
		http.Error(w, "Оценка должна быть от 1 до 5", http.StatusBadRequest)
		return
	}
	if entry.Note == "" {
		http.Error(w, "Заметка не может быть пустой", http.StatusBadRequest)
		return
	}

	err := e.Repo.SaveMood(entry)
	if err != nil {
		http.Error(w, "Ошибка при сохранении в базу", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Данные сохранены!"})
}

func (e *Env) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w, r)

	logs, err := e.Repo.GetAllMoods()
	if err != nil {
		http.Error(w, "Ошибка базы", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
