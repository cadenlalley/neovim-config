package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestSearchRecipes(t *testing.T) {
	testCases := []struct {
		name                  string
		query                 string
		expectedCode          int
		expectedResponse      []recipes.SearchResult
		expectedErrorResponse string
	}{
		{
			name:         "successfully search recipes",
			query:        "pumpkin pie",
			expectedCode: http.StatusOK,
			expectedResponse: []recipes.SearchResult{
				{
					RecipeID:     "rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T",
					KitchenID:    "ktc_2jEx1e1esA5292rBisRGuJwXc14",
					Name:         "Homemade pumpkin pie",
					Cover:        null.StringFrom("uploads/recipes/rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T/2pR9B2cIFxj82GDTDB44lpMzYHu.png"),
					ReviewCount:  4,
					ReviewRating: 3.5,
				},
			},
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := url.QueryEscape(tc.query)

			status, body, err := request(http.MethodGet, "/v1/recipes/search?q="+query, nil)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, status)

			// Check errors
			if tc.expectedErrorResponse != "" {
				assert.JSONEq(t, tc.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual []recipes.SearchResult
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedResponse, actual)
		})
	}
}
