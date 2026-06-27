package handlers

import (
	"net/http"

	"github.com/aarontrelstad/api/internal/services"
	"github.com/aarontrelstad/api/pkg/apperror"
	"github.com/aarontrelstad/api/pkg/response"
	util "github.com/aarontrelstad/api/pkg/httputil"
	"github.com/aarontrelstad/api/pkg/requests"
	"github.com/aarontrelstad/api/pkg/render"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterRequest
	if err := util.Decode(r, &req); err != nil {
		render.HandleError(w, err)
		return
	}

	user, tokens, err := h.authService.Register(r.Context(), req)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	util.SetAuthCookie(w, tokens.JWT)
	util.SetRefreshCookie(w, tokens.RefreshToken)
	render.JSON(w, http.StatusCreated, response.ToUserAuthResponse(user))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req requests.LoginRequest
	if err := util.Decode(r, &req); err != nil {
		render.HandleError(w, err)
		return
	}

	user, tokens, err := h.authService.Login(r.Context(), req)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	util.SetAuthCookie(w, tokens.JWT)
	util.SetRefreshCookie(w, tokens.RefreshToken)
	render.JSON(w, http.StatusOK, response.ToUserAuthResponse(user))
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		util.ClearAuthCookie(w)
		util.ClearRefreshCookie(w)
		render.HandleError(w, apperror.ErrInvalidToken)
		return
	}

	tokens, err := h.authService.Refresh(r.Context(), cookie.Value)
	if err != nil {
		util.ClearAuthCookie(w)
		util.ClearRefreshCookie(w)
		render.HandleError(w, err)
		return
	}

	util.SetAuthCookie(w, tokens.JWT)
	util.SetRefreshCookie(w, tokens.RefreshToken)
	render.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		h.authService.Logout(r.Context(), cookie.Value)
	}
	util.ClearAuthCookie(w)
	util.ClearRefreshCookie(w)
	render.JSON(w, http.StatusOK, map[string]any{"message": "logged out"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDStr, err := util.GetUserID(r)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userIDStr.String())
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, response.ToUserAuthResponse(user))
}
