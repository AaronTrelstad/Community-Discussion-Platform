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
