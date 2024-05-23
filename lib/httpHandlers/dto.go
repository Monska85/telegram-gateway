package httpHandlers

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type SendMessageRequest struct {
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
	Image     string `json:"image,omitempty"`
	ImageName string `json:"image_name,omitempty"` // Optional, used for image file name
}
