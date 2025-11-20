package entity

type UserSetActivityParams struct {
	UserId             string
	NewActivenessValue bool
}

type UserSetActivityResult struct {
	User  User
	Found bool
}

type User struct {
	UserId   string
	Username string
	TeamName string
	IsActive bool
}
