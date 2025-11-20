package model

type Team struct {
	TeamName string
	Members  []TeamMember
}

type TeamMember struct {
	UserId   string
	Username string
	IsActive bool
}
