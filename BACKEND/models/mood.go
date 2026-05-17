package models

import "time"

// MoodLog описывает одну запись в дневнике
type MoodLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Score     int       `json:"score"`     // Оценка от 1 до 5
	Note      string    `json:"note"`      // Текст заметки
	Timestamp time.Time `json:"timestamp"` // Время записи
}
