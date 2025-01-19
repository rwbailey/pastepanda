package models

import (
	"database/sql"
	"time"
)

type Paste struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type PasteModel struct {
	DB *sql.DB
}

func (m *PasteModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func (m *PasteModel) Get(id int) (Paste, error) {
	return Paste{}, nil
}

func (m *PasteModel) Latest() ([]Paste, error) {
	return nil, nil
}
