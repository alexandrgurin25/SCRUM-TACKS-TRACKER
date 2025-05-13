package project

import (
	"context"
	"fmt"
	"gateway/internal/domain/models"
	client "gateway/internal/grpc"

	"github.com/AlexMickh/scrum-protos/pkg/api/projects"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProjectService struct {
	conn   *grpc.ClientConn
	client projects.ProjectsClient
}

func New(addr string) (*ProjectService, error) {
	const op = "grpc.project.New"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := projects.NewProjectsClient(conn)

	return &ProjectService{
		conn:   conn,
		client: client,
	}, nil
}

func (p *ProjectService) CreateProject(ctx context.Context, title, descrition string) (string, error) {
	const op = "grpc.project.CreateProject"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return "", err
	}

	res, err := p.client.CreateProject(ctx, &projects.CreateProjectRequest{
		Title:       title,
		Description: descrition,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetId(), nil
}

func (p *ProjectService) GetProject(ctx context.Context, id string) (*models.Project, error) {
	const op = "grpc.project.GetProject"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return nil, err
	}

	res, err := p.client.GetProject(ctx, &projects.GetProjectRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Project{
		ID:             res.GetId(),
		AuthorID:       res.GetAuthorId(),
		ParticipantsID: res.GetParticipantsId(),
		Title:          res.GetTitle(),
		CreatedAt:      res.GetCreatedAt().AsTime(),
		UpdatedAt:      res.GetUpdatedAt().AsTime(),
	}, nil
}

func (p *ProjectService) AddParticipantToProject(ctx context.Context, id, email string) error {
	const op = "grpc.project.AddParticipantToProject"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return err
	}

	_, err = p.client.AddParticipant(ctx, &projects.AddParticipantRequest{
		ProjectId:        id,
		ParticipantEmail: email,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, id string) error {
	const op = "grpc.project.DeleteProject"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return err
	}

	_, err = p.client.DeleteProject(ctx, &projects.DeleteProjectRequest{Id: id})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ProjectService) Stop() {
	p.conn.Close()
}
