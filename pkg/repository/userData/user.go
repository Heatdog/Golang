package userdata

type User struct {
	ID       string `json:"-"`
	Login    string `json:"username" valid:",required"`
	Password string `json:"password" valid:",required"`
}
