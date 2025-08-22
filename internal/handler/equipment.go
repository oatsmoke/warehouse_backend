package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
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
	userId, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	var request *model.LocationAndRequestLocation
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	id, err := h.EquipmentService.Create(ctx, request.Location.Equipment.SerialNumber, request.Location.Equipment.Profile.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	if err := h.LocationService.AddToStorage(ctx, request.Location.Date, id, userId, request.Location.Company.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	if request.RequestLocation != nil {
		request.RequestLocation[0].EquipmentId = id
		if request.RequestLocation[0].ToDepartment != 0 ||
			request.RequestLocation[0].ToEmployee != 0 ||
			request.RequestLocation[0].ToContract != 0 {
			if err := h.LocationService.TransferTo(ctx, userId, request.RequestLocation); err != nil {
				logger.ErrResponse(ctx, err, http.StatusInternalServerError)
				return
			}
		}
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", request.Location.Equipment.SerialNumber))
	ctx.JSON(http.StatusOK, "")
}

// Update is equipment update
func (h *EquipmentHandler) Update(ctx *gin.Context) {
	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.EquipmentService.Update(ctx, equipment.ID, equipment.SerialNumber, equipment.Profile.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", equipment.SerialNumber))
	ctx.JSON(http.StatusOK, "")
}

// Delete is equipment delete
func (h *EquipmentHandler) Delete(ctx *gin.Context) {
	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.EquipmentService.Delete(ctx, equipment.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s deleted", equipment.SerialNumber))
	ctx.JSON(http.StatusOK, "")
}

// Restore is equipment restore
func (h *EquipmentHandler) Restore(ctx *gin.Context) {
	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.EquipmentService.Restore(ctx, equipment.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s restored", equipment.SerialNumber))
	ctx.JSON(http.StatusOK, "")
}

// GetAll is equipment get all
func (h *EquipmentHandler) GetAll(ctx *gin.Context) {
	res, err := h.EquipmentService.GetAll(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole("equipments list sent")
	ctx.JSON(http.StatusOK, res)
}

// GetByIds is equipment get by id
func (h *EquipmentHandler) GetByIds(ctx *gin.Context) {
	var req []int64
	if err := ctx.BindJSON(&req); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.EquipmentService.GetByIds(ctx, req)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d, get", req))
	ctx.JSON(http.StatusOK, res)
}

// FindBySerialNumber is equipment find by serial number
func (h *EquipmentHandler) FindBySerialNumber(ctx *gin.Context) {
	req := make(map[string]string)
	if err := ctx.BindJSON(&req); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.EquipmentService.FindBySerialNumber(ctx, req["search"])
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s, find", req))
	ctx.JSON(http.StatusOK, res)
}
