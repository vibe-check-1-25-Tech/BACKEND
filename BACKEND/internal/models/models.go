package models

import "time"

type MoodLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Score     int       `json:"score"`
	Note      string    `json:"note"`
	Timestamp time.Time `json:"timestamp"`
}
type ReminderSettings struct {
	UserID       int    `json:"user_id"`
	ReminderTime string `json:"reminder_time"` // Например, "19:00"
	IsEnabled    bool   `json:"is_enabled"`
}
