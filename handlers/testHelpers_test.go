package handlers

import (
	"github.com/geneowak/go-expense-tracker/internal/database"
)

func newTestApiConfig(querier database.Querier) *ApiConfig {
	validate := CreateValidator()

	return &ApiConfig{
		DB:       querier,
		Validate: validate,
	}
}
