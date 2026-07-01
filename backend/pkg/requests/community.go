package requests

type CreateCommunityRequest struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListCommunitiesRequest struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type UpdateCommunityRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
