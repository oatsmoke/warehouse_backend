package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/pkg/model"
)

func (h *Handler) createEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var employee model.Employee
	if err := c.BindJSON(&employee); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.Create(employee.Name, employee.Phone, employee.Email); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) getByIdEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var e model.Employee
	if err := c.BindJSON(&e); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	employee, err := h.service.Employee.GetById(e.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, employee)
}

func (h *Handler) getByDepartmentEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var request model.RequestEmployee
	if err := c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	employees, err := h.service.Employee.GetByDepartment(request.Ids, request.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, employees)
}

func (h *Handler) getAllEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	employees, err := h.service.Employee.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, employees)
}

func (h *Handler) getFreeEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	employees, err := h.service.Employee.GetFree()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, employees)
}

func (h *Handler) getAllButAuthEmployee(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	employees, err := h.service.Employee.GetAllButOne(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, employees)
}

func (h *Handler) getAllButOneEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var e model.Employee
	if err := c.BindJSON(&e); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	employees, err := h.service.Employee.GetAllButOne(e.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, employees)
}

func (h *Handler) addToDepartmentEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var employee model.Employee
	if err := c.BindJSON(&employee); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.AddToDepartment(employee.Id, employee.Department.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) removeFromDepartmentEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var request []int
	if err := c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.RemoveFromDepartment(request[0], request[1]); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) updateEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var employee model.Employee
	if err := c.BindJSON(&employee); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.Update(employee.Id, employee.Name, employee.Phone, employee.Email); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) deleteEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var employee model.Employee
	if err := c.BindJSON(&employee); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.Delete(employee.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) activateEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var employee model.Employee
	if err := c.BindJSON(&employee); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.Activate(employee.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) deactivateEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var employee model.Employee
	if err := c.BindJSON(&employee); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.Deactivate(employee.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) resetPasswordEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var employee model.Employee
	if err := c.BindJSON(&employee); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.ResetPassword(employee.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) changeRoleEmployee(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var employee model.Employee
	if err := c.BindJSON(&employee); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Employee.ChangeRole(employee.Id, employee.Role); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}
