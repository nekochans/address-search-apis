package domain

type AddressRepository interface {
	FindByPostalCode(postalCode string) (*Address, error)
}
