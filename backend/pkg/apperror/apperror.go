package apperror

import "net/http"

type AppError struct {
    Status  int
    Message string
}

func (e *AppError) Error() string {
    return e.Message
}

var (
    ErrBadRequest     = &AppError{http.StatusBadRequest, "Bad request"}
    ErrUnauthorized   = &AppError{http.StatusUnauthorized, "Unauthorized"}
    ErrForbidden      = &AppError{http.StatusForbidden, "Forbidden"}
    ErrNotFound       = &AppError{http.StatusNotFound, "Not found"}
    ErrConflict       = &AppError{http.StatusConflict, "Already exists"}
    ErrInternal       = &AppError{http.StatusInternalServerError, "Something went wrong"}
    ErrInvalidToken   = &AppError{http.StatusUnauthorized, "Invalid or expired token"}
    ErrInvalidCreds   = &AppError{http.StatusUnauthorized, "Invalid credentials"}
)

func New(status int, message string) *AppError {
    return &AppError{status, message}
}
