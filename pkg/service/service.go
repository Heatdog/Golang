package service

import (
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/session"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user userdata.User) (userdata.User, error)
	GenerateToken(login, password string) (string, error)
	GetHash(password string) string
}

type Posts interface {
	CreatePost(post itemdata.CreatePost, userID string) (itemdata.Post, error)
	GetPosts() ([]itemdata.Post, error)
	GetCategory(category string) ([]itemdata.Post, error)
	GetName(login string) ([]itemdata.Post, error)
	GetPostID(id string) (itemdata.Post, error)
	DeletePost(id string, userID string) error
}

type Comments interface {
	CreateComm(postID, userID, comment string) (itemdata.Post, error)
	DeleteComm(postID, userID, commID string) (itemdata.Post, error)
	Upvote(postID, userID string) (itemdata.Post, error)
	Downvote(postID, userID string) (itemdata.Post, error)
	Unvote(postID, userID string) (itemdata.Post, error)
}

type Service struct {
	Authorization
	Posts
	Comments
}

func NewService(userDat userdata.UserData, itemDat itemdata.ItemData, sessionManager session.SesManager, key []byte) *Service {
	return &Service{
		Authorization: NewAuthService(userDat, sessionManager, key),
		Posts:         NewPostService(userDat, itemDat),
		Comments:      NewCommentService(userDat, itemDat),
	}
}
