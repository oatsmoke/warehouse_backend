package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/list_filter"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type DepartmentHandler struct {
	departmentService service.Department
}

func NewDepartmentHandler(departmentService service.Department) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

func (h *DepartmentHandler) Create(ctx *gin.Context) {
	var req *dto.Department
	if err := ctx.BindJSON(&req); err != nil {
		if errors.Is(err, logger.ErrAlreadyExists) {
			logger.ResponseErr(ctx, logger.ErrAlreadyExists.Error(), err, http.StatusConflict)
			return
		}
		logger.ResponseErr(ctx, logger.MsgFailedToInsert, err, http.StatusInternalServerError)
		return
	}

	department := &model.Department{
		Title: req.Title,
	}

	if err := h.departmentService.Create(ctx, department); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToInsert, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, "")
}

func (h *DepartmentHandler) Read(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	res, err := h.departmentService.Read(ctx, id)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *DepartmentHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}
	var req *dto.Department
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	department := &model.Department{
		ID:    id,
		Title: req.Title,
	}

	if err := h.departmentService.Update(ctx, department); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToUpdate, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *DepartmentHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.departmentService.Delete(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToDelete, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *DepartmentHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.departmentService.Restore(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToRestore, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *DepartmentHandler) List(ctx *gin.Context) {
	req := list_filter.ParseQueryParams(ctx)

	res, err := h.departmentService.List(ctx, req)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

//func (h *DepartmentHandler) GetAllButOne(ctx *gin.Context) {
//	employeeId, err := getUserId(ctx)
//	if err != nil {
//		logger.ResponseErr(ctx, err, http.StatusUnauthorized)
//		return
//	}
//
//	var department *model.Department
//	if err := ctx.BindJSON(&department); err != nil {
//		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
//		return
//	}
//	res, err := h.DepartmentService.GetAllButOne(ctx, department.ID, employeeId)
//	if err != nil {
//		logger.ResponseErr(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("departments list sended (except = %d)", department.ID))
//	ctx.JSON(http.StatusOK, res)
//}
