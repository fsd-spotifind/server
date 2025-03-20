package database

import (
	"context"
	"fmt"
	db "fsd-backend/prisma/db"
)

type UserPublic struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type FriendWithUsers struct {
	ID        string          `json:"id"`
	Status    db.FriendStatus `json:"status"`
	UserOne   UserPublic      `json:"userOne"`
	UserTwo   UserPublic      `json:"userTwo"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

func (s *service) AddFriend(ctx context.Context, userId, friendId string) (*db.FriendModel, error) {
	userOne := userId
	userTwo := friendId
	existingFriend, err := s.client.Friend.FindFirst(
		db.Friend.Or(
			db.Friend.And(
				db.Friend.UserOneID.Equals(userOne),
				db.Friend.UserTwoID.Equals(userTwo),
			),
			db.Friend.And(
				db.Friend.UserOneID.Equals(userTwo),
				db.Friend.UserTwoID.Equals(userOne),
			),
		),
	).Exec(ctx)
	if err != nil && err.Error() != "ErrNotFound" {
		return nil, err
	}
	if existingFriend != nil {
		return nil, fmt.Errorf("friendship already exists")
	}
	return s.client.Friend.CreateOne(
		db.Friend.Status.Set(db.FriendStatusPending),
		db.Friend.UserOne.Link(db.User.ID.Equals(userOne)),
		db.Friend.UserTwo.Link(db.User.ID.Equals(userTwo)),
	).Exec(ctx)
}

func (s *service) GetFriendRequests(ctx context.Context, userId string) ([]FriendWithUsers, error) {
	friendRequests, err := s.client.Friend.FindMany(
		db.Friend.UserTwoID.Equals(userId),
		db.Friend.Status.Equals(db.FriendStatusPending),
	).With(
		db.Friend.UserOne.Fetch(),
		db.Friend.UserTwo.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	friendRequestsList := make([]FriendWithUsers, len(friendRequests))
	for i, f := range friendRequests {
		friendRequestsList[i] = FriendWithUsers{
			ID:        f.ID,
			Status:    f.Status,
			UserOne:   UserPublic{ID: f.UserOne().ID, Email: f.UserOne().Email, Username: f.UserOne().Name},
			UserTwo:   UserPublic{ID: f.UserTwo().ID, Email: f.UserTwo().Email, Username: f.UserTwo().Name},
			CreatedAt: f.CreatedAt.String(),
			UpdatedAt: f.UpdatedAt.String(),
		}
	}
	return friendRequestsList, nil
}

func (s *service) AcceptFriendRequest(ctx context.Context, userId, requestId string) (*db.FriendModel, error) {
	friend, err := s.client.Friend.FindUnique(
		db.Friend.ID.Equals(requestId),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("error finding friend request: %v", err)
	}

	if friend == nil || friend.UserTwoID != userId {
		return nil, fmt.Errorf("friend request not found or not intended for this user")
	}

	return s.client.Friend.FindUnique(
		db.Friend.ID.Equals(requestId),
	).Update(
		db.Friend.Status.Set(db.FriendStatusActive),
	).Exec(ctx)
}

func (s *service) GetFriends(ctx context.Context, userId string) ([]FriendWithUsers, error) {
	friends, err := s.client.Friend.FindMany(
		db.Friend.Or(
			db.Friend.UserOneID.Equals(userId),
			db.Friend.UserTwoID.Equals(userId),
		),
		db.Friend.Status.Equals(db.FriendStatusActive),
	).With(
		db.Friend.UserOne.Fetch(),
		db.Friend.UserTwo.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}
	friendList := make([]FriendWithUsers, len(friends))
	for i, f := range friends {
		friendList[i] = FriendWithUsers{
			ID:        f.ID,
			Status:    f.Status,
			UserOne:   UserPublic{ID: f.UserOne().ID, Email: f.UserOne().Email, Username: f.UserOne().Name},
			UserTwo:   UserPublic{ID: f.UserTwo().ID, Email: f.UserTwo().Email, Username: f.UserTwo().Name},
			CreatedAt: f.CreatedAt.String(),
			UpdatedAt: f.UpdatedAt.String(),
		}
	}
	return friendList, nil
}
