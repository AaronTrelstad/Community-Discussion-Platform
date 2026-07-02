package handlers

import (
	"net/http"

	"github.com/aarontrelstad/api/internal/services"
	util "github.com/aarontrelstad/api/pkg/httputil"
	"github.com/aarontrelstad/api/pkg/render"
	"github.com/aarontrelstad/api/pkg/requests"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TeamsHandler struct {
	service *services.TeamService
}

func NewTeamHandler(service *services.TeamService) *TeamsHandler {
	return &TeamsHandler{service: service}
}

func (h *TeamsHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateTeamRequest
	if err := util.Decode(r, &req); err != nil {
		render.HandleError(w, err)
		return
	}

	userID, err := util.GetUserID(r)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	team, err := h.service.CreateTeam(r.Context(), userID, req)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusCreated, team)
}

func (h *TeamsHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	var req requests.ListTeamsRequest
	if err := util.Decode(r, &req); err != nil {
		render.HandleError(w, err)
		return
	}
	
	teams, err := h.service.ListTeams(r.Context(), req)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, teams)
}

func (h *TeamsHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	teamID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.HandleError(w, err)
		return
	}
	
	var req requests.UpdateTeamRequest
	if err := util.Decode(r, &req); err != nil {
		render.HandleError(w, err)
		return
	}

	userID, err := util.GetUserID(r)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	team, err := h.service.UpdateTeam(r.Context(), teamID, userID, req)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, team)
}

func (h *TeamsHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.HandleError(w, err)
		return
	}

	team, err := h.service.GetTeam(r.Context(), teamID)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, team)
}

func (h *TeamsHandler) JoinTeam(w http.ResponseWriter, r *http.Request) {
	teamID, err := uuid.Parse(chi.URLParam(r, "teamID"))
	if err != nil {
		render.HandleError(w, err)
		return
	}

	userID, err := util.GetUserID(r)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	joined, err := h.service.JoinTeam(r.Context(), teamID, userID)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, joined)
}

func (h *TeamsHandler) LeaveTeam(w http.ResponseWriter, r *http.Request) {
	teamID, err := uuid.Parse(chi.URLParam(r, "teamID"))
	if err != nil {
		render.HandleError(w, err)
		return
	}

	userID, err := util.GetUserID(r)
	if err != nil {
		render.HandleError(w, err)
		return
	}
	
	left, err := h.service.LeaveTeam(r.Context(), teamID, userID)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, left)
}
