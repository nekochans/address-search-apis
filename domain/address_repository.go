package domain

import "github.com/pkg/errors"

type AddressRepository interface {
	FindByPostalCode(postalCode string) (*Address, error)
}

var (
	ErrAddressRepositoryNotFound   = errors.New("Address is not found")
	ErrAddressRepositoryUnexpected = errors.New("Unexpected error")
)
