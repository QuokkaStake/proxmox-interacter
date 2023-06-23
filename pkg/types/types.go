package types

type Link struct {
	Name string
	Href string
}

type Resource interface {
	GetID() string
	GetName() string
}
