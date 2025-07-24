package repository

//go:generate mockgen -source=interface.go -destination=../mock/mock_friendship_repository.go

type FriendshipRepository interface {
	CreateFriendship(userId1, userId2 int64) error
	RetrieveFriendsList(userId int64) ([]int64, error)
}
