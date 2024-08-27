package main

import "net/http"

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, http.StatusOK, "index.html", nil)
}
