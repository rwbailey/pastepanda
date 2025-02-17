package main

import (
	"net/http"

	"github.com/justinas/alice"

	"pastepanda/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(
		app.sessionManager.LoadAndSave,
		noSurf,
		app.authenticate,
	)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /paste/view/{id}", dynamic.ThenFunc(app.pasteView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /paste/create", protected.ThenFunc(app.pasteCreate))
	mux.Handle("POST /paste/create", protected.ThenFunc(app.pasteCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		commonHeaders,
	)

	return standard.Then(mux)
}
