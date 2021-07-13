package application

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/nekochans/address-search-apis/domain"
	"github.com/nekochans/address-search-apis/infrastructure/repository"
)

func TestMain(m *testing.M) {
	status := m.Run()

	os.Exit(status)
}

type mockResponse struct {
	Body       string
	StatusCode int
	Result     *repository.FindAddressesResponse
	Error      error
}

type successMockClient struct {
	MockAddress *repository.Address
}

func (c *successMockClient) Do(req *http.Request) (*http.Response, error) {
	mockResBodyData := []*repository.Address{c.MockAddress}

	mockResBody := &repository.FindAddressesResponse{
		Version:   "2021-06-30",
		Addresses: mockResBodyData,
	}

	bodyJson, _ := json.Marshal(mockResBody)

	mockRes := &mockResponse{
		Body:       string(bodyJson),
		StatusCode: http.StatusOK,
		Result:     mockResBody,
	}

	return &http.Response{StatusCode: mockRes.StatusCode, Body: io.NopCloser(strings.NewReader(mockRes.Body))}, nil
}

func createSuccessMockClient(mockAddress *repository.Address) *successMockClient {
	return &successMockClient{MockAddress: mockAddress}
}

func TestHandler(t *testing.T) {
	t.Run("Successful FindByPostalCode", func(t *testing.T) {
		mockAddress := &repository.Address{
			Prefecture: "東京都",
			City:       "新宿区",
			Town:       "市谷加賀町",
		}

		client := createSuccessMockClient(mockAddress)

		repo := &repository.KenallAddressRepository{HttpClient: client}

		scenario := AddressScenario{
			AddressRepository: repo,
		}

		req := &FindByPostalCodeRequest{PostalCode: "1620062"}
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
