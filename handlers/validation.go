package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func handleValidationErrors(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity) // 422 matches Laravel's default

	// Convert errors to a map: {"field": ["error message"]}
	validationErrs := err.(validator.ValidationErrors)
	errorsMap := make(map[string][]string)

	for _, e := range validationErrs {
		field := e.Field()
		msg := getErrorMsg(e)

		// Extract just the error message (e.g., "Email is required")
		// The library provides a way to get a more readable message
		// Often we want "field": ["The email field is required."]
		if _, ok := errorsMap[field]; !ok {
			errorsMap[field] = []string{}
		}
		errorsMap[field] = append(errorsMap[field], msg)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"message": "There were some validation errors",
		"errors":  errorsMap,
	})
}

func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required", err.Field())
	case "email":
		return fmt.Sprintf("The %s must be a valid email.", err.Field())
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters long.", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("The %s field must be at most %s characters long.", err.Field(), err.Param())
	case "alphanum":
		return fmt.Sprintf("The %s field must be alphanumeric.", err.Field())
	default:
		return err.Error()
	}
}
