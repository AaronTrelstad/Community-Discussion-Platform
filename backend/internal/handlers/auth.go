package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aarontrelstad/api/internal/db"
	jwtpkg "github.com/aarontrelstad/api/pkg/jwt"
	"github.com/aarontrelstad/api/pkg/response"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries *db.Queries
}

func NewAuthHandler(queries *db.Queries) *AuthHandler {
	return &AuthHandler{queries: queries}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Eamil    string `json:"email"`
	Password string `json:"password"`
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
	})
}

func setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
}

func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})
}

func clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})
}

func (h *AuthHandler) issueTokens(w http.ResponseWriter, r *http.Request, userIDStr string) error {
	jwt, err := jwtpkg.Generate(userIDStr)
	if err != nil {
		return err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}

	_, err = h.queries.CreateRefreshToken(r.Context(), db.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		return err
	}

	setAuthCookie(w, jwt)
	setRefreshCookie(w, refreshToken)
	return nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := h.queries.CreateUser(r.Context(), db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	})
	if err != nil {
		response.Error(w, http.StatusConflict, "username or email already exists")
		return
	}

	if err := h.issueTokens(w, r, user.ID.String()); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to issue tokens")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), req.Eamil)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := h.issueTokens(w, r, user.ID.String()); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to issue tokens")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		clearAuthCookie(w)
		clearRefreshCookie(w)
		response.Error(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	rt, err := h.queries.GetRefreshToken(r.Context(), cookie.Value)
	if err != nil {
		clearAuthCookie(w)
		clearRefreshCookie(w)
		response.Error(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	h.queries.RevokeRefreshToken(r.Context(), rt.Token)

	if err := h.issueTokens(w, r, rt.UserID.String()); err != nil {
		clearAuthCookie(w)
		clearRefreshCookie(w)
		response.Error(w, http.StatusInternalServerError, "failed to issue tokens")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		h.queries.RevokeRefreshToken(r.Context(), cookie.Value)
	}
	clearAuthCookie(w)
	clearRefreshCookie(w)
	response.JSON(w, http.StatusOK, map[string]any{"message": "logged out"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}
