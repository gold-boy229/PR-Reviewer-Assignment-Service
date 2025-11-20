package dto

type ErrorResponse_Response struct {
	Error Error_Response `json:"error"`
}

type Error_Response struct {
	Code    string `json:"code" validate:"oneof=TEAM_EXISTS PR_EXISTS PR_MERGED NOT_ASSIGNED NO_CANDIDATE NOT_FOUND"`
	Message string `json:"message"`
}

func NewErrorResponse(code, message string) ErrorResponse_Response {
	return ErrorResponse_Response{
		Error: Error_Response{
			Code:    code,
			Message: message,
		},
	}
}
