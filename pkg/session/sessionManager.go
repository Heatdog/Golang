package session

type SesManager interface {
	Check(token string) (string, error)
	Create(token, userID string) error
}
