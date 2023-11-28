package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/ariefsn/go-resik/app/todo/repository/mongo"
	"github.com/ariefsn/go-resik/domain"
	"github.com/ariefsn/go-resik/helper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mockDto := &domain.TodoDto{
		Title:       "Title - 1",
		Description: "Description - 1",
	}

	mt.Run("Success", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateSuccessResponse())

		res, err := mockRepo.Create(context.TODO(), mockDto)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, mockDto.Title, res.Title)
		assert.Equal(t, mockDto.Description, res.Description)
	})

	mt.Run("Failed", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    12,
			Message: "some error",
		}))

		res, err := mockRepo.Create(context.TODO(), mockDto)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestGet(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

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

	mockResultBsonD := []bson.D{}

	for _, v := range mockResult {
		b, _ := helper.ToBsonD(v)
		mockResultBsonD = append(mockResultBsonD, *b)
	}

	mt.Run("Success", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		// Counts
		t.AddMockResponses(mtest.CreateCursorResponse(1, "test.todos", mtest.FirstBatch, bson.D{{Key: "n", Value: len(mockResult)}}))
		t.AddMockResponses(mtest.CreateCursorResponse(0, "test.todos", mtest.FirstBatch, mockResultBsonD...))

		res, total, err := mockRepo.Get(context.TODO(), nil, 0, 10)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.EqualValues(t, len(mockResult), total)
	})

	mt.Run("Success With Filter", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		// Counts
		t.AddMockResponses(mtest.CreateCursorResponse(1, "test.todos", mtest.FirstBatch, bson.D{{Key: "n", Value: 1}}))
		t.AddMockResponses(mtest.CreateCursorResponse(0, "test.todos", mtest.FirstBatch, mockResultBsonD...))

		res, total, err := mockRepo.Get(context.TODO(), bson.M{"title": "Title 1"}, 0, 10)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.EqualValues(t, 1, total)
	})

	mt.Run("Failed - CountDocuments", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    12,
			Message: "some error",
		}))

		res, total, err := mockRepo.Get(context.TODO(), bson.M{"title": "Title 3"}, 0, 10)

		assert.NotNil(t, err)
		assert.Equal(t, []domain.Todo{}, res)
		assert.EqualValues(t, 0, total)
	})

	mt.Run("Failed - Aggregate", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		// Counts
		t.AddMockResponses(mtest.CreateCursorResponse(1, "test.todos", mtest.FirstBatch, bson.D{{Key: "n", Value: 1}}))
		t.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    12,
			Message: "some error",
		}))

		res, total, err := mockRepo.Get(context.TODO(), bson.M{"title": "Title 3"}, 0, 10)

		assert.NotNil(t, err)
		assert.Equal(t, []domain.Todo{}, res)
		assert.EqualValues(t, 0, total)
	})
}
