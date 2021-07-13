package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nekochans/address-search-apis/application"
	"github.com/nekochans/address-search-apis/infrastructure/repository"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const timeout = 10

var client *http.Client

type ResponseOkBody struct {
	Prefecture string `json:"prefecture"`
	Locality   string `json:"locality"`
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

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if val, ok := req.PathParameters["postalCode"]; ok {
		repo := &repository.KenallAddressRepository{HttpClient: client}
		scenario := &application.AddressScenario{
			AddressRepository: repo,
		}

		request := &application.FindByPostalCodeRequest{
			PostalCode: val,
		}

		// TODO 後でエラー処理を行う
		resBody, _ := scenario.FindByPostalCode(request)
		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(http.StatusOK, resBodyJson)

		return res, nil
	}

	// TODO ここに入ってきたらエラーで返すようにする
	resBody := &ResponseOkBody{Prefecture: "東京都", Locality: "新宿区"}
	resBodyJson, _ := json.Marshal(resBody)

	res := createApiGatewayV2Response(http.StatusOK, resBodyJson)

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
