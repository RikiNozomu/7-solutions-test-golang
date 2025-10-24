package handler

import (
	domain "7-solutions-test-golang/core/domains"
	"7-solutions-test-golang/core/service"
	middleware "7-solutions-test-golang/middlewares"
	util "7-solutions-test-golang/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

// Create handles user creation
func (h *UserHandler) Create(c *gin.Context) {
	var user domain.UserCreate
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
		return
	}

	data, err := h.service.Create(user)
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": data})
}

// Get handles fetching a single user
func (h *UserHandler) Get(c *gin.Context) {
	data, err := h.service.Get(c.Param("id"))
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// List handles fetching all users
func (h *UserHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll()
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Update handles user updates
func (h *UserHandler) Update(c *gin.Context) {
	var user domain.UserUpdate
	if !util.IsValid(c) {
		c.Error(gin.Error{
			Err:  util.ErrorAuthenticated,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
		return
	}

	data, err := h.service.Update(c.Param("id"), user)
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Delete handles user deletion
func (h *UserHandler) Delete(c *gin.Context) {
	if !util.IsValid(c) {
		c.Error(gin.Error{
			Err:  util.ErrorAuthenticated,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	err := h.service.Delete(c.Param("id"))
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User has been removed."})
}

// RegisterRoutes registers all user routes
func (h *UserHandler) UserRoutes(router *gin.Engine) {
	// Create a group for user routes with barrier header middleware
	userGroup := router.Group("/user")
	{
		userGroup.POST("", h.Create)
		userGroup.GET("/:id", h.Get)
		userGroup.GET("", h.GetAll)
		userGroup.PUT("/:id", middleware.CheckBarrierHeader(), h.Update)
		userGroup.DELETE("/:id", middleware.CheckBarrierHeader(), h.Delete)
	}
}
