package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type AuthHandler struct {
	AuthService     service.Auth
	EmployeeService service.Employee
}

func NewAuthHandler(authService service.Auth, employeeService service.Employee) *AuthHandler {
	return &AuthHandler{
		AuthService:     authService,
		EmployeeService: employeeService,
	}
}

// SignIn is user authentication for login
func (h *AuthHandler) SignIn(ctx *gin.Context) {
	req := make(map[string]string)
	if err := ctx.BindJSON(&req); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	id, err := h.AuthService.AuthUser(ctx, req["login"], req["password"])
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateToken(id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	hash, err := h.AuthService.GenerateHash(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	employee, err := h.EmployeeService.GetById(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	ctx.Set("token", token)
	ctx.Set("hash", hash)
	if err := setCookie(ctx); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s authenticated", employee.Name))
	ctx.JSON(http.StatusOK, employee)
}

// GetUser is getting user data
func (h *AuthHandler) GetUser(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	employee, err := h.EmployeeService.GetById(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s authenticated", employee.Name))
	ctx.JSON(http.StatusOK, employee)
}

// UserIdentity is user authentication check
func (h *AuthHandler) UserIdentity(ctx *gin.Context) {
	token, err := ctx.Cookie("token")
	if err != nil {
		logger.ErrInConsole(err)
		h.Refresh(ctx)
		return
	}

	userId, err := jwt.ParseToken(token)
	if err != nil {
		logger.ErrInConsole(err)
		h.Refresh(ctx)
		return
	}

	ctx.Set("userId", userId)
	logger.InfoInConsole(fmt.Sprintf("id set to %b", userId))
}

// Refresh is a user authentication update
func (h *AuthHandler) Refresh(ctx *gin.Context) {
	hash, err := ctx.Cookie("hash")
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	userId, err := h.AuthService.FindByHash(ctx, hash)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	newToken, err := jwt.GenerateToken(userId)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	newHash, err := h.AuthService.GenerateHash(ctx, userId)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.Set("token", newToken)
	ctx.Set("hash", newHash)
	ctx.Set("userId", userId)

	if err := setCookie(ctx); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("token and hash for id %b updated", userId))
}

// getUserId is to get the user id from the context
func getUserId(ctx *gin.Context) (int64, error) {
	if userId, ok := ctx.Get("userId"); !ok {
		return 0, errors.New("user id in context not found")
	} else {
		return userId.(int64), nil
	}
}

// setCookie is to set cookies
func setCookie(ctx *gin.Context) error {
	if token, ok := ctx.Get("token"); ok {
		ctx.SetCookie("token", token.(string), 3600, "/", "", true, true)
	} else {
		return errors.New("token in context not found")
	}

	if hash, ok := ctx.Get("hash"); ok {
		ctx.SetCookie("hash", hash.(string), 604800, "/", "", true, true)
	} else {
		return errors.New("hash in context not found")
	}

	return nil
}
