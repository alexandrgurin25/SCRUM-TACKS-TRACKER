package tests

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"projects/internal/models"
	mocks "projects/internal/repository/mocks"
	"projects/internal/repository/projects"
	"testing"
	"time"
)

func TestProjectsRepository_SaveProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjects(ctrl)
	_ = projects.InitRepo(&pgxpool.Pool{})

	ctx := context.Background()
	projectID := "test-project-1"
	authorID := "test-author-1"
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

	t.Run("successful project creation", func(t *testing.T) {
		mockRepo.EXPECT().
			SaveProject(ctx, projectID, authorID, description, title).
			Return(expectedProject, nil)

		createdProject, err := mockRepo.SaveProject(ctx, projectID, authorID, description, title)
		require.NoError(t, err)
		assert.Equal(t, expectedProject, createdProject)
	})

	t.Run("error on duplicate project id", func(t *testing.T) {
		expectedErr := errors.New("duplicate key value violates unique constraint")

		mockRepo.EXPECT().
			SaveProject(ctx, projectID, authorID, description, title).
			Return(nil, expectedErr)

		_, err := mockRepo.SaveProject(ctx, projectID, authorID, description, title)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestProjectsRepository_GetProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjects(ctrl)
	_ = projects.InitRepo(&pgxpool.Pool{})

	ctx := context.Background()
	projectID := "test-project-get-1"
	authorID := "test-author-1"
	title := "Test Project for Get"
	description := "Test Description for Get"

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
		UsersID:   []string{"user-1", "user-2"},
	}

	t.Run("get existing project with participants", func(t *testing.T) {
		mockRepo.EXPECT().
			GetProject(ctx, projectID).
			Return(expectedProject, expectedParticipants, nil)

		project, participants, err := mockRepo.GetProject(ctx, projectID)
		require.NoError(t, err)
		assert.Equal(t, expectedProject, project)
		assert.Equal(t, expectedParticipants, participants)
	})

	t.Run("get non-existent project", func(t *testing.T) {
		expectedErr := errors.New("project not found")

		mockRepo.EXPECT().
			GetProject(ctx, "non-existent-project").
			Return(nil, nil, expectedErr)

		_, _, err := mockRepo.GetProject(ctx, "non-existent-project")
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestProjectsRepository_AddParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjects(ctrl)
	_ = projects.InitRepo(&pgxpool.Pool{})

	ctx := context.Background()
	projectID := "test-project-1"
	participantEmail := "participant@example.com"

	t.Run("successful participant addition", func(t *testing.T) {
		mockRepo.EXPECT().
			AddParticipant(ctx, projectID, participantEmail).
			Return(nil)

		err := mockRepo.AddParticipant(ctx, projectID, participantEmail)
		require.NoError(t, err)
	})

	t.Run("add participant to non-existent project", func(t *testing.T) {
		expectedErr := errors.New("project not found")

		mockRepo.EXPECT().
			AddParticipant(ctx, "non-existent-project", participantEmail).
			Return(expectedErr)

		err := mockRepo.AddParticipant(ctx, "non-existent-project", participantEmail)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("add non-existent participant", func(t *testing.T) {
		expectedErr := errors.New("user not found")

		mockRepo.EXPECT().
			AddParticipant(ctx, projectID, "non-existent@example.com").
			Return(expectedErr)

		err := mockRepo.AddParticipant(ctx, projectID, "non-existent@example.com")
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("add duplicate participant", func(t *testing.T) {
		expectedErr := errors.New("user already participant")

		mockRepo.EXPECT().
			AddParticipant(ctx, projectID, participantEmail).
			Return(expectedErr)

		err := mockRepo.AddParticipant(ctx, projectID, participantEmail)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestProjectsRepository_DeleteProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjects(ctrl)
	_ = projects.InitRepo(&pgxpool.Pool{})

	ctx := context.Background()
	projectID := "test-project-1"
	authorID := "test-author-1"

	t.Run("successful project deletion", func(t *testing.T) {
		mockRepo.EXPECT().
			DeleteProject(ctx, projectID, authorID).
			Return(nil)

		err := mockRepo.DeleteProject(ctx, projectID, authorID)
		require.NoError(t, err)
	})

	t.Run("delete non-existent project", func(t *testing.T) {
		expectedErr := errors.New("project not found")

		mockRepo.EXPECT().
			DeleteProject(ctx, "non-existent-project", authorID).
			Return(expectedErr)

		err := mockRepo.DeleteProject(ctx, "non-existent-project", authorID)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("delete project by non-author", func(t *testing.T) {
		expectedErr := errors.New("unauthorized")

		mockRepo.EXPECT().
			DeleteProject(ctx, projectID, "non-author").
			Return(expectedErr)

		err := mockRepo.DeleteProject(ctx, projectID, "non-author")
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
