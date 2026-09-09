package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/geneowak/go-expense-tracker/internal/auth"
	"github.com/geneowak/go-expense-tracker/internal/database"
)

func (cfg *ApiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,alphanum,min=5"`
	}

	var params loginRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request", err)
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		handleValidationErrors(w, err)
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect email or password", err)
		return
	}

	if matches, err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil || !matches {
		respondWithError(w, http.StatusBadRequest, "Incorrect email or password", err)
		return
	}

	refreshToken, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour), // token valid for 7 days
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}
	jwtToken, err := auth.MakeJWT(user.ID, cfg.JwtSecret, time.Hour) // jwt only valid for 1 hour
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT token", err)
		return
	}
	type response struct {
		database.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		User:         user,
		Token:        jwtToken,
		RefreshToken: refreshToken.Token,
	})
}
