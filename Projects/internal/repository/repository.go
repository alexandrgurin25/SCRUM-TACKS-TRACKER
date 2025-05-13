package repository

import (
	"context"
	"projects/internal/models"
)

type Projects interface {
	SaveProject(ctx context.Context, id string, authorID string, description string, title string) (*models.Project, error)
	GetProject(ctx context.Context, id string) (*models.Project, *models.ParticipantsList, error)
	AddParticipant(ctx context.Context, projectId string, participantEmail string) error
	DeleteProject(ctx context.Context, id string, userID string) error
}
