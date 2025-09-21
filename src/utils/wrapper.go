package utils

import (
	"application-wallet/models"
)

func Data(code int, content interface{}, totalItems int, message string) models.Response {
	return models.Response{
		Code:       code,
		Content:    content,
		TotalItems: totalItems,
		Message:    message,
	}
}
