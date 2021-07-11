package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ResponseOkBody struct {
	Prefecture string `json:"prefecture"`
	Locality string `json:"locality"`
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
	statusCode := 200

	resBody := &ResponseOkBody{Prefecture: "東京都", Locality: "新宿区"}
	resBodyJson, _ := json.Marshal(resBody)

	res := createApiGatewayV2Response(statusCode, resBodyJson)

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
