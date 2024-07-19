package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type ProfileHandler struct {
	ProfileService service.Profile
}

func NewProfileHandler(profileService service.Profile) *ProfileHandler {
	return &ProfileHandler{
		ProfileService: profileService,
	}
}

// Create is profile create
func (h *ProfileHandler) Create(ctx *gin.Context) {
	const fn = "handler.Profile.Create"

	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.ProfileService.Create(ctx, profile.Title, profile.Category.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", profile.Title), fn)
	ctx.JSON(http.StatusOK, "")
}

// Update is profile update
func (h *ProfileHandler) Update(ctx *gin.Context) {
	const fn = "handler.Profile.Update"

	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.ProfileService.Update(ctx, profile.ID, profile.Title, profile.Category.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", profile.Title), fn)
	ctx.JSON(http.StatusOK, "")
}

// Delete is profile delete
func (h *ProfileHandler) Delete(ctx *gin.Context) {
	const fn = "handler.Profile.Delete"

	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.ProfileService.Delete(ctx, profile.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", profile.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// Restore is profile restore
func (h *ProfileHandler) Restore(ctx *gin.Context) {
	const fn = "handler.Profile.Restore"

	var profile *model.Category
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.ProfileService.Restore(ctx, profile.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", profile.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all profiles
func (h *ProfileHandler) GetAll(ctx *gin.Context) {
	const fn = "handler.Profile.GetAll"

	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.ProfileService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("profiles list sended (deleted = %t)", deleted), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get profile by id
func (h *ProfileHandler) GetById(ctx *gin.Context) {
	const fn = "handler.Profile.GetById"

	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.ProfileService.GetById(ctx, profile.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("profile %s found", res.Title), fn)
	ctx.JSON(http.StatusOK, res)
}
