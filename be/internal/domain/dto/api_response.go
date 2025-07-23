package dto

type ApiResponse[T any] struct {
	ResponseMessage string `json:"message"`
	Data            T      `json:"data"`
}

type ApiResponseSuccess[T any] struct {
	Msg  string `json:"message"`
	Data T      `json:"data"`
}

type ApiResponseFail struct {
	Success bool   `json:"success"`
	Msg     string `json:"error"`
}

type ApiResponseSuccessNoData struct {
	Success bool `json:"success"`
}

type ApiResponseSuccessStruct struct {
	Message string  `json:"message" example:"Success"`
	Data    *string `json:"data" example:"null"`
}
