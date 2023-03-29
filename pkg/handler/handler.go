package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"warehouse_backend/pkg/service"
)

type ConfigClient struct {
	Protocol string `json:"protocol"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
}
type errorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(cfg ConfigClient) *gin.Engine {
	router := gin.New()
	router.Use(CORSMiddleware(cfg))
	auth := router.Group("/auth")
	{
		auth.POST("/sing-in", h.signIn)
	}
	api := router.Group("/api", h.userIdentity)
	{
		api.GET("/getUser", h.getUser)
		employee := api.Group("/employee")
		{
			employee.POST("/create", h.createEmployee)
			employee.POST("/getById", h.getByIdEmployee)
			employee.POST("/getByDepartment", h.getByDepartmentEmployee)
			employee.GET("/getAll", h.getAllEmployee)
			employee.GET("/getFree", h.getFreeEmployee)
			employee.GET("/getAllButAuth", h.getAllButAuthEmployee)
			employee.POST("/getAllButOne", h.getAllButOneEmployee)
			employee.POST("/addToDepartment", h.addToDepartmentEmployee)
			employee.POST("/removeFromDepartment", h.removeFromDepartmentEmployee)
			employee.POST("/update", h.updateEmployee)
			employee.POST("/delete", h.deleteEmployee)
			employee.POST("/activate", h.activateEmployee)
			employee.POST("/deactivate", h.deactivateEmployee)
			employee.POST("/resetPassword", h.resetPasswordEmployee)
		}
		department := api.Group("/department")
		{
			department.POST("/create", h.createDepartment)
			department.POST("/getById", h.getByIdDepartment)
			department.GET("/getAll", h.getAllDepartment)
			department.POST("/getAllButOne", h.getAllButOneDepartment)
			department.POST("/update", h.updateDepartment)
			department.POST("/delete", h.deleteDepartment)
		}
		category := api.Group("/category")
		{
			category.POST("/create", h.createCategory)
			category.POST("/getById", h.getByIdCategory)
			category.GET("/getAll", h.getAllCategory)
			category.POST("/update", h.updateCategory)
			category.POST("/delete", h.deleteCategory)
		}
		profile := api.Group("/profile")
		{
			profile.POST("/create", h.createProfile)
			profile.POST("/getById", h.getByIdProfile)
			profile.GET("/getAll", h.getAllProfile)
			profile.POST("/update", h.updateProfile)
			profile.POST("/delete", h.deleteProfile)
		}
		equipment := api.Group("/equipment")
		{
			equipment.POST("/create", h.createEquipment)
			equipment.POST("/getById", h.getByIdEquipment)
			equipment.POST("/getByIds", h.getByIdsEquipment)
			equipment.POST("/getByLocation", h.GetByLocationEquipment)
			equipment.GET("/getAll", h.getAllEquipment)
			equipment.POST("/update", h.updateEquipment)
			equipment.POST("/delete", h.deleteEquipment)
		}
		location := api.Group("/location")
		{
			location.POST("/transferTo", h.transferToLocation)
			location.POST("/getHistory", h.getHistory)
			location.POST("/delete", h.deleteLocation)
		}
		contract := api.Group("/contract")
		{
			contract.POST("/create", h.createContract)
			contract.POST("/getById", h.getByIdContract)
			contract.GET("/getAll", h.getAllContract)
			contract.POST("/update", h.updateContract)
			contract.POST("/delete", h.deleteContract)
		}
	}
	return router
}

func CORSMiddleware(cfg ConfigClient) gin.HandlerFunc {
	clientStr := fmt.Sprintf("%s://%s:%s", cfg.Protocol, cfg.Ip, cfg.Port)
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", clientStr)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
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
