package helper

import (
	"github.com/ariefsn/go-resik/common"
)

func JsonSuccess(data interface{}) common.M {
	return common.ResponseModel{
		Status:  true,
		Data:    data,
		Message: "",
	}.ToM()
}

func JsonError(err error) common.M {
	return common.ResponseModel{
		Status:  false,
		Data:    nil,
		Message: err.Error(),
	}.ToM()
}
