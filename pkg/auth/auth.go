package auth

import (
	"context"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

const ClaimsContextKey = "auth0_claims"
const UserIDContextKey = "auth0_sub"

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	// Scope     string `json:"scope"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Email     string `json:"email"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func NewValidator(domain, audience string, jwkCacheTTL time.Duration) (*validator.Validator, error) {
	issuerURL, err := url.Parse("https://" + domain + "/")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the issuer url")
	}

	provider := jwks.NewCachingProvider(issuerURL, jwkCacheTTL)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up jwt validator")
	}

	return jwtValidator, nil
}
