package repository

import (
	"database/sql"
	"time"
	"vibe-check-backend/internal/models"
)

type MoodRepository struct {
	DB *sql.DB
}

func NewMoodRepository(db *sql.DB) *MoodRepository {
	return &MoodRepository{DB: db}
}

// Сохранение настроения
func (r *MoodRepository) SaveMood(log models.MoodLog) error {
	// ВАЖНО: названия колонок должны совпадать с таблицей напарника (user_id, score, note)
	query := "INSERT INTO Mood_Logs (user_id, score, note, timestamp) VALUES (?, ?, ?, ?)"
	_, err := r.DB.Exec(query, log.UserID, log.Score, log.Note, time.Now())
	return err
}

// Получение всей истории
func (r *MoodRepository) GetAllMoods() ([]models.MoodLog, error) {
	// Используем Mood_Logs (с большой буквы, как в базе)
	rows, err := r.DB.Query("SELECT log_id, user_id, score, note, timestamp FROM Mood_Logs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.MoodLog
	for rows.Next() {
		var l models.MoodLog
		// Сканируем данные в структуру. log_id идет в l.ID
		if err := rows.Scan(&l.ID, &l.UserID, &l.Score, &l.Note, &l.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
