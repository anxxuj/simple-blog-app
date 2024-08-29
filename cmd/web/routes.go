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

	router.HandlerFunc(http.MethodGet, "/", app.index)
	router.HandlerFunc(http.MethodGet, "/post/view/:id", app.postView)
	router.HandlerFunc(http.MethodGet, "/post/add", app.postAdd)
	router.HandlerFunc(http.MethodPost, "/post/add", app.postAddPost)
	router.HandlerFunc(http.MethodGet, "/post/edit/:id", app.postEdit)
	router.HandlerFunc(http.MethodPost, "/post/edit/:id", app.postEditPost)
	router.HandlerFunc(http.MethodGet, "/post/delete/:id", app.postDelete)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
