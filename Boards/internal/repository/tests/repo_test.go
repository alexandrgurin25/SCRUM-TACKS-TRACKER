package tests

import (
	"boards/internal/models"
	"boards/internal/repository/board"
	mocks "boards/internal/repository/mocks"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestBoardRepository_SaveBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBoards(ctrl)
	_ = board.InitRepo(&pgxpool.Pool{})

	ctx := context.Background()
	boardID := "test-board-1"
	authorID := "test-author-1"
	projectID := "test-project-1"
	title := "Test Board"

	expectedBoard := &models.Board{
		ID:        boardID,
		AuthorID:  authorID,
		ProjectID: projectID,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful board creation", func(t *testing.T) {
		mockRepo.EXPECT().
			SaveBoard(ctx, boardID, authorID, projectID, title).
			Return(expectedBoard, nil)

		createdBoard, err := mockRepo.SaveBoard(ctx, boardID, authorID, projectID, title)
		require.NoError(t, err)
		assert.Equal(t, expectedBoard, createdBoard)
	})

	t.Run("error on duplicate board id", func(t *testing.T) {
		expectedErr := errors.New("duplicate key value violates unique constraint")

		mockRepo.EXPECT().
			SaveBoard(ctx, boardID, authorID, projectID, title).
			Return(nil, expectedErr)

		_, err := mockRepo.SaveBoard(ctx, boardID, authorID, projectID, title)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestBoardRepository_GetBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBoards(ctrl)
	_ = board.InitRepo(&pgxpool.Pool{})

	ctx := context.Background()
	boardID := "test-board-get-1"
	authorID := "test-author-1"
	projectID := "test-project-1"
	title := "Test Board for Get"

	expectedBoard := &models.Board{
		ID:        boardID,
		AuthorID:  authorID,
		ProjectID: projectID,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedTasks := &models.TasksList{
		BoardID: boardID,
		Tasks:   []models.Task{},
	}

	t.Run("get existing board", func(t *testing.T) {
		mockRepo.EXPECT().
			GetBoard(ctx, boardID).
			Return(expectedBoard, expectedTasks, nil)

		boardGet, tasks, err := mockRepo.GetBoard(ctx, boardID)
		require.NoError(t, err)
		assert.Equal(t, expectedBoard, boardGet)
		assert.Equal(t, expectedTasks, tasks)
	})

	t.Run("get non-existent board", func(t *testing.T) {
		expectedErr := errors.New("board not found")

		mockRepo.EXPECT().
			GetBoard(ctx, "non-existent-board").
			Return(nil, nil, expectedErr)

		_, _, err := mockRepo.GetBoard(ctx, "non-existent-board")
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestBoardRepository_GetAllBoards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBoards(ctrl)
	_ = board.InitRepo(&pgxpool.Pool{})

	ctx := context.Background()
	projectID := "test-project-get-all"

	expectedBoards := &models.BoardsList{
		ProjectID: projectID,
		Boards: []models.Board{
			{
				ID:        "board-1",
				AuthorID:  "author-1",
				ProjectID: projectID,
				Title:     "Board 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "board-2",
				AuthorID:  "author-1",
				ProjectID: projectID,
				Title:     "Board 2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	t.Run("get all boards for project", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllBoards(ctx, projectID).
			Return(expectedBoards, nil)

		boardsList, err := mockRepo.GetAllBoards(ctx, projectID)
		require.NoError(t, err)
		assert.Equal(t, expectedBoards, boardsList)
	})

	t.Run("get boards for non-existent project", func(t *testing.T) {
		emptyProjectID := "non-existent-project"
		emptyBoards := &models.BoardsList{
			ProjectID: emptyProjectID,
			Boards:    []models.Board{},
		}

		mockRepo.EXPECT().
			GetAllBoards(ctx, emptyProjectID).
			Return(emptyBoards, nil)

		boardsList, err := mockRepo.GetAllBoards(ctx, emptyProjectID)
		require.NoError(t, err)
		assert.Equal(t, emptyBoards, boardsList)
	})
}
