package application

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/nekochans/address-search-apis/domain"
	"github.com/nekochans/address-search-apis/infrastructure/repository"
)

func TestMain(m *testing.M) {
	status := m.Run()

	os.Exit(status)
}

func TestHandler(t *testing.T) {
	t.Run("Successful FindByPostalCode", func(t *testing.T) {
		const timeout = 10
		client := &http.Client{Timeout: timeout * time.Second}

		repo := &repository.KenallAddressRepository{HttpClient: client}

		scenario := AddressScenario{
			AddressRepository: repo,
		}

		req := &FindByPostalCodeRequest{postalCode: "1620062"}
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
}
