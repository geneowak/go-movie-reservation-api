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

func handleJsonDecodeError(w http.ResponseWriter, err error) {
	if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		msg := fmt.Sprintf("Invalid field %s: Expected type %s, but got %s", typeErr.Field, typeErr.Type.String(), typeErr.Value)
		throwAsValidationError(w, typeErr.Field, msg)
		return
	}
	respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
	return
}
