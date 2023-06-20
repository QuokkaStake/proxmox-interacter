package types

type ProxmoxStatusResponse struct {
	Data []map[string]interface{}
}

type Link struct {
	Name string
	Href string
}
