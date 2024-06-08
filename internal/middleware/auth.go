package middleware

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Authorizer struct {
	validator *validator.Validator
}

func NewAuthorizer(validator *validator.Validator) *Authorizer {
	return &Authorizer{
		validator: validator,
	}
}

func (a *Authorizer) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := jwtmiddleware.AuthHeaderTokenExtractor(c.Request())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		validClaims, err := a.validator.ValidateToken(c.Request().Context(), token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(errors.Wrap(err, "error validating token"))
		}

		claims := validClaims.(*validator.ValidatedClaims)

		c.Set(auth.ClaimsContextKey, validClaims)
		c.Set(auth.UserIDContextKey, claims.RegisteredClaims.Subject)

		return next(c)
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
