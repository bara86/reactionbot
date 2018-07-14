package commonstructure

type Storage interface {
	LookupUser(user string) bool
	AddUser(user string, token string) error
	RemoveUser(user string) error
	GetUser(user string) (string, error)
	PopUser(user string) (string, error)
}
