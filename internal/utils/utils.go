package utils

import "github.com/leenzstra/activation_service/internal/responses"

func WrapResponse(result bool, message string, data interface{}) *responses.ResponseBase {
	return &responses.ResponseBase{
		Result:  result,
		Message: message,
		Data:    data,
	}
}