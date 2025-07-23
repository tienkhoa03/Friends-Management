package dto

type CreateFriendshipRequest struct {
	Friends []string `json:"friends" binding:"required,len=2"`
}
