package types

type Link struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

type Resource interface {
	GetID() string
	GetName() string
}
