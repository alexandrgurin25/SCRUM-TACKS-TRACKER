package client

import (
	"context"
	"fmt"
	"gateway/internal/graphql"

	"google.golang.org/grpc/metadata"
)

func AddToken(ctx context.Context, op string) (context.Context, error) {
	token, err := graphql.GetAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", token)

	return ctx, nil
}
