package helper

import (
	"bytes"
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson"
)

type mongoType interface {
	Decode(v interface{}) error
}

func ToBsonM(data interface{}) (bson.M, error) {
	dataMarshalled, err := bson.Marshal(data)

	if err != nil {
		return nil, err
	}

	var bsonM bson.M

	err = bson.Unmarshal(dataMarshalled, &bsonM)

	if err != nil {
		return nil, err
	}

	return bsonM, nil
}

func ToBsonD(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)

	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)

	return
}

func ToJsonBody(v interface{}) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	err := json.NewEncoder(&buff).Encode(v)

	if err != nil {
		return nil, err
	}

	return &buff, nil
}

func FromResponseBody[T any](body io.ReadCloser) (T, error) {
	var target T
	err := json.NewDecoder(body).Decode(&target)

	if err != nil {
		return target, err
	}

	return target, nil
}

func FromJson[T any](from interface{}) (T, error) {
	var to T

	data, err := json.Marshal(from)

	if err != nil {
		return to, err
	}

	err = json.Unmarshal(data, &to)

	if err != nil {
		return to, err
	}

	return to, nil
}
