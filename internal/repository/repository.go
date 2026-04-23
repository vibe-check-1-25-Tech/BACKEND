package repository

import (
	"database/sql"
	"time"
	"vibe-check-backend/internal/models" // Вот это было пропущено!
)

type MoodRepository struct {
	DB *sql.DB
}

func NewMoodRepository(db *sql.DB) *MoodRepository {
	return &MoodRepository{DB: db}
}

// Переименовали в SaveMood, как мы договаривались в хендлерах
func (r *MoodRepository) SaveMood(log models.MoodLog) error {
	query := `INSERT INTO mood_logs (user_id, score, note, timestamp) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, log.UserID, log.Score, log.Note, time.Now())
	return err
}

// Добавили этот метод, чтобы handlers.go не ругался
func (r *MoodRepository) GetAllMoods() ([]models.MoodLog, error) {
	rows, err := r.DB.Query("SELECT id, user_id, score, note, timestamp FROM mood_logs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []models.MoodLog
	for rows.Next() {
		var l models.MoodLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Score, &l.Note, &l.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
