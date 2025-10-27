package handler

import (
	domain "7-solutions-test-golang/core/domains"
	service "7-solutions-test-golang/core/services"
	middleware "7-solutions-test-golang/middlewares"
	util "7-solutions-test-golang/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService // Business logic layer
}

// NewUserHandler initializes a new UserHandler with the given service.
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

// Create handles user creation via POST /user
func (h *UserHandler) Create(c *gin.Context) {
	var user domain.UserCreate

	// Bind and validate incoming JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
		return
	}

	// Create user via service
	data, err := h.service.Create(user)
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	// Respond with created user
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

// Get handles fetching a single user via GET /user/:id
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

// GetAll handles fetching all users via GET /user
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

// Update handles user updates via PUT /user/:id
func (h *UserHandler) Update(c *gin.Context) {
	var user domain.UserUpdate

	// Validate user identity from JWT claims
	if !util.IsValid(c) {
		c.Error(gin.Error{
			Err:  util.ErrorAuthenticated,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	// Bind and validate incoming JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
		return
	}

	// Update user via service
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

// Delete handles user deletion via DELETE /user/:id
func (h *UserHandler) Delete(c *gin.Context) {
	// Validate user identity from JWT claims
	if !util.IsValid(c) {
		c.Error(gin.Error{
			Err:  util.ErrorAuthenticated,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	// Delete user via service
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

// UserRoutes registers all user-related routes under /user
func (h *UserHandler) UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/user")
	userGroup.Use(middleware.RateLimiter())
	{
		userGroup.POST("", h.Create) // Public: create user
		userGroup.GET("/:id", h.Get) // Public: get user by ID
		userGroup.GET("", h.GetAll)  // Public: list all users

		// Protected: require JWT for update/delete
		userGroup.PUT("/:id", middleware.CheckBarrierHeader(), h.Update)
		userGroup.DELETE("/:id", middleware.CheckBarrierHeader(), h.Delete)
	}
}
