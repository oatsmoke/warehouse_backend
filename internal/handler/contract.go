package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type ContractHandler struct {
	ContractService service.Contract
}

func NewContractHandler(contractService service.Contract) *ContractHandler {
	return &ContractHandler{
		ContractService: contractService,
	}
}

// Create is contract create
func (h *ContractHandler) Create(ctx *gin.Context) {
	const fn = "handler.Contract.Create"

	var contract *model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.ContractService.Create(ctx, contract.Number, contract.Address); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", contract.Number), fn)
	ctx.JSON(http.StatusOK, "")
}

// Update is contract update
func (h *ContractHandler) Update(ctx *gin.Context) {
	const fn = "handler.Contract.Update"

	var contract *model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.ContractService.Update(ctx, contract.ID, contract.Number, contract.Address); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", contract.Number), fn)
	ctx.JSON(http.StatusOK, "")
}

// Delete is contract delete
func (h *ContractHandler) Delete(ctx *gin.Context) {
	const fn = "handler.Contract.Delete"

	var contract model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.ContractService.Delete(ctx, contract.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", contract.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// Restore is contract restore
func (h *ContractHandler) Restore(ctx *gin.Context) {
	const fn = "handler.Contract.Restore"

	var contract model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.ContractService.Restore(ctx, contract.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", contract.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all contracts
func (h *ContractHandler) GetAll(ctx *gin.Context) {
	const fn = "handler.Contract.GetAll"

	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.ContractService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("contracts list sended (deleted = %t)", deleted), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get contract by id
func (h *ContractHandler) GetById(ctx *gin.Context) {
	const fn = "handler.Contract.GetById"

	var contract *model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.ContractService.GetById(ctx, contract.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("contract %d sended", contract.ID), fn)
	ctx.JSON(http.StatusOK, res)
}
