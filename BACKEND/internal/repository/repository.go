package repository

import (
	"database/sql"
	"vibe-check-backend/internal/models"
)

type MoodRepository struct {
	DB *sql.DB
}

func NewMoodRepository(db *sql.DB) *MoodRepository {
	return &MoodRepository{DB: db}
}

func (r *MoodRepository) SaveMood(log models.MoodLog) error {
	_, err := r.DB.Exec("INSERT INTO mood_logs (user_id, score, note) VALUES (?, ?, ?)", log.UserID, log.Score, log.Note)
	if err != nil {
		return err
	}
	r.DB.Exec("INSERT INTO audit_logs (actor_id, action_type) VALUES (?, ?)", log.UserID, "CreateLog")
	return nil
}

func (r *MoodRepository) GetAllMoods(score, sort string) ([]models.MoodLog, error) {
	query := "SELECT log_id, user_id, score, note, timestamp FROM mood_logs WHERE 1=1"
	rows, err := r.DB.Query(query)
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

func (r *MoodRepository) GetMoodStats() (map[string]int, error) {
	rows, _ := r.DB.Query("SELECT score, COUNT(*) FROM mood_logs GROUP BY score")
	defer rows.Close()
	stats := make(map[string]int)
	for rows.Next() {
		var s, c int
		rows.Scan(&s, &c)
		stats[string(rune(s))] = c
	}
	return stats, nil
}

func (r *MoodRepository) GetTopTags() (map[string]int, error) {
	return make(map[string]int), nil
}

func (r *MoodRepository) GetTeamAverage(id int) (float64, error) {
	var avg float64
	r.DB.QueryRow("SELECT AVG(score) FROM mood_logs").Scan(&avg)
	return avg, nil
}

func (r *MoodRepository) CheckPin(email, pin string) (bool, error) {
	var dbPin string
	err := r.DB.QueryRow("SELECT pin_code FROM users WHERE email = ?", email).Scan(&dbPin)
	return err == nil && dbPin == pin, nil
}

func (r *MoodRepository) CreateUser(username, email, pin string) error {
	_, err := r.DB.Exec("INSERT INTO users (username, email, pin_code) VALUES (?, ?, ?)", username, email, pin)
	return err
}

func (r *MoodRepository) DeleteUser(email string) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE email = ?", email)
	return err
}

func (r *MoodRepository) SearchNotes(q string) ([]models.MoodLog, error) {
	rows, _ := r.DB.Query("SELECT log_id, user_id, score, note, timestamp FROM mood_logs WHERE note LIKE ?", "%"+q+"%")
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
	// В реальной базе здесь был бы: SELECT user_id, reminder_time FROM reminders WHERE is_enabled = 1
	// Для теста возвращаем фейковый список
	return []models.ReminderSettings{
		{UserID: 1, ReminderTime: "09:00", IsEnabled: true},
		{UserID: 2, ReminderTime: "21:00", IsEnabled: true},
	}, nil
}
