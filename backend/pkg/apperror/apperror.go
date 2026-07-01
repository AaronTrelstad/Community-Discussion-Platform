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
	ErrDecodeRequest         = &AppError{http.StatusBadRequest, "Failed to decode request body"}
	ErrInvalidToken          = &AppError{http.StatusUnauthorized, "Invalid or expired token"}
	ErrHashPassword          = &AppError{http.StatusInternalServerError, "Failed to hash password"}
	ErrGenerateJWTToken      = &AppError{http.StatusInternalServerError, "Failed to generate JWT token"}
	ErrGenerateRefreshToken  = &AppError{http.StatusInternalServerError, "Failed to generate refresh token"}
	ErrCreateRefreshToken    = &AppError{http.StatusInternalServerError, "Failed to create refresh token"}
	ErrCreateUser            = &AppError{http.StatusConflict, "Username or email already exists"}
	ErrCommunityNameConflict = &AppError{http.StatusConflict, "Community name already exists"}
	ErrCreateCommunity       = &AppError{http.StatusInternalServerError, "Failed to create community"}
	ErrJoinCommunity         = &AppError{http.StatusInternalServerError, "Failed to join community"}
	ErrNotFound              = &AppError{http.StatusNotFound, "Not found"}
	ErrForbidden             = &AppError{http.StatusForbidden, "Forbidden"}
	ErrInvalidCredentials    = &AppError{http.StatusUnauthorized, "Invalid credentials"}
	ErrListCommunities       = &AppError{http.StatusInternalServerError, "Failed to list communities"}
	ErrParseURL				 = &AppError{http.StatusBadRequest, "Failed to parse url"}
)
