package tests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"projects/internal/config"
	"projects/internal/models"
	mocks "projects/internal/repository/mocks"
	"projects/internal/service/projects"
	"projects/pkg/logger"
	"testing"
	"time"
)

func getTestContext() context.Context {
	return context.WithValue(context.Background(), logger.Key, logger.NewDefaultLogger())
}

func TestProjectService_CreateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjects(ctrl)
	cfg := &config.Config{}
	serv := projects.InitServ(mockRepo, cfg)

	ctx := context.WithValue(context.Background(), logger.Key, logger.NewDefaultLogger())
	authorID := uuid.NewString()
	title := "Test Project"
	description := "Test Description"
	projectID := uuid.NewString()

	expectedProject := &models.Project{
		ID:          projectID,
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("successful project creation", func(t *testing.T) {
		mockRepo.EXPECT().
			SaveProject(ctx, gomock.Any(), authorID, description, title).
			DoAndReturn(func(_ context.Context, id, authorID, description, title string) (*models.Project, error) {
				return &models.Project{
					ID:          id,
					AuthorID:    authorID,
					Title:       title,
					Description: description,
					CreatedAt:   expectedProject.CreatedAt,
					UpdatedAt:   expectedProject.UpdatedAt,
				}, nil
			})

		id, err := serv.CreateProject(ctx, authorID, description, title)
		require.NoError(t, err)
		assert.NotEmpty(t, id)

		mockRepo.EXPECT().
			GetProject(ctx, id).
			Return(&models.Project{
				ID:          id,
				AuthorID:    authorID,
				Title:       title,
				Description: description,
				CreatedAt:   expectedProject.CreatedAt,
				UpdatedAt:   expectedProject.UpdatedAt,
			}, nil, nil)

		project, _, err := serv.GetProject(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, project.ID)
		assert.Equal(t, authorID, project.AuthorID)
		assert.Equal(t, title, project.Title)
		assert.Equal(t, description, project.Description)
	})

	t.Run("repository error on project creation", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().
			SaveProject(ctx, gomock.Any(), authorID, description, title).
			Return(nil, expectedErr)

		_, err := serv.CreateProject(ctx, authorID, description, title)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to create project")
	})
}

func TestProjectService_GetProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjects(ctrl)
	cfg := &config.Config{}
	serv := projects.InitServ(mockRepo, cfg)

	ctx := getTestContext()
	projectID := uuid.NewString()
	authorID := uuid.NewString()
	title := "Test Project"
	description := "Test Description"

	expectedProject := &models.Project{
		ID:          projectID,
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	expectedParticipants := &models.ParticipantsList{
		ProjectID: projectID,
		UsersID:   []string{uuid.NewString(), uuid.NewString()},
	}

	t.Run("successful project retrieval", func(t *testing.T) {
		mockRepo.EXPECT().
			GetProject(ctx, projectID).
			Return(expectedProject, expectedParticipants, nil)

		project, participants, err := serv.GetProject(ctx, projectID)
		require.NoError(t, err)
		assert.Equal(t, expectedProject, project)
		assert.Equal(t, expectedParticipants, participants)
	})

	t.Run("project not found", func(t *testing.T) {
		expectedErr := errors.New("project not found")

		mockRepo.EXPECT().
			GetProject(ctx, "non-existent-id").
			Return(nil, nil, expectedErr)

		_, _, err := serv.GetProject(ctx, "non-existent-id")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to get project")
	})
}

func TestProjectService_AddParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjects(ctrl)
	cfg := &config.Config{}
	serv := projects.InitServ(mockRepo, cfg)

	ctx := getTestContext()
	projectID := uuid.NewString()
	participantEmail := "participant@example.com"

	t.Run("successful participant addition", func(t *testing.T) {
		mockRepo.EXPECT().
			AddParticipant(ctx, projectID, participantEmail).
			Return(nil)

		err := serv.AddParticipant(ctx, projectID, participantEmail)
		require.NoError(t, err)
	})

	t.Run("repository error on participant addition", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().
			AddParticipant(ctx, projectID, participantEmail).
			Return(expectedErr)

		err := serv.AddParticipant(ctx, projectID, participantEmail)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to add participant")
	})

	t.Run("participant already exists", func(t *testing.T) {
		expectedErr := errors.New("user already participant")

		mockRepo.EXPECT().
			AddParticipant(ctx, projectID, participantEmail).
			Return(expectedErr)

		err := serv.AddParticipant(ctx, projectID, participantEmail)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to add participant")
	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjects(ctrl)
	cfg := &config.Config{}
	serv := projects.InitServ(mockRepo, cfg)

	ctx := getTestContext()
	projectID := uuid.NewString()
	authorID := uuid.NewString()
	otherUserID := uuid.NewString()

	project := &models.Project{
		ID:          projectID,
		AuthorID:    authorID,
		Title:       "Test Project",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("successful project deletion by author", func(t *testing.T) {
		mockRepo.EXPECT().
			GetProject(ctx, projectID).
			Return(project, nil, nil)

		mockRepo.EXPECT().
			DeleteProject(ctx, projectID, authorID).
			Return(nil)

		err := serv.DeleteProject(ctx, projectID, authorID)
		require.NoError(t, err)
	})

	t.Run("project not found for deletion", func(t *testing.T) {
		expectedErr := errors.New("project not found")

		mockRepo.EXPECT().
			GetProject(ctx, "non-existent-id").
			Return(nil, nil, expectedErr)

		err := serv.DeleteProject(ctx, "non-existent-id", authorID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get project")
	})

	t.Run("unauthorized deletion attempt", func(t *testing.T) {
		mockRepo.EXPECT().
			GetProject(ctx, projectID).
			Return(project, nil, nil)

		err := serv.DeleteProject(ctx, projectID, otherUserID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not authorized to delete project")
	})

	t.Run("repository error on project deletion", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().
			GetProject(ctx, projectID).
			Return(project, nil, nil)

		mockRepo.EXPECT().
			DeleteProject(ctx, projectID, authorID).
			Return(expectedErr)

		err := serv.DeleteProject(ctx, projectID, authorID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete project")
	})
}
