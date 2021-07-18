package domain

import (
	"context"

	"github.com/pkg/errors"
)

type AddressRepository interface {
	FindByPostalCode(ctx context.Context, postalCode string) (*Address, error)
}

var (
	ErrAddressRepositoryNotFound   = errors.New("Address is not found")
	ErrAddressRepositoryUnexpected = errors.New("Unexpected error")
)
