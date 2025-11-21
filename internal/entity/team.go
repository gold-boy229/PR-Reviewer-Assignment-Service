package entity

type Team struct {
	TeamName string
	Members  []TeamMember
}

type TeamMember struct {
	UserId   string
	Username string
	IsActive bool
}

type TeamSearchParams struct {
	TeamName string
}

type TeamSearchResult struct {
	Team               Team
	FoundTeam          bool
	ConflictingUserIds []string
}
