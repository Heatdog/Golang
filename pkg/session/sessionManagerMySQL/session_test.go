package sessionmanagermysql

import (
	"errors"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
)

func TestSession_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewSessionManagerMySQL(db)
	type mockBehaviour func(token, userID string)
	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		token         string
		userID        string
		err           error
	}{
		{
			name: "ok",
			mockBehaviour: func(token, userID string) {
				mock.ExpectExec("UPDATE userDB").
					WithArgs(token, userID).
					WillReturnResult(sqlmock.NewErrorResult(nil))
			},
			token:  "111",
			userID: "1",
			err:    nil,
		},
		{
			name: "invalid user id",
			mockBehaviour: func(token, userID string) {
				mock.ExpectExec("UPDATE userDB").
					WithArgs(token, userID).
					WillReturnError(errors.New("invalid user id"))
			},
			token:  "111",
			userID: "1",
			err:    errors.New("invalid user id"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour(testCase.token, testCase.userID)
			err := repo.Create(testCase.token, testCase.userID)
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
			}
		})
	}
}

func TestSession_Check(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewSessionManagerMySQL(db)
	type mockBehaviour func(token, userID string)
	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		token         string
		userID        string
		err           error
	}{
		{
			name: "ok",
			mockBehaviour: func(token, userID string) {
				mock.ExpectQuery("SELECT user_id FROM userDB").
					WithArgs(token).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userID))
			},
			token:  "111",
			userID: "1",
			err:    nil,
		},
		{
			name: "invalid token",
			mockBehaviour: func(token, userID string) {
				mock.ExpectQuery("SELECT user_id FROM userDB").
					WithArgs(token).
					WillReturnError(errors.New("invalid token"))
			},
			token:  "111",
			userID: "1",
			err:    errors.New("invalid token"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour(testCase.token, testCase.userID)
			res, err := repo.Check(testCase.token)
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
				if !reflect.DeepEqual(res, testCase.userID) {
					t.Errorf("results not match, want %v, have %v",
						testCase.userID, res)
					return
				}
			}
		})
	}

}
