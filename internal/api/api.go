package api

import (
	"time"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/ai"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/middleware"
	"github.com/kitchens-io/kitchens-api/internal/search"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

const (
	ENV_DEV  = "development"
	ENV_PROD = "production"
	ENV_TEST = "test"
)

var adminUserIDs = []string{
	"auth0|665e3646139d9f6300bad5e9", // Test Service
	"auth0|670c621fa584f0ad733a5bb8", // kristi
}

type App struct {
	db           *sqlx.DB
	fileManager  *media.S3FileManager
	env          string
	cdnHost      string
	aiClient     *ai.AIClient
	searchClient *search.SearchClient
	API          *echo.Echo
}

type CreateInput struct {
	DB            *sqlx.DB
	FileManager   *media.S3FileManager
	AuthValidator *validator.Validator
	Env           string
	CDNHost       string
	AIClient      *ai.AIClient
	SearchClient  *search.SearchClient
}

// Create will establish an instance of the app with all routes
// and middleware attached.
func Create(input CreateInput) *App {
	app := &App{
		db:           input.DB,
		fileManager:  input.FileManager,
		env:          input.Env,
		cdnHost:      input.CDNHost,
		aiClient:     input.AIClient,
		searchClient: input.SearchClient,
		API:          echo.New(),
	}

	authorizer := middleware.NewAuthorizer(input.AuthValidator, input.Env)
	kitchenAuth := middleware.NewKitchenAuthorizer(input.DB)
	adminAuth := middleware.NewAdminAuthorizer(adminUserIDs)

	// Disable the Echo banners on app start.
	app.API.HideBanner = true
	app.API.HidePort = true

	// Attach middelware and routes to the Echo instance.
	app.API.Use(mw.CORS())
	app.API.Use(mw.Logger())
	app.API.Use(mw.RequestID())
	app.API.Use(mw.BodyLimitWithConfig(mw.BodyLimitConfig{
		Limit: "10M",
	}))
	app.API.Use(mw.TimeoutWithConfig(mw.TimeoutConfig{
		Timeout: 30 * time.Second,
		Skipper: func(c echo.Context) bool {
			// Skip timeout for grocery list generation route
			return c.Request().Method == "GET" && c.Path() == "/v1/account/:account_id/plan/:id/groceries"
		},
	}))

	// Health Handler
	app.API.GET("/health", app.GetHealth)
	app.API.GET("/cdn/*", app.CDN)

	// V1 API routes
	v1 := app.API.Group("/v1")
	v1.Use(authorizer.ValidateToken)

	// Account Routes
	v1.GET("/iam", app.GetIAM)
	v1.POST("/account", app.CreateAccount)
	v1.PATCH("/account", app.UpdateAccount)
	v1.DELETE("/account", app.DeleteAccount)

	// Kitchens
	v1.GET("/kitchen/:kitchen_id", app.GetKitchen)
	v1.PATCH("/kitchen/:kitchen_id", app.UpdateKitchen, kitchenAuth.ValidateWriter)

	// Kitchen Recipes
	mwRecipeWriter := []echo.MiddlewareFunc{
		kitchenAuth.ValidateWriter,
		kitchenAuth.ValidateRecipeWriter,
	}
	v1.GET("/recipes/random", app.GetRandomRecipes)
	v1.GET("/kitchen/:kitchen_id/recipes", app.GetKitchenRecipes)
	v1.GET("/kitchen/:kitchen_id/recipes/:recipe_id", app.GetKitchenRecipe)
	v1.DELETE("/kitchen/:kitchen_id/recipes/:recipe_id", app.DeleteKitchenRecipe, mwRecipeWriter...)
	v1.POST("/kitchen/:kitchen_id/recipes", app.CreateKitchenRecipe, kitchenAuth.ValidateWriter)
	v1.PUT("/kitchen/:kitchen_id/recipes/:recipe_id", app.UpdateKitchenRecipe, mwRecipeWriter...)

	// Kitchen Folders
	mwFolderWriter := []echo.MiddlewareFunc{
		kitchenAuth.ValidateWriter,
		kitchenAuth.ValidateFolderWriter,
	}
	v1.GET("/kitchen/:kitchen_id/folders", app.GetKitchenFolders)
	v1.GET("/kitchen/:kitchen_id/folders/:folder_id", app.GetKitchenFolder)
	v1.DELETE("/kitchen/:kitchen_id/folders/:folder_id", app.DeleteKitchenFolder, mwFolderWriter...)
	v1.POST("/kitchen/:kitchen_id/folders", app.CreateKitchenFolder, kitchenAuth.ValidateWriter)
	v1.PUT("/kitchen/:kitchen_id/folders/:folder_id", app.UpdateKitchenFolder, mwFolderWriter...)

	// Kitchen Recipe Reviews
	v1.GET("/kitchen/:kitchen_id/recipes/:recipe_id/reviews", app.GetRecipeReviews)
	v1.POST("/kitchen/:kitchen_id/recipes/:recipe_id/reviews", app.CreateRecipeReview)
	v1.PUT("/kitchen/:kitchen_id/recipes/:recipe_id/reviews/:review_id", app.UpdateRecipeReview)
	v1.DELETE("/kitchen/:kitchen_id/recipes/:recipe_id/reviews/:review_id", app.DeleteRecipeReview)
	v1.POST("/kitchen/:kitchen_id/recipes/:recipe_id/reviews/:review_id/like", app.LikeRecipeReview)
	v1.DELETE("/kitchen/:kitchen_id/recipes/:recipe_id/reviews/:review_id/like", app.UnlikeRecipeReview)

	// Kitchen Folder Recipes
	v1.POST("/kitchen/:kitchen_id/folders/:folder_id/recipes/add", app.CreateKitchenFolderRecipes, mwFolderWriter...)
	v1.POST("/kitchen/:kitchen_id/folders/:folder_id/recipes/delete", app.DeleteKitchenFolderRecipes, mwFolderWriter...)

	// Kitchen Followers
	v1.GET("/kitchen/:kitchen_id/followers", app.GetKitchenFollowers)
	v1.GET("/kitchen/:kitchen_id/followed", app.GetKitchensFollowing)
	v1.POST("/kitchen/:kitchen_id/follow", app.FollowKitchen, kitchenAuth.ValidateWriter)
	v1.POST("/kitchen/:kitchen_id/unfollow", app.UnfollowKitchen, kitchenAuth.ValidateWriter)

	// Kitchen Saved Recipes
	v1.POST("/kitchen/:kitchen_id/save-recipe", app.SaveRecipe, kitchenAuth.ValidateWriter)
	v1.POST("/kitchen/:kitchen_id/remove-recipe", app.RemoveSavedRecipe, kitchenAuth.ValidateWriter)

	// Search
	v1.GET("/kitchens/search", app.SearchKitchens)
	v1.GET("/recipes/search/web", app.WebSearch)
	v1.GET("/recipes/search", app.RecipeSearch)
	v1.GET("/recipes/search/filters", app.RecipeSearchFilters)

	// Recipe import routes
	v1.POST("/import/url", app.ImportURL)
	v1.POST("/import/image", app.ImportImage)

	// Uploads
	v1.POST("/upload", app.Upload)

	// Meal Planning & grocery list
	v1.POST("/account/:account_id/plan", app.CreatePlan)
	v1.GET("/account/:account_id/plan/:id/categories/order", app.GetCategoryOrder)
	v1.PUT("/account/:account_id/plan/:id/categories/order", app.UpdateCategoryOrder)
	v1.GET("/account/:account_id/plan/:id/recipes", app.GetFullRecipesByPlanID)
	v1.GET("/account/:account_id/plan/:id/groceries", app.GetGroceryListByPlanID)
	v1.POST("/account/:account_id/plan/:id/groceries", app.CreateGroceryListItem)
	v1.DELETE("/account/:account_id/plan/:id/recipe/:recipe_id", app.RemoveRecipeFromPlan)
	v1.DELETE("/account/:account_id/plan/:id/groceries/:item_id", app.DeleteGroceryListItem)
	v1.PUT("/account/:account_id/plan/:id/groceries/:item_id", app.UpdateGroceryListItem)
	v1.POST("/account/:account_id/plan/:id/groceries/:item_id/mark", app.UpdateGroceryListItemMark)
	v1.POST("/account/:account_id/plan/recipes/:id", app.AddRecipesToPlan)
	v1.GET("/account/:account_id/plan/:id", app.GetPlanByID)
	v1.GET("/account/:account_id/plan/:start_date/:end_date", app.GetPlanByAccountIDAndDateRange)
	v1.GET("/account/:account_id/plan", app.GetPlansByUserID)

	// Public
	public := app.API.Group("/public")
	public.GET("/v1/recipes/:recipe_id", app.GetKitchenRecipe)

	// Admin Routes
	admin := app.API.Group("/admin")
	admin.Use(authorizer.ValidateToken)
	admin.Use(adminAuth.ValidateAdmin)
	admin.GET("/accounts", app.AdminListAccounts)
	admin.POST("/recipes/metadata", app.AdminCreateRecipeMetadata)

	// Used only for syncing in development.
	if app.env == ENV_DEV {
		admin.POST("/kitchen/:kitchen_id/recipes", app.CreateKitchenRecipe)
		admin.POST("/kitchen/:kitchen_id/folders/:folder_id/recipes/add", app.CreateKitchenFolderRecipes)
		admin.POST("/kitchen/:kitchen_id/follow", app.FollowKitchen)
	}

	return app
}
