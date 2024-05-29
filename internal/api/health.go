package api

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
)

// Get tests the health of the application.
func (a *App) GetHealth(c echo.Context) error {
	// if err := a.DB.Ping(); err != nil {
	// 	return err
	// }
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

// TODO: Temporary route for checking claims on a JWT.
func (a *App) GetJWT(c echo.Context) error {
	claims := c.Request().Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims) // .CustomClaims.(*CustomClaims)
	return c.JSON(http.StatusOK, claims)
}
