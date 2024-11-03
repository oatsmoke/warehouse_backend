package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type LocationHandler struct {
	LocationService service.Location
}

func NewLocationHandler(locationService service.Location) *LocationHandler {
	return &LocationHandler{
		LocationService: locationService,
	}
}

// TransferTo is equipment transfer to
func (h *LocationHandler) TransferTo(ctx *gin.Context) {
	const fn = "handler.Location.TransferTo"

	userId, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized, fn)
		return
	}

	var request []*model.RequestLocation
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.LocationService.TransferTo(ctx, userId, request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole("transferred", fn)
	ctx.JSON(http.StatusOK, "")
}

// Delete is equipment location delete
func (h *LocationHandler) Delete(ctx *gin.Context) {
	fn := "handler.Location.Delete"

	var location *model.Location
	if err := ctx.BindJSON(&location); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.LocationService.Delete(ctx, location.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d, deleted", location.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// GetById is equipment get by id
func (h *LocationHandler) GetById(ctx *gin.Context) {
	fn := "handler.Location.GetById"

	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.LocationService.GetById(ctx, equipment.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d, get", equipment.ID), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetByIds is equipment get by ids
func (h *LocationHandler) GetByIds(ctx *gin.Context) {
	fn := "handler.Location.GetByIds"

	var request []int64
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.LocationService.GetByIds(ctx, request)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d, get", request), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetHistory is equipment get history
func (h *LocationHandler) GetHistory(ctx *gin.Context) {
	const fn = "handler.Location.GetHistory"

	var location *model.Location
	if err := ctx.BindJSON(&location); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.LocationService.GetHistory(ctx, location.Equipment.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d, get", location.Equipment.ID), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetByLocation is equipment get by location
func (h *LocationHandler) GetByLocation(ctx *gin.Context) {
	const fn = "handler.Location.GetByLocation"

	location := new(model.Location)
	location.ToDepartment = new(model.Department)
	location.ToEmployee = new(model.Employee)
	location.ToContract = new(model.Contract)
	if err := ctx.BindJSON(&location); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.LocationService.GetByLocation(ctx, location.ToDepartment.ID, location.ToEmployee.ID, location.ToContract.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("sent location department: %d, employee: %d, contract: %d", location.ToDepartment.ID, location.ToEmployee.ID, location.ToContract.ID), fn)
	ctx.JSON(http.StatusOK, res)
}

// ReportByCategory is equipment report by category
func (h *LocationHandler) ReportByCategory(ctx *gin.Context) {
	const fn = "handler.Location.ReportByCategory"

	request := struct {
		DepartmentId int64     `json:"departmentId"`
		Date         time.Time `json:"date"`
	}{}
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.LocationService.ReportByCategory(ctx, request.DepartmentId, request.Date)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole("sent report", fn)
	ctx.JSON(http.StatusOK, res)
}
