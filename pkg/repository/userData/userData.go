package userdata

type UserData interface {
	InsertUser(user User) (User, error)
	GetUser(id string) (User, error)
	CheckUser(login string) (string, error)
}
