package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type updateAdminStatusRequest struct {
	IsAdmin *bool `json:"is_admin" validate:"required"`
}

func (cfg *ApiConfig) handleUpdateAdminStatus(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("userId")
	userId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user id", err)
		return
	}

	signedInId, _ := GetUserIdFromContext(r.Context())
	// don't allow user to edit himself so that the system will always have atleast one admin
	if signedInId == userId {
		respondWithError(w, http.StatusBadRequest, "You cannot update your own admin status", nil)
		return
	}

	var req updateAdminStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleJsonDecodeError(w, err)
		return
	}

	if err := cfg.Validate.Struct(req); err != nil {
		handleValidationErrors(w, err)
		return
	}

	updatedUser, err := cfg.DB.UpdateUserAdminStatus(r.Context(), database.UpdateUserAdminStatusParams{
		ID:      userId,
		IsAdmin: *req.IsAdmin,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error updating user admin status", err)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedUser)

}
