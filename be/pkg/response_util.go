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
		ResponseKey:     status,
		ResponseMessage: message,
		Data:            data,
	}
}

func BuildReponseSuccess[T any](status int, responseStatus constant.ResponseStatus, data T) dto.ApiResponseSuccess[T] {
	return dto.ApiResponseSuccess[T]{
		Status: status,
		Msg:    responseStatus.GetResponseMessage(),
		Data:   data,
	}
}

func BuildReponseSuccessNoData(status int, responseStatus constant.ResponseStatus) dto.ApiResponseSuccessNoData {
	return dto.ApiResponseSuccessNoData{
		Status: status,
		Msg:    responseStatus.GetResponseMessage(),
	}
}

func BuildReponseFail(status int, message string) dto.ApiResponseFail {
	return dto.ApiResponseFail{
		Status: status,
		Msg:    message,
	}
}
