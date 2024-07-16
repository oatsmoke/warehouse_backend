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
			employee.POST("/create", h.Employee.createEmployee)
			employee.POST("/getById", h.Employee.getByIdEmployee)
			employee.POST("/getByDepartment", h.Employee.getByDepartmentEmployee)
			employee.GET("/getAll", h.Employee.getAllEmployee)
			employee.GET("/getFree", h.Employee.getFreeEmployee)
			employee.GET("/getAllButAuth", h.Employee.getAllButAuthEmployee)
			employee.POST("/getAllButOne", h.Employee.getAllButOneEmployee)
			employee.POST("/addToDepartment", h.Employee.addToDepartmentEmployee)
			employee.POST("/removeFromDepartment", h.Employee.removeFromDepartmentEmployee)
			employee.POST("/update", h.Employee.updateEmployee)
			employee.POST("/delete", h.Employee.deleteEmployee)
			employee.POST("/activate", h.Employee.activateEmployee)
			employee.POST("/deactivate", h.Employee.deactivateEmployee)
			employee.POST("/resetPassword", h.Employee.resetPasswordEmployee)
			employee.POST("/changeRole", h.Employee.changeRoleEmployee)
		}
		department := api.Group("/department")
		{
			department.POST("/create", h.Department.createDepartment)
			department.POST("/getById", h.Department.getByIdDepartment)
			department.GET("/getAll", h.Department.getAllDepartment)
			department.POST("/getAllButOne", h.Department.getAllButOneDepartment)
			department.POST("/update", h.Department.updateDepartment)
			department.POST("/delete", h.Department.deleteDepartment)
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
			profile.POST("/create", h.Profile.createProfile)
			profile.POST("/getById", h.Profile.getByIdProfile)
			profile.GET("/getAll", h.Profile.getAllProfile)
			profile.POST("/update", h.Profile.updateProfile)
			profile.POST("/delete", h.Profile.deleteProfile)
		}
		equipment := api.Group("/equipment")
		{
			equipment.POST("/create", h.Equipment.createEquipment)
			equipment.POST("/getById", h.Equipment.getByIdEquipment)
			equipment.POST("/getByIds", h.Equipment.getByIdsEquipment)
			equipment.POST("/getByLocation", h.Equipment.GetByLocationEquipment)
			equipment.GET("/getAll", h.Equipment.getAllEquipment)
			equipment.POST("/update", h.Equipment.updateEquipment)
			equipment.POST("/delete", h.Equipment.deleteEquipment)
			equipment.POST("/reportByCategory", h.Equipment.reportByCategory)
		}
		location := api.Group("/location")
		{
			location.POST("/transferTo", h.Location.transferToLocation)
			location.POST("/getHistory", h.Location.getHistory)
			location.POST("/delete", h.Location.deleteLocation)
		}
		contract := api.Group("/contract")
		{
			contract.POST("/create", h.Contract.createContract)
			contract.POST("/getById", h.Contract.getByIdContract)
			contract.GET("/getAll", h.Contract.getAllContract)
			contract.POST("/update", h.Contract.updateContract)
			contract.POST("/delete", h.Contract.deleteContract)
		}
		company := api.Group("/company")
		{
			company.POST("/create", h.Company.createCompany)
			company.POST("/getById", h.Company.getByIdCompany)
			company.GET("/getAll", h.Company.getAllCompany)
			company.POST("/update", h.Company.updateCompany)
			company.POST("/delete", h.Company.deleteCompany)
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