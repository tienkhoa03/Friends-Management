package handler

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/block_relationship"
	"BE_Friends_Management/pkg"
	"BE_Friends_Management/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type BlockRelationshipHandler struct {
	service service.BlockRelationshipService
}

func NewBlockRelationshipHandler(service service.BlockRelationshipService) *BlockRelationshipHandler {
	return &BlockRelationshipHandler{service: service}
}

// Block godoc
// @Summary      Create new block relationship
// @Description  Create new block relationship
// @Tags         BlockRelationship
// @Accept 		json
// @Produce      json
// @Param 		 request body dto.CreateBlockRequest true "Requestor's email and target's email"
// @param Authorization header string true "Authorization"
// @Router       /api/block [POST]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *BlockRelationshipHandler) CreateBlockRelationship(c *gin.Context) {
	defer pkg.PanicHandler(c)
	authUserId := utils.GetAuthUserId(c)
	var request dto.CreateBlockRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	err := h.service.CreateBlockRelationship(authUserId, request.Requestor, request.Target)
	if err != nil {
		log.Error("Happened error when creating new block relationship. Error: ", err)
		switch {
		case errors.Is(err, service.ErrInvalidRequest):
			pkg.PanicExeption(constant.InvalidRequest, err.Error())
		case errors.Is(err, service.ErrUserNotFound):
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		case errors.Is(err, service.ErrAlreadyBlocked):
			pkg.PanicExeption(constant.Conflict, err.Error())
		case errors.Is(err, service.ErrNotSubscribed):
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrNotPermitted):
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when creating new block relationship.")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessNoData())
}
