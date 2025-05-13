package service

import (
	"context"
	"projects/internal/models"
)

type Projects interface {
	CreateProject(ctx context.Context, authorID string, description string, title string) (string, error)
	GetProject(ctx context.Context, id string) (*models.Project, *models.ParticipantsList, error)
	AddParticipant(ctx context.Context, projectId string, participantEmail string) error
	DeleteProject(ctx context.Context, id string, userID string) error
}
