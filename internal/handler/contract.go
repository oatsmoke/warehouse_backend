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

type ContractHandler struct {
	contractService service.Contract
}

func NewContractHandler(contractService service.Contract) *ContractHandler {
	return &ContractHandler{
		contractService: contractService,
	}
}

func (h *ContractHandler) Create(ctx *gin.Context) {
	var req *dto.Contract
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	contract := &model.Contract{
		Number:  req.Number,
		Address: req.Address,
	}

	if err := h.contractService.Create(ctx, contract); err != nil {
		if errors.Is(err, logger.ErrAlreadyExists) {
			logger.ResponseErr(ctx, logger.ErrAlreadyExists.Error(), err, http.StatusConflict)
			return
		}
		logger.ResponseErr(ctx, logger.MsgFailedToInsert, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, "")
}

func (h *ContractHandler) Read(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	res, err := h.contractService.Read(ctx, id)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *ContractHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	var req *dto.Contract
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	contract := &model.Contract{
		ID:      id,
		Number:  req.Number,
		Address: req.Address,
	}

	if err := h.contractService.Update(ctx, contract); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToUpdate, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *ContractHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.contractService.Delete(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToDelete, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *ContractHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.contractService.Restore(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToRestore, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *ContractHandler) List(ctx *gin.Context) {
	req := list_filter.ParseQueryParams(ctx)

	res, err := h.contractService.List(ctx, req)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
