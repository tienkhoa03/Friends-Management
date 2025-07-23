package pkg

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
)

func Null() interface{} {
	return nil
}

func BuildReponse[T any](responseStatus constant.ResponseStatus, data T) dto.ApiResponse[T] {
	return BuildReponse_(responseStatus.GetResponseStatus(), responseStatus.GetResponseMessage(), data)
}

func BuildReponse_[T any](status string, message string, data T) dto.ApiResponse[T] {
	return dto.ApiResponse[T]{
		ResponseMessage: message,
		Data:            data,
	}
}

func BuildReponseSuccess[T any](responseStatus constant.ResponseStatus, data T) dto.ApiResponseSuccess[T] {
	return dto.ApiResponseSuccess[T]{
		Msg:  responseStatus.GetResponseMessage(),
		Data: data,
	}
}

func BuildReponseSuccessNoData() dto.ApiResponseSuccessNoData {
	return dto.ApiResponseSuccessNoData{
		Success: true,
	}
}

func BuildReponseFail(message string) dto.ApiResponseFail {
	return dto.ApiResponseFail{
		Success: false,
		Msg:     message,
	}
}
