package services

import (
	"context"
	"database/sql"

	"github.com/aarontrelstad/api/internal/db"
	"github.com/aarontrelstad/api/pkg/apperror"
	"github.com/aarontrelstad/api/pkg/requests"
	"github.com/google/uuid"
)

type TeamService struct {
	queries *db.Queries
}

func NewTeamService(queries *db.Queries) *TeamService {
	return &TeamService{queries: queries}
}

func (s *TeamService) CreateTeam(ctx context.Context, userID uuid.UUID, req requests.CreateTeamRequest) (db.Team, error) {
	_, err := s.queries.GetTeamByName(ctx, req.Name)
	if err != nil {
		return db.Team{}, apperror.ErrTeamNameConflict
	}

	team, err := s.queries.CreateTeam(ctx, db.CreateTeamParams{
		Name:        req.Name,
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		CreatedBy:   uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return db.Team{}, apperror.ErrCreateTeam
	}

	err = s.queries.JoinTeam(ctx, db.JoinTeamParams{
		UserID:      userID,
		TeamID: team.ID,
	})
	if err != nil {
		return db.Team{}, apperror.ErrJoinTeam
	}

	return team, nil
}

func (s *TeamService) ListTeams(ctx context.Context, req requests.ListTeamsRequest) ([]db.Team, error) {
	teams, err := s.queries.ListTeams(ctx, db.ListTeamsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return []db.Team{}, apperror.ErrListTeams
	}

	return teams, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, teamID, userID uuid.UUID, req requests.UpdateTeamRequest) (db.Team, error) {
	member, err := s.queries.GetTeamMember(ctx, db.GetTeamMemberParams{
		UserID: userID,
		TeamID: teamID,
	})
	if err != nil || (member.Role.String != "moderator" && member.Role.String != "admin") {
		return db.Team{}, apperror.ErrForbidden
	}

	return s.queries.UpdateTeam(ctx, db.UpdateTeamParams{
		ID: teamID,
		Title: req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	})
}
