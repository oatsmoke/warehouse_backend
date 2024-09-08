package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type EquipmentHandler struct {
	EquipmentService service.Equipment
	LocationService  service.Location
}

func NewEquipmentHandler(equipmentService service.Equipment, locationService service.Location) *EquipmentHandler {
	return &EquipmentHandler{
		EquipmentService: equipmentService,
		LocationService:  locationService,
	}
}

// Create is equipment create
func (h *EquipmentHandler) Create(ctx *gin.Context) {
	const fn = "handler.Equipment.Create"

	userId, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized, fn)
		return
	}

	var request *model.LocationAndRequestLocation
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	id, err := h.EquipmentService.Create(ctx, request.Location.Equipment.SerialNumber, request.Location.Equipment.Profile.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	if err := h.LocationService.AddToStorage(ctx, request.Location.Date, id, userId, request.Location.Company.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	if request.RequestLocation != nil {
		request.RequestLocation[0].EquipmentId = id
		if request.RequestLocation[0].ToDepartment != 0 ||
			request.RequestLocation[0].ToEmployee != 0 ||
			request.RequestLocation[0].ToContract != 0 {
			if err := h.LocationService.TransferTo(ctx, userId, request.RequestLocation); err != nil {
				logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
				return
			}
		}
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", request.Location.Equipment.SerialNumber), fn)
	ctx.JSON(http.StatusOK, "")
}

// Update is equipment update
func (h *EquipmentHandler) Update(ctx *gin.Context) {
	const fn = "handler.Equipment.Update"

	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.EquipmentService.Update(ctx, equipment.ID, equipment.SerialNumber, equipment.Profile.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", equipment.SerialNumber), fn)
	ctx.JSON(http.StatusOK, "")
}

// Delete is equipment delete
func (h *EquipmentHandler) Delete(ctx *gin.Context) {
	const fn = "handler.Equipment.Delete"

	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.EquipmentService.Delete(ctx, equipment.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s deleted", equipment.SerialNumber), fn)
	ctx.JSON(http.StatusOK, "")
}

// Restore is equipment restore
func (h *EquipmentHandler) Restore(ctx *gin.Context) {
	const fn = "handler.Equipment.Restore"

	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.EquipmentService.Restore(ctx, equipment.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s restored", equipment.SerialNumber), fn)
	ctx.JSON(http.StatusOK, "")
}

// GetAll is equipment get all
func (h *EquipmentHandler) GetAll(ctx *gin.Context) {
	const fn = "handler.Equipment.GetAll"

	res, err := h.EquipmentService.GetAll(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole("equipments list sent", fn)
	ctx.JSON(http.StatusOK, res)
}

// GetByIds is equipment get by id
func (h *EquipmentHandler) GetByIds(ctx *gin.Context) {
	const fn = "handler.Equipment.GetByIds"

	//req := make(map[string]int64)
	var req []int64
	if err := ctx.BindJSON(&req); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.EquipmentService.GetByIds(ctx, req)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d, get", req), fn)
	ctx.JSON(http.StatusOK, res)
}
