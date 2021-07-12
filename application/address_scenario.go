package application

import (
	"github.com/nekochans/address-search-apis/domain"
)

type AddressScenario struct {
	AddressRepository domain.AddressRepository
}

type FindByPostalCodeRequest struct {
	postalCode string
}

func (s *AddressScenario) FindByPostalCode(req *FindByPostalCodeRequest) (*domain.Address, error) {
	return s.AddressRepository.FindByPostalCode(req.postalCode)
}