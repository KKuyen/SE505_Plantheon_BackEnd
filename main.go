package main

import (
	"log"
	"os"

	"plantheon-backend/common"
	"plantheon-backend/models/activities"
	"plantheon-backend/models/activity_keywords"
	"plantheon-backend/models/blog_tags"
	"plantheon-backend/models/blogs"
	"plantheon-backend/models/comments"
	"plantheon-backend/models/disease_activity_keywords"
	"plantheon-backend/models/diseases"
	"plantheon-backend/models/guide_stages"
	"plantheon-backend/models/noti"
	"plantheon-backend/models/plants"
	"plantheon-backend/models/post_likes"
	"plantheon-backend/models/posts"
	"plantheon-backend/models/scan_history"
	"plantheon-backend/models/sub_guide_stages"
	"plantheon-backend/models/users"

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
	err := db.AutoMigrate(
		&users.User{},
		&posts.Post{},
		&comments.Comments{},
		&comments.CommentLike{},
		&diseases.Disease{},
		&scan_history.ScanHistory{},
		&activities.Activity{},
		&activity_keywords.ActivityKeyword{},
		&disease_activity_keywords.DiseaseActivityKeyword{},
		&post_likes.PostLike{},
		&blogs.Blog{},
		&blog_tags.BlogTag{},
		&plants.Plant{},
		&guide_stages.GuideStage{},
		&sub_guide_stages.SubGuideStage{},
		&noti.Notification{},
	)
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
				"status":  "OK",
				"message": "Plantheon Backend API is running",
			})
		})

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", users.Register)
			auth.POST("/login", users.Login)
		}

		// Public routes (no auth required)
		publicRoutes := api.Group("/public")
		{
			publicRoutes.GET("/users/:userId", posts.GetPublicUserProfileWithPostsHandler)
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
			diseaseRoutes.GET("/all", diseases.GetAllDiseasesHandler)
			diseaseRoutes.GET("/count", diseases.GetDiseasesCountHandler)
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

		// Activity routes
		activityRoutes := api.Group("/activities")
		{
			// Public routes (anyone can view activities)
			activityRoutes.GET("", activities.GetActivities)
			activityRoutes.GET("/all", activities.GetAllActivitiesHandler)
			activityRoutes.GET("/count", activities.GetActivitiesCountHandler)
			activityRoutes.GET("/get-activites-by-month", activities.GetActivitiesCalendarByMonthHandler)
			activityRoutes.GET("/by-day", activities.GetActivitiesByDayHandler)
			activityRoutes.GET("/financial/monthly", activities.GetMonthlyFinancialSummaryHandler)
			activityRoutes.GET("/financial/annual", activities.GetAnnualFinancialSummaryHandler)
			activityRoutes.GET("/financial/multi-year", activities.GetMultiYearFinancialSummaryHandler)
			activityRoutes.GET("/:id", activities.GetActivity)
			activityRoutes.POST("", activities.CreateActivityHandler)
			activityRoutes.PUT("/:id", activities.UpdateActivityHandler)
			activityRoutes.DELETE("", activities.DeleteActivitiesHandler)
			activityRoutes.DELETE("/:id", activities.DeleteActivityHandler)
		}

		// Admin-only activity routes (require admin role)
		adminActivityRoutes := api.Group("/activities")
		adminActivityRoutes.Use(users.RequireAdmin())
		{

		}

		// Activity Keywords routes
		activityKeywordRoutes := api.Group("/activity-keywords")
		{
			// Public routes (anyone can view)
			activityKeywordRoutes.GET("", activity_keywords.GetAllActivityKeywordsHandler)
			activityKeywordRoutes.GET("/:id", activity_keywords.GetActivityKeywordByIDHandler)
			activityKeywordRoutes.GET("/search", activity_keywords.SearchActivityKeywordsHandler)
			activityKeywordRoutes.GET("/disease/:disease_id", disease_activity_keywords.GetActivityKeywordsForDiseaseHandler)
		}

		// Admin-only Activity Keywords routes
		adminActivityKeywordRoutes := api.Group("/activity-keywords")
		adminActivityKeywordRoutes.Use(users.RequireAdmin())
		{
			adminActivityKeywordRoutes.POST("", activity_keywords.CreateActivityKeywordHandler)
			adminActivityKeywordRoutes.PUT("/:id", activity_keywords.UpdateActivityKeywordHandler)
			adminActivityKeywordRoutes.DELETE("/:id", activity_keywords.DeleteActivityKeywordHandler)
		}

		// Disease-Activity Keywords relationship routes
		diseaseActivityKeywordRoutes := api.Group("/disease-activity-keywords")
		{
			// Public routes - anyone can view
			diseaseActivityKeywordRoutes.GET("/disease/:disease_id", disease_activity_keywords.GetActivityKeywordsForDiseaseHandler)
			diseaseActivityKeywordRoutes.GET("/keyword/:keyword_id", disease_activity_keywords.GetDiseasesForKeywordHandler)
		}

		// Admin-only Disease-Activity Keywords routes
		adminDiseaseActivityKeywordRoutes := api.Group("/disease-activity-keywords")
		adminDiseaseActivityKeywordRoutes.Use(users.RequireAdmin())
		{
			adminDiseaseActivityKeywordRoutes.POST("/disease/:disease_id", disease_activity_keywords.AddActivityKeywordToDiseaseHandler)
			adminDiseaseActivityKeywordRoutes.DELETE("/disease/:disease_id/keyword/:keyword_id", disease_activity_keywords.RemoveActivityKeywordFromDiseaseHandler)
			adminDiseaseActivityKeywordRoutes.PUT("/disease/:disease_id", disease_activity_keywords.SetActivityKeywordsForDiseaseHandler)
			adminDiseaseActivityKeywordRoutes.POST("/import-csv", disease_activity_keywords.ImportActivityKeywordsFromCSVHandler)
		}

		// Scan History routes (protected - users can manage their own scan history)
		scanHistoryRoutes := api.Group("/scan-history")
		scanHistoryRoutes.Use(users.AuthMiddleware())
		{
			scanHistoryRoutes.GET("", scan_history.GetScanHistoriesHandler)
			scanHistoryRoutes.GET("/:id", scan_history.GetScanHistoryByIDHandler)
			scanHistoryRoutes.POST("", scan_history.CreateScanHistoryHandler)
			scanHistoryRoutes.DELETE("", scan_history.DeleteAllScanHistoriesHandler)
			scanHistoryRoutes.DELETE("/:id", scan_history.DeleteScanHistoryByIDHandler)
		}

		postRoutes := api.Group("/posts")
		postRoutes.Use(users.AuthMiddleware())
		{
			postRoutes.GET("", posts.GetPostsHandler)
			postRoutes.GET("/my", posts.GetMyPostsHandler)
			postRoutes.GET("/:id", posts.GetPostByIDHandler)
			postRoutes.POST("", posts.CreatePostHandler)
			postRoutes.PUT("/:id", posts.UpdatePostHandler)
			postRoutes.DELETE("/:id", posts.DeletePostByIDHandler)
			postRoutes.PUT("/:id/like", posts.LikePostHandler)
			postRoutes.PUT("/:id/unlike", posts.UnlikePostHandler)
			postRoutes.PUT("/:id/share", posts.SharePostHandler)
			postRoutes.GET("/:id/comments", comments.GetCommentsByPostIDHandler)
			postRoutes.POST("/:id/comments", comments.AddCommentHandler)
			postRoutes.PUT("/:id/comments", comments.UpdateCommentHandler)
			postRoutes.DELETE("/comments/:commentId", comments.DeleteCommentHandler)
			postRoutes.PUT("/comments/:commentId/like", comments.LikeCommentHandler)
			postRoutes.PUT("/comments/:commentId/unlike", comments.UnlikeCommentHandler)
			postRoutes.GET("/user/:userId", posts.GetPostsByUserIDHandler)
		}
		

		// Guide Stage routes
		guideStageRoutes := api.Group("/guide-stages")
		{

			// Public routes
			guideStageRoutes.GET("/:id", guide_stages.GetGuideStageByIDHandler)
			guideStageRoutes.GET("/plant/:plant_id", guide_stages.GetGuideStagesByPlantIDHandler)
		}

		// Admin-only Guide Stage routes
		adminGuideStageRoutes := api.Group("/guide-stages")
		adminGuideStageRoutes.Use(users.RequireAdmin())
		{
			adminGuideStageRoutes.POST("", guide_stages.CreateGuideStageHandler)
			adminGuideStageRoutes.PUT("/:id", guide_stages.UpdateGuideStageHandler)
			adminGuideStageRoutes.DELETE("/:id", guide_stages.DeleteGuideStageHandler)
		}

		// Sub Guide Stage routes
		subGuideStageRoutes := api.Group("/sub-guide-stages")
		{
			// Public routes
			subGuideStageRoutes.GET("/:id", sub_guide_stages.GetSubGuideStageByIDHandler)
			subGuideStageRoutes.GET("/guide-stage/:guide_stage_id", sub_guide_stages.GetSubGuideStagesByGuideStageIDHandler)
		}

		// Admin-only Sub Guide Stage routes
		adminSubGuideStageRoutes := api.Group("/sub-guide-stages")
		adminSubGuideStageRoutes.Use(users.RequireAdmin())
		{
			adminSubGuideStageRoutes.POST("", sub_guide_stages.CreateSubGuideStageHandler)
			adminSubGuideStageRoutes.PUT("/:id", sub_guide_stages.UpdateSubGuideStageHandler)
			adminSubGuideStageRoutes.DELETE("/:id", sub_guide_stages.DeleteSubGuideStageHandler)
		}

		// Plant routes
		plantRoutes := api.Group("/plants")
		{
			// Public routes
			plantRoutes.GET("", plants.GetAllPlantsHandler)
			plantRoutes.GET("/:id", plants.GetPlantByIDHandler)
		}

		// Admin-only Plant routes
		adminPlantRoutes := api.Group("/plants")
		adminPlantRoutes.Use(users.RequireAdmin())
		{
			adminPlantRoutes.POST("", plants.CreatePlantHandler)
			adminPlantRoutes.PUT("/:id", plants.UpdatePlantHandler)
			adminPlantRoutes.DELETE("/:id", plants.DeletePlantHandler)
		}

		newsRoutes := api.Group("/news")
		newsRoutes.Use(users.AuthMiddleware())
		{
			newsRoutes.GET("", blogs.GetNewsHandler)
			newsRoutes.GET("/:id", blogs.GetNewsByIDHandler)
		}

		// Admin-only news routes (require admin role)
		adminNewsRoutes := api.Group("/news")
		adminNewsRoutes.Use(users.RequireAdmin())
		{
			adminNewsRoutes.POST("", blogs.CreateNewsHandler)
			adminNewsRoutes.PUT("/:id", blogs.UpdateNewsHandler)
			adminNewsRoutes.DELETE("/:id", blogs.DeleteNewsHandler)
		}
		// Public news tag routes
		api.GET("/news-tags", blog_tags.GetAllBlogTagsHandler)

		// news tag routes (admin only)
		blogTagRoutes := api.Group("/news-tags")
		blogTagRoutes.Use(users.RequireAdmin())
		{
			blogTagRoutes.POST("", blog_tags.CreateBlogTagHandler)
			blogTagRoutes.PUT("/:id", blog_tags.UpdateBlogTagHandler)
			blogTagRoutes.DELETE("/:id", blog_tags.DeleteBlogTagHandler)
		}

		notificationRoutes := api.Group("/notification")
		notificationRoutes.Use(users.AuthMiddleware())
		{
			notificationRoutes.POST("", noti.CreateNotificationHandler)
			notificationRoutes.GET("", noti.GetNotificationsHandler)
			notificationRoutes.GET("/:id/seen", noti.MarkNotificationAsSeenHandler)
			notificationRoutes.DELETE("/:id", noti.DeleteNotificationHandler)
			notificationRoutes.DELETE("", noti.DeleteAllNotificationsHandler)
		}
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
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
	log.Printf("  GET  /api/diseases/all - Xem tất cả bệnh (không pagination)")
	log.Printf("  GET  /api/diseases/count - Xem số lượng bệnh")
	log.Printf("  GET  /api/diseases/:id - Xem chi tiết bệnh")
	log.Printf("  GET  /api/diseases/class/:className - Xem bệnh theo class name")
	log.Printf("Disease routes (cần admin role):")
	log.Printf("  POST /api/diseases - Tạo bệnh mới")
	log.Printf("  POST /api/diseases/import-excel - Import nhiều bệnh từ Excel")
	log.Printf("  PUT  /api/diseases/:id - Cập nhật bệnh")
	log.Printf("  DELETE /api/diseases/:ClassName - Xóa bệnh")
	log.Printf("Activity routes (public):")
	log.Printf("  GET  /api/activities - Xem danh sách hoạt động (có pagination, search, filter)")
	log.Printf("  GET  /api/activities/all - Xem tất cả hoạt động (không pagination)")
	log.Printf("  GET  /api/activities/count - Xem số lượng hoạt động")
	log.Printf("  GET  /api/activities/:id - Xem chi tiết hoạt động")
	log.Printf("Activity routes (cần admin role):")
	log.Printf("  POST /api/activities - Tạo hoạt động mới")
	log.Printf("  PUT  /api/activities/:id - Cập nhật hoạt động")
	log.Printf("  DELETE /api/activities/:id - Xóa hoạt động")
	log.Printf("Scan History routes (cần token):")
	log.Printf("  GET  /api/scan-history - Lấy tất cả lịch sử quét")
	log.Printf("  GET  /api/scan-history/:id - Lấy chi tiết lịch sử quét theo ID")
	log.Printf("  POST /api/scan-history - Tạo lịch sử quét mới")
	log.Printf("  DELETE /api/scan-history - Xóa tất cả lịch sử quét")
	log.Printf("  DELETE /api/scan-history/:id - Xóa lịch sử quét theo ID")
	log.Printf("Admin routes (cần admin role):")
	log.Printf("  /api/admin/users/* - Quản lý người dùng (commented out)")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
