package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/list_filter"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type EquipmentHandler struct {
	equipmentService service.Equipment
}

func NewEquipmentHandler(equipmentService service.Equipment) *EquipmentHandler {
	return &EquipmentHandler{
		equipmentService: equipmentService,
	}
}

func (h *EquipmentHandler) Create(ctx *gin.Context) {
	//userId, err := getUserId(ctx)
	//if err != nil {
	//	logger.ResponseErr(ctx, err, http.StatusUnauthorized)
	//	return
	//}

	var req *dto.Equipment
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	equipment := &model.Equipment{
		SerialNumber: req.SerialNumber,
		Profile: &model.Profile{
			ID: req.ProfileID,
		},
	}

	if err := h.equipmentService.Create(ctx, equipment); err != nil {
		if errors.Is(err, logger.ErrAlreadyExists) {
			logger.ResponseErr(ctx, logger.ErrAlreadyExists.Error(), err, http.StatusConflict)
			return
		}
		logger.ResponseErr(ctx, logger.MsgFailedToInsert, err, http.StatusInternalServerError)
		return
	}

	//if err := h.LocationService.AddToStorage(ctx, request.Location.Date, id, userId, request.Location.Company.ID); err != nil {
	//	logger.ResponseErr(ctx, err, http.StatusInternalServerError)
	//	return
	//}
	//
	//if request.RequestLocation != nil {
	//	request.RequestLocation[0].EquipmentId = id
	//	if request.RequestLocation[0].ToDepartment != 0 ||
	//		request.RequestLocation[0].ToEmployee != 0 ||
	//		request.RequestLocation[0].ToContract != 0 {
	//		if err := h.LocationService.TransferTo(ctx, userId, request.RequestLocation); err != nil {
	//			logger.ResponseErr(ctx, err, http.StatusInternalServerError)
	//			return
	//		}
	//	}
	//}

	ctx.JSON(http.StatusCreated, "")
}

func (h *EquipmentHandler) Read(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	res, err := h.equipmentService.Read(ctx, id)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *EquipmentHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	var req *dto.Equipment
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	equipment := &model.Equipment{
		ID:           id,
		SerialNumber: req.SerialNumber,
		Profile: &model.Profile{
			ID: req.ProfileID,
		},
	}

	if err := h.equipmentService.Update(ctx, equipment); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToUpdate, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *EquipmentHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.equipmentService.Delete(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToDelete, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *EquipmentHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.equipmentService.Restore(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToRestore, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *EquipmentHandler) List(ctx *gin.Context) {
	req := list_filter.ParseQueryParams(ctx)

	res, err := h.equipmentService.List(ctx, req)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
