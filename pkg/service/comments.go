package service

import (
	"errors"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
	"time"
)

var _ Comments = (*CommentService)(nil)

type CommentService struct {
	dbUser  userdata.UserData
	dbItems itemdata.ItemData
}

func NewCommentService(dbUser userdata.UserData, dbItems itemdata.ItemData) *CommentService {
	return &CommentService{
		dbUser:  dbUser,
		dbItems: dbItems,
	}
}

func (cmServ *CommentService) CreateComm(postID, userID, comment string) (itemdata.Post, error) {
	user, err := cmServ.dbUser.GetUser(userID)
	if err != nil {
		return itemdata.Post{}, err
	}
	comm := itemdata.Comment{
		Ath: itemdata.Author{
			Username: user.Login,
			ID:       user.ID,
		},
		Body:    comment,
		Created: time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		ID:      utils.RandomHex(),
	}
	post, err := cmServ.dbItems.GetPostID(postID)
	if err != nil {
		return itemdata.Post{}, err
	}
	post.Comments = append(post.Comments, comm)
	if err = cmServ.dbItems.SetPost(post); err != nil {
		return post, err
	}
	return post, nil
}

func (cmServ *CommentService) DeleteComm(postID, userID, commID string) (itemdata.Post, error) {
	post, err := cmServ.dbItems.GetPostID(postID)
	if err != nil {
		return itemdata.Post{}, err
	}
	comment, ok := cmServ.findComment(post, commID)
	if !ok {
		return itemdata.Post{}, errors.New("invalid comment id")
	}
	if comment.Ath.ID != userID {
		return itemdata.Post{}, errors.New("invalid user id")
	}
	post, ok = cmServ.deleteComment(comment, post)
	if !ok {
		return itemdata.Post{}, errors.New("invalid comment id")
	}
	if err = cmServ.dbItems.SetPost(post); err != nil {
		return post, err
	}
	return post, nil
}

func (cmServ *CommentService) Upvote(postID, userID string) (itemdata.Post, error) {
	return cmServ.vote(postID, userID, 1)
}

func (cmServ *CommentService) Downvote(postID, userID string) (itemdata.Post, error) {
	return cmServ.vote(postID, userID, -1)
}

func (cmServ *CommentService) Unvote(postID, userID string) (itemdata.Post, error) {
	post, err := cmServ.dbItems.GetPostID(postID)
	if err != nil {
		return itemdata.Post{}, err
	}
	old, ok := cmServ.findVote(post, userID)
	if !ok {
		return itemdata.Post{}, errors.New("invalid vote")
	}
	post = cmServ.deleteVote(post, userID)
	post.Score -= old
	post.UpvotePercentage = cmServ.getPercentage(post)
	if err = cmServ.dbItems.SetPost(post); err != nil {
		return post, err
	}
	return post, nil
}

func (cmServ CommentService) findComment(post itemdata.Post, commID string) (itemdata.Comment, bool) {
	for _, el := range post.Comments {
		if el.ID == commID {
			return el, true
		}
	}
	return itemdata.Comment{}, false
}

func (cmServ CommentService) deleteComment(comment itemdata.Comment, post itemdata.Post) (itemdata.Post, bool) {
	for i, el := range post.Comments {
		if el.ID == comment.ID {
			post.Comments = append(post.Comments[:i], post.Comments[i+1:]...)
			return post, true
		}
	}
	return itemdata.Post{}, false
}

func (cmServ CommentService) findVote(post itemdata.Post, userID string) (int, bool) {
	for _, el := range post.Vote {
		if el.User == userID {
			return el.Vote, true
		}
	}
	return 0, false
}

func (cmServ CommentService) deleteVote(post itemdata.Post, userID string) itemdata.Post {
	for i, el := range post.Vote {
		if el.User == userID {
			post.Vote = append(post.Vote[:i], post.Vote[i+1:]...)
			return post
		}
	}
	return itemdata.Post{}
}

func (cmServ CommentService) getPercentage(post itemdata.Post) int {
	pos := 0
	var i int
	var el itemdata.Votes
	for i, el = range post.Vote {
		if el.Vote == 1 {
			pos++
		}
	}
	return (pos * 100) / (i + 1)
}

func (cmServ CommentService) vote(postID, userID string, diff int) (itemdata.Post, error) {
	post, err := cmServ.dbItems.GetPostID(postID)
	if err != nil {
		return itemdata.Post{}, err
	}
	old, ok := cmServ.findVote(post, userID)
	if ok {
		post = cmServ.deleteVote(post, userID)
		post.Score -= old
	}
	post.Vote = append(post.Vote, itemdata.Votes{
		User: userID,
		Vote: diff,
	})
	post.Score += diff
	post.UpvotePercentage = cmServ.getPercentage(post)
	if err = cmServ.dbItems.SetPost(post); err != nil {
		return post, err
	}
	return post, nil
}
