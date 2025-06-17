package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/kitchens-io/kitchens-api/internal/fixtures"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/stretchr/testify/assert"
)

func TestSearchRecipes(t *testing.T) {
	testCases := []struct {
		name                  string
		parameters            map[string]string
		expectedCode          int
		expectedResponse      []recipes.SearchResult
		expectedErrorResponse string
	}{
		{
			name:             "successfully search recipes",
			parameters:       map[string]string{"q": "pumpkin pie"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"pie"}),
		},
		{
			name:             "successfully search recipes with no results",
			parameters:       map[string]string{"q": "invalid"},
			expectedCode:     http.StatusOK,
			expectedResponse: []recipes.SearchResult{},
		},
		// Filters
		// =================
		{
			name:             "successfully filter recipes by kitchen id",
			parameters:       map[string]string{"q": "classic", "kitchenId": "ktc_2jEx1eCS13KMS8udlPoK12e5KPW"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"bolognese"}),
		},
		{
			name:             "successfully filter recipes by course",
			parameters:       map[string]string{"q": "classic", "course": "dessert"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"cookie"}),
		},
		{
			name:             "successfully filter recipes by class",
			parameters:       map[string]string{"q": "classic", "class": "dessert"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"cookie"}),
		},
		{
			name:             "successfully filter recipes by cuisine",
			parameters:       map[string]string{"q": "homemade", "cuisine": "American"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"pie", "cookie"}),
		},
		{
			name:             "successfully filter recipes by difficulty",
			parameters:       map[string]string{"q": "classic", "difficulty": "1"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"cookie"}),
		},
		{
			name:             "successfully filter recipes by rating",
			parameters:       map[string]string{"q": "homemade", "rating": "3"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"pie"}),
		},
		{
			name:             "successfully filter recipes by time",
			parameters:       map[string]string{"q": "homemade", "course": "dessert", "time": "30"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"cookie"}),
		},
		// Pagination
		// =================
		{
			name:             "successfully paginate recipes with limit",
			parameters:       map[string]string{"q": "homemade", "limit": "1"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"pie"}),
		},
		{
			name:             "successfully paginate recipes with limit and offset",
			parameters:       map[string]string{"q": "homemade", "limit": "1", "offset": "1"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"cookie"}),
		},
		// Sorting
		// =================
		{
			name:             "successfully sort recipes by rating",
			parameters:       map[string]string{"q": "love", "sort": "top"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"pie", "bolognese", "cookie"}),
		},
		{
			name:             "successfully sort recipes by newest",
			parameters:       map[string]string{"q": "love", "sort": "new"},
			expectedCode:     http.StatusOK,
			expectedResponse: fixtures.GetSearchRecipes([]string{"cookie", "bolognese", "pie"}),
		},
		// Bad requests
		// =================
		{
			name:                  "responds bad request with missing query parameter 'q'",
			parameters:            map[string]string{},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"missing query parameter 'q'"}`,
		},
		{
			name:                  "responds bad request with invalid value for query parameter 'difficulty'",
			parameters:            map[string]string{"q": "classic", "difficulty": "invalid"},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"invalid value for query parameter 'difficulty'"}`,
		},
		{
			name:                  "responds bad request with invalid value for query parameter 'limit'",
			parameters:            map[string]string{"q": "classic", "limit": "invalid"},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"invalid value for query parameter 'limit'"}`,
		},
		{
			name:                  "responds bad request with invalid value for query parameter 'offset'",
			parameters:            map[string]string{"q": "classic", "offset": "invalid"},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"invalid value for query parameter 'offset'"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := make(url.Values)
			for k, v := range tc.parameters {
				params.Set(k, url.QueryEscape(v))
			}

			status, body, err := request(http.MethodGet, "/v1/recipes/search?"+params.Encode(), nil)
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

func TestRecipeSearchFilters(t *testing.T) {
	testCases := []struct {
		name                  string
		expectedCode          int
		expectedResponse      RecipeSearchFiltersResponse
		expectedErrorResponse string
	}{
		{
			name:         "successfully get recipe search filters",
			expectedCode: http.StatusOK,
			expectedResponse: RecipeSearchFiltersResponse{
				Courses:    stringMapKeys(recipes.ValidCourses),
				Classes:    stringMapKeys(recipes.ValidClasses),
				Difficulty: []int{1, 2, 3, 4, 5},
				Rating:     []int{1, 2, 3, 4, 5},
				Sort:       stringMapKeys(recipes.ValidSort),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, body, err := request(http.MethodGet, "/v1/recipes/search/filters", nil)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, status)

			// Actual
			var actual RecipeSearchFiltersResponse
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			assert.ElementsMatch(t, tc.expectedResponse.Courses, actual.Courses)
			assert.ElementsMatch(t, tc.expectedResponse.Classes, actual.Classes)
			assert.ElementsMatch(t, tc.expectedResponse.Difficulty, actual.Difficulty)
			assert.ElementsMatch(t, tc.expectedResponse.Rating, actual.Rating)
			assert.ElementsMatch(t, tc.expectedResponse.Sort, actual.Sort)
		})
	}
}
