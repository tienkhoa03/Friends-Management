package handler

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/block_relationship"
	"BE_Friends_Management/pkg"
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

// User godoc
// @Summary      Create new block relationship
// @Description  Create new block relationship
// @Tags         BlockRelationship
// @Accept 		json
// @Produce      json
// @Param 		 request body dto.CreateBlockRequest true "Requestor's email and target's email"
// @Router       /api/block [POST]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
func (h *BlockRelationshipHandler) CreateBlockRelationship(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.CreateBlockRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	err := h.service.CreateBlockRelationship(request.Requestor, request.Target)
	if err != nil {
		log.Error("Happened error when creating new block relationship. Error: ", err)
		switch err {
		case service.ErrInvalidRequest:
			pkg.PanicExeption(constant.InvalidRequest, err.Error())
		case service.ErrUserNotFound:
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		case service.ErrAlreadyBlocked:
			pkg.PanicExeption(constant.Conflict, err.Error())
		case service.ErrNotSubscribed:
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when creating new block relationship.")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessNoData())
}
