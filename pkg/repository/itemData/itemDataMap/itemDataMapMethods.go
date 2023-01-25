package itemdatamap

import (
	"errors"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
	"sort"
)

func (dt *itemDataMap) CreatePost(post itemdata.Post) (itemdata.Post, error) {
	post.ID = utils.RandomHex()
	dt.mux.Lock()
	defer dt.mux.Unlock()
	dt.data[post.ID] = post
	return post, nil
}

func (dt *itemDataMap) GetPosts() ([]itemdata.Post, error) {
	res := make([]itemdata.Post, 0, len(dt.data))
	dt.mux.RLock()
	for _, el := range dt.data {
		res = append(res, el)
	}
	dt.mux.RUnlock()
	res = dt.sort(res)
	return res, nil
}

func (dt *itemDataMap) GetCategory(category string) ([]itemdata.Post, error) {
	res := make([]itemdata.Post, 0, len(dt.data))
	dt.mux.RLock()
	for _, el := range dt.data {
		if el.Cat == category {
			res = append(res, el)
		}
	}
	dt.mux.RUnlock()
	res = dt.sort(res)
	return res, nil
}

func (dt *itemDataMap) GetName(login string) ([]itemdata.Post, error) {
	res := make([]itemdata.Post, 0, len(dt.data))
	dt.mux.RLock()
	for _, el := range dt.data {
		if el.Ath.Username == login {
			res = append(res, el)
		}
	}
	dt.mux.RUnlock()
	res = dt.sort(res)
	return res, nil
}

func (dt *itemDataMap) GetPostID(id string) (itemdata.Post, error) {
	dt.mux.RLock()
	defer dt.mux.RUnlock()
	post, ok := dt.data[id]
	if !ok {
		return post, errors.New("invalid post id")
	}
	return post, nil
}

func (dt *itemDataMap) SetPost(post itemdata.Post) error {
	dt.mux.Lock()
	defer dt.mux.Unlock()
	dt.data[post.ID] = post
	return nil
}

func (dt *itemDataMap) DeletePost(postID string) error {
	dt.mux.Lock()
	defer dt.mux.Unlock()
	delete(dt.data, postID)
	return nil
}

func (dt *itemDataMap) sort(posts []itemdata.Post) []itemdata.Post {
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Score == posts[j].Score {
			return posts[i].Created < posts[j].Created
		}
		return posts[i].Score > posts[j].Score
	})
	return posts
}
