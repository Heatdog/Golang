package server

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/service"
	mockservice "gitlab.com/vk-go/lectures-2022-2/pkg/service/mocks"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer_Register(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthorization, user userdata.User)
	testingTable := []struct {
		name              string
		inputBody         string
		inputUser         userdata.User
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name:      "ok",
			inputBody: `{"username": "test", "password": "12345678"}`,
			inputUser: userdata.User{
				ID:       "",
				Login:    "test",
				Password: "12345678",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, user userdata.User) {
				s.EXPECT().CreateUser(user).Return(userdata.User{
					ID:       "1",
					Login:    user.Login,
					Password: "abcdef",
				}, nil)
				s.EXPECT().GenerateToken(user.Login, "abcdef").Return("123", nil)
			},
			expectStatusCode:  http.StatusCreated,
			expectRequestBody: []byte(`"token":"123"`),
		},
		{
			name:              "json error",
			inputBody:         `{231321}`,
			inputUser:         userdata.User{},
			mockBehavior:      func(s *mockservice.MockAuthorization, user userdata.User) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid json input"`),
		},
		{
			name:              "fields error",
			inputBody:         `{"213213": "test", "312321321": "12345678"}`,
			inputUser:         userdata.User{},
			mockBehavior:      func(s *mockservice.MockAuthorization, user userdata.User) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid struct fields"`),
		},
		{
			name:      "create existed user",
			inputBody: `{"username": "test", "password": "12345678"}`,
			inputUser: userdata.User{
				ID:       "",
				Login:    "test",
				Password: "12345678",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, user userdata.User) {
				s.EXPECT().CreateUser(user).Return(userdata.User{}, errors.New("already exists"))
			},
			expectStatusCode:  400,
			expectRequestBody: []byte(`{"errors":[{"location":"body","param":"username","value":"test","msg":"already exists"}]}`),
		},
		{
			name:      "generate token error",
			inputBody: `{"username": "test", "password": "12345678"}`,
			inputUser: userdata.User{
				ID:       "",
				Login:    "test",
				Password: "12345678",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, user userdata.User) {
				s.EXPECT().CreateUser(user).Return(userdata.User{
					ID:       "1",
					Login:    user.Login,
					Password: "abcdef",
				}, nil)
				s.EXPECT().GenerateToken(user.Login, "abcdef").Return("", errors.New("invalid token generation"))
			},
			expectStatusCode:  200,
			expectRequestBody: []byte(`{"message":"invalid token generation"}`),
		},
	}
	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))

			r := httptest.NewRequest("POST", "/register", bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()

			handler.Register(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}

func TestServer_Login(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthorization, user userdata.User)
	testingTable := []struct {
		name              string
		inputBody         string
		inputUser         userdata.User
		mockBehavior      mockBehavior
		expectStatusCode  int
		expectRequestBody []byte
	}{
		{
			name:      "ok",
			inputBody: `{"username": "test", "password": "12345678"}`,
			inputUser: userdata.User{
				ID:       "",
				Login:    "test",
				Password: "12345678",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, user userdata.User) {
				s.EXPECT().GetHash("12345678").Return("abcd")
				s.EXPECT().GenerateToken("test", "abcd").Return("123", nil)
			},
			expectStatusCode:  200,
			expectRequestBody: []byte(`"token":"123"`),
		},
		{
			name:              "json error",
			inputBody:         `{231321}`,
			inputUser:         userdata.User{},
			mockBehavior:      func(s *mockservice.MockAuthorization, user userdata.User) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid json input"`),
		},
		{
			name:              "fields error",
			inputBody:         `{"213213": "test", "312321321": "12345678"}`,
			inputUser:         userdata.User{},
			mockBehavior:      func(s *mockservice.MockAuthorization, user userdata.User) {},
			expectStatusCode:  400,
			expectRequestBody: []byte(`"message":"invalid struct fields"`),
		},
		{
			name:      "incorrect password",
			inputBody: `{"username": "test", "password": "12345678"}`,
			inputUser: userdata.User{
				ID:       "",
				Login:    "test",
				Password: "12345678",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, user userdata.User) {
				s.EXPECT().GetHash("12345678").Return("abcd")
				s.EXPECT().GenerateToken("test", "abcd").Return("", errors.New("invalid password"))
			},
			expectStatusCode:  401,
			expectRequestBody: []byte(`{"message":"invalid password"}`),
		},
		{
			name:      "incorrect login",
			inputBody: `{"username":"test", "password":"12345678"}`,
			inputUser: userdata.User{
				ID:       "",
				Login:    "test",
				Password: "12345678",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, user userdata.User) {
				s.EXPECT().GetHash("12345678").Return("abcd")
				s.EXPECT().GenerateToken("test", "abcd").Return("", errors.New("invalid user"))
			},
			expectStatusCode:  401,
			expectRequestBody: []byte(`{"message":"invalid user"}`),
		},
	}
	for _, testCase := range testingTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewServer(services, log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile))

			r := httptest.NewRequest("POST", "/login", bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()

			handler.Login(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if !bytes.Contains(body, testCase.expectRequestBody) {
				t.Errorf("no text found")
			}
		})
	}
}
