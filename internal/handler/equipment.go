package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type EquipmentHandler struct {
	serviceEquipment service.Equipment
	serviceLocation  service.Location
}

func NewEquipmentHandler(serviceEquipment service.Equipment, serviceLocation service.Location) *EquipmentHandler {
	return &EquipmentHandler{
		serviceEquipment: serviceEquipment,
		serviceLocation:  serviceLocation,
	}
}

func (h *EquipmentHandler) createEquipment(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var request *model.LocationAndRequestLocation
	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	equipmentId, err := h.serviceEquipment.Create(ctx, request.Location.Date, request.Location.Company.ID, request.Location.Equipment.SerialNumber, request.Location.Equipment.Profile.ID, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	request.RequestLocation[0].EquipmentId = equipmentId
	if request.RequestLocation[0].ToDepartment != 0 ||
		request.RequestLocation[0].ToEmployee != 0 ||
		request.RequestLocation[0].ToContract != 0 {
		if err := h.serviceLocation.TransferTo(ctx, userId, request.RequestLocation); err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *EquipmentHandler) getByIdEquipment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var e *model.Employee
	if err := ctx.BindJSON(&e); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	equipment, err := h.serviceEquipment.GetById(ctx, e.ID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) getByIdsEquipment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var request map[string][]int64
	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	equipments, err := h.serviceEquipment.GetByIds(ctx, request["ids"])
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, equipments)
}

func (h *EquipmentHandler) GetByLocationEquipment(ctx *gin.Context) {
	if _, err := getUserId(ctx); err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var l *model.Location
	if err := ctx.BindJSON(&l); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	equipments, err := h.serviceEquipment.GetByLocation(ctx, l.ToDepartment.ID, l.ToEmployee.ID, l.ToContract.ID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, equipments)
}

func (h *EquipmentHandler) getAllEquipment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	equipments, err := h.serviceEquipment.GetAll(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, equipments)
}

func (h *EquipmentHandler) updateEquipment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceEquipment.Update(ctx, equipment.ID, equipment.SerialNumber, equipment.Profile.ID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *EquipmentHandler) deleteEquipment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var equipment *model.Equipment
	if err := ctx.BindJSON(&equipment); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceEquipment.Delete(ctx, equipment.ID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *EquipmentHandler) reportByCategory(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	request := struct {
		DepartmentId int64 `json:"departmentId"`
		Date         int64 `json:"date"`
	}{}
	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	report, err := h.serviceEquipment.ReportByCategory(ctx, request.DepartmentId, request.Date)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, report)
}
