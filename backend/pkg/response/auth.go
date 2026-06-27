package response

import (
	"github.com/aarontrelstad/api/internal/db"
	"github.com/google/uuid"
)

type UserAuthResponse struct {
    ID       uuid.UUID `json:"id"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
}

func ToUserAuthResponse(user db.User) UserAuthResponse {
    return UserAuthResponse{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
    }
}
