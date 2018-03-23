package cberr

import (
	"encoding/json"

	cc "github.com/clearblade/Go-SDK/cberr/constants/category"
	uuid "github.com/clearblade/Go-SDK/uuid"
)

type Info struct {
	Id            uuid.Uuid           `json:"id"`
	Code          int                 `json:"code"`
	Level         int8                `json:"level"`
	Category      cc.CategoryConstant `json:"category"`
	Message       string              `json:"message"`
	Detail        string              `json:"detail"`
	LowLevelError error               `json:"lowLevelError"`
	Line          string              `json:"line"`
}

type Response struct {
	Info       Info `json:"error"`
	StatusCode int  `json:"statusCode"`
}

func (e *Response) Error() string {
	marsha, _ := json.Marshal(e)
	return string(marsha)
}

func (e *Response) HTTPStatusCode() int {
	return e.StatusCode
}

func (e *Response) FromNet(body interface{}) *Response {
	var errorMap Response
	err := json.Unmarshal(body.([]byte), &errorMap)
	if err != nil {
		return nil
	}
	e.Info = errorMap.Info
	e.StatusCode = errorMap.StatusCode
	return e
}

func CreateResponseFromNet(body interface{}) *Response {
	var errResp Response
	return errResp.FromNet(body)
}
