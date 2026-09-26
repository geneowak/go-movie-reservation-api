package handlers

import (
	"github.com/geneowak/cinehold/internal/database"
)

func newTestApiConfig(querier database.Querier) *ApiConfig {
	cfg := ApiConfig{
		DB: querier,
	}

	validate := CreateValidator(&cfg)

	cfg.Validate = validate

	return &cfg
}
