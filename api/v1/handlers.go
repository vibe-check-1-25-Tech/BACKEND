package v1

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"vibe-check-backend/internal/models"
	"vibe-check-backend/internal/repository"
)

type Env struct {
	Repo *repository.MoodRepository
}

// --- Обработчики настроения ---

// CreateMoodHandler сохраняет новое настроение и проверяет необходимость поддержки
func (e *Env) CreateMoodHandler(w http.ResponseWriter, r *http.Request) {
	var log models.MoodLog
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 1. Сохраняем настроение в базу данных
	e.Repo.SaveMood(log)

	// 2. Проверяем состояние пользователя (автоматическая поддержка)
	needsSupport, content := e.Repo.CheckIfUserNeedsSupport(log.UserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status": "success",
	}

	// Если логика репозитория решила, что пользователю грустно, прикрепляем мем сразу
	if needsSupport {
		response["support"] = content
	}

	json.NewEncoder(w).Encode(response)
}

// GetMoodsHandler возвращает историю всех записей
func (e *Env) GetMoodsHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := e.Repo.GetAllMoods("", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// --- Аналитика и статистика ---

func (e *Env) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := e.Repo.GetMoodStats()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (e *Env) GetTopTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := e.Repo.GetTopTags()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (e *Env) GetTeamStatsHandler(w http.ResponseWriter, r *http.Request) {
	avg, err := e.Repo.GetTeamAverage(1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	ok, err := e.Repo.CheckPin(c.Email, c.Pin)
	if err != nil || !ok {
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
	logs, err := e.Repo.SearchNotes(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetSupportContent — хендлер для ручного запроса мема (например, кнопка "Поддержать")
func (e *Env) GetSupportContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Используем нашу обновленную функцию из models.go
	content := models.GetRandomSupport()
	json.NewEncoder(w).Encode(content)
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

	// 1. Обработка PDF (заглушка)
	if strings.Contains(r.URL.Path, "pdf") {
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=vibe_check_report.pdf")
		w.Write([]byte("Отчет Vibe Check: PDF в разработке. Используйте CSV для проверки данных."))
		return
	}

	// 2. Обработка CSV
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=vibe_check_export.csv")

	// Добавляем UTF-8 BOM для корректного открытия русского языка в Excel
	w.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(w)
	writer.Comma = ';'
	defer writer.Flush()

	// Заголовки таблицы
	writer.Write([]string{"Оценка", "Заметка"})

	// Наполнение данными
	for _, l := range logs {
		row := []string{
			fmt.Sprintf("%d", l.Score),
			l.Note,
		}
		writer.Write(row)
	}
}

// --- Служебные функции ---

func (e *Env) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Страница не найдена"))
}

func (e *Env) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "pong"}`))
}
