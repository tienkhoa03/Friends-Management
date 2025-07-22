package dto

type ApiResponse[T any] struct {
	ResponseKey     string `json:"status"`
	ResponseMessage string `json:"message"`
	Data            T      `json:"data"`
}

type ApiResponseSuccess[T any] struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
	Data   T      `json:"data"`
}

type ApiResponseFail struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
}
type ApiResponseSuccessNoData struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
}
type ApiResponseSuccessStruct struct {
	Code    int     `json:"code" example:"200"`
	Message string  `json:"message" example:"Success"`
	Data    *string `json:"data" example:"null"`
}
