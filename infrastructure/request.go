package infrastructure

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/google/uuid"
)

type contextKey string

const lambdaRequestIdContextKey contextKey = "lambdaRequestId"

const httpRequestIdContextKey contextKey = "httpRequestId"

func CreateContextWithRequestId(ctx context.Context, r string) context.Context {
	newCtx := context.Background()

	u, _ := uuid.NewRandom()

	var httpReqId string
	if r == "" {
		httpReqId = u.String()
	} else {
		httpReqId = r
	}

	newCtx = context.WithValue(newCtx, httpRequestIdContextKey, httpReqId)

	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		newCtx = context.WithValue(newCtx, lambdaRequestIdContextKey, lc.AwsRequestID)
	} else {
		uu := u.String()
		newCtx = context.WithValue(newCtx, lambdaRequestIdContextKey, uu)
	}

	return newCtx
}

func ExtractLambdaRequestIdFromContext(ctx context.Context) string {
	requestId, ok := ctx.Value(lambdaRequestIdContextKey).(string)

	if ok {
		return requestId
	}

	return ""
}

func ExtractHttpRequestIdFromContext(ctx context.Context) string {
	requestId, ok := ctx.Value(httpRequestIdContextKey).(string)

	if ok {
		return requestId
	}

	return ""
}
