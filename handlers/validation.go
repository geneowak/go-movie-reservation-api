package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func CreateValidator(cfg *ApiConfig) *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	// update the validator to return the json field name instead of the struct name
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})

	// register custom validators
	validate.RegisterValidation("cinema-exists", cfg.ValidateCinemaId)
	validate.RegisterValidation("is-valid-seat", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		// seat no needs to be in the format of Row:Seat eg A:1
		elements := strings.Split(val, ":")
		if len(elements) != 2 {
			return false
		}
		// the second element must be a number
		if _, err := strconv.Atoi(elements[1]); err != nil {
			return false
		}
		return true
	})

	return validate
}

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

func throwAsValidationError(w http.ResponseWriter, field, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity) // 422 matches Laravel's default

	errorsMap := map[string][]string{
		field: {errMsg},
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
	case "url":
		return fmt.Sprintf("The %s must be a valid URL.", err.Field())
	case "uuid_rfc4122", "uuid":
		return fmt.Sprintf("The %s must be a valid UUID.", err.Field())
	case "datetime":
		return fmt.Sprintf("The %s must be a valid datetime of the format %s.", err.Field(), err.Param())
	case "gtfield", "gtcsfield":
		return fmt.Sprintf("The %s field must be a greater than the %s field.", err.Field(), err.Param())
	case "gtefield", "gtecsfield":
		return fmt.Sprintf("The %s field must be a greater than or equal to the %s field.", err.Field(), err.Param())
	case "min":
		switch err.Kind() {
		case reflect.Slice, reflect.Array:
			return fmt.Sprintf("The %s field must contain atleast %s items.", err.Field(), err.Param())
		case reflect.Map:
			return fmt.Sprintf("The %s field must have at least %s entries.", err.Field(), err.Param())
		default:
			return fmt.Sprintf("The %s field must be at least %s characters long.", err.Field(), err.Param())
		}
	case "max":
		switch err.Kind() {
		case reflect.Slice, reflect.Array:
			return fmt.Sprintf("The %s field must not contain more than %s items.", err.Field(), err.Param())
		case reflect.Map:
			return fmt.Sprintf("The %s field must not have more than %s entries.", err.Field(), err.Param())
		default:
			return fmt.Sprintf("The %s field must be at most %s characters long.", err.Field(), err.Param())
		}
	case "alphanum":
		return fmt.Sprintf("The %s field must be alphanumeric.", err.Field())
	// custom validator errors
	case "cinema-exists":
		return fmt.Sprintf("The %s field must be a valid ID of an existing cinema.", err.Field())
	case "is-valid-seat":
		return fmt.Sprintf("The %s field is not a valid cinema seat number. Expected format is Row:SeatNo eg A:1", err.Field())
	default:
		return err.Error()
	}
}
