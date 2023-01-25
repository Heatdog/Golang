package userdatamap

import (
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"sync"
)

var _ userdata.UserData = (*userDataMap)(nil)

type userDataMap struct {
	data map[string]userdata.User
	mux  *sync.RWMutex
}

func NewUserDataMap() *userDataMap {
	return &userDataMap{
		data: make(map[string]userdata.User, 10),
		mux:  &sync.RWMutex{},
	}
}
