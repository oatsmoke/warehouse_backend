package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"warehouse_backend/pkg/model"
)

func (h *Handler) signIn(c *gin.Context) {
	var input model.SignInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.Employee.FindUser(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := h.service.Employee.GenerateToken(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	hash, err := h.service.Employee.GenerateHash(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Set("token", token)
	c.Set("hash", hash)
	setCookie(c)
	c.JSON(http.StatusOK, map[string]interface{}{"token": token, "hash": hash})
}

func (h *Handler) userIdentity(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		userId, err := h.service.Employee.ParseToken(token)
		if err == nil {
			c.Set("userId", userId)
			return
		}
	} else {
		hash, err := c.Cookie("hash")
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, "no authorization hash")
			return
		}
		id, err := h.service.Employee.FindByHash(hash)
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		newToken, err := h.service.Employee.GenerateToken(id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		newHash, err := h.service.Employee.GenerateHash(id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.Set("token", newToken)
		c.Set("hash", newHash)
		c.Set("userId", strconv.Itoa(id))
	}
}

func (h *Handler) getUser(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	employee, err := h.service.Employee.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, employee)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")
	if !ok {
		return 0, errors.New("user id not found")
	}
	idInt, err := strconv.Atoi(id.(string))
	if err != nil {
		return 0, errors.New("user id is invalid of type")
	}
	return idInt, nil
}

func setCookie(c *gin.Context) {
	token, ok := c.Get("token")
	if ok {
		c.SetCookie("token", token.(string), 3600, "/", "", true, true)
	}
	hash, ok := c.Get("hash")
	if ok {
		c.SetCookie("hash", hash.(string), 604800, "/", "", true, true)
	}
}
