package service

import (
	"errors"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"time"
)

var _ Posts = (*PostService)(nil)

type PostService struct {
	dbUser  userdata.UserData
	dbPosts itemdata.ItemData
}

func NewPostService(dbUser userdata.UserData, dbPosts itemdata.ItemData) *PostService {
	return &PostService{
		dbUser:  dbUser,
		dbPosts: dbPosts,
	}
}

func (postServ *PostService) CreatePost(post itemdata.CreatePost, userID string) (itemdata.Post, error) {
	us, err := postServ.dbUser.GetUser(userID)
	if err != nil {
		return itemdata.Post{}, err
	}
	resp := itemdata.Post{
		Ath: itemdata.Author{
			ID:       userID,
			Username: us.Login,
		},
		Comments:         make([]itemdata.Comment, 0, 10),
		Cat:              post.Cat,
		Score:            1,
		Type:             post.Type,
		Title:            post.Title,
		Created:          time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpvotePercentage: 100,
		Views:            0,
		Text:             post.Text,
		Vote: []itemdata.Votes{
			{
				Vote: 1,
				User: us.ID,
			},
		},
	}
	return postServ.dbPosts.CreatePost(resp)
}

func (postServ *PostService) GetPosts() ([]itemdata.Post, error) {
	return postServ.dbPosts.GetPosts()
}

func (postServ *PostService) GetCategory(category string) ([]itemdata.Post, error) {
	return postServ.dbPosts.GetCategory(category)
}

func (postServ *PostService) GetName(login string) ([]itemdata.Post, error) {
	return postServ.dbPosts.GetName(login)
}

func (postServ *PostService) GetPostID(id string) (itemdata.Post, error) {
	post, err := postServ.dbPosts.GetPostID(id)
	if err != nil {
		return itemdata.Post{}, err
	}
	post.Views++
	if err = postServ.dbPosts.SetPost(post); err != nil {
		return post, err
	}
	return post, nil
}

func (postServ *PostService) DeletePost(id string, userID string) error {
	post, err := postServ.dbPosts.GetPostID(id)
	if err != nil {
		return err
	}
	if post.Ath.ID != userID {
		return errors.New("invalid user id")
	}
	return postServ.dbPosts.DeletePost(id)
}
