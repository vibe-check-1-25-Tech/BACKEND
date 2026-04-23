package handlers

import (
	"encoding/json"
	"net/http"
	"vibe-check-backend/internal/models"
	"vibe-check-backend/internal/repository"
)

// Определение структуры Env должно быть здесь!
type Env struct {
	Repo *repository.MoodRepository
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok", "message": "Сервер Аэлиты готов!"}`))
}

func (e *Env) CreateLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var entry models.MoodLog
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Ошибка в формате JSON", http.StatusBadRequest)
		return
	}

	// Валидация
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

// Добавим функцию истории, чтобы закрыть все задачи
func (e *Env) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := e.Repo.GetAllMoods()
	if err != nil {
		http.Error(w, "Ошибка базы", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
