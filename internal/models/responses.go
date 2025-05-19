package models

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type MessageResponse struct {
	Details string `json:"details"`
}
