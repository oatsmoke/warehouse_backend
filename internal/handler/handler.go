package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

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

func New(service *service.Service) *Handler {
	return &Handler{
		Auth:       NewAuthHandler(service.Auth, service.Employee),
		Category:   NewCategoryHandler(service.Category),
		Company:    NewCompanyHandler(service.Company),
		Contract:   NewContractHandler(service.Contract),
		Department: NewDepartmentHandler(service.Department),
		Employee:   NewEmployeeHandler(service.Employee),
		Equipment:  NewEquipmentHandler(service.Equipment),
		Location:   NewLocationHandler(service.Location),
		Profile:    NewProfileHandler(service.Profile),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	if err := router.SetTrustedProxies(nil); err != nil {
		return nil
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{env.GetClientUrl()},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/singIn", h.Auth.SignIn)
	}

	api := router.Group("/api") //, h.Auth.UserIdentity
	{
		api.GET("/getUser", h.Auth.GetUser)

		employee := api.Group("/employee")
		{
			employee.POST("/create", h.Employee.Create)
			employee.POST("/update", h.Employee.Update)
			employee.POST("/delete", h.Employee.Delete)
			employee.POST("/restore", h.Employee.Restore)
			employee.POST("/getAll", h.Employee.GetAll)
			employee.POST("/getAllShort", h.Employee.GetAllShort)
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

		category := api.Group("/categories")
		{
			category.POST("", h.Category.Create)
			category.GET("/:id", h.Category.Read)
			category.PUT("/:id", h.Category.Update)
			category.DELETE("/:id", h.Category.Delete)
			category.PUT("/:id/restore", h.Category.Restore)
			category.GET("", h.Category.List)
		}

		profile := api.Group("/profiles")
		{
			profile.POST("", h.Profile.Create)
			profile.GET("/:id", h.Profile.Read)
			profile.PUT("/:id", h.Profile.Update)
			profile.DELETE("/:id", h.Profile.Delete)
			profile.PUT("/:id/restore", h.Profile.Restore)
			profile.GET("", h.Profile.List)
		}

		equipment := api.Group("/equipments")
		{
			equipment.POST("", h.Equipment.Create)
			equipment.GET("/:id", h.Equipment.Read)
			equipment.PUT("/:id", h.Equipment.Update)
			equipment.DELETE("/:id", h.Equipment.Delete)
			equipment.PUT("/:id/restore", h.Equipment.Restore)
			equipment.GET("", h.Equipment.List)
		}

		location := api.Group("/location")
		{
			location.POST("/transferTo", h.Location.TransferTo)
			location.POST("/delete", h.Location.Delete)
			location.POST("/getById", h.Location.GetById)
			location.POST("/getByIds", h.Location.GetByIds)
			location.POST("/getHistory", h.Location.GetHistory)
			location.POST("/getByLocation", h.Location.GetByLocation)
			location.POST("/reportByCategory", h.Location.ReportByCategory)
		}

		contract := api.Group("/contracts")
		{
			contract.POST("", h.Contract.Create)
			contract.GET("/:id", h.Contract.Read)
			contract.PUT("/:id", h.Contract.Update)
			contract.DELETE("/:id", h.Contract.Delete)
			contract.PUT("/:id/restore", h.Contract.Restore)
			contract.GET("", h.Contract.List)
		}

		company := api.Group("/companies")
		{
			company.POST("", h.Company.Create)
			company.GET("/:id", h.Company.Read)
			company.PUT("/:id", h.Company.Update)
			company.DELETE("/:id", h.Company.Delete)
			company.PUT("/:id/restore", h.Company.Restore)
			company.GET("", h.Company.List)
		}
	}

	return router
}
