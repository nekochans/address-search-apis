package test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/nekochans/address-search-apis/infrastructure/repository"
)

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

func CreateFindAddressesSuccessMockClient(mockAddress *repository.Address) *successMockClient {
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

func CreateFindAddressesErrorMockClient(statusCode int, mockResBody interface{}) *errorMockClient {
	return &errorMockClient{StatusCode: statusCode, MockResBody: mockResBody}
}
