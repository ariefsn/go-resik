package api_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ariefsn/go-resik/app/todo/delivery/api"
	"github.com/ariefsn/go-resik/common"
	"github.com/ariefsn/go-resik/domain"
	"github.com/ariefsn/go-resik/domain/mocks"
	"github.com/ariefsn/go-resik/helper"
	"github.com/stretchr/testify/assert"
)

var MOCK_CTX = context.Background()

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

var MOCK_DTO_M, _ = helper.FromJson[common.M](MOCK_DTO)
var MOCK_DTO_UPDATE_M, _ = helper.FromJson[common.M](MOCK_DTO_UPDATE)
var MOCK_DATA_SINGLE_M, _ = helper.FromJson[map[string]interface{}](MOCK_DATA_SINGLE)

var svc = new(mocks.TodoService)

func TestCreate(t *testing.T) {
	cases := []struct {
		name    string
		success bool
		payload interface{}
	}{
		{
			name:    "Success",
			success: true,
			payload: MOCK_DTO,
		},
		{
			name:    "Failed",
			success: false,
			payload: MOCK_DTO,
		},
		{
			name:    "Failed - DTO",
			success: false,
			payload: nil,
		},
	}

	app := api.NewTodoApi(svc)

	for _, c := range cases {
		if c.success {
			svc.On("Create", MOCK_CTX, MOCK_DTO).Return(&MOCK_DATA_SINGLE, nil).Once()
		} else {
			if c.payload == nil {
				svc.On("Create", MOCK_CTX, nil).Return(nil, errors.New("unexpected end of JSON input")).Once()
			} else {
				svc.On("Create", MOCK_CTX, MOCK_DTO).Return(nil, errors.New("create failed")).Once()
			}
		}

		body := new(bytes.Buffer)
		if c.payload != nil {
			body, _ = helper.ToJsonBody(MOCK_DTO)
		}

		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		result, _ := helper.FromResponseBody[common.ResponseModel](res.Body)

		if c.success {
			assert.True(t, result.Status)
			assert.Equal(t, "", result.Message)
			assert.Equal(t, MOCK_DATA_SINGLE_M, result.Data)
		} else {
			assert.False(t, result.Status)
			assert.Nil(t, nil, result.Data)

			if c.payload == nil {
				assert.Equal(t, "unexpected end of JSON input", result.Message)
			} else {
				assert.Equal(t, "create failed", result.Message)
			}
		}
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		name    string
		success bool
		filter  common.M
	}{
		{
			name:    "Success",
			success: true,
			filter: common.M{
				"title":       "Title 1",
				"description": "Description 1",
			},
		},
		{
			name:    "Failed",
			success: false,
			filter:  common.M{},
		},
	}

	app := api.NewTodoApi(svc)

	for _, c := range cases {
		if c.success {
			svc.On("Get", MOCK_CTX, c.filter, int64(0), int64(10)).Return([]domain.Todo{MOCK_DATA_SINGLE}, int64(1), nil).Once()
		} else {
			svc.On("Get", MOCK_CTX, c.filter, int64(0), int64(10)).Return([]domain.Todo{}, int64(0), errors.New("some error")).Once()
		}

		params := url.Values{}
		params.Add("skip", "0")
		params.Add("limit", "10")

		for k, v := range c.filter {
			params.Add(k, v.(string))
		}

		req := httptest.NewRequest(http.MethodGet, "/?"+params.Encode(), nil)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		result, _ := helper.FromResponseBody[common.ResponseModel](res.Body)

		if c.success {
			assert.True(t, result.Status)
			assert.Equal(t, "", result.Message)
			assert.EqualValues(t, map[string]interface{}{
				"items": []interface{}{MOCK_DATA_SINGLE_M},
				"total": float64(1),
			}, result.Data)
		} else {
			assert.False(t, result.Status)
			assert.Nil(t, nil, result.Data)
			assert.Equal(t, "some error", result.Message)
		}
	}
}

func TestGetByID(t *testing.T) {
	cases := []struct {
		name    string
		success bool
		filter  common.M
	}{
		{
			name:    "Success",
			success: true,
			filter: common.M{
				"id": "1",
			},
		},
		{
			name:    "Failed",
			success: false,
			filter: common.M{
				"id": "1",
			},
		},
	}

	app := api.NewTodoApi(svc)

	for _, c := range cases {
		if c.success {
			svc.On("GetByID", MOCK_CTX, "1").Return(&MOCK_DATA_SINGLE, nil).Once()
		} else {
			svc.On("GetByID", MOCK_CTX, "1").Return(nil, errors.New("some error")).Once()
		}

		params := url.Values{}

		for k, v := range c.filter {
			params.Add(k, v.(string))
		}

		req := httptest.NewRequest(http.MethodGet, "/"+params.Get("id"), nil)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		result, _ := helper.FromResponseBody[common.ResponseModel](res.Body)

		if c.success {
			assert.True(t, result.Status)
			assert.Equal(t, "", result.Message)
			assert.EqualValues(t, MOCK_DATA_SINGLE_M, result.Data)
		} else {
			assert.False(t, result.Status)
			assert.Nil(t, nil, result.Data)
			assert.Equal(t, "some error", result.Message)
		}
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name    string
		success bool
		id      string
		payload interface{}
	}{
		{
			name:    "Success",
			success: true,
			id:      "1",
			payload: MOCK_DTO_UPDATE_M,
		},
		{
			name:    "Failed",
			success: false,
			id:      "1",
			payload: MOCK_DTO_UPDATE_M,
		},
		{
			name:    "Failed - Invalid Payload",
			success: false,
			id:      "1",
			payload: nil,
		},
	}

	app := api.NewTodoApi(svc)

	for _, c := range cases {
		if c.success {
			svc.On("Update", MOCK_CTX, "1", &MOCK_DTO_UPDATE).Return(&MOCK_DATA_SINGLE, nil).Once()
		} else {
			if c.payload == nil {
				svc.On("Update", MOCK_CTX, "1", nil).Return(nil, errors.New("unexpected end of JSON input")).Once()
			} else {
				svc.On("Update", MOCK_CTX, "1", &MOCK_DTO_UPDATE).Return(nil, errors.New("some error")).Once()
			}
		}

		body := new(bytes.Buffer)
		if c.payload != nil {
			body, _ = helper.ToJsonBody(MOCK_DTO_UPDATE)
		}

		req := httptest.NewRequest(http.MethodPut, "/"+c.id, body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		result, _ := helper.FromResponseBody[common.ResponseModel](res.Body)

		if c.success {
			assert.True(t, result.Status)
			assert.Equal(t, "", result.Message)
			assert.EqualValues(t, MOCK_DATA_SINGLE_M, result.Data)
		} else {
			assert.False(t, result.Status)
			assert.Nil(t, nil, result.Data)
			msg := ""
			if c.payload == nil {
				msg = "unexpected end of JSON input"
			} else {
				msg = "some error"
			}
			assert.Equal(t, msg, result.Message)
		}
	}
}

func TestUpdateStatus(t *testing.T) {
	cases := []struct {
		name    string
		success bool
		id      string
		payload interface{}
	}{
		{
			name:    "Success",
			success: true,
			id:      "1",
			payload: common.M{
				"isCompleted": true,
			},
		},
		{
			name:    "Failed",
			success: false,
			id:      "1",
			payload: common.M{
				"isCompleted": true,
			},
		},
		{
			name:    "Failed - Invalid Payload",
			success: false,
			id:      "1",
			payload: nil,
		},
	}

	app := api.NewTodoApi(svc)

	for _, c := range cases {
		if c.success {
			svc.On("UpdateStatus", MOCK_CTX, "1", true).Return(&MOCK_DATA_SINGLE, nil).Once()
		} else {
			if c.payload == nil {
				svc.On("UpdateStatus", MOCK_CTX, "1", nil).Return(nil, errors.New("unexpected end of JSON input")).Once()
			} else {
				svc.On("UpdateStatus", MOCK_CTX, "1", true).Return(nil, errors.New("some error")).Once()
			}
		}

		body := new(bytes.Buffer)
		if c.payload != nil {
			body, _ = helper.ToJsonBody(c.payload)
		}

		req := httptest.NewRequest(http.MethodPatch, "/"+c.id, body)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		result, _ := helper.FromResponseBody[common.ResponseModel](res.Body)

		if c.success {
			assert.True(t, result.Status)
			assert.Equal(t, "", result.Message)
			assert.EqualValues(t, MOCK_DATA_SINGLE_M, result.Data)
		} else {
			assert.False(t, result.Status)
			assert.Nil(t, nil, result.Data)
			msg := ""
			if c.payload == nil {
				msg = "unexpected end of JSON input"
			} else {
				msg = "some error"
			}
			assert.Equal(t, msg, result.Message)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name    string
		success bool
		id      string
	}{
		{
			name:    "Success",
			success: true,
			id:      "1",
		},
		{
			name:    "Failed",
			success: false,
			id:      "1",
		},
	}

	app := api.NewTodoApi(svc)

	for _, c := range cases {
		if c.success {
			svc.On("Delete", MOCK_CTX, "1").Return(nil).Once()
		} else {
			svc.On("Delete", MOCK_CTX, "1").Return(errors.New("some error")).Once()
		}

		req := httptest.NewRequest(http.MethodDelete, "/"+c.id, nil)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		result, _ := helper.FromResponseBody[common.ResponseModel](res.Body)

		if c.success {
			assert.True(t, result.Status)
			assert.Equal(t, "", result.Message)
			assert.EqualValues(t, "1", result.Data)
		} else {
			assert.False(t, result.Status)
			assert.Nil(t, nil, result.Data)
			msg := "some error"
			assert.Equal(t, msg, result.Message)
		}
	}
}
