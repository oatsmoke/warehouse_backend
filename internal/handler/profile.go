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

type ProfileHandler struct {
	profileService service.Profile
}

func NewProfileHandler(profileService service.Profile) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (h *ProfileHandler) Create(ctx *gin.Context) {
	var req *dto.Profile
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	profile := &model.Profile{
		Title: req.Title,
		Category: &model.Category{
			ID: req.CategoryID,
		},
	}

	if err := h.profileService.Create(ctx, profile); err != nil {
		if errors.Is(err, logger.ErrAlreadyExists) {
			logger.ResponseErr(ctx, logger.ErrAlreadyExists.Error(), err, http.StatusConflict)
			return
		}
		logger.ResponseErr(ctx, logger.MsgFailedToInsert, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, "")
}

func (h *ProfileHandler) Read(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	res, err := h.profileService.Read(ctx, id)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *ProfileHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	var req *dto.Profile
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	profile := &model.Profile{
		ID:    id,
		Title: req.Title,
		Category: &model.Category{
			ID: req.CategoryID,
		},
	}

	if err := h.profileService.Update(ctx, profile); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToUpdate, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *ProfileHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.profileService.Delete(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToDelete, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *ProfileHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.profileService.Restore(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToRestore, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *ProfileHandler) List(ctx *gin.Context) {
	req := list_filter.ParseQueryParams(ctx)

	res, err := h.profileService.List(ctx, req)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
