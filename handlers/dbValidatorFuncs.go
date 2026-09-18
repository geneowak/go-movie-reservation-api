package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) ValidateCinemaId(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	id, err := uuid.Parse(val)
	if err != nil {
		return false
	}
	exists, err := cfg.DB.CheckCinemaById(context.Background(), id)
	if err != nil {
		return false
	}

	return exists
}
