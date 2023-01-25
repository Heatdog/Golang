package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

func TestServer_CreateComment(t *testing.T) {
	type mockBehavior func(s *mockservice.MockComments, userID string, postID string, comment string)
	testingTable := []struct {
		name              string
		inputBody         string
		inputUserID       string
		inputPostID       string
		inputComment      string
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name:         "ok",
			inputBody:    `{"comment":"123"}`,
			inputUserID:  "1",
			inputPostID:  "abcd",
			inputComment: "123",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string, comment string) {
				s.EXPECT().CreateComm(postID, userID, comment).Return(itemdata.Post{
					ID: postID,
					Ath: itemdata.Author{
						ID:       userID,
						Username: "test",
					},
					Comments: []itemdata.Comment{
						{
							Ath: itemdata.Author{
								ID:       userID,
								Username: "test",
							},
							Body:    comment,
							Created: "2022-11-04T17:55:14Z",
							ID:      "111",
						},
					},
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
							User: userID,
							Vote: 1,
						},
					},
				}, nil)
			},
			expectStatusCode:  201,
			expectRequestBody: []byte(`{"id":"abcd","author":{"id":"1","username":"test"},"comments":[{"author":{"id":"1","username":"test"},"body":"123","created":"2022-11-04T17:55:14Z","id":"111"}],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name:              "invalid user id",
			inputBody:         "",
			inputUserID:       "",
			inputPostID:       "",
			inputComment:      "",
			mockBehavior:      func(s *mockservice.MockComments, userID string, postID string, comment string) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid user id"}`),
		},
		{
			name:              "invalid json",
			inputBody:         `{12321321}`,
			inputUserID:       "",
			inputPostID:       "",
			inputComment:      "",
			mockBehavior:      func(s *mockservice.MockComments, userID string, postID string, comment string) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid json input"}`),
		},
		{
			name:         "invalid user in DB",
			inputBody:    `{"comment":"123"}`,
			inputUserID:  "1",
			inputPostID:  "123",
			inputComment: "123",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string, comment string) {
				s.EXPECT().CreateComm(postID, userID, comment).Return(itemdata.Post{}, errors.New("invalid user id"))
			},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid user id"}`),
		},
	}
	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comments := mockservice.NewMockComments(c)
			testCase.mockBehavior(comments, testCase.inputUserID, testCase.inputPostID, testCase.inputComment)

			services := &service.Service{Comments: comments}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))
			r := httptest.NewRequest("POST", "/post/"+testCase.inputPostID, bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"post_id": testCase.inputPostID,
			})
			var ctx context.Context
			if testCase.name == "invalid user id" {
				ctx = r.Context()
			} else {
				ctx = context.WithValue(r.Context(), "id", testCase.inputUserID)
			}
			handler.CreateComment(w, r.WithContext(ctx))
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_DeleteComment(t *testing.T) {
	type mockBehavior func(s *mockservice.MockComments, userID string, postID string, commentID string)
	testingTable := []struct {
		name              string
		inputBody         string
		inputUserID       string
		inputPostID       string
		inputCommentID    string
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name:           "ok",
			inputBody:      "",
			inputUserID:    "123",
			inputPostID:    "1",
			inputCommentID: "1",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string, commentID string) {
				s.EXPECT().DeleteComm(postID, userID, commentID).Return(itemdata.Post{
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
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"id":"1","author":{"id":"1","username":"123"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name:              "invalid user id",
			inputBody:         "",
			inputUserID:       "123",
			inputPostID:       "111",
			inputCommentID:    "123",
			mockBehavior:      func(s *mockservice.MockComments, userID string, postID string, commentID string) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid user id"}`),
		},
		{
			name:           "invalid post id",
			inputBody:      "",
			inputUserID:    "123",
			inputPostID:    "111",
			inputCommentID: "123",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string, commentID string) {
				s.EXPECT().DeleteComm(postID, userID, commentID).Return(itemdata.Post{}, errors.New("invalid post id"))
			},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid post id"}`),
		},
	}
	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comments := mockservice.NewMockComments(c)
			testCase.mockBehavior(comments, testCase.inputUserID, testCase.inputPostID, testCase.inputCommentID)

			services := &service.Service{Comments: comments}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))
			r := httptest.NewRequest("DELETE", fmt.Sprintf("/post/%s/%s", testCase.inputPostID, testCase.inputCommentID), bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"post_id":    testCase.inputPostID,
				"comment_id": testCase.inputCommentID,
			})
			var ctx context.Context
			if testCase.name == "invalid user id" {
				ctx = r.Context()
			} else {
				ctx = context.WithValue(r.Context(), "id", testCase.inputUserID)
			}
			handler.DeleteComment(w, r.WithContext(ctx))
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_Upvote(t *testing.T) {
	type mockBehavior func(s *mockservice.MockComments, userID string, postID string)
	testingTable := []struct {
		name              string
		inputBody         string
		inputUserID       string
		inputPostID       string
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name:        "ok",
			inputBody:   "",
			inputUserID: "123",
			inputPostID: "1",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string) {
				s.EXPECT().Upvote(postID, userID).Return(itemdata.Post{
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
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"id":"1","author":{"id":"1","username":"123"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name:              "invalid user id",
			inputBody:         "",
			inputUserID:       "123",
			inputPostID:       "123",
			mockBehavior:      func(s *mockservice.MockComments, userID string, postID string) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid user id"}`),
		},
		{
			name:        "invalid post id",
			inputBody:   "",
			inputUserID: "123",
			inputPostID: "111",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string) {
				s.EXPECT().Upvote(postID, userID).Return(itemdata.Post{}, errors.New("invalid post id"))
			},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid post id"}`),
		},
	}
	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comments := mockservice.NewMockComments(c)
			testCase.mockBehavior(comments, testCase.inputUserID, testCase.inputPostID)

			services := &service.Service{Comments: comments}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))
			r := httptest.NewRequest("GET", fmt.Sprintf("/post/%s/upvote", testCase.inputPostID), bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"post_id": testCase.inputPostID,
			})
			var ctx context.Context
			if testCase.name == "invalid user id" {
				ctx = r.Context()
			} else {
				ctx = context.WithValue(r.Context(), "id", testCase.inputUserID)
			}
			handler.Upvote(w, r.WithContext(ctx))
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_Downvote(t *testing.T) {
	type mockBehavior func(s *mockservice.MockComments, userID string, postID string)
	testingTable := []struct {
		name              string
		inputBody         string
		inputUserID       string
		inputPostID       string
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name:        "ok",
			inputBody:   "",
			inputUserID: "123",
			inputPostID: "1",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string) {
				s.EXPECT().Downvote(postID, userID).Return(itemdata.Post{
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
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"id":"1","author":{"id":"1","username":"123"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name:              "invalid user id",
			inputBody:         "",
			inputUserID:       "123",
			inputPostID:       "123",
			mockBehavior:      func(s *mockservice.MockComments, userID string, postID string) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid user id"}`),
		},
		{
			name:        "invalid post id",
			inputBody:   "",
			inputUserID: "123",
			inputPostID: "111",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string) {
				s.EXPECT().Downvote(postID, userID).Return(itemdata.Post{}, errors.New("invalid post id"))
			},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid post id"}`),
		},
	}
	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comments := mockservice.NewMockComments(c)
			testCase.mockBehavior(comments, testCase.inputUserID, testCase.inputPostID)

			services := &service.Service{Comments: comments}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))
			r := httptest.NewRequest("GET", fmt.Sprintf("/post/%s/downvote", testCase.inputPostID), bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"post_id": testCase.inputPostID,
			})
			var ctx context.Context
			if testCase.name == "invalid user id" {
				ctx = r.Context()
			} else {
				ctx = context.WithValue(r.Context(), "id", testCase.inputUserID)
			}
			handler.Downvote(w, r.WithContext(ctx))
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_Unvote(t *testing.T) {
	type mockBehavior func(s *mockservice.MockComments, userID string, postID string)
	testingTable := []struct {
		name              string
		inputBody         string
		inputUserID       string
		inputPostID       string
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name:        "ok",
			inputBody:   "",
			inputUserID: "123",
			inputPostID: "1",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string) {
				s.EXPECT().Unvote(postID, userID).Return(itemdata.Post{
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
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"id":"1","author":{"id":"1","username":"123"},"comments":[],"category":"music","score":1,"type":"text","title":"123","created":"2022-11-04T17:55:14Z","upvotePercentage":100,"views":1,"votes":[{"user":"1","vote":1}],"text":"123"}`),
		},
		{
			name:              "invalid user id",
			inputBody:         "",
			inputUserID:       "123",
			inputPostID:       "123",
			mockBehavior:      func(s *mockservice.MockComments, userID string, postID string) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid user id"}`),
		},
		{
			name:        "invalid post id",
			inputBody:   "",
			inputUserID: "123",
			inputPostID: "111",
			mockBehavior: func(s *mockservice.MockComments, userID string, postID string) {
				s.EXPECT().Unvote(postID, userID).Return(itemdata.Post{}, errors.New("invalid post id"))
			},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"message":"invalid post id"}`),
		},
	}
	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comments := mockservice.NewMockComments(c)
			testCase.mockBehavior(comments, testCase.inputUserID, testCase.inputPostID)

			services := &service.Service{Comments: comments}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))
			r := httptest.NewRequest("GET", fmt.Sprintf("/post/%s/unvote", testCase.inputPostID), bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()
			r = mux.SetURLVars(r, map[string]string{
				"post_id": testCase.inputPostID,
			})
			var ctx context.Context
			if testCase.name == "invalid user id" {
				ctx = r.Context()
			} else {
				ctx = context.WithValue(r.Context(), "id", testCase.inputUserID)
			}
			handler.Unvote(w, r.WithContext(ctx))
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}
