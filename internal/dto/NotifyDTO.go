package dto

type NotifyUserRequest struct {
	To       string `json:"target"`
	Sub      string `json:"subject"`
	Channel  string `json:"channel"`
	BodyType string `json:"body_type"`
	Body     string `json:"body"`
}
