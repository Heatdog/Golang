package sessionmanagermysql

import (
	"database/sql"
	"gitlab.com/vk-go/lectures-2022-2/pkg/session"
)

var _ session.SesManager = (*SessionManagerMySQL)(nil)

type SessionManagerMySQL struct {
	db *sql.DB
}

func NewSessionManagerMySQL(db *sql.DB) *SessionManagerMySQL {
	return &SessionManagerMySQL{db: db}
}
