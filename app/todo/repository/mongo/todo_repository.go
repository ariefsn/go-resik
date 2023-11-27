package mongo

import (
	"context"
	"time"

	"github.com/ariefsn/go-resik/common"
	"github.com/ariefsn/go-resik/domain"
	"github.com/ariefsn/go-resik/helper"
	"github.com/ariefsn/go-resik/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTodoRepository struct {
	Db *mongo.Database
}

// Create implements domain.TodoRepository.
func (r *mongoTodoRepository) Create(ctx context.Context, payload *domain.TodoDto) (*domain.Todo, error) {
	data := domain.Todo{
		ID:          primitive.NewObjectID().Hex(),
		Title:       payload.Title,
		Description: payload.Description,
		IsCompleted: false,
		Audit: &domain.Audit{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	_, err := r.Db.Collection(data.TableName()).InsertOne(ctx, data)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &data, nil
}

// Delete implements domain.TodoRepository.
func (r *mongoTodoRepository) Delete(ctx context.Context, id string) error {
	res := r.Db.Collection(domain.Todo{}.TableName()).FindOneAndDelete(ctx, bson.D{
		bson.E{
			Key:   "_id",
			Value: id,
		},
	})

	if res.Err() != nil {
		logger.Error(res.Err())
	}

	return res.Err()
}

// Get implements domain.TodoRepository.
func (r *mongoTodoRepository) Get(ctx context.Context, skip int64, limit int64) ([]domain.Todo, int64, error) {

	result := []domain.Todo{}
	coll := r.Db.Collection(domain.Todo{}.TableName())

	count, err := coll.CountDocuments(ctx, bson.D{})

	if err != nil && err != mongo.ErrNilDocument {
		logger.Error(err)
		return result, 0, err
	}

	pipe := helper.MongoPipe(helper.MongoAggregate{
		Match: nil,
		Sort:  nil,
		Skip:  &skip,
		Limit: &limit,
	})

	cur, err := coll.Aggregate(ctx, pipe)

	if err != nil {
		logger.Error(err)
		return result, 0, err
	}

	for cur.Next(ctx) {
		var row domain.Todo

		err = cur.Decode(&row)
		if err != nil {
			break
		}

		result = append(result, row)
	}

	return result, count, nil
}

// GetByID implements domain.TodoRepository.
func (r *mongoTodoRepository) GetByID(ctx context.Context, id string) (*domain.Todo, error) {
	result := domain.Todo{}
	err := r.Db.Collection(result.TableName()).FindOne(ctx, bson.D{
		bson.E{
			Key:   "_id",
			Value: id,
		},
	}).Decode(&result)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &result, nil
}

// GetByTitle implements domain.TodoRepository.
func (r *mongoTodoRepository) GetByTitle(ctx context.Context, title string) (*domain.Todo, error) {
	panic("unimplemented")
}

// Update implements domain.TodoRepository.
func (r *mongoTodoRepository) Update(ctx context.Context, id string, payload *domain.TodoDto) (*domain.Todo, error) {
	var data domain.Todo
	res := r.Db.Collection(domain.Todo{}.TableName()).FindOneAndUpdate(ctx, bson.D{
		bson.E{
			Key:   "_id",
			Value: id,
		},
	}, common.M{
		"title":           payload.Title,
		"description":     payload.Description,
		"audit.updatedAt": time.Now(),
	})

	if res.Err() != nil {
		logger.Error(res.Err())
		return nil, res.Err()
	}

	res.Decode(&data)

	return &data, nil
}

// UpdateStatus implements domain.TodoRepository.
func (r *mongoTodoRepository) UpdateStatus(ctx context.Context, id string, isCompleted bool) (*domain.Todo, error) {
	panic("unimplemented")
}

// NewMongoTodoRepository will create an object that represent the todo.Repository interface
func NewMongoTodoRepository(database *mongo.Database) domain.TodoRepository {
	return &mongoTodoRepository{
		Db: database,
	}
}
