package dto

type TeamAdd_Response Team_Response

type Team_Response struct {
	TeamName string                `json:"team_name"`
	Members  []TeamMember_Response `json:"members"`
}

type TeamMember_Response struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
