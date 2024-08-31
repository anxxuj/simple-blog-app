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

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.index))
	router.Handler(http.MethodGet, "/post/view/:id", dynamic.ThenFunc(app.postView))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/register", dynamic.ThenFunc(app.userRegister))
	router.Handler(http.MethodPost, "/user/register", dynamic.ThenFunc(app.userRegisterPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/post/add", protected.ThenFunc(app.postAdd))
	router.Handler(http.MethodPost, "/post/add", protected.ThenFunc(app.postAddPost))
	router.Handler(http.MethodGet, "/post/edit/:id", protected.ThenFunc(app.postEdit))
	router.Handler(http.MethodPost, "/post/edit/:id", protected.ThenFunc(app.postEditPost))
	router.Handler(http.MethodGet, "/post/delete/:id", protected.ThenFunc(app.postDelete))
	router.Handler(http.MethodGet, "/user/logout", protected.ThenFunc(app.userLogout))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
