package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/auth"
	"github.com/geneowak/go-expense-tracker/internal/database"
)

func (cfg *ApiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,alphanum,min=5"`
	}

	var params createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if err := cfg.Validate.Struct(params); err != nil {
		handleValidationErrors(w, err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}
