package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"warehouse_backend/internal/lib/config"
	"warehouse_backend/internal/service"
)

type errorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	Auth       *AuthHandler
	Category   *CategoryHandler
	Company    *CompanyHandler
	Contract   *ContractHandler
	Department *DepartmentHandler
	Employee   *EmployeeHandler
	Equipment  *EquipmentHandler
	Location   *LocationHandler
	Profile    *ProfileHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Auth:       NewAuthHandler(service.Auth, service.Employee),
		Category:   NewCategoryHandler(service.Category),
		Company:    NewCompanyHandler(service.Company),
		Contract:   NewContractHandler(service.Contract),
		Department: NewDepartmentHandler(service.Department),
		Employee:   NewEmployeeHandler(service.Employee),
		Equipment:  NewEquipmentHandler(service.Equipment, service.Location),
		Location:   NewLocationHandler(service.Location),
		Profile:    NewProfileHandler(service.Profile),
	}
}

func (h *Handler) InitRoutes(cfg *config.Client) *gin.Engine {
	router := gin.New()
	router.ForwardedByClientIP = true
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil
	}
	router.Use(CORSMiddleware(cfg))
	auth := router.Group("/auth")
	{
		auth.POST("/sing-in", h.Auth.SignIn)
	}
	api := router.Group("/api", h.Auth.UserIdentity)
	{
		api.GET("/getUser", h.Auth.GetUser)
		employee := api.Group("/employee")
		{
			employee.POST("/create", h.Employee.Create)
			employee.POST("/update", h.Employee.Update)
			employee.POST("/delete", h.Employee.Delete)
			employee.POST("/restore", h.Employee.Restore)
			employee.POST("/getAll", h.Employee.GetById)
			employee.POST("/getAllButAuth", h.Employee.GetAllButAuth)
			employee.POST("/getAllButOne", h.Employee.GetAllButOne)
			employee.POST("/getById", h.Employee.GetById)
			employee.POST("/getFree", h.Employee.GetFree)
			employee.POST("/getByDepartment", h.Employee.GetByDepartment)
			employee.POST("/addToDepartment", h.Employee.AddToDepartment)
			employee.POST("/removeFromDepartment", h.Employee.RemoveFromDepartment)
			employee.POST("/activate", h.Employee.Activate)
			employee.POST("/deactivate", h.Employee.Deactivate)
			employee.POST("/resetPassword", h.Employee.ResetPassword)
			employee.POST("/changeRole", h.Employee.ChangeRole)
		}
		department := api.Group("/department")
		{
			department.POST("/create", h.Department.Create)
			department.POST("/update", h.Department.Update)
			department.POST("/delete", h.Department.Delete)
			department.POST("/restore", h.Department.Restore)
			department.POST("/getAll", h.Department.GetAll)
			department.POST("/getById", h.Department.GetById)
			department.POST("/getAllButOne", h.Department.GetAllButOne)
		}
		category := api.Group("/category")
		{
			category.POST("/create", h.Category.Create)
			category.POST("/update", h.Category.Update)
			category.POST("/delete", h.Category.Delete)
			category.POST("/restore", h.Category.Restore)
			category.POST("/getAll", h.Category.GetAll)
			category.POST("/getById", h.Category.GetById)
		}
		profile := api.Group("/profile")
		{
			profile.POST("/create", h.Profile.Create)
			profile.POST("/update", h.Profile.Update)
			profile.POST("/delete", h.Profile.Delete)
			profile.POST("/restore", h.Profile.Restore)
			profile.POST("/getAll", h.Profile.GetAll)
			profile.POST("/getById", h.Profile.GetById)
		}
		equipment := api.Group("/equipment")
		{
			equipment.POST("/create", h.Equipment.Create)
			equipment.POST("/update", h.Equipment.Update)
			equipment.POST("/delete", h.Equipment.Delete)
			equipment.POST("/restore", h.Equipment.Restore)
			equipment.POST("/getAll", h.Equipment.GetAll)
		}
		location := api.Group("/location")
		{
			location.POST("/transferTo", h.Location.transferToLocation)
			location.POST("/getHistory", h.Location.getHistory)
			location.POST("/delete", h.Location.deleteLocation)
			location.POST("/getById", h.Equipment.getByIdEquipment)
			location.POST("/getByIds", h.Equipment.getByIdsEquipment)
			location.POST("/getByLocation", h.Equipment.GetByLocationEquipment)
			location.POST("/reportByCategory", h.Equipment.reportByCategory)
		}
		contract := api.Group("/contract")
		{
			contract.POST("/create", h.Contract.Create)
			contract.POST("/update", h.Contract.Update)
			contract.POST("/delete", h.Contract.Delete)
			contract.POST("/getAll", h.Contract.GetAll)
			contract.POST("/getById", h.Contract.GetById)
		}
		company := api.Group("/company")
		{
			company.POST("/create", h.Company.Create)
			company.POST("/update", h.Company.Update)
			company.POST("/delete", h.Company.Delete)
			company.POST("/restore", h.Company.Restore)
			company.POST("/getAll", h.Company.GetAll)
			company.POST("/getById", h.Company.GetById)
		}
	}
	return router
}

func CORSMiddleware(cfg *config.Client) gin.HandlerFunc {
	clientStr := fmt.Sprintf("%s://%s:%s", cfg.Protocol, cfg.Ip, cfg.Port)
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", clientStr)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(statusCode, " - ", message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
