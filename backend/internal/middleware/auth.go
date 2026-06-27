package middleware

import (
	"context"
	"net/http"

	"github.com/aarontrelstad/api/pkg/response"
	jwtpkg "github.com/aarontrelstad/api/pkg/jwt"

)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "missing auth token")
			return
		}

		claims, err := jwtpkg.Verify(cookie.Value)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid auth token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
