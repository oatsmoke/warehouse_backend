package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
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
	const fn = "handler.Department.Create"

	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.DepartmentService.Create(ctx, department.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", department.Title), fn)
	ctx.JSON(http.StatusOK, "")
}

// Update is department update
func (h *DepartmentHandler) Update(ctx *gin.Context) {
	const fn = "handler.Department.Update"

	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.DepartmentService.Update(ctx, department.ID, department.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", department.Title), fn)
	ctx.JSON(http.StatusOK, "")
}

// Delete is department delete
func (h *DepartmentHandler) Delete(ctx *gin.Context) {
	const fn = "handler.Department.Delete"

	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.DepartmentService.Delete(ctx, department.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", department.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// Restore is department restore
func (h *DepartmentHandler) Restore(ctx *gin.Context) {
	const fn = "handler.Department.Restore"

	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.DepartmentService.Restore(ctx, department.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", department.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all departments
func (h *DepartmentHandler) GetAll(ctx *gin.Context) {
	const fn = "handler.Department.GetAll"

	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.DepartmentService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("departments list sent (deleted = %t)", deleted), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get department by id
func (h *DepartmentHandler) GetById(ctx *gin.Context) {
	const fn = "handler.Department.GetById"

	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.DepartmentService.GetById(ctx, department.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("department %s found", res.Title), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetAllButOne is to get all departments but one
func (h *DepartmentHandler) GetAllButOne(ctx *gin.Context) {
	const fn = "handler.Department.GetAllButOne"

	employeeId, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized, fn)
		return
	}

	var department *model.Department
	if err := ctx.BindJSON(&department); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}
	res, err := h.DepartmentService.GetAllButOne(ctx, department.ID, employeeId)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("departments list sended (except = %d)", department.ID), fn)
	ctx.JSON(http.StatusOK, res)
}
