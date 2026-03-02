package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	ae "skyfox/error"
)

func HandleStructValidationError(err error) interface{} {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		// Not a field-validation error (e.g. malformed JSON body).
		return ae.BadRequestError("BadRequest", err.Error(), err)
	}

	for _, fieldErr := range validationErrors {
		msg := fieldError{fieldErr}.String()
		return ae.BadRequestError("ValidationFailed", msg, fieldErr)
	}
	return nil
}

type fieldError struct {
	err validator.FieldError
}

func (e fieldError) String() string {
	var sb strings.Builder

	sb.WriteString("field '" + e.err.Field() + "'")
	sb.WriteString(", condition: " + validationErrorToText(e.err))

	if e.err.Value() != nil && e.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", provided: %v", e.err.Value()))
	}
	return sb.String()
}

func validationErrorToText(fieldErr validator.FieldError) string {
	switch fieldErr.ActualTag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldErr.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", fieldErr.Field(), fieldErr.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", fieldErr.Field(), fieldErr.Param())
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", fieldErr.Field(), fieldErr.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than %s", fieldErr.Field(), fieldErr.Param())
	case "phoneNumber":
		return fmt.Sprintf("%s must be exactly 10 digits", fieldErr.Field())
	case "passwordStrength":
		return fmt.Sprintf("%s must be at least 8 characters and contain an uppercase letter, a digit, and a special character", fieldErr.Field())
	}
	return fmt.Sprintf("%s is not valid", fieldErr.Field())
}
