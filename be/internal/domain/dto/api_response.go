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

type ApiResponseSuccessWithFriendsList struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int64    `json:"count"`
}

type ApiResponseSuccessWithRecipients struct {
	Success    bool     `json:"success"`
	Recipients []string `json:"recipients"`
}

type ApiResponseSuccessWithTokens struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ApiResponseSuccessWithAccessToken struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token"`
}
