package service

import (
	"BE_Friends_Management/internal/domain/entity"
	block_relationship "BE_Friends_Management/internal/repository/block_relationship"
	friendship "BE_Friends_Management/internal/repository/friendship"
	user "BE_Friends_Management/internal/repository/users"
	"errors"
	"strings"

	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=../mock/mock_friendship_service.go

type FriendshipService interface {
	CreateFriendship(authUserId int64, email1, email2 string) error
	RetrieveFriendsList(authUserId int64, email string) ([]*entity.User, error)
	RetrieveCommonFriends(authUserId int64, email1, email2 string) ([]*entity.User, error)
	CountFriends(friendsList []*entity.User) int64
}

type friendshipService struct {
	repo                  friendship.FriendshipRepository
	userRepo              user.UserRepository
	blockRelationshipRepo block_relationship.BlockRelationshipRepository
}

func NewFriendshipService(repo friendship.FriendshipRepository, userRepo user.UserRepository, blockRelationshipRepo block_relationship.BlockRelationshipRepository) FriendshipService {
	return &friendshipService{
		repo:                  repo,
		userRepo:              userRepo,
		blockRelationshipRepo: blockRelationshipRepo,
	}
}

func (service *friendshipService) CreateFriendship(authUserId int64, email1, email2 string) error {
	user1, err := service.userRepo.GetUserByEmail(email1)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	user2, err := service.userRepo.GetUserByEmail(email2)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if user1.Id != authUserId && user2.Id != authUserId {
		return ErrNotPermitted
	}
	if user1.Id == user2.Id {
		return ErrInvalidRequest
	}
	_, err = service.blockRelationshipRepo.GetBlockRelationship(user1.Id, user2.Id)
	if err == nil {
		return ErrIsBlocked
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	_, err = service.blockRelationshipRepo.GetBlockRelationship(user2.Id, user1.Id)
	if err == nil {
		return ErrIsBlocked
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if user1.Id < user2.Id {
		err = service.repo.CreateFriendship(user1.Id, user2.Id)
	} else {
		err = service.repo.CreateFriendship(user2.Id, user1.Id)
	}
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return ErrAlreadyFriend
	}
	return err
}

func (service *friendshipService) RetrieveFriendsList(authUserId int64, email string) ([]*entity.User, error) {
	user, err := service.userRepo.GetUserByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if authUserId != user.Id {
		return nil, ErrNotPermitted
	}
	friendIds, err := service.repo.RetrieveFriendIds(user.Id)
	if err != nil {
		return nil, err
	}
	friends := make([]*entity.User, len(friendIds))
	for i, id := range friendIds {
		friend, err := service.userRepo.GetUserById(id)
		if err != nil {
			return nil, err
		}
		friends[i] = friend
	}
	return friends, nil
}

func (service *friendshipService) RetrieveCommonFriends(authUserId int64, email1, email2 string) ([]*entity.User, error) {
	user1, err := service.userRepo.GetUserByEmail(email1)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user2, err := service.userRepo.GetUserByEmail(email2)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if authUserId != user1.Id && authUserId != user2.Id {
		return nil, ErrNotPermitted
	}
	if user1.Id == user2.Id {
		return nil, ErrInvalidRequest
	}
	friendIdsOfUser1, err := service.repo.RetrieveFriendIds(user1.Id)
	if err != nil {
		return nil, err
	}
	friendIdsOfUser2, err := service.repo.RetrieveFriendIds(user2.Id)
	if err != nil {
		return nil, err
	}

	set := make(map[int64]bool)
	commonFriendIds := []int64{}
	for _, id1 := range friendIdsOfUser1 {
		set[id1] = true
	}
	for _, id2 := range friendIdsOfUser2 {
		if set[id2] {
			commonFriendIds = append(commonFriendIds, id2)
		}
	}

	commonFriends := make([]*entity.User, len(commonFriendIds))
	for i, id := range commonFriendIds {
		commenFriend, err := service.userRepo.GetUserById(id)
		if err != nil {
			return nil, err
		}
		commonFriends[i] = commenFriend
	}
	return commonFriends, nil
}

func (service *friendshipService) CountFriends(friends []*entity.User) int64 {
	var count int64 = 0
	for _, friend := range friends {
		if friend != nil {
			count++
		}
	}
	return count
}
