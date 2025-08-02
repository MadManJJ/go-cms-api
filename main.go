package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MadManJJ/cms-api/config"
	_ "github.com/MadManJJ/cms-api/docs"
	appHandler "github.com/MadManJJ/cms-api/handlers/app"
	cmsHandler "github.com/MadManJJ/cms-api/handlers/cms"
	commonHandler "github.com/MadManJJ/cms-api/handlers/common"
	"github.com/MadManJJ/cms-api/middleware"
	"github.com/MadManJJ/cms-api/repositories"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Swagger documentation
// @title Swagger Example API
// @description This API supports both Web and admin (CMS) operations, including authentication, page management.
// @version 1.0
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	// Load configuration
	cfg := config.New()

	// convert port number from string to int
	Port, err := strconv.Atoi(cfg.Database.Port)
	if err != nil {
		panic("failed to convert port number")
	}

	// Database connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DatabaseName)

	// New logger for detailed SQL logging
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold: time.Second,     // Slow SQL threshold
			LogLevel:      gormlogger.Info, // Log level
			Colorful:      true,            // Enable color
		},
	)

	// Initialize GORM with PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.Server.AppName,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	// Set up static file serving
	app.Static(cfg.App.StaticFilePrefix, cfg.App.UploadPath, fiber.Static{
		Compress: true,
		// Browse: cfg.App.Environment == "development",
		MaxAge: 3600,
	})

	// Setup Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Use middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(middleware.HPP())
	// app.Use(limiter.New(limiter.Config{
	// 	Max:        200,              // Max 200 requests
	// 	Expiration: 10 * time.Minute, // Per 10 minutes
	// 	LimitReached: func(c *fiber.Ctx) error {
	// 		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
	// 			"message": "Too many requests",
	// 		})
	// 	},
	// }))
	if cfg.App.Environment == "production" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: cfg.App.FrontendURLS,
			AllowHeaders: "Content-Type, Authorization",
			AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		}))
	} else {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Content-Type, Authorization",
			AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		}))
	}

	// Initialize repositories
	appRepo := repositories.NewMockAppRepository()
	appLandingPageRepo := repositories.NewAppLandingPageRepository(db)
	appPartnerPageRepo := repositories.NewAppPartnerPageRepository(db)
	appFaqPageRepo := repositories.NewAppFaqPageRepository(db)
	cmsRepo := repositories.NewMockCMSRepository(db)
	cmsAuthRepo := repositories.NewCMSAuthCMSAuthRepository(db)
	cmsCategoryTypeRepo := repositories.NewCMSCategoryTypeRepository(db)
	cmsCategoryRepo := repositories.NewCMSCategoryRepository(db)
	cmsFaqPageRepo := repositories.NewCMSFaqPageRepository(db)
	cmsLandingPageRepo := repositories.NewCMSLandingPageRepository(db)
	cmsPartnerPageRepo := repositories.NewCMSPartnerPageRepository(db)
	emailCategoryRepo := repositories.NewEmailCategoryRepository(db)
	emailContentRepo := repositories.NewEmailContentRepository(db)
	mediaFileRepo := repositories.NewMediaFileRepository(db)
	formRepo := repositories.NewFormRepository(db)
	formSubmissionRepo := repositories.NewFormSubmissionRepository(db)

	// Initialize services
	appService := services.NewAppService(appRepo)
	appLandingPageService := services.NewAppLandingPageService(appLandingPageRepo)
	appPartnerPageService := services.NewAppPartnerPageService(appPartnerPageRepo)
	appFaqPageService := services.NewAppFaqPageService(appFaqPageRepo)
	cmsService := services.NewCMSService(cmsRepo)
	cmsAuthService := services.NewCMSAuthService(cmsAuthRepo)
	categoryService := services.NewCMSCategoryService(cmsCategoryRepo, cmsCategoryTypeRepo)
	cmsCategoryTypeService := services.NewCMSCategoryTypeService(cmsCategoryTypeRepo, cmsCategoryRepo, categoryService)
	cmsFaqPageService := services.NewCMSFaqPageService(cmsFaqPageRepo, cfg)
	emailCategoryService := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)
	emailContentService := services.NewEmailContentService(emailContentRepo, emailCategoryRepo)
	emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	mediaFileService := services.NewMediaFileService(cfg, mediaFileRepo)
	cmsLandingPageService := services.NewCMSLandingPageService(cmsLandingPageRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)
	cmsPartnerPageService := services.NewCMSPartnerPageService(cmsPartnerPageRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)
	cmsFormService := services.NewCMSFormService(db, formRepo, emailCategoryRepo, cfg)
	commonLineLoginService := services.NewLineLoginService(cfg, cmsAuthRepo)
	cmsFormSubmissionService := services.NewCMSFormSubmissionService(formSubmissionRepo, emailSendingService)

	// Initialize handlers
	healthHandler := commonHandler.NewHealthHandler()
	commonLineLoginHandler := commonHandler.NewLineLoginHandler(commonLineLoginService)
	testMiddlewareHanlder := commonHandler.NewTestMiddlewareHandler()
	appLandingPageHandler := appHandler.NewAppLandingPageHandler(appLandingPageService)
	appPartnerPageHandler := appHandler.NewAppPartnerPageHandler(appPartnerPageService)
	appFaqPageHandler := appHandler.NewAppFaqPageHandler(appFaqPageService)
	appHandler := appHandler.NewAppHandler(appService)
	cmsCategoryTypeHandler := cmsHandler.NewCMSCategoryTypeHandler(cmsCategoryTypeService)
	cmsCategoryHandler := cmsHandler.NewCMSCategoryHandler(categoryService)
	cmsFaqPageHandler := cmsHandler.NewCMSFaqPageHandler(cmsFaqPageService)
	cmsLandingPageHandler := cmsHandler.NewCMSLandingPageHandler(cmsLandingPageService)
	cmsPartnerPageHandler := cmsHandler.NewCMSPartnerPageHandler(cmsPartnerPageService)
	cmsAuthHandler := cmsHandler.NewAuthCMSHandler(cmsAuthService)
	cmsFormHandler := cmsHandler.NewCMSFormHandler(cmsFormService)
	cmsFormSubmissionHandler := cmsHandler.NewCMSFormSubmissionHandler(cmsFormSubmissionService)
	emailCategoryCMSHandler := cmsHandler.NewEmailCategoryHandler(emailCategoryService)
	emailContentCMSHandler := cmsHandler.NewEmailContentHandler(emailContentService)
	emailSendingHandler := commonHandler.NewEmailSendingHandler(emailSendingService)
	mediaFileCMSHandler := cmsHandler.NewMediaFileHandler(mediaFileService)
	cmsHandler := cmsHandler.NewCMSHandler(cmsService)

	// Setup routes directly in main.go
	// Health check endpoint
	app.Get("/health", healthHandler.HandleHealthCheck)

	// API route group
	apiGroup := app.Group("/api/v1")

	// App routes under v1
	appGroup := apiGroup.Group("/app")
	appGroup.Get("/test", appHandler.HandleTest)
	appGroup.Get("/additional", appHandler.HandleAdditional)

	appLandingGroup := appGroup.Group("/landingpages")
	appLandingGroup.Get("/:languageCode/by-alias", appLandingPageHandler.HandleGetLandingPageByUrlAlias)
	appLandingGroup.Get("/previews/:id", appLandingPageHandler.HandleGetLandingContentPreview)

	appPartnerGroup := appGroup.Group("/partnerpages")
	appPartnerGroup.Get("/:languageCode/by-alias", appPartnerPageHandler.HandleGetPartnerPageByAlias)
	appPartnerGroup.Get("/:languageCode/by-url", appPartnerPageHandler.HandleGetPartnerPageByUrl)
	appPartnerGroup.Get("/previews/:id", appPartnerPageHandler.HandleGetPartnerContentPreview)

	appFaqGroup := appGroup.Group("/faqpages")
	appFaqGroup.Get("/:languageCode/by-alias", appFaqPageHandler.HandleGetFaqPageByAlias)
	appFaqGroup.Get("/:languageCode/by-url", appFaqPageHandler.HandleGetFaqPageByUrl)
	appFaqGroup.Get("/previews/:id", appFaqPageHandler.HandleGetFaqContentPreview)

	// CMS routes under v1
	cmsGroup := apiGroup.Group("/cms")
	cmsGroup.Get("/test", cmsHandler.HandleTest)
	cmsGroup.Get("/additional", cmsHandler.HandleAdditional)
	cmsAuthGroup := cmsGroup.Group("/auth")
	cmsAuthGroup.Post("/register", cmsAuthHandler.HandleRegister)
	cmsAuthGroup.Post("/login", cmsAuthHandler.HandleLogin)

	cmsFaqPageGroup := cmsGroup.Group("/faqpages")
	cmsFaqPageGroup.Post("/", cmsFaqPageHandler.HandleCreateFaqPage)
	cmsFaqPageGroup.Get("/", cmsFaqPageHandler.HandleGetFaqPages)
	cmsFaqPageGroup.Get("/:pageId", cmsFaqPageHandler.HandleGetFaqPageById)
	cmsFaqPageGroup.Get("/:pageId/latestcontents/:languageCode", cmsFaqPageHandler.HandleGetLatestContentByFaqPageId)
	cmsFaqPageGroup.Post("/:revisionId/revisions", cmsFaqPageHandler.HandleRevertFaqContent)
	cmsFaqPageGroup.Put("/:contentId/contents", cmsFaqPageHandler.HandleUpdateFaqContent)
	cmsFaqPageGroup.Delete("/:pageId", cmsFaqPageHandler.HandleDeleteFaqPage)
	cmsFaqPageGroup.Get("/:pageId/contents/:languageCode", cmsFaqPageHandler.HandleGetContentByFaqPageId)
	cmsFaqPageGroup.Delete("/:pageId/contents/:languageCode", cmsFaqPageHandler.HandleDeleteFaqContentByPageId)
	cmsFaqPageGroup.Post("/duplicate/:pageId/pages", cmsFaqPageHandler.HandleDuplicateFaqPage)
	cmsFaqPageGroup.Post("/duplicate/:contentId/contents", cmsFaqPageHandler.HandleDuplicateFaqContentToAnotherLanguage)
	cmsFaqPageGroup.Get("/category/:categoryTypeCode/:pageId/:languageCode", cmsFaqPageHandler.HandleGetCategory)
	cmsFaqPageGroup.Get("/revisions/:languageCode/:pageId", cmsFaqPageHandler.HandleGetRevisions)
	cmsFaqPageGroup.Post("/previews/:pageId", cmsFaqPageHandler.HandlePreviewFaqContent)

	cmsLandingPageGroup := cmsGroup.Group("/landingpages")
	cmsLandingPageGroup.Post("/", cmsLandingPageHandler.HandleCreateLandingPage)
	cmsLandingPageGroup.Get("/", cmsLandingPageHandler.HandleGetLandingPages)
	cmsLandingPageGroup.Get("/:pageId", cmsLandingPageHandler.HandleGetLandingPageById)
	cmsLandingPageGroup.Get("/:pageId/latestcontents/:languageCode", cmsLandingPageHandler.HandleGetLatestContentByLandingPageId)
	cmsLandingPageGroup.Post("/:revisionId/revisions", cmsLandingPageHandler.HandleRevertLandingContent)
	cmsLandingPageGroup.Put("/:contentId/contents", cmsLandingPageHandler.HandleUpdateLandingContent)
	cmsLandingPageGroup.Delete("/:pageId", cmsLandingPageHandler.HandleDeleteLandingPage)
	cmsLandingPageGroup.Get("/:pageId/contents/:languageCode", cmsLandingPageHandler.HandleGetContentByLandingPageId)
	cmsLandingPageGroup.Delete("/:pageId/contents/:languageCode", cmsLandingPageHandler.HandleDeleteLandingContentByPageId)
	cmsLandingPageGroup.Post("/duplicate/:pageId/pages", cmsLandingPageHandler.HandleDuplicateLandingPage)
	cmsLandingPageGroup.Post("/duplicate/:contentId/contents", cmsLandingPageHandler.HandleDuplicateLandingContentToAnotherLanguage)
	cmsLandingPageGroup.Get("/category/:categoryTypeCode/:pageId/:languageCode", cmsLandingPageHandler.HandleGetCategory)
	cmsLandingPageGroup.Get("/revisions/:languageCode/:pageId", cmsLandingPageHandler.HandleGetRevisions)
	cmsLandingPageGroup.Post("/previews/:pageId", cmsLandingPageHandler.HandlePreviewLandingContent)

	cmsPartnerPageGroup := cmsGroup.Group("/partnerpages")
	cmsPartnerPageGroup.Post("/", cmsPartnerPageHandler.HandleCreatePartnerPage)
	cmsPartnerPageGroup.Get("/", cmsPartnerPageHandler.HandleGetPartnerPages)
	cmsPartnerPageGroup.Get("/:pageId", cmsPartnerPageHandler.HandleGetPartnerPageById)
	cmsPartnerPageGroup.Get("/:pageId/latestcontents/:languageCode", cmsPartnerPageHandler.HandleGetLatestContentByPartnerPageId)
	cmsPartnerPageGroup.Post("/:revisionId/revisions", cmsPartnerPageHandler.HandleRevertPartnerContent)
	cmsPartnerPageGroup.Put("/:contentId/contents", cmsPartnerPageHandler.HandleUpdatePartnerContent)
	cmsPartnerPageGroup.Delete("/:pageId", cmsPartnerPageHandler.HandleDeletePartnerPage)
	cmsPartnerPageGroup.Get("/:pageId/contents/:languageCode", cmsPartnerPageHandler.HandleGetContentByPartnerPageId)
	cmsPartnerPageGroup.Delete("/:pageId/contents/:languageCode", cmsPartnerPageHandler.HandleDeletePartnerContentByPageId)
	cmsPartnerPageGroup.Post("/duplicate/:pageId/pages", cmsPartnerPageHandler.HandleDuplicatePartnerPage)
	cmsPartnerPageGroup.Post("/duplicate/:contentId/contents", cmsPartnerPageHandler.HandleDuplicatePartnerContentToAnotherLanguage)
	cmsPartnerPageGroup.Get("/category/:categoryTypeCode/:pageId/:languageCode", cmsPartnerPageHandler.HandleGetCategory)
	cmsPartnerPageGroup.Get("/revisions/:languageCode/:pageId", cmsPartnerPageHandler.HandleGetRevisions)
	cmsPartnerPageGroup.Post("/previews/:pageId", cmsPartnerPageHandler.HandlePreviewPartnerContent)

	cmsCategoryTypesGroup := cmsGroup.Group("/category-types")
	cmsCategoryTypesGroup.Post("/", cmsCategoryTypeHandler.HandleCreateCategoryType)
	cmsCategoryTypesGroup.Get("/", cmsCategoryTypeHandler.HandleListCategoryTypes)
	cmsCategoryTypesGroup.Get("/:id", cmsCategoryTypeHandler.HandleGetCategoryType)
	cmsCategoryTypesGroup.Patch("/:id", cmsCategoryTypeHandler.HandleUpdateCategoryType)
	cmsCategoryTypesGroup.Delete("/:id", cmsCategoryTypeHandler.HandleDeleteCategoryType)
	cmsCategoryTypesGroup.Get("/:categoryTypeId/categories", cmsCategoryTypeHandler.HandleListCategoriesForType)

	categoriesGroup := cmsGroup.Group("/categories")
	categoriesGroup.Post("/", cmsCategoryHandler.HandleCreateCategory)
	categoriesGroup.Get("/", cmsCategoryHandler.HandleListAllCategories)
	categoriesGroup.Get("/:categoryUuid", cmsCategoryHandler.HandleGetCategoryByUUID) 
	categoriesGroup.Patch("/:categoryUuid", cmsCategoryHandler.HandleUpdateCategory)
	categoriesGroup.Delete("/:categoryUuid", cmsCategoryHandler.HandleDeleteCategory)

	emailCategoriesCMSGroup := cmsGroup.Group("/email-categories")
	emailCategoriesCMSGroup.Post("/", emailCategoryCMSHandler.HandleCreateEmailCategory)
	emailCategoriesCMSGroup.Get("/", emailCategoryCMSHandler.HandleListEmailCategories)
	emailCategoriesCMSGroup.Get("/:id", emailCategoryCMSHandler.HandleGetEmailCategory)
	emailCategoriesCMSGroup.Patch("/:id", emailCategoryCMSHandler.HandleUpdateEmailCategory)
	emailCategoriesCMSGroup.Delete("/:id", emailCategoryCMSHandler.HandleDeleteEmailCategory)

	emailContentsCMSGroup := cmsGroup.Group("/email-contents")
	emailContentsCMSGroup.Post("/", emailContentCMSHandler.HandleCreateEmailContent)
	emailContentsCMSGroup.Get("/", emailContentCMSHandler.HandleListEmailContents)
	emailContentsCMSGroup.Get("/category/:email_category_id/language/:language", emailContentCMSHandler.HandleGetEmailContentByCategoryAndLanguage)
	emailContentsCMSGroup.Get("", emailContentCMSHandler.HandleGetEmailContent)
	emailContentsCMSGroup.Patch("/:id", emailContentCMSHandler.HandleUpdateEmailContent)
	emailContentsCMSGroup.Delete("/:id", emailContentCMSHandler.HandleDeleteEmailContent)

	// Media files CMS routes
	mediaFilesCMSGroup := cmsGroup.Group("/media-files")
	mediaFilesCMSGroup.Post("/", mediaFileCMSHandler.HandleUploadMediaFile)    
	mediaFilesCMSGroup.Get("/", mediaFileCMSHandler.HandleListMediaFiles)      
	mediaFilesCMSGroup.Get("/:id", mediaFileCMSHandler.HandleGetMediaFileByID) 
	mediaFilesCMSGroup.Delete("/:id", mediaFileCMSHandler.HandleDeleteMediaFile)

	// EMAIL SENDING ROUTE (can be under /api/v1 or /api/v1/common etc.)
	emailSendingGroup := apiGroup.Group("/emails")
	emailSendingGroup.Post("/send", emailSendingHandler.HandleSendEmail)

	// Common routes
	commonGroup := apiGroup.Group("/common")
	commonGroup.Get("/login-link", commonLineLoginHandler.HandleGetLoginLink)
	commonGroup.Post("/authenticate", commonLineLoginHandler.HandleAuthenticate)
	commonGroup.Post("/refresh-token", commonLineLoginHandler.HandleRefreshToken)

	cmsFormBuilderGroup := cmsGroup.Group("/forms")
	cmsFormBuilderGroup.Post("/", cmsFormHandler.HandleCreateForm)
	cmsFormBuilderGroup.Get("/", cmsFormHandler.HandleListForms)
	cmsFormBuilderGroup.Get("/:formId", cmsFormHandler.HandleGetForm)
	cmsFormBuilderGroup.Put("/:formId", cmsFormHandler.HandleUpdateForm)
	cmsFormBuilderGroup.Delete("/:formId", cmsFormHandler.HandleDeleteForm)
	cmsFormBuilderGroup.Post("/:formId/submissions", cmsFormSubmissionHandler.HandleCreateFormSubmission)
	cmsFormBuilderGroup.Get("/:formId/submissions", cmsFormSubmissionHandler.HandleGetFormSubmissions)
	cmsFormBuilderGroup.Get("/submissions/:submissionId", cmsFormSubmissionHandler.HandleGetFormSubmission)

	appFormGroup := appGroup.Group("/forms")
	appFormGroup.Get("/:formId/structure", cmsFormHandler.HandleGetFormStructure)

	// Test routes for middleware
	testGroup := apiGroup.Group("/middleware")
	testGroup.Get("/test", middleware.CheckAnyTokenMiddleware(cfg.SecretKey.LineKey, cfg.SecretKey.NormalKey, cmsAuthRepo), testMiddlewareHanlder.HandleTestMiddleware)

	// Start the server
	log.Printf("Starting server on port %s in %s mode", cfg.Server.Port, cfg.App.Environment)
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
