package repository

type FriendshipRepository interface {
	CreateFriendship(userId1, userId2 int64) error
}
