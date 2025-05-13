package client

import (
	"context"
	"errors"
	"fmt"

	auth "github.com/AlexMickh/scrum-protos/pkg/api/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client auth.AuthClient
}

var ErrUserNotFound = errors.New("user not found")

func New(addr string) (*Client, error) {
	const op = "grpc.client.New"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := auth.NewAuthClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) VerifyToken(ctx context.Context, token string) (string, error) {
	const op = "grpc.client.GetUserID"

	res, err := c.client.VerifyToken(ctx, &auth.VerifyTokenRequest{AccessToken: token})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetUserId(), nil
}

func (c *Client) Close() {
	c.conn.Close()
}
