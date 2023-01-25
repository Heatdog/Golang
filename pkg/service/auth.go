package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/session"
	"time"
)

var _ Authorization = (*AuthService)(nil)

type AuthService struct {
	db        userdata.UserData
	sessionDB session.SesManager
	key       []byte
}

func NewAuthService(bd userdata.UserData, sessionDB session.SesManager, key []byte) *AuthService {
	return &AuthService{db: bd, sessionDB: sessionDB, key: key}
}

func (ser *AuthService) CreateUser(user userdata.User) (userdata.User, error) {
	user.Password = ser.GetHash(user.Password)
	_, err := ser.db.CheckUser(user.Login)
	if err == nil {
		return userdata.User{}, err
	}
	user, err = ser.db.InsertUser(user)
	return user, err
}

func (ser *AuthService) GetHash(password string) string {
	hashPass := md5.New()
	hashPass.Write([]byte(password))
	return hex.EncodeToString(hashPass.Sum(nil))
}

func (ser *AuthService) GenerateToken(login, password string) (string, error) {
	id, err := ser.db.CheckUser(login)
	if err != nil {
		return "", err
	}
	userDB, err := ser.db.GetUser(id)
	if err != nil {
		return "", err
	}
	if userDB.Password != password {
		return "", errors.New("invalid password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": struct {
			ID    string `json:"id"`
			Login string `json:"username"`
		}{
			userDB.ID, userDB.Login,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().AddDate(0, 0, 7).Unix(),
	})
	resToken, err := token.SignedString(ser.key)
	if err != nil {
		return "", err
	}
	err = ser.sessionDB.Create(resToken, userDB.ID)
	return resToken, err
}
