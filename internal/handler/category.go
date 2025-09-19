package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type CategoryHandler struct {
	categoryService service.Category
}

func NewCategoryHandler(categoryService service.Category) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Create(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.categoryService.Create(ctx, category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, "")
}

func (h *CategoryHandler) Read(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.categoryService.Read(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("category %d found", id))
	ctx.JSON(http.StatusOK, res)
}

func (h *CategoryHandler) Update(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	category.ID = id
	if err := h.categoryService.Update(ctx, category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *CategoryHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.categoryService.Delete(ctx, id); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *CategoryHandler) Restore(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.categoryService.Restore(ctx, id); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *CategoryHandler) List(ctx *gin.Context) {
	var withDeleted bool
	if ctx.Query("deleted") == "true" {
		withDeleted = true
	}

	res, err := h.categoryService.List(ctx, withDeleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
