package mocks

import (
	"pastepanda/internal/models"
	"time"
)

var mockPaste = models.Paste{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type PasteModel struct{}

func (m *PasteModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *PasteModel) Get(id int) (models.Paste, error) {
	switch id {
	case 1:
		return mockPaste, nil
	default:
		return models.Paste{}, models.ErrNoRecord
	}
}

func (m *PasteModel) Latest() ([]models.Paste, error) {
	return []models.Paste{mockPaste}, nil
}
