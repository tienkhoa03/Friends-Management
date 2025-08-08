package middleware

import (
	"BE_Friends_Management/constant"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"

	"BE_Friends_Management/pkg"
	"BE_Friends_Management/pkg/utils"
	"github.com/gin-gonic/gin"
)

var (
	ErrAccessTokenExpires = errors.New("access token has expired")
	ErrNotPermitted       = errors.New("action not permitted")
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func TimeoutMiddleware(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func ValidateAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer pkg.PanicHandler(c)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Error("Happened error when validating access token. Error: Missing access token")
			pkg.PanicExeption(constant.Unauthorized, "Missing access token.")
		}
		rawAccessToken := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseAccessToken(rawAccessToken)
		if err != nil {
			log.Error("Happened error when validating access token. Error: ", err)
			switch {
			case errors.Is(err, utils.ErrInvalidAccessToken):
				pkg.PanicExeption(constant.Unauthorized, err.Error())
			case errors.Is(err, utils.ErrInvalidSigningMethod):
				pkg.PanicExeption(constant.Unauthorized, err.Error())
			default:
				pkg.PanicExeption(constant.UnknownError, "Invalid Access Token.")
			}
		}
		if claims.ExpiresAt.Time.Before(time.Now()) {
			log.Error("Happened error when validating access token. Error: ", ErrAccessTokenExpires)
			pkg.PanicExeption(constant.Unauthorized, ErrAccessTokenExpires.Error())
		}
		c.Set("authUserId", claims.UserId)
		c.Set("authUserRole", claims.Role)
		c.Next()
	}
}

func RequireAnyRole(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer pkg.PanicHandler(c)
		rawAuthUserRole, exists := c.Get("authUserRole")
		authUserRole := fmt.Sprint(rawAuthUserRole)
		if !exists {
			log.Error("Happened error when validating access token. Error: ", ErrNotPermitted)
			pkg.PanicExeption(constant.StatusForbidden, ErrNotPermitted.Error())
		}
		for _, role := range roles {
			if role == authUserRole {
				c.Next()
				return
			}
		}
		log.Error("Happened error when validating access token. Error: ", ErrNotPermitted)
		pkg.PanicExeption(constant.StatusForbidden, ErrNotPermitted.Error())
	}
}
