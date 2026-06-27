package handlers

import (
	"net/http"

	"github.com/aarontrelstad/api/internal/services"
	"github.com/aarontrelstad/api/pkg/render"
	"github.com/aarontrelstad/api/pkg/requests"
	util "github.com/aarontrelstad/api/pkg/httputil"
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
