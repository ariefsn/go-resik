package service

import (
	"context"

	"github.com/ariefsn/go-resik/domain"
)

type todoService struct {
	todoRepo domain.TodoRepository
}

// Create implements domain.TodoService.
func (s *todoService) Create(ctx context.Context, payload *domain.TodoDto) (*domain.Todo, error) {
	return s.todoRepo.Create(ctx, payload)
}

// Delete implements domain.TodoService.
func (s *todoService) Delete(ctx context.Context, id string) error {
	return s.todoRepo.Delete(ctx, id)
}

// Get implements domain.TodoService.
func (s *todoService) Get(ctx context.Context, filter interface{}, skip int64, limit int64) ([]domain.Todo, int64, error) {
	return s.todoRepo.Get(ctx, filter, skip, limit)
}

// GetByID implements domain.TodoService.
func (s *todoService) GetByID(ctx context.Context, id string) (*domain.Todo, error) {
	return s.todoRepo.GetByID(ctx, id)
}

// Update implements domain.TodoService.
func (s *todoService) Update(ctx context.Context, id string, payload *domain.TodoDto) (*domain.Todo, error) {
	return s.todoRepo.Update(ctx, id, payload)
}

// UpdateStatus implements domain.TodoService.
func (s *todoService) UpdateStatus(ctx context.Context, id string, isCompleted bool) (*domain.Todo, error) {
	return s.todoRepo.UpdateStatus(ctx, id, isCompleted)
}

// NewTodoService will create new an todoService object representation of domain.TodoService interface
func NewTodoService(todoRepo domain.TodoRepository) domain.TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}
