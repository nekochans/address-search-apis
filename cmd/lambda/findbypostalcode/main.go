package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/nekochans/address-search-apis/infrastructure"

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

func createApiGatewayV2Response(
	ctx context.Context,
	statusCode int,
	resBodyJson []byte,
) events.APIGatewayV2HTTPResponse {
	res := events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":        "application/json",
			"X-Request-Id":        infrastructure.ExtractHttpRequestIdFromContext(ctx),
			"X-Lambda-Request-Id": infrastructure.ExtractLambdaRequestIdFromContext(ctx),
		},
		Body:            string(resBodyJson),
		IsBase64Encoded: false,
	}

	return res
}

func createErrorResponse(ctx context.Context, statusCode int, message string) events.APIGatewayV2HTTPResponse {
	resBody := &ResponseErrorBody{Message: message}
	resBodyJson, _ := json.Marshal(resBody)

	res := events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":        "application/json",
			"X-Request-Id":        infrastructure.ExtractHttpRequestIdFromContext(ctx),
			"X-Lambda-Request-Id": infrastructure.ExtractLambdaRequestIdFromContext(ctx),
		},
		Body:            string(resBodyJson),
		IsBase64Encoded: false,
	}

	return res
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	newCtx := infrastructure.CreateContextWithRequestId(ctx, req.Headers["x-request-id"])

	logger := infrastructure.CreateLogger(
		infrastructure.ExtractLambdaRequestIdFromContext(newCtx),
		infrastructure.ExtractHttpRequestIdFromContext(newCtx),
	)

	if val, ok := req.PathParameters["postalCode"]; ok {
		repo := &repository.KenallAddressRepository{HttpClient: client}
		scenario := &application.AddressScenario{
			AddressRepository: repo,
		}

		postalCode := strings.ReplaceAll(val, "-", "")

		request := &application.FindByPostalCodeRequest{
			Ctx:        newCtx,
			PostalCode: postalCode,
		}

		resBody, err := scenario.FindByPostalCode(request)
		if err != nil {
			var statusCode int
			var message string

			switch err.Error() {
			case domain.ErrAddressRepositoryNotFound.Error():
				statusCode = http.StatusNotFound
				message = "???????????????????????????????????????"
			case application.ErrFindByPostalCodeValidation.Error():
				statusCode = http.StatusUnprocessableEntity
				message = "????????????????????????????????????????????????????????????"
			default:
				statusCode = http.StatusInternalServerError
				message = "??????????????????????????????????????????"

				logger.Error(err.Error())
			}

			return createErrorResponse(newCtx, statusCode, message), nil
		}

		resBodyJson, _ := json.Marshal(resBody)

		res := createApiGatewayV2Response(newCtx, http.StatusOK, resBodyJson)

		logger.Info(string(resBodyJson))

		return res, nil
	}

	return createErrorResponse(
		newCtx,
		http.StatusInternalServerError,
		"??????????????????????????????????????????",
	), nil
}

func main() {
	lambda.Start(Handler)
}
