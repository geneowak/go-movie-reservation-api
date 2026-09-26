package handlers

import (
	"errors"
	"net/http"

	"github.com/geneowak/cinehold/internal/auth"
	"github.com/jackc/pgx/v5"
)

func (cfg *ApiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	refreshToken, err := cfg.DB.GetRefreshToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

	err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error refreshing token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
