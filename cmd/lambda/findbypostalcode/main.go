package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nekochans/address-search-apis/domain"

	"github.com/nekochans/address-search-apis/application"
	"github.com/nekochans/address-search-apis/infrastructure/repository"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const timeout = 10

var client *http.Client

type ResponseErrorBody struct {
	Message string `json:"message"`
}

//nolint:gochecknoinits
func init() {
	client = &http.Client{Timeout: timeout * time.Second}
}

func createApiGatewayV2Response(statusCode int, resBodyJson []byte) events.APIGatewayV2HTTPResponse {
	res := events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(resBodyJson),
		IsBase64Encoded: false,
	}

	return res
}

func createErrorResponse(statusCode int, message string) events.APIGatewayV2HTTPResponse {
	resBody := &ResponseErrorBody{Message: message}
	resBodyJson, _ := json.Marshal(resBody)

	res := events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(resBodyJson),
		IsBase64Encoded: false,
	}

	return res
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if val, ok := req.PathParameters["postalCode"]; ok {
		repo := &repository.KenallAddressRepository{HttpClient: client}
		scenario := &application.AddressScenario{
			AddressRepository: repo,
		}

		request := &application.FindByPostalCodeRequest{
			PostalCode: val,
		}

		resBody, err := scenario.FindByPostalCode(request)
		if err != nil {
			var statusCode int
			var message string

			switch err.Error() {
			case domain.ErrAddressRepositoryNotFound.Error():
				statusCode = http.StatusNotFound
				message = "住所が見つかりませんでした"
			default:
				statusCode = http.StatusInternalServerError
				message = "予期せぬエラーが発生しました"
			}

			return createErrorResponse(statusCode, message), nil
		}

		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(http.StatusOK, resBodyJson)

		return res, nil
	}

	return createErrorResponse(http.StatusInternalServerError, "予期せぬエラーが発生しました"), nil
}

func main() {
	lambda.Start(Handler)
}
