package dto

type UsersSetIsActive_Response struct {
	User User_Response `json:"user"`
}

type User_Response struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
