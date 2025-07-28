package dto

type CreateBlockRequest struct {
	Requestor string `json:"requestor" binding:"required"`
	Target    string `json:"target" binding:"required"`
}
