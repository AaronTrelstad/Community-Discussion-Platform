package render

import (
	"encoding/json"
	"net/http"

	"github.com/aarontrelstad/api/pkg/apperror"
)

func JSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string) {
    JSON(w, status, map[string]string{"error": message})
}

func HandleError(w http.ResponseWriter, err error) {
    if appErr, ok := err.(*apperror.AppError); ok {
        Error(w, appErr.Status, appErr.Message)
        return
    }
    Error(w, http.StatusInternalServerError, "something went wrong")
}
