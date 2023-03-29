package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/pkg/model"
)

func (h *Handler) createProfile(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var profile model.Profile
	if err := c.BindJSON(&profile); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Profile.Create(profile.Title, profile.Category.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) getByIdProfile(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var p model.Profile
	if err := c.BindJSON(&p); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	profile, err := h.service.Profile.GetById(p.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, profile)
}

func (h *Handler) getAllProfile(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	profiles, err := h.service.Profile.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, profiles)
}

func (h *Handler) updateProfile(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var profile model.Profile
	if err := c.BindJSON(&profile); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Profile.Update(profile.Id, profile.Title, profile.Category.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) deleteProfile(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var profile model.Profile
	if err := c.BindJSON(&profile); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Profile.Delete(profile.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}
