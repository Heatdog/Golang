package session

type Session struct {
	Token  string
	UserID string
}

func NewSession(token, userID string) *Session {
	return &Session{
		UserID: userID,
		Token:  token,
	}
}
