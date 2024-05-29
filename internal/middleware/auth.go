package middleware

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Authorizer struct {
	validator *validator.Validator
}

func NewAuthorizer(validator *validator.Validator) *Authorizer {
	return &Authorizer{
		validator: validator,
	}
}

// ValidateToken is a middleware that will check the validity of our JWT.
func (a *Authorizer) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Err(err).Msg("encountered error while validating JWT")
	}

	middleware := jwtmiddleware.New(
		a.validator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(ctx echo.Context) error {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.SetRequest(r)
			next(ctx)
		}

		middleware.CheckJWT(handler).ServeHTTP(ctx.Response(), ctx.Request())

		if encounteredError {
			ctx.JSON(http.StatusUnauthorized, map[string]string{
				"message": "JWT is invalid",
			})
		}
		return nil
	}
}

// HasScope checks whether our claims have a specific scope.
// func (c CustomClaims) HasScope(expectedScope string) bool {
// 	result := strings.Split(c.Scope, " ")
// 	for i := range result {
// 		if result[i] == expectedScope {
// 			return true
// 		}
// 	}

// 	return false
// }
