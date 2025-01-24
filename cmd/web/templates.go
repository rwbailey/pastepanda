package main

import "pastepanda/internal/models"

type templateData struct {
	Paste  models.Paste
	Pastes []models.Paste
}
