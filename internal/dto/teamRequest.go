package dto

type TeamAdd_Request Team_Request

type Team_Request struct {
	TeamName string               `json:"team_name" validate:"required"`
	Members  []TeamMember_Request `json:"members" validate:"required"`
}

type TeamMember_Request struct {
	UserId   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type TeamNameQuery_Request struct {
	TeamName string `query:"team_name" validate:"required"`
}
