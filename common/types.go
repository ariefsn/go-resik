package common

import "encoding/json"

type M map[string]interface{}

type ResponseModel struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (m ResponseModel) ToM() M {
	res := M{}

	jByte, _ := json.Marshal(m)
	json.Unmarshal(jByte, &res)

	return res
}
