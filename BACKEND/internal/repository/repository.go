package repository

import (
	"database/sql"
	"fmt"
	"vibe-check-backend/internal/models"
)

type MoodRepository struct {
	DB *sql.DB
}

func NewMoodRepository(db *sql.DB) *MoodRepository {
	return &MoodRepository{DB: db}
}

func (r *MoodRepository) SaveMood(log models.MoodLog) error {
	_, err := r.DB.Exec("INSERT INTO Mood_Logs (user_id, score, note) VALUES (?, ?, ?)", log.UserID, log.Score, log.Note)
	if err != nil {
		return err
	}
	// Запись в аудит
	r.DB.Exec("INSERT INTO Audit_Logs (actor_id, action_type) VALUES (?, ?)", log.UserID, "CreateLog")
	return nil
}

func (r *MoodRepository) GetAllMoods(score, sort string) ([]models.MoodLog, error) {
	query := "SELECT log_id, user_id, score, note, timestamp FROM Mood_Logs"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.MoodLog
	for rows.Next() {
		var l models.MoodLog
		err := rows.Scan(&l.ID, &l.UserID, &l.Score, &l.Note, &l.Timestamp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *MoodRepository) GetMoodStats() (map[string]int, error) {
	rows, err := r.DB.Query("SELECT score, COUNT(*) FROM Mood_Logs GROUP BY score")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var s, c int
		rows.Scan(&s, &c)
		stats[fmt.Sprintf("%d", s)] = c
	}
	return stats, nil
}

func (r *MoodRepository) GetTopTags() (map[string]int, error) {
	return make(map[string]int), nil
}

func (r *MoodRepository) GetTeamAverage(id int) (float64, error) {
	var avg float64
	err := r.DB.QueryRow("SELECT AVG(score) FROM Mood_Logs").Scan(&avg)
	if err != nil {
		return 0, err
	}
	return avg, nil
}

func (r *MoodRepository) CheckPin(email, pin string) (bool, error) {
	var dbPin string
	err := r.DB.QueryRow("SELECT password_hash FROM Users WHERE email = ?", email).Scan(&dbPin)
	return err == nil && dbPin == pin, nil
}

func (r *MoodRepository) CreateUser(username, email, pin string) error {
	_, err := r.DB.Exec("INSERT INTO Users (username, email, password_hash) VALUES (?, ?, ?)", username, email, pin)
	return err
}

func (r *MoodRepository) DeleteUser(email string) error {
	_, err := r.DB.Exec("DELETE FROM Users WHERE email = ?", email)
	return err
}

func (r *MoodRepository) SearchNotes(q string) ([]models.MoodLog, error) {
	rows, err := r.DB.Query("SELECT log_id, user_id, score, note, timestamp FROM Mood_Logs WHERE note LIKE ?", "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.MoodLog
	for rows.Next() {
		var l models.MoodLog
		rows.Scan(&l.ID, &l.UserID, &l.Score, &l.Note, &l.Timestamp)
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *MoodRepository) GetScheduledReminders() ([]models.ReminderSettings, error) {
	return []models.ReminderSettings{
		{UserID: 1, ReminderTime: "09:00", IsEnabled: true},
		{UserID: 2, ReminderTime: "21:00", IsEnabled: true},
	}, nil
}
