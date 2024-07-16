package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type ProfileHandler struct {
	serviceProfile service.Profile
}

func NewProfileHandler(serviceProfile service.Profile) *ProfileHandler {
	return &ProfileHandler{
		serviceProfile: serviceProfile,
	}
}

func (h *ProfileHandler) createProfile(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceProfile.Create(ctx, profile.Title, profile.Category.ID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *ProfileHandler) getByIdProfile(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var p *model.Profile
	if err := ctx.BindJSON(&p); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	profile, err := h.serviceProfile.GetById(ctx, p.ID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) getAllProfile(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	profiles, err := h.serviceProfile.GetAll(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, profiles)
}

func (h *ProfileHandler) updateProfile(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceProfile.Update(ctx, profile.ID, profile.Title, profile.Category.ID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *ProfileHandler) deleteProfile(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceProfile.Delete(ctx, profile.ID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}
