package application

import (
	"context"

	"github.com/nekochans/address-search-apis/domain"
)

type AddressScenario struct {
	AddressRepository domain.AddressRepository
}

type FindByPostalCodeRequest struct {
	Ctx        context.Context
	PostalCode string
}

func (s *AddressScenario) FindByPostalCode(req *FindByPostalCodeRequest) (*domain.Address, error) {
	return s.AddressRepository.FindByPostalCode(req.Ctx, req.PostalCode)
}
