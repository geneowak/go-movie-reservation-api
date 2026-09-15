package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/auth"
)

func (cfg *ApiConfig) middlewareAdminAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check that the token exists
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		// extract user id from the token
		userId, err := auth.ValidateJWT(token, cfg.JwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		// get the user and check if they are an admin
		user, err := cfg.DB.GetUserById(r.Context(), userId)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		if !user.IsAdmin {
			respondWithError(w, http.StatusForbidden, "You don't have permissions to perform this action", fmt.Errorf("User not admin"))
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
