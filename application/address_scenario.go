package application

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/nekochans/address-search-apis/domain"
	"github.com/pkg/errors"
)

type AddressScenario struct {
	AddressRepository domain.AddressRepository
}

type FindByPostalCodeRequest struct {
	Ctx        context.Context
	PostalCode string `validate:"required,len=7,numeric"`
}

var (
	ErrFindByPostalCodeValidation = errors.New("postalCode format is invalid")
)

func (s *AddressScenario) FindByPostalCode(req *FindByPostalCodeRequest) (*domain.Address, error) {
	validate := validator.New()
	err := validate.Struct(req)

	if err != nil {
		return nil, ErrFindByPostalCodeValidation
	}

	return s.AddressRepository.FindByPostalCode(req.Ctx, req.PostalCode)
}
