package userdatamysql

import (
	"database/sql"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
)

var _ userdata.UserData = (*UserDataMySQL)(nil)

type UserDataMySQL struct {
	db *sql.DB
}

func NewUserDataMySql(db *sql.DB) *UserDataMySQL {
	return &UserDataMySQL{db: db}
}
