package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type LocationHandler struct {
	serviceLocation service.Location
}

func NewLocationHandler(serviceLocation service.Location) *LocationHandler {
	return &LocationHandler{
		serviceLocation: serviceLocation,
	}
}

func (h *LocationHandler) transferToLocation(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var request []*model.RequestLocation
	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceLocation.TransferTo(ctx, userId, request); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}
func (h *LocationHandler) getHistory(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var location *model.Location
	if err := ctx.BindJSON(&location); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	history, err := h.serviceLocation.GetHistory(ctx, location.Equipment.ID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, history)
}

func (h *LocationHandler) deleteLocation(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var location *model.Location
	if err := ctx.BindJSON(&location); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceLocation.Delete(ctx, location.ID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}
