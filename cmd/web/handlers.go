package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anxxuj/simple-blog-app/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderTemplate(w, http.StatusOK, "index.html", &tempateData{
		Posts: posts,
	})
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	post, err := app.posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.renderTemplate(w, http.StatusOK, "post.html", &tempateData{
		Post: post,
	})
}

func (app *application) postAdd(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, http.StatusOK, "post_form.html", &tempateData{
		Form: &PostForm{Name: "Add Post"},
	})
}

func (app *application) postAddPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &PostForm{
		Name:    "Add Post",
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}

	if !form.Validate() {
		app.renderTemplate(w, http.StatusUnprocessableEntity, "post_form.html", &tempateData{
			Form: form,
		})
		return
	}

	id, err := app.posts.Insert(form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func (app *application) postEdit(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	post, err := app.posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	form := &PostForm{
		Name:    "Edit Post",
		Title:   post.Title,
		Content: post.Content,
	}

	app.renderTemplate(w, http.StatusOK, "post_form.html", &tempateData{
		Form: form,
	})
}

func (app *application) postEditPost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &PostForm{
		Name:    "Edit Post",
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}

	if !form.Validate() {
		app.renderTemplate(w, http.StatusUnprocessableEntity, "post_form.html", &tempateData{
			Form: form,
		})
		return
	}

	err = app.posts.Update(id, form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func (app *application) postDelete(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.posts.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
