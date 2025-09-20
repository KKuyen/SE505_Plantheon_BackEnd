package main

import (
	"log"
	"os"

	"plantheon-backend/common"
	"plantheon-backend/diseases"
	"plantheon-backend/users"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	db := common.Init()

	// Auto migrate database tables
	err := db.AutoMigrate(&users.User{}, &diseases.Disease{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Set up Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "OK",
				"message": "Plantheon Backend API is running",
			})
		})

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", users.Register)
			auth.POST("/login", users.Login)
		}

		// User routes (protected)
		userRoutes := api.Group("/users")
		userRoutes.Use(users.AuthMiddleware())
		{
			userRoutes.GET("/profile", users.GetProfile)
			userRoutes.PUT("/profile", users.UpdateProfile)
		}

		// Admin-only user management routes
		adminUserRoutes := api.Group("/admin/users")
		adminUserRoutes.Use(users.RequireAdmin())
		{
			// Add admin-only user management endpoints here
			// adminUserRoutes.GET("", users.GetAllUsers)
			// adminUserRoutes.GET("/:id", users.GetUserByIDHandler)  
			// adminUserRoutes.PUT("/:id/role", users.UpdateUserRole)
			// adminUserRoutes.DELETE("/:id", users.DeleteUserHandler)
		}

		// Disease routes
		diseaseRoutes := api.Group("/diseases")
		{
			// Public routes (anyone can view diseases)
			diseaseRoutes.GET("", diseases.GetDiseases)
			// diseaseRoutes.GET("/:id", diseases.GetDisease)
			diseaseRoutes.GET("/:ClassName", diseases.GetDiseaseByClassNameHandler)
		}

		// Admin-only disease routes (require admin role)
		adminDiseaseRoutes := api.Group("/diseases")
		adminDiseaseRoutes.Use(users.RequireAdmin())
		{
			adminDiseaseRoutes.POST("", diseases.CreateDiseaseHandler)
			adminDiseaseRoutes.POST("/import-excel", diseases.ImportDiseasesFromExcelHandler)
			adminDiseaseRoutes.PUT("/:id", diseases.UpdateDiseaseHandler)
			adminDiseaseRoutes.DELETE("/:ClassName", diseases.DeleteDiseaseHandler)
		}
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Health check: http://localhost:%s/api/health", port)
	log.Printf("API documentation:")
	log.Printf("Auth routes:")
	log.Printf("  POST /api/auth/register - Đăng ký tài khoản")
	log.Printf("  POST /api/auth/login - Đăng nhập")
	log.Printf("User routes (cần token):")
	log.Printf("  GET  /api/users/profile - Xem profile")
	log.Printf("  PUT  /api/users/profile - Cập nhật profile")
	log.Printf("Disease routes (public):")
	log.Printf("  GET  /api/diseases - Xem danh sách bệnh (có pagination, search, filter)")
	log.Printf("  GET  /api/diseases/:id - Xem chi tiết bệnh")
	log.Printf("  GET  /api/diseases/class/:className - Xem bệnh theo class name")
	log.Printf("Disease routes (cần admin role):")
	log.Printf("  POST /api/diseases - Tạo bệnh mới")
	log.Printf("  POST /api/diseases/import-excel - Import nhiều bệnh từ Excel")
	log.Printf("  PUT  /api/diseases/:id - Cập nhật bệnh")
	log.Printf("  DELETE /api/diseases/:ClassName - Xóa bệnh")
	log.Printf("Admin routes (cần admin role):")
	log.Printf("  /api/admin/users/* - Quản lý người dùng (commented out)")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}