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
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) pasteCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.pastes.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/view/%d", id), http.StatusSeeOther)
}
