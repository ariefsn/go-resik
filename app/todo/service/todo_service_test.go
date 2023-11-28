package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ariefsn/go-resik/app/todo/service"
	"github.com/ariefsn/go-resik/domain"
	"github.com/ariefsn/go-resik/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	mockTodoRepo := new(mocks.TodoRepository)

	payload := &domain.TodoDto{
		Title:       "Title 1",
		Description: "Description 1",
	}
	result := &domain.Todo{
		ID:          "1",
		Title:       payload.Title,
		Description: payload.Description,
		IsCompleted: false,
		Audit: &domain.Audit{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockError := errors.New("some error")

	t.Run("Success", func(t *testing.T) {
		mockTodoRepo.On("Create", mock.Anything, payload).Return(result, nil).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, err := svc.Create(context.TODO(), payload)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, result, res)
	})

	t.Run("Failed", func(t *testing.T) {
		mockTodoRepo.On("Create", mock.Anything, payload).Return(nil, mockError).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, err := svc.Create(context.TODO(), payload)

		assert.NotNil(t, err)
		assert.Equal(t, mockError, err)
		assert.Nil(t, res)
	})
}

func TestGet(t *testing.T) {
	mockTodoRepo := new(mocks.TodoRepository)

	mockResult := []domain.Todo{
		{
			ID:          "1",
			Title:       "Title 1",
			Description: "Description 1",
			IsCompleted: false,
			Audit: &domain.Audit{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			ID:          "2",
			Title:       "Title 2",
			Description: "Description 2",
			IsCompleted: false,
			Audit: &domain.Audit{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	mockError := errors.New("some error")

	t.Run("Success", func(t *testing.T) {
		mockTodoRepo.On("Get", context.TODO(), nil, int64(0), int64(10)).Return(mockResult, int64(len(mockResult)), nil).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, total, err := svc.Get(context.TODO(), nil, int64(0), int64(10))

		assert.Nil(t, err)
		assert.Equal(t, len(mockResult), len(res))
		assert.EqualValues(t, 2, total)
	})

	t.Run("Failed", func(t *testing.T) {
		mockTodoRepo.On("Get", context.TODO(), nil, int64(0), int64(10)).Return([]domain.Todo{}, int64(0), mockError).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, total, err := svc.Get(context.TODO(), nil, int64(0), int64(10))

		assert.NotNil(t, err)
		assert.Equal(t, 0, len(res))
		assert.EqualValues(t, 0, total)
	})
}

func TestGetByID(t *testing.T) {
	mockTodoRepo := new(mocks.TodoRepository)

	mockResult := &domain.Todo{
		ID:          "1",
		Title:       "Title 1",
		Description: "Description 1",
		IsCompleted: false,
		Audit: &domain.Audit{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockError := errors.New("some error")

	t.Run("Success", func(t *testing.T) {
		mockTodoRepo.On("GetByID", context.TODO(), "1").Return(mockResult, nil).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, err := svc.GetByID(context.TODO(), "1")

		assert.Nil(t, err)
		assert.Equal(t, mockResult, res)
	})

	t.Run("Failed", func(t *testing.T) {
		mockTodoRepo.On("GetByID", context.TODO(), "").Return(nil, mockError).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, err := svc.GetByID(context.TODO(), "")

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestUpdate(t *testing.T) {
	mockTodoRepo := new(mocks.TodoRepository)

	mockDto := &domain.TodoDto{
		Title:       "Title 1 - Updated",
		Description: "Description 1 - Updated",
	}

	mockResult := &domain.Todo{
		ID:          "1",
		Title:       mockDto.Title,
		Description: mockDto.Description,
		IsCompleted: false,
		Audit: &domain.Audit{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockError := errors.New("some error")

	t.Run("Success", func(t *testing.T) {
		mockTodoRepo.On("Update", context.TODO(), mockResult.ID, mockDto).Return(mockResult, nil).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, err := svc.Update(context.TODO(), "1", mockDto)

		assert.Nil(t, err)
		assert.Equal(t, mockResult, res)
	})

	t.Run("Failed", func(t *testing.T) {
		mockTodoRepo.On("Update", context.TODO(), mockResult.ID, mockDto).Return(nil, mockError).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, err := svc.Update(context.TODO(), "1", mockDto)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestUpdateStatus(t *testing.T) {
	mockTodoRepo := new(mocks.TodoRepository)

	mockResult := &domain.Todo{
		ID:          "1",
		Title:       "Title 1",
		Description: "Description 1",
		IsCompleted: true,
		Audit: &domain.Audit{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockError := errors.New("some error")

	t.Run("Success", func(t *testing.T) {
		mockTodoRepo.On("UpdateStatus", context.TODO(), mockResult.ID, true).Return(mockResult, nil).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, err := svc.UpdateStatus(context.TODO(), "1", true)

		assert.Nil(t, err)
		assert.Equal(t, mockResult, res)
		assert.True(t, res.IsCompleted)
	})

	t.Run("Failed", func(t *testing.T) {
		mockTodoRepo.On("UpdateStatus", context.TODO(), mockResult.ID, true).Return(nil, mockError).Once()

		svc := service.NewTodoService(mockTodoRepo)
		res, err := svc.UpdateStatus(context.TODO(), "1", true)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestDelete(t *testing.T) {
	mockTodoRepo := new(mocks.TodoRepository)

	mockError := errors.New("some error")

	t.Run("Success", func(t *testing.T) {
		mockTodoRepo.On("Delete", context.TODO(), "1").Return(nil).Once()

		svc := service.NewTodoService(mockTodoRepo)
		err := svc.Delete(context.TODO(), "1")

		assert.Nil(t, err)
	})

	t.Run("Failed", func(t *testing.T) {
		mockTodoRepo.On("Delete", context.TODO(), "1").Return(mockError).Once()

		svc := service.NewTodoService(mockTodoRepo)
		err := svc.Delete(context.TODO(), "1")

		assert.NotNil(t, err)
	})
}
