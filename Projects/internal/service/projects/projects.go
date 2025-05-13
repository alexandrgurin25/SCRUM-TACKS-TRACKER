package projects

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"projects/internal/config"
	"projects/internal/models"
	"projects/internal/repository"
	"projects/internal/service"
	"projects/pkg/logger"
)

type Serv struct {
	repo repository.Projects
	cfg  *config.Config
}

func InitServ(projectsRepo repository.Projects, cfg *config.Config) service.Projects {
	return &Serv{repo: projectsRepo, cfg: cfg}
}

func (s *Serv) CreateProject(ctx context.Context, authorID string, description string, title string) (string, error) {

	id := uuid.New().String()

	_, err := s.repo.SaveProject(ctx, id, authorID, description, title)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "unable to create project", zap.Error(err))
		return "", fmt.Errorf("unable to create project: %w", err)
	}
	return id, nil
}

func (s *Serv) GetProject(ctx context.Context, id string) (*models.Project, *models.ParticipantsList, error) {
	project, participants, err := s.repo.GetProject(ctx, id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "unable to get project", zap.Error(err))
		return nil, nil, fmt.Errorf("unable to get project: %w", err)
	}
	return project, participants, nil
}

func (s *Serv) AddParticipant(ctx context.Context, projectId string, participantEmail string) error {
	err := s.repo.AddParticipant(ctx, projectId, participantEmail)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "unable to add participant to project", zap.Error(err))
		return fmt.Errorf("unable to add participant to project: %w", err)
	}
	return nil
}

func (s *Serv) DeleteProject(ctx context.Context, id string, userID string) error {
	project, _, err := s.repo.GetProject(ctx, id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get project for deletion", zap.Error(err))
		return fmt.Errorf("failed to get project: %w", err)
	}

	if project.AuthorID != userID {
		logger.GetLoggerFromCtx(ctx).Warn(ctx, "unauthorized project deletion attempt",
			zap.String("projectID", id),
			zap.String("userID", userID),
			zap.String("authorID", project.AuthorID))
		return fmt.Errorf("user %s is not authorized to delete project %s", userID, id)
	}

	delErr := s.repo.DeleteProject(ctx, id, userID)
	if delErr != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to delete project", zap.Error(delErr))
		return fmt.Errorf("failed to delete project: %w", delErr)
	}

	return nil
}
