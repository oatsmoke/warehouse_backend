package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
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
	var contract *model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.ContractService.Create(ctx, contract.Number, contract.Address); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", contract.Number))
	ctx.JSON(http.StatusOK, "")
}

// Update is contract update
func (h *ContractHandler) Update(ctx *gin.Context) {
	var contract *model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.ContractService.Update(ctx, contract.ID, contract.Number, contract.Address); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", contract.Number))
	ctx.JSON(http.StatusOK, "")
}

// Delete is contract delete
func (h *ContractHandler) Delete(ctx *gin.Context) {
	var contract model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.ContractService.Delete(ctx, contract.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", contract.ID))
	ctx.JSON(http.StatusOK, "")
}

// Restore is contract restore
func (h *ContractHandler) Restore(ctx *gin.Context) {
	var contract model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.ContractService.Restore(ctx, contract.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", contract.ID))
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all contracts
func (h *ContractHandler) GetAll(ctx *gin.Context) {
	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.ContractService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("contracts list sended (deleted = %t)", deleted))
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get contract by id
func (h *ContractHandler) GetById(ctx *gin.Context) {
	var contract *model.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.ContractService.GetById(ctx, contract.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("contract %d sended", contract.ID))
	ctx.JSON(http.StatusOK, res)
}
