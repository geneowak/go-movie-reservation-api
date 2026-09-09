package handlers

import (
	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/go-playground/validator/v10"
)

type ApiConfig struct {
	DB        *database.Queries
	Platform  string
	JwtSecret string
	Validate  *validator.Validate
}
