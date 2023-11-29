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
	mdb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var MOCK_DTO = &domain.TodoDto{
	Title:       "Title - 1",
	Description: "Description - 1",
}

var MOCK_DTO_UPDATE = domain.TodoDto{
	Title:       "Title 1 - Updated",
	Description: "Description 1 - Updated",
}

var MOCK_DATA_LIST = []domain.Todo{
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

var MOCK_DATA_SINGLE = domain.Todo{
	ID:          "1",
	Title:       "Title 1",
	Description: "Description 1",
	IsCompleted: false,
	Audit: &domain.Audit{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

var MOCK_DATA_SINGLE_UPDATED = domain.Todo{
	ID:          "1",
	Title:       "Title 1 - Updated",
	Description: "Description 1 - Updated",
	IsCompleted: false,
	Audit: &domain.Audit{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

var MOCK_DATA_SINGLE_STATUS_UPDATED = domain.Todo{
	ID:          "1",
	Title:       "Title 1 - Updated",
	Description: "Description 1 - Updated",
	IsCompleted: true,
	Audit: &domain.Audit{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

var MOCK_DATA_SINGLE_BSOND, _ = helper.ToBsonD(MOCK_DATA_SINGLE)

var MOCK_DATA_SINGLE_UPDATED_BSOND, _ = helper.ToBsonD(MOCK_DATA_SINGLE_UPDATED)

var MOCK_DATA_SINGLE_STATUS_UPDATED_BSOND, _ = helper.ToBsonD(MOCK_DATA_SINGLE_STATUS_UPDATED)

func TestCreate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateSuccessResponse())

		res, err := mockRepo.Create(context.TODO(), MOCK_DTO)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, MOCK_DTO.Title, res.Title)
		assert.Equal(t, MOCK_DTO.Description, res.Description)
	})

	mt.Run("Failed", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    12,
			Message: "some error",
		}))

		res, err := mockRepo.Create(context.TODO(), MOCK_DTO)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestGet(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mockResultBsonD := []bson.D{}

	for _, v := range MOCK_DATA_LIST {
		b, _ := helper.ToBsonD(v)
		mockResultBsonD = append(mockResultBsonD, *b)
	}

	mt.Run("Success", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		// Counts
		t.AddMockResponses(mtest.CreateCursorResponse(1, "test.todos", mtest.FirstBatch, bson.D{{Key: "n", Value: len(MOCK_DATA_LIST)}}))
		t.AddMockResponses(mtest.CreateCursorResponse(0, "test.todos", mtest.FirstBatch, mockResultBsonD...))

		res, total, err := mockRepo.Get(context.TODO(), nil, 0, 10)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.EqualValues(t, len(MOCK_DATA_LIST), total)
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

func TestGetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateCursorResponse(1, "test.todos", mtest.FirstBatch, *MOCK_DATA_SINGLE_BSOND))

		res, err := mockRepo.GetByID(context.TODO(), "1")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, MOCK_DATA_SINGLE.Title, res.Title)
		assert.Equal(t, MOCK_DATA_SINGLE.Description, res.Description)
	})

	mt.Run("Failed", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: mdb.ErrNoDocuments.Error()}))

		res, err := mockRepo.GetByID(context.TODO(), "1")

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(t *mtest.T) {
		db := t.Client.Database("mock-db")
		mockRepo := mongo.NewMongoTodoRepository(db)

		t.AddMockResponses(bson.D{
			{
				Key:   "ok",
				Value: 1,
			},
			{
				Key:   "value",
				Value: MOCK_DATA_SINGLE_UPDATED_BSOND,
			},
		})

		res, err := mockRepo.Update(context.TODO(), "1", &MOCK_DTO_UPDATE)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, MOCK_DATA_SINGLE_UPDATED.Title, res.Title)
		assert.Equal(t, MOCK_DATA_SINGLE_UPDATED.Description, res.Description)
	})

	mt.Run("Failed", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: mdb.ErrNoDocuments.Error()}))

		res, err := mockRepo.Update(context.TODO(), "1", &MOCK_DTO_UPDATE)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestUpdateStatus(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(t *mtest.T) {
		db := t.Client.Database("mock-db")
		mockRepo := mongo.NewMongoTodoRepository(db)

		t.AddMockResponses(bson.D{
			{
				Key:   "ok",
				Value: 1,
			},
			{
				Key:   "value",
				Value: MOCK_DATA_SINGLE_STATUS_UPDATED_BSOND,
			},
		})

		res, err := mockRepo.UpdateStatus(context.TODO(), "1", true)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, MOCK_DATA_SINGLE_STATUS_UPDATED.Title, res.Title)
		assert.Equal(t, MOCK_DATA_SINGLE_STATUS_UPDATED.Description, res.Description)
		assert.Equal(t, MOCK_DATA_SINGLE_STATUS_UPDATED.IsCompleted, res.IsCompleted)
	})

	mt.Run("Failed", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: mdb.ErrNoDocuments.Error()}))

		res, err := mockRepo.UpdateStatus(context.TODO(), "1", true)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(t *mtest.T) {
		db := t.Client.Database("mock-db")
		mockRepo := mongo.NewMongoTodoRepository(db)

		t.AddMockResponses(bson.D{
			{
				Key:   "ok",
				Value: 1,
			},
			{
				Key:   "value",
				Value: MOCK_DATA_SINGLE_BSOND,
			},
		})

		err := mockRepo.Delete(context.TODO(), "1")

		assert.Nil(t, err)
	})

	mt.Run("Failed", func(t *mtest.T) {
		mockRepo := mongo.NewMongoTodoRepository(t.Client.Database("mock-db"))

		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: mdb.ErrNoDocuments.Error()}))

		err := mockRepo.Delete(context.TODO(), "1")

		assert.NotNil(t, err)
	})
}
