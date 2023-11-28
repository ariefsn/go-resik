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
	"go.mongodb.org/mongo-driver/mongo/options"
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
	res := r.Db.Collection(domain.Todo{}.TableName()).FindOneAndDelete(ctx, bson.M{"_id": id})

	if res.Err() != nil {
		err := helper.ParseMongoError(res.Err())
		logger.Error(err)
	}

	return helper.ParseMongoError(res.Err())
}

// Get implements domain.TodoRepository.
func (r *mongoTodoRepository) Get(ctx context.Context, filter interface{}, skip int64, limit int64) ([]domain.Todo, int64, error) {
	result := []domain.Todo{}
	coll := r.Db.Collection(domain.Todo{}.TableName())

	if filter == nil {
		filter = common.M{}
	}

	_bson, _ := helper.ToBsonM(filter)

	filterBson := bson.M{}

	for k, v := range _bson {
		switch t := v.(type) {
		case string:
			_f := helper.MongoFilter(helper.FoContains, k, t)
			filterBson[k] = _f[k]
		}
	}

	count, err := coll.CountDocuments(ctx, filterBson)

	if err != nil && err != mongo.ErrNilDocument {
		logger.Error(err)
		return result, 0, err
	}

	pipe := helper.MongoPipe(helper.MongoAggregate{
		Match: filterBson,
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
	err := r.Db.Collection(result.TableName()).FindOne(ctx, bson.M{"_id": id}).Decode(&result)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &result, nil
}

func (r *mongoTodoRepository) update(ctx context.Context, filter bson.M, payload bson.M) *mongo.SingleResult {
	upsert := true
	returnDoc := options.After

	return r.Db.Collection(domain.Todo{}.TableName()).FindOneAndUpdate(ctx, filter, bson.M{
		"$set": payload,
	}, &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDoc,
		Upsert:         &upsert,
	})
}

// Update implements domain.TodoRepository.
func (r *mongoTodoRepository) Update(ctx context.Context, id string, payload *domain.TodoDto) (*domain.Todo, error) {
	var data domain.Todo

	res := r.update(ctx, bson.M{"_id": id}, bson.M{
		"title":           payload.Title,
		"description":     payload.Description,
		"audit.updatedAt": time.Now(),
	})

	if res.Err() != nil {
		err := helper.ParseMongoError(res.Err())
		logger.Error(err)
		return nil, err
	}

	res.Decode(&data)

	return &data, nil
}

// UpdateStatus implements domain.TodoRepository.
func (r *mongoTodoRepository) UpdateStatus(ctx context.Context, id string, isCompleted bool) (*domain.Todo, error) {
	var data domain.Todo

	res := r.update(ctx, bson.M{"_id": id}, bson.M{
		"isCompleted":     isCompleted,
		"audit.updatedAt": time.Now(),
	})

	if res.Err() != nil {
		err := helper.ParseMongoError(res.Err())
		logger.Error(err)
		return nil, err
	}

	res.Decode(&data)

	return &data, nil
}

// NewMongoTodoRepository will create an object that represent the todo.Repository interface
func NewMongoTodoRepository(database *mongo.Database) domain.TodoRepository {
	return &mongoTodoRepository{
		Db: database,
	}
}
