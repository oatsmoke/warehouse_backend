package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/lib/logger"
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
func (h *EquipmentHandler) getByIdEquipment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized, fn)
		return
	}
	var e *model.Employee
	if err := ctx.BindJSON(&e); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}
	equipment, err := h.EquipmentService.GetById(ctx, e.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) getByIdsEquipment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized, fn)
		return
	}
	var request map[string][]int64
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}
	equipments, err := h.EquipmentService.GetByIds(ctx, request["ids"])
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, equipments)
}

func (h *EquipmentHandler) GetByLocationEquipment(ctx *gin.Context) {
	if _, err := getUserId(ctx); err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized, fn)
		return
	}
	var l *model.Location
	if err := ctx.BindJSON(&l); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}
	equipments, err := h.EquipmentService.GetByLocation(ctx, l.ToDepartment.ID, l.ToEmployee.ID, l.ToContract.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, equipments)
}

func (h *EquipmentHandler) reportByCategory(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized, fn)
		return
	}
	request := struct {
		DepartmentId int64 `json:"departmentId"`
		Date         int64 `json:"date"`
	}{}
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}
	report, err := h.EquipmentService.ReportByCategory(ctx, request.DepartmentId, request.Date)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, report)
}
