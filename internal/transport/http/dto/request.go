package dto

type CreateChatRequest struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
}

type CreateMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=5000"`
}
