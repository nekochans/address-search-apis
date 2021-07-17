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
	Result     interface{}
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

type errorMockClient struct {
	StatusCode  int
	MockResBody interface{}
}

func (c *errorMockClient) Do(req *http.Request) (*http.Response, error) {
	bodyJson, _ := json.Marshal(c.MockResBody)

	mockRes := &mockResponse{
		Body:       string(bodyJson),
		StatusCode: c.StatusCode,
		Result:     c.MockResBody,
	}

	return &http.Response{StatusCode: mockRes.StatusCode, Body: io.NopCloser(strings.NewReader(mockRes.Body))}, nil
}

func createErrorMockClient(statusCode int, mockResBody interface{}) *errorMockClient {
	return &errorMockClient{StatusCode: statusCode, MockResBody: mockResBody}
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

	t.Run("Error FindByPostalCode Address is not found", func(t *testing.T) {
		mockResBody := &repository.Address{}

		client := createErrorMockClient(404, mockResBody)

		repo := &repository.KenallAddressRepository{HttpClient: client}

		scenario := AddressScenario{
			AddressRepository: repo,
		}

		req := &FindByPostalCodeRequest{PostalCode: "4040000"}
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

		client := createErrorMockClient(500, mockResBody)

		repo := &repository.KenallAddressRepository{HttpClient: client}

		scenario := AddressScenario{
			AddressRepository: repo,
		}

		req := &FindByPostalCodeRequest{PostalCode: "1000000"}
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
}
