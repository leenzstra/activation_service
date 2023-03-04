package utils

import "github.com/leenzstra/activation_service/internal/responses"

// Обёртка для ответа от сервера формата:
// {
// 	"result":bool,
// 	"data":interface{},
// 	"message":string
// }
func WrapResponse(result bool, message string, data interface{}) *responses.ResponseBase {
	return &responses.ResponseBase{
		Result:  result,
		Message: message,
		Data:    data,
	}
}