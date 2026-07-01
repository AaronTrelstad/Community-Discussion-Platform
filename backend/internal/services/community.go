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
		Name:        req.Name,
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		CreatedBy:   uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return db.Community{}, apperror.ErrCreateCommunity
	}

	err = s.queries.JoinCommunity(ctx, db.JoinCommunityParams{
		UserID:      userID,
		CommunityID: community.ID,
	})
	if err != nil {
		return db.Community{}, apperror.ErrJoinCommunity
	}

	return community, nil
}

func (s *CommunityService) ListCommunities(ctx context.Context, req requests.ListCommunitiesRequest) ([]db.Community, error) {
	communities, err := s.queries.ListCommunities(ctx, db.ListCommunitiesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return []db.Community{}, apperror.ErrListCommunities
	}

	return communities, nil
}

func (s *CommunityService) UpdateCommunity(ctx context.Context, communityID, userID uuid.UUID, req requests.UpdateCommunityRequest) (db.Community, error) {
	member, err := s.queries.GetCommunityMember(ctx, db.GetCommunityMemberParams{
		UserID: userID,
		CommunityID: communityID,
	})
	if err != nil || (member.Role.String != "moderator" && member.Role.String != "admin") {
		return db.Community{}, apperror.ErrForbidden
	}

	return s.queries.UpdateCommunity(ctx, db.UpdateCommunityParams{
		ID: communityID,
		Title: req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	})
}
