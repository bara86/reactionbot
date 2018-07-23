package commonstructure

type Storage interface {
	LookupUserToken(id string) (bool, error)
	AddUserToken(id string, token string) error
	RemoveUserToken(id string) error
	GetUserToken(id string) (string, error)
	PopUserToken(id string) (string, error)
}
