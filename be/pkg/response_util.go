package pkg

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
)

func Null() interface{} {
	return nil
}

func BuildResponse[T any](responseStatus constant.ResponseStatus, data T) dto.ApiResponse[T] {
	return BuildResponse_(responseStatus.GetResponseStatus(), responseStatus.GetResponseMessage(), data)
}

func BuildResponse_[T any](status string, message string, data T) dto.ApiResponse[T] {
	return dto.ApiResponse[T]{
		ResponseMessage: message,
		Data:            data,
	}
}

func BuildResponseSuccess[T any](responseStatus constant.ResponseStatus, data T) dto.ApiResponseSuccess[T] {
	return dto.ApiResponseSuccess[T]{
		Msg:  responseStatus.GetResponseMessage(),
		Data: data,
	}
}

func BuildResponseSuccessNoData() dto.ApiResponseSuccessNoData {
	return dto.ApiResponseSuccessNoData{
		Success: true,
	}
}

func BuildResponseFail(message string) dto.ApiResponseFail {
	return dto.ApiResponseFail{
		Success: false,
		Msg:     message,
	}
}

func BuildResponseSuccessWithFriendsList(friends []string, count int64) dto.ApiResponseSuccessWithFriendsList {
	return dto.ApiResponseSuccessWithFriendsList{
		Success: true,
		Friends: friends,
		Count:   count,
	}
}

func BuildResponseSuccessWithRecipients(recipients []string) dto.ApiResponseSuccessWithRecipients {
	return dto.ApiResponseSuccessWithRecipients{
		Success:    true,
		Recipients: recipients,
	}
}

func BuildResponseSuccessWithTokens(accessToken, refreshToken string) dto.ApiResponseSuccessWithTokens {
	return dto.ApiResponseSuccessWithTokens{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func BuildResponseSuccessWithAccessToken(accessToken string) dto.ApiResponseSuccessWithAccessToken {
	return dto.ApiResponseSuccessWithAccessToken{
		Success:     true,
		AccessToken: accessToken,
	}
}
