package http

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

type MessageResponse struct {
	Message string `json:"message" example:"operation completed successfully"`
}
