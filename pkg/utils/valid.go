package utils

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
)

func UnmarshalCreatPost(buf []byte, post *itemdata.CreatePost) error {
	if err := json.Unmarshal(buf, &post); err != nil {
		return err
	}
	if post.Type == "link" {
		urlBody := struct {
			URL string `json:"url"`
		}{}
		if err := json.Unmarshal(buf, &urlBody); err != nil {
			return err
		}
		post.Text = urlBody.URL
		if !govalidator.IsURL(post.Text) {
			return errors.New("invalid URL")
		}
	} else {
		textBody := struct {
			URL string `json:"text"`
		}{}
		if err := json.Unmarshal(buf, &textBody); err != nil {
			return err
		}
		post.Text = textBody.URL
	}
	return nil
}

func MarshalPost(post itemdata.Post) ([]byte, error) {
	var resPost interface{}
	if post.Type == "link" {
		resPost = struct {
			itemdata.Post
			URL string `json:"url"`
		}{
			Post: post,
			URL:  post.Text,
		}
	} else {
		resPost = struct {
			itemdata.Post
			Text string `json:"text"`
		}{
			Post: post,
			Text: post.Text,
		}
	}
	resp, err := json.Marshal(&resPost)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func MarshalSlice(posts []itemdata.Post) ([]byte, error) {
	res := make([]byte, 0, 50)
	res = append(res, []byte("[")...)
	for i, el := range posts {
		r, err := MarshalPost(el)
		if err != nil {
			return nil, err
		}
		res = append(res, r...)
		if i != len(posts)-1 {
			res = append(res, []byte(",")...)
		}
	}
	res = append(res, []byte("]")...)
	return res, nil
}
