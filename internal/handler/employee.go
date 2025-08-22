package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
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
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.Create(ctx, employee.Name, employee.Phone, employee.Email); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", employee.Name))
	ctx.JSON(http.StatusOK, "")
}

// Update is employee update
func (h *EmployeeHandler) Update(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.Update(ctx, employee.ID, employee.Name, employee.Phone, employee.Email); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", employee.Name))
	ctx.JSON(http.StatusOK, "")
}

// Delete is employee delete
func (h *EmployeeHandler) Delete(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.Delete(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", employee.ID))
	ctx.JSON(http.StatusOK, "")
}

// Restore is employee restore
func (h *EmployeeHandler) Restore(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.Restore(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restore", employee.ID))
	ctx.JSON(http.StatusOK, "")
}

// GetAll is employee get all
func (h *EmployeeHandler) GetAll(ctx *gin.Context) {
	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	employees, err := h.serviceEmployee.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("employees list sended (deleted = %t)", deleted))
	ctx.JSON(http.StatusOK, employees)
}

// GetAllShort is employee get all short
func (h *EmployeeHandler) GetAllShort(ctx *gin.Context) {
	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	employees, err := h.serviceEmployee.GetAllShort(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("employees list sended (deleted = %t)", deleted))
	ctx.JSON(http.StatusOK, employees)
}

// GetAllButAuth is employee get all but auth
func (h *EmployeeHandler) GetAllButAuth(ctx *gin.Context) {
	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	id, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	employees, err := h.serviceEmployee.GetAllButOne(ctx, id, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole("employees list sent")
	ctx.JSON(http.StatusOK, employees)
}

// GetAllButOne is employee get all but one
func (h *EmployeeHandler) GetAllButOne(ctx *gin.Context) {
	var employees *model.Employee
	if err := ctx.BindJSON(&employees); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.serviceEmployee.GetAllButOne(ctx, employees.ID, employees.Deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole("employees list sent")
	ctx.JSON(http.StatusOK, res)
}

// GetById is employee get by id
func (h *EmployeeHandler) GetById(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.serviceEmployee.GetById(ctx, employee.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("employee %s found", res.Name))
	ctx.JSON(http.StatusOK, res)
}

// GetFree is employee get free
func (h *EmployeeHandler) GetFree(ctx *gin.Context) {
	employees, err := h.serviceEmployee.GetFree(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole("employees list sent ")
	ctx.JSON(http.StatusOK, employees)
}

// GetByDepartment is employee get by department
func (h *EmployeeHandler) GetByDepartment(ctx *gin.Context) {
	var request *model.RequestEmployee
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	employees, err := h.serviceEmployee.GetByDepartment(ctx, request.Ids, request.DepartmentId)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("employees get by department %d found", request.DepartmentId))
	ctx.JSON(http.StatusOK, employees)
}

// AddToDepartment is employee add to department
func (h *EmployeeHandler) AddToDepartment(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.AddToDepartment(ctx, employee.ID, employee.Department.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d added to %d", employee.ID, employee.Department.ID))
	ctx.JSON(http.StatusOK, "")
}

// RemoveFromDepartment is employee remove from department
func (h *EmployeeHandler) RemoveFromDepartment(ctx *gin.Context) {
	var request []int64
	if err := ctx.BindJSON(&request); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.RemoveFromDepartment(ctx, request[0], request[1]); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d removed from %d", request[0], request[1]))
	ctx.JSON(http.StatusOK, "")
}

// Activate is employee activate
func (h *EmployeeHandler) Activate(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.Activate(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d activated", employee.ID))
	ctx.JSON(http.StatusOK, "")
}

// Deactivate is employee deactivate
func (h *EmployeeHandler) Deactivate(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.Deactivate(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deactivated", employee.ID))
	ctx.JSON(http.StatusOK, "")
}

// ResetPassword is employee reset password
func (h *EmployeeHandler) ResetPassword(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.ResetPassword(ctx, employee.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d reset password", employee.ID))
	ctx.JSON(http.StatusOK, "")
}

// ChangeRole is employee change role
func (h *EmployeeHandler) ChangeRole(ctx *gin.Context) {
	var employee *model.Employee
	if err := ctx.BindJSON(&employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.ChangeRole(ctx, employee.ID, employee.Role); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d change role", employee.ID))
	ctx.JSON(http.StatusOK, "")
}
