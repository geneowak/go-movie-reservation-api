package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) middlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do the auth stuff here
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		userId, err := auth.ValidateJWT(token, cfg.JwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		// ensure that the user id exists in our db
		if exists, err := cfg.DB.CheckUserId(r.Context(), userId); err != nil || !exists {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value("user_id")
	if val == nil {
		return uuid.Nil, errors.New("User ID not found in context")
	}

	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("ID value not of type uuid.UUID")
	}

	return id, nil
}
