package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type EmployeeHandler struct {
	serviceEmployee service.Employee
}

func NewEmployeeHandler(serviceEmployee service.Employee) *EmployeeHandler {
	return &EmployeeHandler{
		serviceEmployee: serviceEmployee,
	}
}

// Create is employee create
func (h *EmployeeHandler) Create(ctx *gin.Context) {
	const fn = "handler.Employee.Create"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.Create(ctx, employee.Name, employee.Phone, employee.Email); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", employee.Name), fn)
	ctx.JSON(http.StatusOK, "")
}

// Update is employee update
func (h *EmployeeHandler) Update(ctx *gin.Context) {
	const fn = "handler.Employee.Update"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.Update(ctx, employee.ID, employee.Name, employee.Phone, employee.Email); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", employee.Name), fn)
	ctx.JSON(http.StatusOK, "")
}

// Delete is employee delete
func (h *EmployeeHandler) Delete(ctx *gin.Context) {
	const fn = "handler.Employee.Delete"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.Delete(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", employee.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// Restore is employee restore
func (h *EmployeeHandler) Restore(ctx *gin.Context) {
	const fn = "handler.Employee.Restore"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.Restore(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restore", employee.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// GetAll is employee get all
func (h *EmployeeHandler) GetAll(ctx *gin.Context) {
	const fn = "handler.Employee.GetAll"

	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	employees, err := h.serviceEmployee.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("employees list sended (deleted = %t)", deleted), fn)
	ctx.JSON(http.StatusOK, employees)
}

// GetAllButAuth is employee get all but auth
func (h *EmployeeHandler) GetAllButAuth(ctx *gin.Context) {
	const fn = "handler.Employee.GetAllButAuth"

	id, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized, fn)
		return
	}

	employees, err := h.serviceEmployee.GetAllButOne(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole("employees list sent", fn)
	ctx.JSON(http.StatusOK, employees)
}

// GetAllButOne is employee get all but one
func (h *EmployeeHandler) GetAllButOne(ctx *gin.Context) {
	const fn = "handler.Employee.GetAllButOne"

	var employees *model.Employee
	if err := ctx.BindJSON(&employees); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.serviceEmployee.GetAllButOne(ctx, employees.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole("employees list sent", fn)
	ctx.JSON(http.StatusOK, res)
}

// GetById is employee get by id
func (h *EmployeeHandler) GetById(ctx *gin.Context) {
	const fn = "handler.Employee.GetById"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.serviceEmployee.GetById(ctx, employee.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("employee %s found", res.Name), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetFree is employee get free
func (h *EmployeeHandler) GetFree(ctx *gin.Context) {
	const fn = "handler.Employee.GetFree"

	employees, err := h.serviceEmployee.GetFree(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole("employees list sent ", fn)
	ctx.JSON(http.StatusOK, employees)
}

// GetByDepartment is employee get by department
func (h *EmployeeHandler) GetByDepartment(ctx *gin.Context) {
	const fn = "handler.Employee.GetByDepartment"

	var request *model.RequestEmployee
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	employees, err := h.serviceEmployee.GetByDepartment(ctx, request.Ids, request.DepartmentId)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("employees get by department %d found", request.DepartmentId), fn)
	ctx.JSON(http.StatusOK, employees)
}

// AddToDepartment is employee add to department
func (h *EmployeeHandler) AddToDepartment(ctx *gin.Context) {
	const fn = "handler.Employee.AddToDepartment"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.AddToDepartment(ctx, employee.ID, employee.Department.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d added to %d", employee.ID, employee.Department.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// RemoveFromDepartment is employee remove from department
func (h *EmployeeHandler) RemoveFromDepartment(ctx *gin.Context) {
	const fn = "handler.Employee.RemoveFromDepartment"

	var request []int64
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.RemoveFromDepartment(ctx, request[0], request[1]); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d removed from %d", request[0], request[1]), fn)
	ctx.JSON(http.StatusOK, "")
}

// Activate is employee activate
func (h *EmployeeHandler) Activate(ctx *gin.Context) {
	const fn = "handler.Employee.Activate"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.Activate(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d activated", employee.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// Deactivate is employee deactivate
func (h *EmployeeHandler) Deactivate(ctx *gin.Context) {
	const fn = "handler.Employee.Deactivate"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.Deactivate(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deactivated", employee.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// ResetPassword is employee reset password
func (h *EmployeeHandler) ResetPassword(ctx *gin.Context) {
	const fn = "handler.Employee.ResetPassword"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.ResetPassword(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d reset password", employee.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// ChangeRole is employee change role
func (h *EmployeeHandler) ChangeRole(ctx *gin.Context) {
	const fn = "handler.Employee.ChangeRole"

	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.serviceEmployee.ChangeRole(ctx, employee.ID, employee.Role); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d change role", employee.ID), fn)
	ctx.JSON(http.StatusOK, "")
}
