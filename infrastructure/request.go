package infrastructure

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/google/uuid"
)

type contextKey string

const requestIdContextKey contextKey = "requestId"

func CreateContextWithRequestId(ctx context.Context) context.Context {
	newCtx := context.Background()

	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		newCtx = context.WithValue(newCtx, requestIdContextKey, lc.AwsRequestID)
	} else {
		u, _ := uuid.NewRandom()
		uu := u.String()
		newCtx = context.WithValue(newCtx, requestIdContextKey, uu)
	}

	return newCtx
}

func ExtractRequestIdFromContext(ctx context.Context) string {
	requestId, ok := ctx.Value(requestIdContextKey).(string)

	if ok {
		return requestId
	}

	return ""
}
