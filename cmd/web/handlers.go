package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"pastepanda/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	pastes, err := app.pastes.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Pastes = pastes

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	paste, err := app.pastes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)

		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Paste = paste

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) pasteCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new paste..."))
}

func (app *application) pasteCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.pastes.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/view/%d", id), http.StatusSeeOther)
}
