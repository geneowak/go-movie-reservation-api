package handlers

import (
	"github.com/geneowak/cinehold/internal/database"
	"github.com/go-playground/validator/v10"
)

type ApiConfig struct {
	DB        database.Querier
	Platform  string
	JwtSecret string
	Validate  *validator.Validate
}
