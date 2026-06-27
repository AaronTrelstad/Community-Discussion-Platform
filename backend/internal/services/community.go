package services

import (
	"context"
	"database/sql"

	"github.com/aarontrelstad/api/internal/db"
	"github.com/aarontrelstad/api/pkg/apperror"
	"github.com/aarontrelstad/api/pkg/requests"
	"github.com/google/uuid"
)

type CommunityService struct {
	queries *db.Queries 
}

func NewCommunityService(queries *db.Queries) *CommunityService {
	return &CommunityService{queries: queries}
}

func (s *CommunityService) CreateCommunity(ctx context.Context, userID uuid.UUID, req requests.CreateCommunityRequest) (db.Community, error) {
	_, err := s.queries.GetCommunityByName(ctx, req.Name)
	if err != nil {
		return db.Community{}, apperror.ErrCommunityNameConflict
	}

	community, err := s.queries.CreateCommunity(ctx, db.CreateCommunityParams{
		Name: req.Name,
		Title: req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		CreatedBy:   uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return db.Community{}, apperror.ErrCreateCommunity
	}

	err = s.queries.JoinCommunity(ctx, db.JoinCommunityParams{
		UserID: userID,
		CommunityID: community.ID,
	})
	if err != nil {
		return db.Community{}, apperror.ErrJoinCommunity
	}

	return community, nil
}
