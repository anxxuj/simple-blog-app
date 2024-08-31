package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.index))
	router.Handler(http.MethodGet, "/post/view/:id", dynamic.ThenFunc(app.postView))
	router.Handler(http.MethodGet, "/post/add", dynamic.ThenFunc(app.postAdd))
	router.Handler(http.MethodPost, "/post/add", dynamic.ThenFunc(app.postAddPost))
	router.Handler(http.MethodGet, "/post/edit/:id", dynamic.ThenFunc(app.postEdit))
	router.Handler(http.MethodPost, "/post/edit/:id", dynamic.ThenFunc(app.postEditPost))
	router.Handler(http.MethodGet, "/post/delete/:id", dynamic.ThenFunc(app.postDelete))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
