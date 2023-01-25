package itemdatamongo

import (
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dt *itemDataMongo) CreatePost(post itemdata.Post) (itemdata.Post, error) {
	post.ID = utils.RandomHex()
	_, err := dt.collection.InsertOne(dt.ctx, post)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (dt *itemDataMongo) GetPosts() ([]itemdata.Post, error) {
	return dt.sort(bson.D{})
}

func (dt *itemDataMongo) GetCategory(category string) ([]itemdata.Post, error) {
	return dt.sort(bson.D{{"category", category}})
}

func (dt *itemDataMongo) GetName(login string) ([]itemdata.Post, error) {
	return dt.sort(bson.D{{"username", login}})
}

func (dt *itemDataMongo) sort(m bson.D) ([]itemdata.Post, error) {
	res := make([]itemdata.Post, 0, 10)
	opts := options.Find().SetSort(bson.D{{"score", -1}, {"created", 1}})
	posts, err := dt.collection.Find(dt.ctx, m, opts)
	if err != nil {
		return nil, err
	}
	if err = posts.All(dt.ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (dt *itemDataMongo) GetPostID(id string) (itemdata.Post, error) {
	var post itemdata.Post
	if err := dt.collection.FindOne(dt.ctx, bson.M{"_id": id}).Decode(&post); err != nil {
		return post, err
	}
	return post, nil
}

func (dt *itemDataMongo) SetPost(post itemdata.Post) error {
	_, err := dt.collection.ReplaceOne(dt.ctx, bson.M{"_id": post.ID}, post)
	if err != nil {
		return err
	}
	return nil
}

func (dt *itemDataMongo) DeletePost(postID string) error {
	_, err := dt.collection.DeleteOne(dt.ctx, bson.M{"_id": postID})
	return err
}
