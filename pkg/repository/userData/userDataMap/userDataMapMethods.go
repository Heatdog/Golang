package userdatamap

import (
	"errors"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
)

func (usData *userDataMap) CheckUser(login string) (string, error) {
	usData.mux.RLock()
	defer usData.mux.RUnlock()
	for _, el := range usData.data {
		if el.Login == login {
			return el.ID, nil
		}
	}
	return "", errors.New("invalid login")
}

func (usData *userDataMap) InsertUser(user userdata.User) (userdata.User, error) {
	id := utils.RandomHex()
	user.ID = id
	usData.mux.Lock()
	defer usData.mux.Unlock()
	usData.data[id] = user
	return user, nil
}

func (usData *userDataMap) GetUser(id string) (userdata.User, error) {
	usData.mux.RLock()
	defer usData.mux.RUnlock()
	user, ok := usData.data[id]
	if !ok {
		return user, errors.New("invalid id")
	}
	return user, nil
}
