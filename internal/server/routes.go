package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"instashop/internal/controller"
	"instashop/internal/database"
	"instashop/internal/middleware"
	"instashop/internal/service"
)

// RegisterRoutes sets up all the routes for the application
func (s *Server) RegisterRoutes() http.Handler {

	dbService := database.New()

	// Pass GORM DB to the user service
	userService := service.NewUserService(dbService.GetGORM())
	userController := controller.NewUserController(userService)
	productService := service.NewProductService(dbService.GetGORM())
	productController := controller.NewProductController(productService)

	// Initialize router
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandlerMiddleware)

	// Basic routes
	r.GET("/", s.HelloWorldHandler)

	// API routes
	r.POST("/v1/auth/users/create", userController.CreateUser)
	r.POST("/v1/auth/admin/create", userController.CreateAdmin)

	// Verify email route
	r.POST("/v1/auth/verify-email", userController.VerifyEmailHandler)

	r.POST("/v1/auth/send-email", userController.SendEmailHandler)

	r.POST("/v1/auth/login", userController.LoginHandler)

	// Product routes
	authorized := r.Group("/v1")
	authorized.Use(middleware.VerifyToken())
	{
		authorized.POST("/products", productController.CreateProduct)
		authorized.GET("/products/:productID", productController.GetProduct)
		authorized.GET("/products", productController.GetAllProductsByUserID)
		authorized.PATCH("/products/:productID", productController.UpdateProduct)
		authorized.DELETE("/products/:productID", productController.DeletePendingProduct)
	}

	// Handle not found routes
	r.NoRoute(middleware.HandleNotFound)

	// Assign the DB close method to the Server's cleanup logic
	s.cleanupFunc = func() {
		s.db.Close()
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := map[string]string{
		"message": "Welcome to INSTASHOP Service",
	}
	c.JSON(http.StatusOK, resp)
}
