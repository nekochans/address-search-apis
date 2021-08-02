package application

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/nekochans/address-search-apis/domain"
	"github.com/nekochans/address-search-apis/infrastructure"
	"github.com/nekochans/address-search-apis/infrastructure/repository"
	"github.com/nekochans/address-search-apis/test"
)

func TestMain(m *testing.M) {
	status := m.Run()

	os.Exit(status)
}

// nolint:funlen
func TestHandler(t *testing.T) {
	t.Run("Successful FindByPostalCode", func(t *testing.T) {
		mockAddress := &repository.Address{
			Prefecture: "東京都",
			City:       "新宿区",
			Town:       "市谷加賀町",
		}

		client := test.CreateFindAddressesSuccessMockClient(mockAddress)

		repo := &repository.KenallAddressRepository{HttpClient: client}

		scenario := AddressScenario{
			AddressRepository: repo,
		}

		ctx := context.Background()
		req := &FindByPostalCodeRequest{
			Ctx:        infrastructure.CreateContextWithRequestId(ctx, "aaaaaaaaaaaaa-bbbb-cccc-ddddddddddd1"),
			PostalCode: "1620062",
		}
		res, err := scenario.FindByPostalCode(req)

		if err != nil {
			t.Fatal("Error failed to FindByPostalCode", err)
		}

		expected := &domain.Address{
			PostalCode: "1620062",
			Prefecture: "東京都",
			Locality:   "新宿区市谷加賀町",
		}

		if reflect.DeepEqual(res, expected) == false {
			t.Error("\nActually: ", res, "\nExpected: ", expected)
		}
	})

	t.Run("Error FindByPostalCode Address is not found", func(t *testing.T) {
		mockResBody := &repository.Address{}

		client := test.CreateFindAddressesErrorMockClient(404, mockResBody)

		repo := &repository.KenallAddressRepository{HttpClient: client}

		scenario := AddressScenario{
			AddressRepository: repo,
		}

		ctx := context.Background()
		req := &FindByPostalCodeRequest{Ctx: infrastructure.CreateContextWithRequestId(
			ctx,
			"aaaaaaaaaaaaa-bbbb-cccc-ddddddddddd1"),
			PostalCode: "4040000",
		}
		res, err := scenario.FindByPostalCode(req)

		expected := "Address is not found"
		if res != nil {
			t.Error("\nActually: ", res, "\nExpected: ", expected)
		}

		if err != nil {
			if err.Error() != expected {
				t.Error("\nActually: ", err, "\nExpected: ", expected)
			}
		}
	})

	t.Run("Error FindByPostalCode Unexpected error", func(t *testing.T) {
		mockResBody := &repository.Address{}

		client := test.CreateFindAddressesErrorMockClient(500, mockResBody)

		repo := &repository.KenallAddressRepository{HttpClient: client}

		scenario := AddressScenario{
			AddressRepository: repo,
		}

		ctx := context.Background()
		req := &FindByPostalCodeRequest{
			Ctx:        infrastructure.CreateContextWithRequestId(ctx, "aaaaaaaaaaaaa-bbbb-cccc-ddddddddddd1"),
			PostalCode: "1000000",
		}
		res, err := scenario.FindByPostalCode(req)

		expected := "Unexpected error"
		if res != nil {
			t.Error("\nActually: ", res, "\nExpected: ", expected)
		}

		if err != nil {
			if err.Error() != expected {
				t.Error("\nActually: ", err, "\nExpected: ", expected)
			}
		}
	})

	t.Run("Error FindByPostalCode Validation error", func(t *testing.T) {
		mockAddress := &repository.Address{
			Prefecture: "東京都",
			City:       "新宿区",
			Town:       "市谷加賀町",
		}

		client := test.CreateFindAddressesSuccessMockClient(mockAddress)

		repo := &repository.KenallAddressRepository{HttpClient: client}

		scenario := AddressScenario{
			AddressRepository: repo,
		}

		ctx := context.Background()
		req := &FindByPostalCodeRequest{
			Ctx:        infrastructure.CreateContextWithRequestId(ctx, "aaaaaaaaaaaaa-bbbb-cccc-ddddddddddd1"),
			PostalCode: "16200621",
		}
		res, err := scenario.FindByPostalCode(req)

		expected := "postalCode format is invalid"
		if res != nil {
			t.Error("\nActually: ", res, "\nExpected: ", expected)
		}

		if err != nil {
			if err.Error() != expected {
				t.Error("\nActually: ", err, "\nExpected: ", expected)
			}
		}
	})
}
