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
