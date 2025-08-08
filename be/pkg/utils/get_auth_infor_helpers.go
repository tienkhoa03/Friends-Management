package utils

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/pkg"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAuthUserId(c *gin.Context) int64 {
	rawAuthUserId, exists := c.Get("authUserId")
	if !exists {
		log.Error("Happened error when creating new friendship. Error: ", ErrInvalidAccessToken)
		pkg.PanicExeption(constant.InvalidRequest, ErrInvalidAccessToken.Error())
	}
	authUserIdStr := fmt.Sprint(rawAuthUserId)
	authUserId, err := strconv.ParseInt(authUserIdStr, 10, 64)
	if err != nil {
		log.Error("Happened error when creating new friendship. Error: ", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when creating new friendship.")
	}
	return authUserId
}

func GetAuthUserRole(c *gin.Context) string {
	rawAuthUserRole, exists := c.Get("authUserRole")
	if !exists {
		log.Error("Happened error when creating new friendship. Error: ", ErrInvalidAccessToken)
		pkg.PanicExeption(constant.InvalidRequest, ErrInvalidAccessToken.Error())
	}
	authUserRole := fmt.Sprint(rawAuthUserRole)
	return authUserRole
}
