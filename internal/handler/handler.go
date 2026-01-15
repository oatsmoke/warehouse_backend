package handler

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
	"github.com/oatsmoke/warehouse_backend/internal/lib/websocket"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type Handler struct {
	Auth       *AuthHandler
	User       *UserHandler
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
		Auth:       NewAuthHandler(service.Auth, service.User),
		User:       NewUserHandler(service.User),
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
	hub := websocket.NewHub()
	go hub.Run()

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
		auth.POST("/login", h.Auth.Login)
	}

	api := router.Group("/api", h.Auth.UserIdentity)
	{
		api.GET("/ws", func(ctx *gin.Context) {
			websocket.NewClient(ctx.Writer, ctx.Request, hub)
		})
		api.GET("/roles", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, role.AllRole())
		})
		api.GET("/user", h.Auth.GetUser)

		user := api.Group("/users")
		{
			user.POST("", h.User.Create)
			user.GET("/:id", h.User.Read)
			user.PUT("/:id", h.User.Update)
			user.DELETE("/:id", h.User.Delete)
			user.GET("", h.User.List)
			user.PUT("/:id/set_password", h.User.SetPassword)
			user.PUT("/:id/reset_password", h.User.ResetPassword)
			user.PUT("/:id/set_enabled", h.User.SetEnabled)
		}

		employee := api.Group("/employees")
		{
			employee.POST("", h.Employee.Create)
			employee.GET("/:id", h.Employee.Read)
			employee.PUT("/:id", h.Employee.Update)
			employee.DELETE("/:id", h.Employee.Delete)
			employee.PUT("/:id/restore", h.Employee.Restore)
			employee.GET("", h.Employee.List)
			employee.PUT("/:id/set_department", h.Employee.SetDepartment)
			//employee.POST("/getAllShort", h.Employee.GetAllShort)
			//employee.POST("/getAllButAuth", h.Employee.GetAllButAuth)
			//employee.POST("/getAllButOne", h.Employee.GetAllButOne)
			//employee.POST("/getFree", h.Employee.GetFree)
			//employee.POST("/getByDepartment", h.Employee.GetByDepartment)
			//employee.POST("/addToDepartment", h.Employee.AddToDepartment)
			//employee.POST("/removeFromDepartment", h.Employee.RemoveFromDepartment)
			//employee.POST("/activate", h.Employee.Activate)
			//employee.POST("/deactivate", h.Employee.Deactivate)
			//employee.POST("/resetPassword", h.Employee.ResetPassword)
			//employee.POST("/changeRole", h.Employee.ChangeRole)
		}

		department := api.Group("/departments")
		{
			department.POST("", h.Department.Create)
			department.GET("/:id", h.Department.Read)
			department.PUT("/:id", h.Department.Update)
			department.DELETE("/:id", h.Department.Delete)
			department.PUT("/:id/restore", h.Department.Restore)
			department.GET("", h.Department.List)
			//department.POST("/getAllButOne", h.Department.GetAllButOne)
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
			company.POST("", h.Auth.AdminAccess, h.Company.Create)
			company.GET("/:id", h.Auth.AdminAccess, h.Company.Read)
			company.PUT("/:id", h.Auth.AdminAccess, h.Company.Update)
			company.DELETE("/:id", h.Auth.AdminAccess, h.Company.Delete)
			company.PUT("/:id/restore", h.Auth.AdminAccess, h.Company.Restore)
			company.GET("", h.Company.List)
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

		location := api.Group("/locations")
		{
			location.GET("", h.Location.List)
			//location.POST("/transferTo", h.Location.TransferTo)
			//location.POST("/delete", h.Location.Delete)
			//location.POST("/getById", h.Location.GetById)
			//location.POST("/getByIds", h.Location.GetByIds)
			//location.POST("/getHistory", h.Location.GetHistory)
			//location.POST("/getByLocation", h.Location.GetByLocation)
			//location.POST("/reportByCategory", h.Location.ReportByCategory)
		}
	}

	return router
}
