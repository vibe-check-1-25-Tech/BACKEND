package v1

import (
	"encoding/json"
	"net/http"
	"vibe-check-backend/internal/models"
	"vibe-check-backend/internal/repository"
)

type Env struct {
	Repo *repository.MoodRepository
}

// ======================
// MOODS (у тебя уже есть)
// ======================

func (e *Env) MoodsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		logs, err := e.Repo.GetAllMoods("", "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to fetch moods",
			})
			return
		}
		json.NewEncoder(w).Encode(logs)

	case http.MethodPost:
		defer r.Body.Close()

		var mood models.MoodLog

		if err := json.NewDecoder(r.Body).Decode(&mood); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid request body",
			})
			return
		}

		err := e.Repo.SaveMood(mood)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to save mood",
			})
			return
		}

		needsSupport, content := e.Repo.CheckIfUserNeedsSupport(mood.UserID)

		response := map[string]interface{}{
			"status": "success",
		}

		if needsSupport {
			response["support"] = content
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "method not allowed",
		})
	}
}

// ======================
// AUTH (заглушки)
// ======================

func (e *Env) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "register ok",
	})
}

func (e *Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "login ok",
	})
}

// ======================
// SEARCH / SUPPORT
// ======================

func (e *Env) SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "search ok",
	})
}

func (e *Env) GetSupportContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "support ok",
	})
}

// ======================
// STATS
// ======================

func (e *Env) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "stats ok",
	})
}

func (e *Env) GetTopTagsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "top tags ok",
	})
}

func (e *Env) GetTeamStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "team stats ok",
	})
}

// ======================
// REMINDERS
// ======================

func (e *Env) SetReminderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "reminder ok",
	})
}

// ======================
// EXPORT
// ======================

func (e *Env) ExportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "export ok",
	})
}

// ======================
// SYSTEM
// ======================

func (e *Env) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (e *Env) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
