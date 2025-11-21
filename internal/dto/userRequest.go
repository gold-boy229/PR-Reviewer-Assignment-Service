package dto

type UsersSetIsActive_Request struct {
	UserId             string `json:"user_id" validate:"required"`
	NewActivenessValue *bool  `json:"is_active" validate:"required"`
}

type UserGetReview_Request struct {
	UserId string `query:"user_id" validate:"required"`
}
