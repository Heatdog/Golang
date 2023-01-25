package userdatamysql

import (
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"strconv"
)

func (usData *UserDataMySQL) CheckUser(login string) (string, error) {
	var userID int
	row := usData.db.QueryRow("SELECT user_id FROM userDB  WHERE login = ?", login)
	if err := row.Scan(&userID); err != nil {
		return "", err
	}
	return strconv.Itoa(userID), nil
}

func (usData *UserDataMySQL) InsertUser(user userdata.User) (userdata.User, error) {
	result, err := usData.db.Exec("INSERT INTO userDB (login, password) VALUE (?, ?)", user.Login, user.Password)
	if err != nil {
		return user, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return user, err
	}
	user.ID = strconv.Itoa(int(id))
	return user, nil
}

func (usData *UserDataMySQL) GetUser(id string) (userdata.User, error) {
	usID, _ := strconv.Atoi(id)
	var login, password string
	var user userdata.User
	row := usData.db.QueryRow("SELECT login, password FROM userDB WHERE user_id = ?", usID)
	if err := row.Scan(&login, &password); err != nil {
		return user, err
	}
	user.Login = login
	user.Password = password
	user.ID = id
	return user, nil
}
