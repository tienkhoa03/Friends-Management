package dto

type CreateSubscriptionRequest struct {
	Requestor string `json:"requestor" binding:"required"`
	Target    string `json:"target" binding:"required"`
}
