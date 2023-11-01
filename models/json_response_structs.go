package models

type ErrorResponseWrapper struct {
	Error ErrorResponse `json:"error"`
}

type SuccessResponseWrapper struct {
	Success SuccessResponse `json:"success"`
}

type GitHubClientIdResponseWrapper struct {
	GitHubClientId GitHubClientIdResponse `json:"github_client_id"`
}

type AccessCodeResponseWrapper struct {
	AccessCode AccessCodeResponse `json:"access_code"`
}
