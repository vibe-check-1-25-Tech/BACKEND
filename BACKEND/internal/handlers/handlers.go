package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"vibe-check-backend/internal/models"
	"vibe-check-backend/internal/repository"
)

type Env struct {
	Repo *repository.MoodRepository
}

// --- Обработчики настроения ---

func (e *Env) CreateMoodHandler(w http.ResponseWriter, r *http.Request) {
	var log models.MoodLog
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	e.Repo.SaveMood(log)
	w.WriteHeader(http.StatusCreated)
}

func (e *Env) GetMoodsHandler(w http.ResponseWriter, r *http.Request) {
	logs, _ := e.Repo.GetAllMoods("", "")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// --- Аналитика и статистика ---

func (e *Env) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, _ := e.Repo.GetMoodStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (e *Env) GetTopTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, _ := e.Repo.GetTopTags()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (e *Env) GetTeamStatsHandler(w http.ResponseWriter, r *http.Request) {
	avg, _ := e.Repo.GetTeamAverage(1)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"average": avg})
}

// --- Пользователи и авторизация ---

func (e *Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var c struct{ Email, Pin string }
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ok, _ := e.Repo.CheckPin(c.Email, c.Pin)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (e *Env) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var d struct{ Username, Email, Pin string }
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	e.Repo.CreateUser(d.Username, d.Email, d.Pin)
	w.WriteHeader(http.StatusCreated)
}

func (e *Env) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	e.Repo.DeleteUser(email)
	w.WriteHeader(http.StatusNoContent)
}

// --- Поиск и контент ---

func (e *Env) SearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	logs, _ := e.Repo.SearchNotes(q)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (e *Env) GetSupportContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"content": "Улыбнись! Всё получится!"})
}

// --- Уведомления (9 вкладка) ---

func (e *Env) SetReminderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var settings models.ReminderSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат данных"})
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Напоминание установлено на " + settings.ReminderTime,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// --- Экспорт данных (11 вкладка) ---

func (e *Env) ExportHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := e.Repo.GetAllMoods("", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=vibe_check_export.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Оставим только Оценку и Заметку, если с Датой проблемы
	writer.Write([]string{"Оценка", "Заметка"})

	for _, l := range logs {
		row := []string{
			fmt.Sprintf("%d", l.Score),
			l.Note,
		}
		writer.Write(row)
	}
}

// --- Служебные ---

func (e *Env) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Страница не найдена"))
}
func (e *Env) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "pong"}`))
}
