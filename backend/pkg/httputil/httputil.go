// pkg/httputil/httputil.go
package httputil

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/aarontrelstad/api/pkg/apperror"
    "github.com/google/uuid"
)

func Decode(r *http.Request, dst any) error {
    if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
        return apperror.ErrDecodeRequest
    }
    return nil
}

func GetUserID(r *http.Request) (uuid.UUID, error) {
    userIDStr, ok := r.Context().Value("user_id").(string)
    if !ok {
        return uuid.UUID{}, apperror.ErrInvalidToken
    }
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return uuid.UUID{}, apperror.ErrInvalidToken
    }
    return userID, nil
}

func SetAuthCookie(w http.ResponseWriter, token string) {
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

func SetRefreshCookie(w http.ResponseWriter, token string) {
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

func ClearAuthCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    "",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
        Path:     "/",
        Expires:  time.Unix(0, 0),
        MaxAge:   -1,
    })
}

func ClearRefreshCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
        Path:     "/",
        Expires:  time.Unix(0, 0),
        MaxAge:   -1,
    })
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := "http://localhost:3000"
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
