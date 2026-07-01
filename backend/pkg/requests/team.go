package requests

type CreateTeamRequest struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListTeamsRequest struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type UpdateTeamRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
