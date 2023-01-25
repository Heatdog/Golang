package server

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/service"
	mockservice "gitlab.com/vk-go/lectures-2022-2/pkg/service/mocks"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer_CreatePost(t *testing.T) {
	type mockBehavior func(s *mockservice.MockPosts, userID string, post itemdata.CreatePost)
	testingTable := []struct {
		name              string
		inputBody         string
		inputUserID       string
		inputPost         itemdata.CreatePost
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name:        "ok",
			inputBody:   `{"category":"music","type":"text","title":"123","text":"123"}`,
			inputUserID: "1",
			inputPost: itemdata.CreatePost{
				Cat:   "music",
				Title: "123",
				Type:  "text",
				Text:  "123",
			},
			mockBehavior: func(s *mockservice.MockPosts, userID string, post itemdata.CreatePost) {
				s.EXPECT().CreatePost(post, userID).Return(itemdata.Post{
					ID: "abcd",
					Ath: itemdata.Author{
						ID:       userID,
						Username: "test",
					},
					Comments:         make([]itemdata.Comment, 0, 10),
					Cat:              post.Cat,
					Score:            1,
					Type:             post.Type,
					Title:            post.Title,
					Created:          "2022-11-04T17:55:14Z",
					UpvotePercentage: 100,
					Views:            0,
					Text:             post.Text,
					Vote: []itemdata.Votes{
						{
							User: userID,
							Vote: 1,
						},
					},
				}, nil)
			},
			expectStatusCode:  201,
			expectRequestBody: []byte(`{"id":"abcd","author":{"id":"1","username":"test"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":0,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name:              "invalid user id",
			inputBody:         ``,
			inputUserID:       "",
			inputPost:         itemdata.CreatePost{},
			mockBehavior:      func(s *mockservice.MockPosts, userID string, post itemdata.CreatePost) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid user id"}`),
		},
		{
			name:              "invalid json",
			inputBody:         "`{231321}`",
			inputUserID:       "1",
			inputPost:         itemdata.CreatePost{},
			mockBehavior:      func(s *mockservice.MockPosts, userID string, post itemdata.CreatePost) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid json input"}`),
		},
		{
			name:              "invalid fields",
			inputBody:         `{"type":"123","title":"123","text":"123"}`,
			inputUserID:       "1",
			inputPost:         itemdata.CreatePost{},
			mockBehavior:      func(s *mockservice.MockPosts, userID string, post itemdata.CreatePost) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid struct fields"`),
		},
		{
			name:        "invalid user id BD",
			inputBody:   `{"category":"music","type":"text","title":"123","text":"123"}`,
			inputUserID: "123",
			inputPost: itemdata.CreatePost{
				Cat:   "music",
				Title: "123",
				Type:  "text",
				Text:  "123",
			},
			mockBehavior: func(s *mockservice.MockPosts, userID string, post itemdata.CreatePost) {
				s.EXPECT().CreatePost(post, userID).Return(itemdata.Post{}, errors.New("invalid user id"))
			},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid user id"}`),
		},
	}
	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			posts := mockservice.NewMockPosts(c)
			testCase.mockBehavior(posts, testCase.inputUserID, testCase.inputPost)

			services := &service.Service{Posts: posts}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))

			r := httptest.NewRequest("POST", "/posts", bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			var ctx context.Context
			if testCase.name == "invalid user id" {
				ctx = r.Context()
			} else {
				ctx = context.WithValue(r.Context(), "id", testCase.inputUserID)
			}
			handler.CreatePost(w, r.WithContext(ctx))
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_GetPosts(t *testing.T) {
	type mockBehavior func(s *mockservice.MockPosts)
	testingTable := []struct {
		name              string
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name: "ok",
			mockBehavior: func(s *mockservice.MockPosts) {
				s.EXPECT().GetPosts().Return([]itemdata.Post{
					{
						ID: "1",
						Ath: itemdata.Author{
							ID:       "1",
							Username: "123",
						},
						Comments:         []itemdata.Comment{},
						Cat:              "music",
						Score:            1,
						Type:             "text",
						Title:            "123",
						Created:          "2022-11-04T17:55:14Z",
						UpvotePercentage: 100,
						Views:            1,
						Text:             "123",
						Vote: []itemdata.Votes{
							{
								User: "1",
								Vote: 1,
							},
						},
					},
				}, nil)
			},
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"id":"1","author":{"id":"1","username":"123"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name: "get server problems",
			mockBehavior: func(s *mockservice.MockPosts) {
				s.EXPECT().GetPosts().Return(nil, errors.New("invalid collection"))
			},
			expectStatusCode:  500,
			expectRequestBody: []byte(`"message":"invalid collection"`),
		},
	}

	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			posts := mockservice.NewMockPosts(c)
			testCase.mockBehavior(posts)

			services := &service.Service{Posts: posts}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))

			r := httptest.NewRequest("GET", "/posts", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			handler.GetPosts(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_GetCategory(t *testing.T) {
	type mockBehavior func(s *mockservice.MockPosts, category string)
	testingTable := []struct {
		name              string
		mockBehavior      mockBehavior
		category          string
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name: "ok",
			mockBehavior: func(s *mockservice.MockPosts, category string) {
				s.EXPECT().GetCategory(category).Return([]itemdata.Post{
					{
						ID: "1",
						Ath: itemdata.Author{
							ID:       "1",
							Username: "123",
						},
						Comments:         []itemdata.Comment{},
						Cat:              "music",
						Score:            1,
						Type:             "text",
						Title:            "123",
						Created:          "2022-11-04T17:55:14Z",
						UpvotePercentage: 100,
						Views:            1,
						Text:             "123",
						Vote: []itemdata.Votes{
							{
								User: "1",
								Vote: 1,
							},
						},
					},
				}, nil)
			},
			category:          "music",
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"id":"1","author":{"id":"1","username":"123"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name:              "invalid category",
			mockBehavior:      func(s *mockservice.MockPosts, category string) {},
			category:          "humans",
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid category"}`),
		},
		{
			name: "invalid collection",
			mockBehavior: func(s *mockservice.MockPosts, category string) {
				s.EXPECT().GetCategory(category).Return(nil, errors.New("invalid collection"))
			},
			category:          "music",
			expectStatusCode:  500,
			expectRequestBody: []byte(`{"message":"invalid collection"}`),
		},
	}

	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			posts := mockservice.NewMockPosts(c)
			testCase.mockBehavior(posts, testCase.category)

			services := &service.Service{Posts: posts}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))

			r := httptest.NewRequest("GET", "/posts/", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"category": testCase.category,
			})
			handler.GetCategory(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_GetUser(t *testing.T) {
	type mockBehavior func(s *mockservice.MockPosts, login string)
	testingTable := []struct {
		name              string
		mockBehavior      mockBehavior
		login             string
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name: "ok",
			mockBehavior: func(s *mockservice.MockPosts, login string) {
				s.EXPECT().GetName(login).Return([]itemdata.Post{
					{
						ID: "1",
						Ath: itemdata.Author{
							ID:       "1",
							Username: "123",
						},
						Comments:         []itemdata.Comment{},
						Cat:              "music",
						Score:            1,
						Type:             "text",
						Title:            "123",
						Created:          "2022-11-04T17:55:14Z",
						UpvotePercentage: 100,
						Views:            1,
						Text:             "123",
						Vote: []itemdata.Votes{
							{
								User: "1",
								Vote: 1,
							},
						},
					},
				}, nil)
			},
			login:             "123",
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"id":"1","author":{"id":"1","username":"123"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name: "invalid user login",
			mockBehavior: func(s *mockservice.MockPosts, login string) {
				s.EXPECT().GetName(login).Return(nil, errors.New("invalid user login"))
			},
			login:             "123",
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid user login"`),
		},
	}

	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			posts := mockservice.NewMockPosts(c)
			testCase.mockBehavior(posts, testCase.login)

			services := &service.Service{Posts: posts}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))

			r := httptest.NewRequest("GET", "/user/", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"user_login": testCase.login,
			})
			handler.GetUser(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_GetPostID(t *testing.T) {
	type mockBehavior func(s *mockservice.MockPosts, postID string)
	testingTable := []struct {
		name              string
		mockBehavior      mockBehavior
		postID            string
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name: "ok",
			mockBehavior: func(s *mockservice.MockPosts, postID string) {
				s.EXPECT().GetPostID(postID).Return(itemdata.Post{
					ID: postID,
					Ath: itemdata.Author{
						ID:       "1",
						Username: "123",
					},
					Comments:         []itemdata.Comment{},
					Cat:              "music",
					Score:            1,
					Type:             "text",
					Title:            "123",
					Created:          "2022-11-04T17:55:14Z",
					UpvotePercentage: 100,
					Views:            1,
					Text:             "123",
					Vote: []itemdata.Votes{
						{
							User: "1",
							Vote: 1,
						},
					},
				}, nil)
			},
			postID:            "1",
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"id":"1","author":{"id":"1","username":"123"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name: "invalid post id",
			mockBehavior: func(s *mockservice.MockPosts, postID string) {
				s.EXPECT().GetPostID(postID).Return(itemdata.Post{}, errors.New("invalid post id"))
			},
			postID:            "1",
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid post id"`),
		},
	}

	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			posts := mockservice.NewMockPosts(c)
			testCase.mockBehavior(posts, testCase.postID)

			services := &service.Service{Posts: posts}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))

			r := httptest.NewRequest("GET", "/post/", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"post_id": testCase.postID,
			})
			handler.GetPostID(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_DeletePost(t *testing.T) {
	type mockBehavior func(s *mockservice.MockPosts, postID, userID string)
	testingTable := []struct {
		name              string
		mockBehavior      mockBehavior
		postID            string
		userID            string
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name: "ok",
			mockBehavior: func(s *mockservice.MockPosts, postID, userID string) {
				s.EXPECT().DeletePost(postID, userID).Return(nil)
			},
			postID:            "123",
			userID:            "111",
			expectStatusCode:  200,
			expectRequestBody: []byte(`"message":"success"`),
		},
		{
			name:              "invalid user id",
			mockBehavior:      func(s *mockservice.MockPosts, postID, userID string) {},
			postID:            "213",
			userID:            "123",
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid user id"`),
		},
		{
			name: "invalid auth user",
			mockBehavior: func(s *mockservice.MockPosts, postID, userID string) {
				s.EXPECT().DeletePost(postID, userID).Return(errors.New("invalid user id"))
			},
			postID:            "123",
			userID:            "213",
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid user id"`),
		},
	}

	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			posts := mockservice.NewMockPosts(c)
			testCase.mockBehavior(posts, testCase.postID, testCase.userID)

			services := &service.Service{Posts: posts}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))

			r := httptest.NewRequest("GET", "/post/", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"post_id": testCase.postID,
			})
			var ctx context.Context
			if testCase.name == "invalid user id" {
				ctx = r.Context()
			} else {
				ctx = context.WithValue(r.Context(), "id", testCase.userID)
			}
			handler.DeletePost(w, r.WithContext(ctx))
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}
