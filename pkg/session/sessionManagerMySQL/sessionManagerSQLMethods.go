package sessionmanagermysql

import (
	"strconv"
)

func (manager *SessionManagerMySQL) Create(token, userID string) error {
	_, err := manager.db.Exec("UPDATE userDB SET token=? WHERE user_id=?", token, userID)
	return err
}

func (manager *SessionManagerMySQL) Check(token string) (string, error) {
	var userID int
	res := manager.db.QueryRow("SELECT user_id FROM userDB WHERE token=?", token)
	if err := res.Scan(&userID); err != nil {
		return "", err
	}
	return strconv.Itoa(userID), nil
}
