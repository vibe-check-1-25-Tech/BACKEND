package models

import (
	"math/rand"
	"time"
)

// MoodLog представляет запись о настроении в базе данных
type MoodLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Score     int       `json:"score"`
	Note      string    `json:"note"`
	Photo     string    `json:"photo"`
	Tags      string    `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
}

// ReminderSettings для настройки уведомлений (ДОБАВЛЕНО)
type ReminderSettings struct {
	UserID       int    `json:"user_id"`
	ReminderTime string `json:"reminder_time"`
	IsEnabled    bool   `json:"is_enabled"`
}

// MoodSaveResponse — структура, которую ждет фронтенд (JS)
type MoodSaveResponse struct {
	Status  string          `json:"status"`
	Support *SupportContent `json:"support,omitempty"`
}

type SupportContent struct {
	Type    string `json:"type"` // "image" или "joke"
	Content string `json:"content"`
}

// Твой список контента
var SupportList = []SupportContent{
	{Type: "joke", Content: "— Куда идешь?\n— В спортзал.\n— А зачем?\n— Стать сильным.\n— И что ты будешь делать с этой силой?\n— Шкаф передвину, а то зарядка не достает."},
	{Type: "joke", Content: "Программист ставит на тумбочку два стакана. Один с водой — на случай, если захочет пить. Второй пустой — на случай, если не захочет."},
	{Type: "meme", Content: "/assets/cat.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem1.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem2.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem3.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem4.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem5.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem6.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem7.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem8.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem9.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem10.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem11.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem12.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem13.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem14.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem15.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem16.jpg.jpeg"},
	{Type: "meme", Content: "/assets/mem18.jpeg"},
	{Type: "meme", Content: "/assets/mem19.jpeg"},
}

func GetRandomSupport() SupportContent {
	rand.Seed(time.Now().UnixNano())

	if len(SupportList) == 0 {
		return SupportContent{Type: "joke", Content: "Улыбнись! Всё получится!"}
	}

	content := SupportList[rand.Intn(len(SupportList))]

	// ВАЖНО: JS на фронтенде ждет тип "image", а в списке "meme".
	if content.Type == "meme" {
		content.Type = "image"
	}

	return content
}
