package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /paste/view/{id}", dynamic.ThenFunc(app.pasteView))
	mux.Handle("GET /paste/create", dynamic.ThenFunc(app.pasteCreate))
	mux.Handle("POST /paste/create", dynamic.ThenFunc(app.pasteCreatePost))

	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		commonHeaders,
	)

	return standard.Then(mux)
}
