package handlers

import (
	"github.com/go-playground/validator/v10"
)

type Validatable interface {
	Validate() validator.ValidationErrors
}

type createUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Email     string `json:"email"      validate:"required,email"`
}

func (r createUserRequest) Validate() validator.ValidationErrors {
	validate := validator.New()
	err := validate.Struct(r)
	validationErrors, ok := err.(validator.ValidationErrors)
	if ok && len(validationErrors) > 0 {
		return err.(validator.ValidationErrors)
	}

	return nil
}

type createUserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
