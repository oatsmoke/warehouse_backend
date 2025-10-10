package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/list_filter"
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
	var req *dto.Employee
	if err := ctx.BindJSON(&req); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	employee := &model.Employee{
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		Phone:      req.Phone,
	}

	if err := h.serviceEmployee.Create(ctx, employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, "")
}

func (h *EmployeeHandler) Read(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.serviceEmployee.Read(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *EmployeeHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	var req *dto.Employee
	if err := ctx.BindJSON(&req); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	employee := &model.Employee{
		ID:         id,
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		Phone:      req.Phone,
		Department: &model.Department{
			ID: req.DepartmentID,
		},
	}

	if err := h.serviceEmployee.Update(ctx, employee); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *EmployeeHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.Delete(ctx, id); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *EmployeeHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.serviceEmployee.Restore(ctx, id); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *EmployeeHandler) List(ctx *gin.Context) {
	req := list_filter.ParseQueryParams(ctx)

	res, err := h.serviceEmployee.List(ctx, req)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// GetAllShort is employee get all short
//func (h *EmployeeHandler) GetAllShort(ctx *gin.Context) {
//	var deleted bool
//	if err := ctx.BindJSON(&deleted); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	employees, err := h.serviceEmployee.GetAllShort(ctx, deleted)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("employees list sended (deleted = %t)", deleted))
//	ctx.JSON(http.StatusOK, employees)
//}
//
//// GetAllButAuth is employee get all but auth
//func (h *EmployeeHandler) GetAllButAuth(ctx *gin.Context) {
//	var deleted bool
//	if err := ctx.BindJSON(&deleted); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	id, err := getUserId(ctx)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
//		return
//	}
//
//	employees, err := h.serviceEmployee.GetAllButOne(ctx, id, deleted)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole("employees list sent")
//	ctx.JSON(http.StatusOK, employees)
//}
//
//// GetAllButOne is employee get all but one
//func (h *EmployeeHandler) GetAllButOne(ctx *gin.Context) {
//	var employees *model.Employee
//	if err := ctx.BindJSON(&employees); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	res, err := h.serviceEmployee.GetAllButOne(ctx, employees.ID, employees.Deleted)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole("employees list sent")
//	ctx.JSON(http.StatusOK, res)
//}
//
//// GetById is employee get by id
//
//// GetFree is employee get free
//func (h *EmployeeHandler) GetFree(ctx *gin.Context) {
//	employees, err := h.serviceEmployee.GetFree(ctx)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole("employees list sent ")
//	ctx.JSON(http.StatusOK, employees)
//}
//
//// GetByDepartment is employee get by department
//func (h *EmployeeHandler) GetByDepartment(ctx *gin.Context) {
//	var request *model.RequestEmployee
//	if err := ctx.BindJSON(&request); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	employees, err := h.serviceEmployee.GetByDepartment(ctx, request.Ids, request.DepartmentId)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("employees get by department %d found", request.DepartmentId))
//	ctx.JSON(http.StatusOK, employees)
//}
//
//// AddToDepartment is employee add to department
//func (h *EmployeeHandler) AddToDepartment(ctx *gin.Context) {
//	var employee *model.Employee
//	if err := ctx.BindJSON(&employee); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	if err := h.serviceEmployee.AddToDepartment(ctx, employee.ID, employee.Department.ID); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("%d added to %d", employee.ID, employee.Department.ID))
//	ctx.JSON(http.StatusOK, "")
//}
//
//// RemoveFromDepartment is employee remove from department
//func (h *EmployeeHandler) RemoveFromDepartment(ctx *gin.Context) {
//	var request []int64
//	if err := ctx.BindJSON(&request); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	if err := h.serviceEmployee.RemoveFromDepartment(ctx, request[0], request[1]); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("%d removed from %d", request[0], request[1]))
//	ctx.JSON(http.StatusOK, "")
//}
//
//// Activate is employee activate
//func (h *EmployeeHandler) Activate(ctx *gin.Context) {
//	var employee *model.Employee
//	if err := ctx.BindJSON(&employee); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	if err := h.serviceEmployee.Activate(ctx, employee.ID); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("%d activated", employee.ID))
//	ctx.JSON(http.StatusOK, "")
//}
//
//// Deactivate is employee deactivate
//func (h *EmployeeHandler) Deactivate(ctx *gin.Context) {
//	var employee *model.Employee
//	if err := ctx.BindJSON(&employee); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	if err := h.serviceEmployee.Deactivate(ctx, employee.ID); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("%d deactivated", employee.ID))
//	ctx.JSON(http.StatusOK, "")
//}
//
//// ResetPassword is employee reset password
//func (h *EmployeeHandler) ResetPassword(ctx *gin.Context) {
//	var employee *model.Employee
//	if err := ctx.BindJSON(&employee); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	if err := h.serviceEmployee.ResetPassword(ctx, employee.ID); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("%d reset password", employee.ID))
//	ctx.JSON(http.StatusOK, "")
//}
//
//// ChangeRole is employee change role
//func (h *EmployeeHandler) ChangeRole(ctx *gin.Context) {
//	var employee *model.Employee
//	if err := ctx.BindJSON(&employee); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusBadRequest)
//		return
//	}
//
//	if err := h.serviceEmployee.ChangeRole(ctx, employee.ID, employee.Role); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("%d change role", employee.ID))
//	ctx.JSON(http.StatusOK, "")
//}
