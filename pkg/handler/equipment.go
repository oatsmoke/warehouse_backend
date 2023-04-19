package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/pkg/model"
)

func (h *Handler) createEquipment(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var request model.LocationAndRequestLocation
	if err := c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	equipmentId, err := h.service.Equipment.Create(
		request.Location.Date,
		request.Location.Company.Id,
		request.Location.Equipment.SerialNumber,
		request.Location.Equipment.Profile.Id,
		userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	request.RequestLocation[0].EquipmentId = equipmentId
	if request.RequestLocation[0].ToDepartment != 0 ||
		request.RequestLocation[0].ToEmployee != 0 ||
		request.RequestLocation[0].ToContract != 0 {
		if err := h.service.Location.TransferTo(userId, request.RequestLocation); err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) getByIdEquipment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var e model.Employee
	if err := c.BindJSON(&e); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	equipment, err := h.service.Equipment.GetById(e.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, equipment)
}

func (h *Handler) getByIdsEquipment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var request map[string][]int
	if err := c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	equipments, err := h.service.Equipment.GetByIds(request["ids"])
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, equipments)
}

func (h *Handler) GetByLocationEquipment(c *gin.Context) {
	if _, err := getUserId(c); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var l model.Location
	if err := c.BindJSON(&l); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	equipments, err := h.service.Equipment.GetByLocation(l.ToDepartment.Id, l.ToEmployee.Id, l.ToContract.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, equipments)
}

func (h *Handler) getAllEquipment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	equipments, err := h.service.Equipment.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, equipments)
}

func (h *Handler) updateEquipment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var equipment model.Equipment
	if err := c.BindJSON(&equipment); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Equipment.Update(equipment.Id, equipment.SerialNumber, equipment.Profile.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) deleteEquipment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var equipment model.Equipment
	if err := c.BindJSON(&equipment); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Equipment.Delete(equipment.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}
