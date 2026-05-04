package models

import (
	"math/rand"
	"time"
)

// MoodLog представляет запись о настроении в базе данных
type MoodLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Score     int       `json:"score"` // Оценка от 1 до 5
	Note      string    `json:"note"`  // Текстовая заметка
	Photo     string    `json:"photo"` // Base64 строка изображения
	Tags      string    `json:"tags"`  // Теги через запятую
	Timestamp time.Time `json:"timestamp"`
}

// ReminderSettings для настройки уведомлений
type ReminderSettings struct {
	UserID       int    `json:"user_id"`
	ReminderTime string `json:"reminder_time"`
	IsEnabled    bool   `json:"is_enabled"`
}

// SupportContent структура для мемов и анекдотов
type SupportContent struct {
	Type    string `json:"type"`    // "meme" или "joke"
	Content string `json:"content"` // Путь к файлу или текст анекдота
}

// SupportList — ПОЛНЫЙ список на основе твоего скриншота
var SupportList = []SupportContent{
	// Анекдоты
	{Type: "joke", Content: "— Куда идешь?\n— В спортзал.\n— А зачем?\n— Стать сильным.\n— И что ты будешь делать с этой силой?\n— Шкаф передвину, а то зарядка не достает."},
	{Type: "joke", Content: "Программист ставит на тумбочку два стакана. Один с водой — на случай, если захочет пить. Второй пустой — на случай, если не захочет."},

	// ПЕРВЫЕ МЕМЫ (по твоему списку из assets)
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

	// ОСТАЛЬНЫЕ МЕМЫ
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

// GetRandomSupport выбирает случайный контент из списка
func GetRandomSupport() SupportContent {
	rand.Seed(time.Now().UnixNano())

	if len(SupportList) == 0 {
		return SupportContent{Type: "joke", Content: "Улыбнись! Всё получится!"}
	}

	return SupportList[rand.Intn(len(SupportList))]
}
