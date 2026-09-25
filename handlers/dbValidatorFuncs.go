package handlers

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (cfg *ApiConfig) IsUniqueEmail(fl validator.FieldLevel) bool {
	// we don't need to verify if this is actually an email because we have validation rules for that
	email := fl.Field().String()
	exists, err := cfg.DB.CheckUserEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true
		}
		return false
	}

	return !exists
}
