package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"projects/pkg/api/auth"
)

type Client struct {
	conn *grpc.ClientConn
	auth auth.AuthClient
}

func InitClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthClient(conn)

	return &Client{conn: conn, auth: client}, nil
}

func (c *Client) VerifyToken(ctx context.Context, token string) (string, error) {

	res, err := c.auth.VerifyToken(ctx, &auth.VerifyTokenRequest{AccessToken: token})
	if err != nil {
		return "", err
	}
	if res.GetUserId() == "" {
		return "", err
	}

	return res.GetUserId(), nil
}

func (c *Client) Close() {
	c.conn.Close()
}
