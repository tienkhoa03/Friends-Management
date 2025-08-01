package pkg

import (
	"BE_Friends_Management/constant"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var ErrorMessages = map[string]string{
	"uni_users_email":    "Email was used",
	"uni_users_username": "Username was used",
	"idx_name_userID":    "Name was used",
	"idx_name_projectID": "Name was used",
	"idx_name_stepID":    "Name was used",
	"record not found":   "Something went wrong",
	"chk_steps_name":     "Name was short",
}

func PanicExeption_(key string, message string) {
	err := errors.New(message)
	err = fmt.Errorf("%s: %w", key, err)
	if err != nil {
		panic(err)
	}
}

func PanicExeption(responseKey constant.ResponseStatus, customMessage ...string) {
	message := responseKey.GetResponseMessage()
	if len(customMessage) > 0 {
		message = customMessage[0]
	}
	for key, msg := range ErrorMessages {
		if strings.Contains(message, key) {
			message = msg
		}
	}
	PanicExeption_(responseKey.GetResponseStatus(), message)
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func PanicHandler(c *gin.Context) {
	if err := recover(); err != nil {
		str := fmt.Sprint(err)
		strArr := strings.SplitN(str, ":", 2)

		key := strArr[0]
		msg := capitalizeFirst(strings.Trim(strArr[1], " "))
		switch key {
		case constant.InvalidRequest.GetResponseStatus():
			c.JSON(http.StatusBadRequest, BuildResponseFail(msg))
			c.Abort()
		case constant.DataNotFound.GetResponseStatus():
			c.JSON(http.StatusNotFound, BuildResponseFail(msg))
			c.Abort()
		case constant.Unauthorized.GetResponseStatus():
			c.JSON(http.StatusUnauthorized, BuildResponseFail(msg))
			c.Abort()
		case constant.StatusForbidden.GetResponseStatus():
			c.JSON(http.StatusForbidden, BuildResponseFail(msg))
			c.Abort()
		case constant.Conflict.GetResponseStatus():
			c.JSON(http.StatusConflict, BuildResponseFail(msg))
			c.Abort()
		default:
			c.JSON(http.StatusInternalServerError, BuildResponseFail(msg))
			c.Abort()
		}
	}
}
