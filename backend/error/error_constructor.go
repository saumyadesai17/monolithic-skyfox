package error

import "net/http"

func NotFoundError(code string, message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, code, message, err)
}

func UnProcessableError(code string, message string, err error) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, code, message, err)
}

func InternalServerError(code string, message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, code, message, err)
}

func BadRequestError(code string, message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, code, message, err)
}

func NewAppError(httpCode int, code string, message string, err error) *AppError {
	return &AppError{httpCode: httpCode, Code: code, Message: message, error: err}
}

func InvalidCredentialsError(code string, message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, code, message, err)
}
