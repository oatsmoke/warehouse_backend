package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type DepartmentHandler struct {
	DepartmentService service.Department
}

func NewDepartmentHandler(departmentService service.Department) *DepartmentHandler {
	return &DepartmentHandler{
		DepartmentService: departmentService,
	}
}

// Create is department create
func (h *DepartmentHandler) Create(ctx *gin.Context) {
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.DepartmentService.Create(ctx, department.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", department.Title))
	ctx.JSON(http.StatusOK, "")
}

// Update is department update
func (h *DepartmentHandler) Update(ctx *gin.Context) {
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.DepartmentService.Update(ctx, department.ID, department.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", department.Title))
	ctx.JSON(http.StatusOK, "")
}

// Delete is department delete
func (h *DepartmentHandler) Delete(ctx *gin.Context) {
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.DepartmentService.Delete(ctx, department.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", department.ID))
	ctx.JSON(http.StatusOK, "")
}

// Restore is department restore
func (h *DepartmentHandler) Restore(ctx *gin.Context) {
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.DepartmentService.Restore(ctx, department.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", department.ID))
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all departments
func (h *DepartmentHandler) GetAll(ctx *gin.Context) {
	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.DepartmentService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("departments list sent (deleted = %t)", deleted))
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get department by id
func (h *DepartmentHandler) GetById(ctx *gin.Context) {
	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.DepartmentService.GetById(ctx, department.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("department %s found", res.Title))
	ctx.JSON(http.StatusOK, res)
}

// GetAllButOne is to get all departments but one
func (h *DepartmentHandler) GetAllButOne(ctx *gin.Context) {
	employeeId, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}
	res, err := h.DepartmentService.GetAllButOne(ctx, department.ID, employeeId)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("departments list sended (except = %d)", department.ID))
	ctx.JSON(http.StatusOK, res)
}
