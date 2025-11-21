package dto

type UsersSetIsActive_Response struct {
	User User_Response `json:"user"`
}

type UserGetReview_Response struct {
	UserId       string                      `json:"user_id"`
	PullRequests []PullRequestShort_Response `json:"pull_requests"`
}

type User_Response struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
