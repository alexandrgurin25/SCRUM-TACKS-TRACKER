package comment

import (
	"context"
	"fmt"
	"gateway/internal/graphql/graph/model"
	client "gateway/internal/grpc"
	"io"

	"github.com/AlexMickh/scrum-protos/pkg/api/comments"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CommentClient struct {
	conn   *grpc.ClientConn
	client comments.CommentsClient
}

func New(addr string) (*CommentClient, error) {
	const op = "grpc.comment.New"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := comments.NewCommentsClient(conn)

	return &CommentClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *CommentClient) GetClient() comments.CommentsClient {
	return c.client
}

func (c *CommentClient) CreateComment(ctx context.Context, taskId, title, description string) (string, error) {
	const op = "grpc.comment.CreateComment"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return "", err
	}

	res, err := c.client.CreateComment(ctx, &comments.CreateCommentRequest{
		TaskId:      taskId,
		Title:       title,
		Description: description,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetId(), nil
}

func (c *CommentClient) GetComments(
	ctx context.Context,
	taskId string,
	ch chan *model.Comment,
	done chan struct{},
) {
	const op = "grpc.comment.GetComment"
	fmt.Println("yes")

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return
	}
	fmt.Println("yes3")

	stream, err := c.client.GetComments(ctx, &comments.GetCommentsRequest{TaskId: taskId})
	if err != nil {
		fmt.Println(err)
		return
	}

	// done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				// done <- true //close(done)
				done <- struct{}{}
				return
			}
			if err != nil {
				fmt.Println(err)
				return
			}

			com := &model.Comment{
				ID:          resp.GetId(),
				AuthorID:    resp.GetAuthorId(),
				TaskID:      resp.GetTaskId(),
				Title:       resp.GetTitle(),
				Description: resp.GetDescription(),
			}
			fmt.Println(*com)
			ch <- com
		}
	}()

	// <-done
}

func (c *CommentClient) Close() {
	c.conn.Close()
}
