package handlers

import (
	"net/http"

	"github.com/aarontrelstad/api/internal/services"
	"github.com/aarontrelstad/api/pkg/apperror"
	util "github.com/aarontrelstad/api/pkg/httputil"
	"github.com/aarontrelstad/api/pkg/render"
	"github.com/aarontrelstad/api/pkg/requests"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CommunitiesHandler struct {
	service *services.CommunityService
}

func NewCommunitiesHandler(service *services.CommunityService) *CommunitiesHandler {
	return &CommunitiesHandler{service: service}
}

func (h *CommunitiesHandler) CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateCommunityRequest
	if err := util.Decode(r, &req); err != nil {
		render.HandleError(w, err)
		return
	}

	userID, err := util.GetUserID(r)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	community, err := h.service.CreateCommunity(r.Context(), userID, req)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusCreated, community)
}

func (h *CommunitiesHandler) ListCommunities(w http.ResponseWriter, r *http.Request) {
	var req requests.ListCommunitiesRequest
	if err := util.Decode(r, &req); err != nil {
		render.HandleError(w, err)
		return
	}
	
	communities, err := h.service.ListCommunities(r.Context(), req)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, communities)
}

func (h *CommunitiesHandler) UpdateCommunity(w http.ResponseWriter, r *http.Request) {
	communityID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.HandleError(w, apperror.ErrParseURL)
		return
	}
	
	var req requests.UpdateCommunityRequest
	if err := util.Decode(r, &req); err != nil {
		render.HandleError(w, err)
		return
	}

	userID, err := util.GetUserID(r)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	community, err := h.service.UpdateCommunity(r.Context(), communityID, userID, req)
	if err != nil {
		render.HandleError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, community)
}
