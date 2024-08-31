package main

import "github.com/anxxuj/simple-blog-app/internal/validator"

type postForm struct {
	Name    string
	Title   string
	Content string
	validator.Validator
}

func (form *postForm) Validate() bool {
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be empty")
	form.CheckField(validator.MaxChars(form.Title, 140), "title", "This field cannot be more than 140 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be empty")

	return form.Valid()
}

type registerForm struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
	validator.Validator
}

func (form *registerForm) Validate() bool {
	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be empty")
	form.CheckField(validator.Matches(form.Username, validator.UsernameRX), "username", "This field must be a valid username")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be empty")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be atleast 8 characters long")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirmPassword", "This field cannot be empty")
	form.CheckField(validator.EqualTo(form.ConfirmPassword, form.Password), "confirmPassword", "This field should be equal to password")

	return form.Valid()
}

type loginForm struct {
	Username string
	Password string
	validator.Validator
}

func (form *loginForm) Validate() bool {
	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be empty")
	form.CheckField(validator.Matches(form.Username, validator.UsernameRX), "username", "This field must be a valid username")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be empty")

	return form.Valid()
}
