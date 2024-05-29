package api

import (
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/kitchens-io/kitchens-api/internal/middleware"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

type App struct {
	API *echo.Echo
}

type CreateInput struct {
	AuthValidator *validator.Validator
}

// Create will establish an instance of the app with all routes
// and middleware attached.
func Create(input CreateInput) *App {
	app := &App{
		API: echo.New(),
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

	// TODO: temporary route for checking claims on a JWT.
	v1.GET("/jwt", app.GetJWT)

	// Auth Handler
	// app.API.POST("/sign-in", app.PostSignIn)
	// app.API.POST("/sign-up", app.PostSignUp)

	// V1 API routes
	// v1 := app.API.Group("/v1")
	// v1.Use(
	// 	mw.JWT([]byte(app.JWTSignature)),
	// 	web.MiddlewareSetUserID(),
	// )

	// user := CreateUserResource(db)
	// v1.GET("/user", user.Get)
	// v1.PATCH("/user", user.Patch)

	return app
}
