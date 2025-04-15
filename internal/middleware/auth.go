package middleware

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const (
	ENV_TEST = "test"
)

type Authorizer struct {
	validator *validator.Validator
	env       string
}

func NewAuthorizer(validator *validator.Validator, env string) *Authorizer {
	return &Authorizer{
		validator: validator,
		env:       env,
	}
}

func (a *Authorizer) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := jwtmiddleware.AuthHeaderTokenExtractor(c.Request())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var validClaims interface{}
		if a.env == ENV_TEST {
			validClaims = a.SetTestClaims()
		} else {
			validClaims, err = a.validator.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(errors.Wrap(err, "error validating token"))
			}
		}

		claims := validClaims.(*validator.ValidatedClaims)

		c.Set(auth.ClaimsContextKey, validClaims)
		c.Set(auth.UserIDContextKey, claims.RegisteredClaims.Subject)
		c.Set(auth.KitchenIDContextKey, c.Request().Header.Get("X-Kitchen-ID"))

		return next(c)
	}
}

func (a *Authorizer) SetTestClaims() interface{} {
	return &validator.ValidatedClaims{
		RegisteredClaims: validator.RegisteredClaims{
			Issuer:   "https://dev-b7ndo7iy1caliyo4.us.auth0.com/",
			Subject:  "auth0|665e3646139d9f6300bad5e9",
			Audience: []string{"Fqb3Ypo9LTR98NoC9BKYpNaOi2dCt18r"},
			ID:       "665e3646139d9f6300bad5e9",
		},
		CustomClaims: auth.CustomClaims{
			Email:         "test-service@kitchens-app.com",
			EmailVerified: true,
		},
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
