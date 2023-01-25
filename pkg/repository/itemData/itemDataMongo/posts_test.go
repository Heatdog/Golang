package itemdatamongo

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"strings"
	"testing"
)

func marshalPost(post itemdata.Post) bson.D {
	var bsonData []byte
	bsonData, _ = bson.Marshal(post)
	var bsonD bson.D
	_ = bson.Unmarshal(bsonData, &bsonD)
	return bsonD
}

func marshalPosts(posts []itemdata.Post) []bson.D {
	docs := make([]bson.D, 0)
	for _, post := range posts {
		bsonData, _ := bson.Marshal(post)
		var bsonD bson.D
		_ = bson.Unmarshal(bsonData, &bsonD)
		docs = append(docs, bsonD)
	}
	return docs
}

func TestPosts_Create(t *testing.T) {
	post := itemdata.Post{
		ID: "",
		Ath: itemdata.Author{
			ID:       "1",
			Username: "123",
		},
		Comments:         []itemdata.Comment{},
		Cat:              "music",
		Score:            1,
		Type:             "text",
		Title:            "213",
		Created:          "2022-11-04T17:55:14Z",
		UpvotePercentage: 100,
		Views:            1,
		Text:             "123",
		Vote:             []itemdata.Votes{},
	}
	t.Parallel()
	testCases := []struct {
		name      string
		inputPost itemdata.Post
		mongoRes  bson.D
		wantErr   error
	}{

		{
			name:      "ok",
			inputPost: post,
			mongoRes:  mtest.CreateSuccessResponse(),
			wantErr:   nil,
		},

		{
			name:      "insert error",
			inputPost: post,
			mongoRes: mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    123,
				Message: "invalid insert",
			}),
			wantErr: errors.New("invalid insert"),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tc := range testCases {
		tc := tc
		mt.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.mongoRes)
			mongo := NewItemDataMongo(mt.DB.Collection("postDB"), context.Background())
			res, gotErr := mongo.CreatePost(tc.inputPost)
			if tc.wantErr == nil {
				if gotErr != nil {
					t.Errorf("unexpected error")
				}
				if res.ID == "" {
					t.Errorf("no text found")
				}
			} else {
				if !strings.Contains(gotErr.Error(), tc.wantErr.Error()) {
					t.Errorf("unexpected error")
				}
			}
		})
	}
}

func TestPosts_GetPosts(t *testing.T) {
	posts := []itemdata.Post{
		{
			ID: "123",
			Ath: itemdata.Author{
				ID:       "1",
				Username: "123",
			},
			Comments:         []itemdata.Comment{},
			Cat:              "music",
			Score:            1,
			Type:             "text",
			Title:            "213",
			Created:          "2022-11-04T17:55:14Z",
			UpvotePercentage: 100,
			Views:            1,
			Text:             "123",
			Vote:             []itemdata.Votes{},
		},
	}
	testCases := []struct {
		name     string
		postsRes []itemdata.Post
		mongoRes func(mt *mtest.T, posts []itemdata.Post) []bson.D
		wantErr  error
	}{
		{
			name: "ok",
			mongoRes: func(mt *mtest.T, posts []itemdata.Post) []bson.D {
				bsonD := marshalPosts(posts)
				return []bson.D{mtest.CreateCursorResponse(1, "postDB.postDB", mtest.FirstBatch, bsonD...),
					mtest.CreateCursorResponse(0, "postDB.postDB", mtest.NextBatch)}
			},
			postsRes: posts,
			wantErr:  nil,
		},
		{
			name:     "get error",
			postsRes: nil,
			mongoRes: func(mt *mtest.T, posts []itemdata.Post) []bson.D {
				return []bson.D{
					mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "invalid server posts",
					}),
				}
			},
			wantErr: errors.New("invalid server posts"),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tc := range testCases {
		tc := tc
		mt.Run(tc.name, func(mt *mtest.T) {
			mongo := NewItemDataMongo(mt.DB.Collection("postDB"), context.Background())
			mt.AddMockResponses(tc.mongoRes(mt, tc.postsRes)...)
			res, err := mongo.GetPosts()
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error")
				}
				for i := range res {
					require.EqualValues(mt, res[i].ID, tc.postsRes[i].ID)
					require.EqualValues(mt, res[i].Text, tc.postsRes[i].Text)
					require.EqualValues(mt, res[i].Cat, tc.postsRes[i].Cat)
					require.EqualValues(mt, res[i].Type, tc.postsRes[i].Type)
					require.EqualValues(mt, res[i].Title, tc.postsRes[i].Title)
					require.EqualValues(mt, res[i].Created, tc.postsRes[i].Created)
				}
			} else {
				if !strings.Contains(err.Error(), tc.wantErr.Error()) {
					t.Errorf("unexpected error")
				}
			}
		})
	}
}

func TestPosts_GetCategory(t *testing.T) {
	posts := []itemdata.Post{
		{
			ID: "123",
			Ath: itemdata.Author{
				ID:       "1",
				Username: "123",
			},
			Comments:         []itemdata.Comment{},
			Cat:              "music",
			Score:            1,
			Type:             "text",
			Title:            "213",
			Created:          "2022-11-04T17:55:14Z",
			UpvotePercentage: 100,
			Views:            1,
			Text:             "123",
			Vote:             []itemdata.Votes{},
		},
	}
	testCases := []struct {
		name     string
		postsRes []itemdata.Post
		mongoRes func(mt *mtest.T, posts []itemdata.Post) []bson.D
		wantErr  error
		category string
	}{
		{
			name: "ok",
			mongoRes: func(mt *mtest.T, posts []itemdata.Post) []bson.D {
				bsonD := marshalPosts(posts)
				return []bson.D{mtest.CreateCursorResponse(1, "postDB.postDB", mtest.FirstBatch, bsonD...),
					mtest.CreateCursorResponse(0, "postDB.postDB", mtest.NextBatch)}
			},
			postsRes: posts,
			wantErr:  nil,
		},
		{
			name:     "get error",
			postsRes: nil,
			mongoRes: func(mt *mtest.T, posts []itemdata.Post) []bson.D {
				return []bson.D{
					mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "invalid server posts",
					}),
				}
			},
			wantErr:  errors.New("invalid server posts"),
			category: "music",
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tc := range testCases {
		tc := tc
		mt.Run(tc.name, func(mt *mtest.T) {
			mongo := NewItemDataMongo(mt.DB.Collection("postDB"), context.Background())
			mt.AddMockResponses(tc.mongoRes(mt, tc.postsRes)...)
			res, err := mongo.GetCategory(tc.category)
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error")
				}
				for i := range res {
					require.EqualValues(mt, res[i].ID, tc.postsRes[i].ID)
					require.EqualValues(mt, res[i].Text, tc.postsRes[i].Text)
					require.EqualValues(mt, res[i].Cat, tc.postsRes[i].Cat)
					require.EqualValues(mt, res[i].Type, tc.postsRes[i].Type)
					require.EqualValues(mt, res[i].Title, tc.postsRes[i].Title)
					require.EqualValues(mt, res[i].Created, tc.postsRes[i].Created)
				}
			} else {
				if !strings.Contains(err.Error(), tc.wantErr.Error()) {
					t.Errorf("unexpected error")
				}
			}
		})
	}
}

func TestPosts_GetName(t *testing.T) {
	posts := []itemdata.Post{
		{
			ID: "123",
			Ath: itemdata.Author{
				ID:       "1",
				Username: "123",
			},
			Comments:         []itemdata.Comment{},
			Cat:              "music",
			Score:            1,
			Type:             "text",
			Title:            "213",
			Created:          "2022-11-04T17:55:14Z",
			UpvotePercentage: 100,
			Views:            1,
			Text:             "123",
			Vote:             []itemdata.Votes{},
		},
	}
	testCases := []struct {
		name     string
		postsRes []itemdata.Post
		mongoRes func(mt *mtest.T, posts []itemdata.Post) []bson.D
		wantErr  error
		nameUser string
	}{
		{
			name: "ok",
			mongoRes: func(mt *mtest.T, posts []itemdata.Post) []bson.D {
				bsonD := marshalPosts(posts)
				return []bson.D{mtest.CreateCursorResponse(1, "postDB.postDB", mtest.FirstBatch, bsonD...),
					mtest.CreateCursorResponse(0, "postDB.postDB", mtest.NextBatch)}
			},
			postsRes: posts,
			wantErr:  nil,
			nameUser: "123",
		},
		{
			name:     "get error",
			postsRes: nil,
			mongoRes: func(mt *mtest.T, posts []itemdata.Post) []bson.D {
				return []bson.D{
					mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "invalid server posts",
					}),
				}
			},
			wantErr:  errors.New("invalid server posts"),
			nameUser: "123",
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tc := range testCases {
		tc := tc
		mt.Run(tc.name, func(mt *mtest.T) {
			mongo := NewItemDataMongo(mt.DB.Collection("postDB"), context.Background())
			mt.AddMockResponses(tc.mongoRes(mt, tc.postsRes)...)
			res, err := mongo.GetName(tc.nameUser)
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error")
				}
				for i := range res {
					require.EqualValues(mt, res[i].ID, tc.postsRes[i].ID)
					require.EqualValues(mt, res[i].Text, tc.postsRes[i].Text)
					require.EqualValues(mt, res[i].Cat, tc.postsRes[i].Cat)
					require.EqualValues(mt, res[i].Type, tc.postsRes[i].Type)
					require.EqualValues(mt, res[i].Title, tc.postsRes[i].Title)
					require.EqualValues(mt, res[i].Created, tc.postsRes[i].Created)
				}
			} else {
				if !strings.Contains(err.Error(), tc.wantErr.Error()) {
					t.Errorf("unexpected error")
				}
			}
		})
	}
}

func TestPosts_GetPostID(t *testing.T) {
	posts := itemdata.Post{
		ID: "123",
		Ath: itemdata.Author{
			ID:       "1",
			Username: "123",
		},
		Comments:         []itemdata.Comment{},
		Cat:              "music",
		Score:            1,
		Type:             "text",
		Title:            "213",
		Created:          "2022-11-04T17:55:14Z",
		UpvotePercentage: 100,
		Views:            1,
		Text:             "123",
		Vote:             []itemdata.Votes{},
	}
	testCases := []struct {
		name     string
		postsRes itemdata.Post
		mongoRes func(mt *mtest.T, posts itemdata.Post) []bson.D
		wantErr  error
		postID   string
	}{
		{
			name: "ok",
			mongoRes: func(mt *mtest.T, posts itemdata.Post) []bson.D {
				return []bson.D{mtest.CreateCursorResponse(1, "postDB.postDB", mtest.FirstBatch, marshalPost(posts)),
					mtest.CreateCursorResponse(0, "postDB.postDB", mtest.NextBatch)}
			},
			postsRes: posts,
			wantErr:  nil,
			postID:   "123",
		},
		{
			name:     "get error",
			postsRes: itemdata.Post{},
			mongoRes: func(mt *mtest.T, posts itemdata.Post) []bson.D {
				return []bson.D{
					mtest.CreateWriteErrorsResponse(mtest.WriteError{
						Index:   1,
						Code:    123,
						Message: "invalid post id",
					}),
				}
			},
			wantErr: errors.New("invalid post id"),
			postID:  "123",
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tc := range testCases {
		tc := tc
		mt.Run(tc.name, func(mt *mtest.T) {
			mongo := NewItemDataMongo(mt.DB.Collection("postDB"), context.Background())
			mt.AddMockResponses(tc.mongoRes(mt, tc.postsRes)...)
			res, err := mongo.GetPostID(tc.postID)
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error")
				}
				require.EqualValues(mt, res.ID, tc.postsRes.ID)
				require.EqualValues(mt, res.Text, tc.postsRes.Text)
				require.EqualValues(mt, res.Cat, tc.postsRes.Cat)
				require.EqualValues(mt, res.Type, tc.postsRes.Type)
				require.EqualValues(mt, res.Title, tc.postsRes.Title)
				require.EqualValues(mt, res.Created, tc.postsRes.Created)
			} else {
				if !strings.Contains(err.Error(), tc.wantErr.Error()) {
					t.Errorf("unexpected error")
				}
			}
		})
	}
}

func TestPosts_SetPost(t *testing.T) {
	post := itemdata.Post{
		ID: "123",
		Ath: itemdata.Author{
			ID:       "1",
			Username: "123",
		},
		Comments:         []itemdata.Comment{},
		Cat:              "music",
		Score:            1,
		Type:             "text",
		Title:            "213",
		Created:          "2022-11-04T17:55:14Z",
		UpvotePercentage: 100,
		Views:            1,
		Text:             "123",
		Vote:             []itemdata.Votes{},
	}
	t.Parallel()
	testCases := []struct {
		name      string
		inputPost itemdata.Post
		mongoRes  bson.D
		wantErr   error
	}{

		{
			name:      "ok",
			inputPost: post,
			mongoRes:  mtest.CreateSuccessResponse(),
			wantErr:   nil,
		},

		{
			name:      "replace error",
			inputPost: post,
			mongoRes: mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    123,
				Message: "invalid post id",
			}),
			wantErr: errors.New("invalid post id"),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tc := range testCases {
		tc := tc
		mt.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.mongoRes)
			mongo := NewItemDataMongo(mt.DB.Collection("postDB"), context.Background())
			gotErr := mongo.SetPost(tc.inputPost)
			if tc.wantErr == nil {
				if gotErr != nil {
					t.Errorf("unexpected error")
				}
			} else {
				if !strings.Contains(gotErr.Error(), tc.wantErr.Error()) {
					t.Errorf("unexpected error")
				}
			}
		})
	}
}

func TestPosts_DeletePost(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		inputID  string
		mongoRes bson.D
		wantErr  error
	}{

		{
			name:     "ok",
			inputID:  "123",
			mongoRes: mtest.CreateSuccessResponse(),
			wantErr:  nil,
		},

		{
			name:    "replace error",
			inputID: "123",
			mongoRes: mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    123,
				Message: "invalid post id",
			}),
			wantErr: errors.New("invalid post id"),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, tc := range testCases {
		tc := tc
		mt.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.mongoRes)
			mongo := NewItemDataMongo(mt.DB.Collection("postDB"), context.Background())
			gotErr := mongo.DeletePost(tc.inputID)
			if tc.wantErr == nil {
				if gotErr != nil {
					t.Errorf("unexpected error")
				}
			} else {
				if !strings.Contains(gotErr.Error(), tc.wantErr.Error()) {
					t.Errorf("unexpected error")
				}
			}
		})
	}
}
