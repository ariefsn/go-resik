package domain

import (
	"context"
)

// Todo: Todo model struct
type Todo struct {
	ID          string `json:"id" bson:"_id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	IsCompleted bool   `json:"isCompleted" bson:"isCompleted"`
	*Audit
}

func (t Todo) TableName() string {
	return "todos"
}

// TodoDto: TodoDto model struct
type TodoDto struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// TodoService represent the todo's usecases
type TodoService interface {
	Get(ctx context.Context, filter interface{}, skip, limit int64) ([]Todo, int64, error)
	GetByID(ctx context.Context, id string) (*Todo, error)
	Update(ctx context.Context, id string, payload *TodoDto) (*Todo, error)
	UpdateStatus(ctx context.Context, id string, isCompleted bool) (*Todo, error)
	Create(ctx context.Context, payload *TodoDto) (*Todo, error)
	Delete(ctx context.Context, id string) error
}

// TodoRepository represent the todo's repository contract
type TodoRepository interface {
	Get(ctx context.Context, filter interface{}, skip, limit int64) ([]Todo, int64, error)
	GetByID(ctx context.Context, id string) (*Todo, error)
	Update(ctx context.Context, id string, payload *TodoDto) (*Todo, error)
	UpdateStatus(ctx context.Context, id string, isCompleted bool) (*Todo, error)
	Create(ctx context.Context, payload *TodoDto) (*Todo, error)
	Delete(ctx context.Context, id string) error
}
