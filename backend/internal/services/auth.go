package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/aarontrelstad/api/internal/db"
	"github.com/aarontrelstad/api/pkg/apperror"
	jwtpkg "github.com/aarontrelstad/api/pkg/jwt"
	"github.com/aarontrelstad/api/pkg/requests"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries *db.Queries
}

type TokenPair struct {
	JWT          string
	RefreshToken string
}

func NewAuthService(queries *db.Queries) *AuthService {
	return &AuthService{queries: queries}
}

func (s *AuthService) Register(ctx context.Context, req requests.RegisterRequest) (db.User, TokenPair, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, TokenPair{}, apperror.ErrHashPassword
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	})
	if err != nil {
		return db.User{}, TokenPair{}, apperror.ErrCreateUser
	}

	tokens, err := s.issueTokens(ctx, user.ID.String())
	if err != nil {
		return db.User{}, TokenPair{}, err
	}

	return user, tokens, nil
}

func (s *AuthService) Login(ctx context.Context, req requests.LoginRequest) (db.User, TokenPair, error) {
	user, err := s.queries.GetUserByEmail(ctx, req.Eamil)
	if err != nil {
		return db.User{}, TokenPair{}, apperror.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return db.User{}, TokenPair{}, apperror.ErrInvalidCredentials
	}

	tokens, err := s.issueTokens(ctx, user.ID.String())
	if err != nil {
		return db.User{}, TokenPair{}, err
	}

	return user, tokens, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	rt, err := s.queries.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return TokenPair{}, apperror.ErrInvalidToken
	}

	s.queries.RevokeRefreshToken(ctx, rt.Token)

	return s.issueTokens(ctx, rt.UserID.String())
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.queries.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (db.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return db.User{}, apperror.ErrInvalidToken
	}
	return s.queries.GetUserByID(ctx, id)
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) issueTokens(ctx context.Context, userIDStr string) (TokenPair, error) {
	jwt, err := jwtpkg.Generate(userIDStr)
	if err != nil {
		return TokenPair{}, apperror.ErrGenerateJWTToken
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, apperror.ErrGenerateRefreshToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return TokenPair{}, apperror.ErrInvalidToken
	}

	_, err = s.queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		return TokenPair{}, apperror.ErrCreateRefreshToken
	}

	return TokenPair{JWT: jwt, RefreshToken: refreshToken}, nil
}
