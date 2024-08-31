package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anxxuj/microblog/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = posts
	app.renderTemplate(w, http.StatusOK, "index.html", data)
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

	data := app.newTemplateData(r)
	data.Post = post
	app.renderTemplate(w, http.StatusOK, "post.html", data)
}

func (app *application) postAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &postForm{Name: "Add Post"}
	app.renderTemplate(w, http.StatusOK, "post_form.html", data)
}

func (app *application) postAddPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &postForm{
		Name:    "Add Post",
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}

	if !form.Validate() {
		data := app.newTemplateData(r)
		data.Form = form
		app.renderTemplate(w, http.StatusUnprocessableEntity, "post_form.html", data)
		return
	}

	id, err := app.posts.Insert(form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Post created successfully")

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

	form := &postForm{
		Name:    "Edit Post",
		Title:   post.Title,
		Content: post.Content,
	}

	data := app.newTemplateData(r)
	data.Form = form
	app.renderTemplate(w, http.StatusOK, "post_form.html", data)
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

	form := &postForm{
		Name:    "Edit Post",
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}

	if !form.Validate() {
		data := app.newTemplateData(r)
		data.Form = form
		app.renderTemplate(w, http.StatusUnprocessableEntity, "post_form.html", data)
		return
	}

	err = app.posts.Update(id, form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Post updated successfully")

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

	app.sessionManager.Put(r.Context(), "flash", "Post deleted successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userRegister(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Form = &registerForm{}
	app.renderTemplate(w, http.StatusOK, "register.html", data)
}

func (app *application) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &registerForm{
		Username:        r.PostForm.Get("username"),
		Email:           r.PostForm.Get("email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm-password"),
	}

	if !form.Validate() {
		data := app.newTemplateData(r)
		data.Form = form
		app.renderTemplate(w, http.StatusUnprocessableEntity, "register.html", data)
		return
	}

	err = app.users.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUsername) {
			form.AddFieldError("username", "Username is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.renderTemplate(w, http.StatusUnprocessableEntity, "register.html", data)
		} else if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.renderTemplate(w, http.StatusUnprocessableEntity, "register.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User registered successfully")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Form = &loginForm{}
	app.renderTemplate(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &loginForm{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
	}

	if !form.Validate() {
		data := app.newTemplateData(r)
		data.Form = form
		app.renderTemplate(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Username or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.renderTemplate(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "flash", "User logged in successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "User logged out successfully")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
