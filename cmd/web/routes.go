package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /paste/view/{id}", app.pasteView)
	mux.HandleFunc("GET /paste/create", app.pasteCreate)
	mux.HandleFunc("POST /paste/create", app.pasteCreatePost)

	return app.recoverPanic(
		app.logRequest(
			commonHeaders(mux),
		),
	)
}
