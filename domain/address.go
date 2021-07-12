package domain

type Address struct {
	PostalCode string `json:"postalCode"`
	Prefecture string `json:"prefecture"`
	Locality   string `json:"locality"`
}
