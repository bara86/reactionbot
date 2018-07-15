package commonstructure

type Storage interface {
	Lookup(code string) bool
	Add(code string, value string) error
	Remove(code string) error
	Get(code string) (string, error)
	Pop(code string) (string, error)
}