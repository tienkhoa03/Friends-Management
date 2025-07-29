package dto

type GetUpdateRecipientsRequest struct {
	Sender string `json:"sender" binding:"required"`
	Text   string `json:"text" binding:"required"`
}
