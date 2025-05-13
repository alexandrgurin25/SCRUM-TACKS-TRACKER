package tests

import (
	"boards/internal/config"
	"boards/internal/models"
	"boards/internal/service/board"
	"boards/pkg/logger"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveBoard(ctx context.Context, id, authorID, projectID, title string) (*models.Board, error) {
	args := m.Called(ctx, id, authorID, projectID, title)
	return args.Get(0).(*models.Board), args.Error(1)
}

func (m *MockRepository) GetBoard(ctx context.Context, id string) (*models.Board, *models.TasksList, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Board), args.Get(1).(*models.TasksList), args.Error(2)
}

func (m *MockRepository) GetAllBoards(ctx context.Context, projectID string) (*models.BoardsList, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).(*models.BoardsList), args.Error(1)
}

func getTestContext() context.Context {
	return context.WithValue(context.Background(), logger.Key, zap.NewNop())
}

func getTestConfig() *config.Config {
	return &config.Config{
		Postgres: config.PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "1111",
			Database: "boards",
			MinConns: 3,
			MaxConns: 5,
		},
	}
}

func TestBoardService_CreateBoard(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := board.InitServ(mockRepo, getTestConfig())
		ctx := getTestContext()

		mockRepo.On("SaveBoard", ctx, mock.AnythingOfType("string"), "author1", "project1", "Test Board").
			Return(&models.Board{ID: "generated-id"}, nil)

		id, err := service.CreateBoard(ctx, "author1", "project1", "Test Board")

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
		assert.Regexp(t, `^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := board.InitServ(mockRepo, getTestConfig())
		ctx := getTestContext()

		mockRepo.On("SaveBoard", ctx, mock.AnythingOfType("string"), "author1", "project1", "Test Board").
			Return(&models.Board{}, errors.New("repository error"))

		id, err := service.CreateBoard(ctx, "author1", "project1", "Test Board")

		assert.Error(t, err)
		assert.Empty(t, id)
		assert.Contains(t, err.Error(), "failed to create board")
		mockRepo.AssertExpectations(t)
	})
}

// Остальные тесты остаются аналогичными, с использованием getTestConfig()
func TestBoardService_GetBoard(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := board.InitServ(mockRepo, getTestConfig())
		ctx := getTestContext()

		expectedBoard := &models.Board{ID: "board1"}
		expectedTasks := &models.TasksList{Tasks: []models.Task{}}
		mockRepo.On("GetBoard", ctx, "board1").
			Return(expectedBoard, expectedTasks, nil)

		boardGet, tasks, err := service.GetBoard(ctx, "board1")

		assert.NoError(t, err)
		assert.Equal(t, expectedBoard, boardGet)
		assert.Equal(t, expectedTasks, tasks)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := board.InitServ(mockRepo, getTestConfig())
		ctx := getTestContext()

		mockRepo.On("GetBoard", ctx, "nonexistent").
			Return(&models.Board{}, &models.TasksList{}, errors.New("not found"))

		boardGet, tasks, err := service.GetBoard(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, boardGet)
		assert.Nil(t, tasks)
		assert.Contains(t, err.Error(), "failed to get board")
		mockRepo.AssertExpectations(t)
	})
}

func TestBoardService_GetAllBoards(t *testing.T) {
	t.Run("successful get all", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := board.InitServ(mockRepo, getTestConfig())
		ctx := getTestContext()

		expectedBoards := &models.BoardsList{
			Boards: []models.Board{
				{ID: "board1"},
				{ID: "board2"},
			},
		}
		mockRepo.On("GetAllBoards", ctx, "project1").
			Return(expectedBoards, nil)

		boards, err := service.GetAllBoards(ctx, "project1")

		assert.NoError(t, err)
		assert.Equal(t, expectedBoards, boards)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty project", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := board.InitServ(mockRepo, getTestConfig())
		ctx := getTestContext()

		emptyBoards := &models.BoardsList{Boards: []models.Board{}}
		mockRepo.On("GetAllBoards", ctx, "empty project").
			Return(emptyBoards, nil)

		boards, err := service.GetAllBoards(ctx, "empty project")

		assert.NoError(t, err)
		assert.Empty(t, boards.Boards)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := board.InitServ(mockRepo, getTestConfig())
		ctx := getTestContext()

		mockRepo.On("GetAllBoards", ctx, "error project").
			Return(&models.BoardsList{}, errors.New("repository error"))

		boards, err := service.GetAllBoards(ctx, "error project")

		assert.Error(t, err)
		assert.Nil(t, boards)
		assert.Contains(t, err.Error(), "failed to get boards by project_id")
		mockRepo.AssertExpectations(t)
	})
}
