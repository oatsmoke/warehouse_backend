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
		Equipment:  NewEquipmentHandler(service.Equipment, service.Location),
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

	api := router.Group("/api", h.Auth.UserIdentity)
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

		category := api.Group("/category")
		{
			category.POST("/create", h.Category.Create)
			category.PUT("/update/:id", h.Category.Update)
			category.DELETE("/delete/:id", h.Category.Delete)
			category.PUT("/restore/:id", h.Category.Restore)
			category.POST("/getAll", h.Category.GetAll)
			category.GET("/getById/:id", h.Category.GetById)
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
			equipment.POST("/getByIds", h.Equipment.GetByIds)
			equipment.POST("/findBySerialNumber", h.Equipment.FindBySerialNumber)
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
