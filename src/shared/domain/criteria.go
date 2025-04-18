package domain

type Criteria interface {
	Validate() error
	Debug() string
}
