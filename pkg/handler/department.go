package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/pkg/model"
)

func (h *Handler) createDepartment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var department model.Department
	if err := c.BindJSON(&department); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Department.Create(department.Title); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) getByIdDepartment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var d model.Department
	if err := c.BindJSON(&d); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	department, err := h.service.Department.GetById(d.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, department)
}

func (h *Handler) getAllDepartment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	departments, err := h.service.Department.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, departments)
}

func (h *Handler) getAllButOneDepartment(c *gin.Context) {
	employeeId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var department model.Department
	if err := c.BindJSON(&department); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	departments, err := h.service.Department.GetAllButOne(department.Id, employeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, departments)
}

func (h *Handler) updateDepartment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var department model.Department
	if err := c.BindJSON(&department); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Department.Update(department.Id, department.Title); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) deleteDepartment(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var department model.Department
	if err := c.BindJSON(&department); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Department.Delete(department.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}
