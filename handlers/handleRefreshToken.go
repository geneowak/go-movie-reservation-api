package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/geneowak/go-expense-tracker/internal/auth"
)

func (cfg *ApiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	refreshToken, err := cfg.DB.GetRefreshToken(r.Context(), token)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching refresh token", err)
		return
	}

	if refreshToken.RevokedAt != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has already been revoked", nil)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has already expired", nil)
		return
	}

	jwtToken, err := auth.MakeJWT(refreshToken.UserID, cfg.JwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error refreshing token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{Token: jwtToken})
}
