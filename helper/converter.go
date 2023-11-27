package helper

import (
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