package api

import (
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/middleware"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

type App struct {
	db          *sqlx.DB
	fileManager *media.S3FileManager
	API         *echo.Echo
}

type CreateInput struct {
	DB            *sqlx.DB
	FileManager   *media.S3FileManager
	AuthValidator *validator.Validator
}

// Create will establish an instance of the app with all routes
// and middleware attached.
func Create(input CreateInput) *App {
	app := &App{
		db:          input.DB,
		fileManager: input.FileManager,
		API:         echo.New(),
	}

	authorizer := middleware.NewAuthorizer(input.AuthValidator)

	// Disable the Echo banners on app start.
	app.API.HideBanner = true
	app.API.HidePort = true

	// Attach middelware and routes to the Echo instance.
	app.API.Use(mw.Logger())
	app.API.Use(mw.RequestID())

	// Health Handler
	app.API.GET("/health", app.GetHealth)

	// V1 API routes
	v1 := app.API.Group("/v1")
	v1.Use(authorizer.ValidateToken)

	// Account Routes
	v1.GET("/iam", app.GetIAM)
	v1.POST("/account", app.CreateAccount)
	v1.PATCH("/account", app.UpdateAccount)

	// Kitchens
	v1.GET("/kitchen/:kitchen_id", app.GetKitchen)
	v1.PATCH("/kitchen/:kitchen_id", app.UpdateKitchen)

	// Kitchen Recipes
	v1.GET("/kitchen/:kitchen_id/recipes", app.GetKitchenRecipes)
	v1.GET("/kitchen/:kitchen_id/recipes/:recipe_id", app.GetKitchenRecipe)

	// Uploads
	v1.POST("/upload", app.Upload)

	return app
}
