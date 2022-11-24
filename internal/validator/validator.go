package db

import (
	"github.com/go-playground/validator"
)

func GetValidator() *validator.Validate {
	return validator.New()
}
