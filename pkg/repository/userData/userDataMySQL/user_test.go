package userdatamysql

import (
	"errors"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"strconv"
	"testing"
)

func TestUser_Check(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserDataMySql(db)
	type mockBehaviour func(login string, userID int)
	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		id            int
		login         string
		err           error
	}{
		{
			name: "ok",
			mockBehaviour: func(login string, userID int) {
				rows := sqlmock.NewRows([]string{"user_id"}).AddRow(userID)
				mock.ExpectQuery("SELECT user_id FROM userDB WHERE").
					WithArgs(login).
					WillReturnRows(rows)
			},
			id:    2,
			login: "123",
			err:   nil,
		},
		{
			name: "invalid user",
			mockBehaviour: func(login string, userID int) {
				rows := sqlmock.NewRows([]string{"user_id"}).RowError(1, errors.New(""))
				mock.ExpectQuery("SELECT user_id FROM userDB WHERE").
					WithArgs(login).
					WillReturnRows(rows)
			},
			id:    1,
			login: "123",
			err:   errors.New("sql: no rows in result set"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour(testCase.login, testCase.id)
			got, err := repo.CheckUser(testCase.login)
			if testCase.err != nil {
				if !reflect.DeepEqual(err.Error(), testCase.err.Error()) {
					t.Errorf("results not match, want %v, have %v", testCase.err.Error(), err.Error())
					return
				}
			} else {
				if err != nil {
					t.Errorf("unexpected err: %s", err)
					return
				}
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
					return
				}
				if !reflect.DeepEqual(got, strconv.Itoa(testCase.id)) {
					t.Errorf("results not match, want %v, have %v", testCase.id, got)
					return
				}
			}
		})
	}
}

func TestUser_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserDataMySql(db)
	type mockBehaviour func(user userdata.User, userID int64)
	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		user          userdata.User
		userID        int64
		err           error
	}{
		{
			name: "ok",
			mockBehaviour: func(user userdata.User, userID int64) {
				mock.ExpectExec("INSERT INTO userDB").
					WithArgs(user.Login, user.Password).
					WillReturnResult(sqlmock.NewResult(userID, 1))
			},
			user: userdata.User{
				ID:       "",
				Login:    "123",
				Password: "456",
			},
			userID: 1,
			err:    nil,
		},
		{
			name: "invalid insert",
			mockBehaviour: func(user userdata.User, userID int64) {
				mock.ExpectExec("INSERT INTO userDB").
					WithArgs(user.Login, user.Password).
					WillReturnError(errors.New("invalid insert"))
			},
			user: userdata.User{
				ID:       "",
				Login:    "213",
				Password: "45456",
			},
			userID: 1,
			err:    errors.New("invalid insert"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour(testCase.user, testCase.userID)
			got, err := repo.InsertUser(testCase.user)
			if testCase.err != nil {
				if !reflect.DeepEqual(err.Error(), testCase.err.Error()) {
					t.Errorf("results not match, want %v, have %v", testCase.err.Error(), err.Error())
					return
				}
			} else {
				if err != nil {
					t.Errorf("unexpected err: %s", err)
					return
				}
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
					return
				}
				if !reflect.DeepEqual(got, userdata.User{
					ID:       strconv.Itoa(int(testCase.userID)),
					Login:    testCase.user.Login,
					Password: testCase.user.Password,
				}) {
					t.Errorf("results not match, want %v, have %v", userdata.User{
						ID:       strconv.Itoa(int(testCase.userID)),
						Login:    testCase.user.Login,
						Password: testCase.user.Password,
					}, got)
					return
				}
			}
		})
	}
}

func TestUser_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserDataMySql(db)
	type mockBehaviour func(login, password string, userID int)
	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		login         string
		password      string
		userID        int
		err           error
	}{
		{
			name: "ok",
			mockBehaviour: func(login, password string, userID int) {
				rows := sqlmock.NewRows([]string{"login", "password"}).AddRow(login, password)
				mock.ExpectQuery("SELECT login, password FROM userDB WHERE").
					WithArgs(userID).
					WillReturnRows(rows)
			},
			login:    "123",
			password: "456",
			userID:   1,
			err:      nil,
		},
		{
			name: "invalid user id",
			mockBehaviour: func(login, password string, userID int) {
				rows := sqlmock.NewRows([]string{"login", "password"}).RowError(1, errors.New("invalid user id"))
				mock.ExpectQuery("SELECT login, password FROM userDB WHERE").
					WithArgs(userID).
					WillReturnRows(rows)
			},
			login:    "123",
			password: "456",
			userID:   1,
			err:      errors.New("invalid user id"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour(testCase.login, testCase.password, testCase.userID)
			got, err := repo.GetUser(strconv.Itoa(testCase.userID))
			if testCase.err != nil {
				if !reflect.DeepEqual(err.Error(), testCase.err.Error()) {
					t.Errorf("results not match, want %v, have %v", testCase.err.Error(), err.Error())
					return
				}
			} else {
				if err != nil {
					t.Errorf("unexpected err: %s", err)
					return
				}
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
					return
				}
				if !reflect.DeepEqual(got, userdata.User{
					ID:       strconv.Itoa(testCase.userID),
					Login:    testCase.login,
					Password: testCase.password,
				}) {
					t.Errorf("results not match, want %v, have %v", userdata.User{
						ID:       strconv.Itoa(testCase.userID),
						Login:    testCase.login,
						Password: testCase.password,
					}, got)
					return
				}
			}
		})
	}
}
