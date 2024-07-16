package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type DepartmentHandler struct {
	serviceDepartment service.Department
}

func NewDepartmentHandler(serviceDepartment service.Department) *DepartmentHandler {
	return &DepartmentHandler{
		serviceDepartment: serviceDepartment,
	}
}

func (h *DepartmentHandler) createDepartment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceDepartment.Create(ctx, department.Title); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *DepartmentHandler) getByIdDepartment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var d *model.Department
	if err := ctx.BindJSON(&d); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	department, err := h.serviceDepartment.GetById(ctx, d.ID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, department)
}

func (h *DepartmentHandler) getAllDepartment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	departments, err := h.serviceDepartment.GetAll(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, departments)
}

func (h *DepartmentHandler) getAllButOneDepartment(ctx *gin.Context) {
	employeeId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	departments, err := h.serviceDepartment.GetAllButOne(ctx, department.ID, employeeId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, departments)
}

func (h *DepartmentHandler) updateDepartment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceDepartment.Update(ctx, department.ID, department.Title); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *DepartmentHandler) deleteDepartment(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceDepartment.Delete(ctx, department.ID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}
