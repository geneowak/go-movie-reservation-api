package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithValidationErrors(w http.ResponseWriter, errorsMap map[string][]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "There were some validation errors",
		"errors":  errorsMap,
	})
}

func handleJsonDecodeError(w http.ResponseWriter, err error) {
	if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		msg := fmt.Sprintf("Invalid field %s: Expected type %s, but got %s", typeErr.Field, typeErr.Type.String(), typeErr.Value)
		respondWithValidationErrors(w, map[string][]string{
			typeErr.Field: {msg},
		})
		return
	}
	respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
	return
}
