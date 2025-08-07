package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/kitchens-io/kitchens-api/internal/fixtures"
	"github.com/kitchens-io/kitchens-api/pkg/ptr"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestKitchenRecipesAuthorization(t *testing.T) {
	testKitchenID := "ktc_2jEx1j3CVPIIAaOwGIORKqHfK89"
	testRecipeID := "rcp_2oSUH6fs0iCWGNP1AF2XemKYClo"

	testCases := []struct {
		name     string
		methods  []string
		endpoint string
	}{
		{
			name:     "prevents recipe creation if not kitchen owner",
			methods:  []string{"POST"},
			endpoint: "/v1/kitchen/" + testKitchenID + "/recipes",
		},
		{
			name:     "prevents recipe modification if not kitchen owner",
			methods:  []string{"PUT", "DELETE"},
			endpoint: "/v1/kitchen/" + testKitchenID + "/recipes/" + testRecipeID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, method := range tc.methods {
				status, body, err := request(method, tc.endpoint, nil)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusForbidden, status)
				assert.JSONEq(t, `{"message":"Forbidden"}`, body.String())
			}
		})
	}
}

func TestCreateKitchenRecipe(t *testing.T) {
	testPayload := recipes.Recipe{
		Name:       "Tacos",
		Summary:    null.StringFrom("Protein, tortilla, and toppings."),
		PrepTime:   ptr.Int(10),
		CookTime:   ptr.Int(15),
		Servings:   ptr.Int(4),
		Difficulty: 1,
		Course:     null.StringFrom("dinner"),
		Class:      null.StringFrom("main"),
		Cuisine:    null.StringFrom("Mexican"),
		Steps: []recipes.RecipeStep{
			{
				StepID:      1,
				Instruction: "Buy tacos from the store.",
			},
		},
		Ingredients: []recipes.RecipeIngredient{
			{
				IngredientID: 1,
				Name:         "Taco",
				Quantity:     null.FloatFrom(1),
				Unit:         null.StringFrom("piece"),
			},
		},
	}

	testCases := []struct {
		name                  string
		kitchenID             string
		payload               recipes.Recipe
		expectedCode          int
		expectedResponse      recipes.Recipe
		expectedErrorResponse string
	}{
		{
			name:             "successfully create recipe",
			kitchenID:        "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			payload:          testPayload,
			expectedCode:     http.StatusOK,
			expectedResponse: testPayload,
		},
		{
			name:                  "prevents recipe creation with bad data",
			kitchenID:             "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			payload:               recipes.Recipe{},
			expectedCode:          http.StatusBadRequest,
			expectedResponse:      recipes.Recipe{},
			expectedErrorResponse: `{"message":"missing required field 'recipe.name'"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, body, err := request(http.MethodPost, "/v1/kitchen/"+tc.kitchenID+"/recipes", tc.payload)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, status)

			// Check errors
			if tc.expectedErrorResponse != "" {
				assert.JSONEq(t, tc.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual recipes.Recipe
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// TODO: Should these come back as nil?
			tc.expectedResponse.Steps = nil
			tc.expectedResponse.Ingredients = nil

			// Handle non-deterministic fields
			tc.expectedResponse.KitchenID = tc.kitchenID
			tc.expectedResponse.RecipeID = actual.RecipeID

			actual.ShareURL = ""
			actual.CreatedAt = time.Time{}
			actual.UpdatedAt = time.Time{}
			actual.DeletedAt = null.Time{}

			assert.Equal(t, tc.expectedResponse, actual)
		})
	}
}

func TestUpdateRecipe(t *testing.T) {
	// Setup expected output
	testRecipe := fixtures.GetTestRecipe()
	testRecipe.Name = "Store bought pumpkin pie"
	testRecipe.ShareURL = "/2jbgfAMKOCnKrWQroRBkXPIRI6T/store-bought-pumpkin-pie"
	testRecipe.Difficulty = 1

	// These wont be set on the response.
	testRecipe.ReviewCount = 0
	testRecipe.ReviewRating = 0

	testCases := []struct {
		name                  string
		recipeID              string
		kitchenID             string
		payload               recipes.Recipe
		expectedCode          int
		expectedResponse      recipes.Recipe
		expectedErrorResponse string
	}{
		{
			name:             "successfully update recipe",
			recipeID:         testRecipe.RecipeID,
			kitchenID:        testRecipe.KitchenID,
			payload:          testRecipe,
			expectedCode:     http.StatusOK,
			expectedResponse: testRecipe,
		},
		{
			name:                  "prevents recipe update with bad data",
			recipeID:              testRecipe.RecipeID,
			kitchenID:             testRecipe.KitchenID,
			payload:               recipes.Recipe{},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"missing required field 'recipe.name'"}`,
		},
		{
			name:                  "recipe not found",
			recipeID:              "rcp_test",
			kitchenID:             testRecipe.KitchenID,
			payload:               recipes.Recipe{Name: "Cookie Wookies"},
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"recipe not found"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, body, err := request(http.MethodPut, "/v1/kitchen/"+tc.kitchenID+"/recipes/"+tc.recipeID, tc.payload)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, status)

			// Check errors
			if tc.expectedErrorResponse != "" {
				assert.JSONEq(t, tc.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual recipes.Recipe
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// TODO: Should these come back as nil?
			tc.expectedResponse.Steps = nil
			tc.expectedResponse.Ingredients = nil

			// Handle time
			actual.CreatedAt = time.Time{}
			actual.UpdatedAt = time.Time{}

			assert.Equal(t, tc.expectedResponse, actual)
		})
	}
}

func TestGetKitchenRecipe(t *testing.T) {
	// Setup expected output
	testRecipe := fixtures.GetTestRecipe()

	testCases := []struct {
		name                  string
		kitchenID             string
		recipeID              string
		expectedCode          int
		expectedResponse      recipes.Recipe
		expectedErrorResponse string
	}{
		{
			name:             "successfully get recipe",
			kitchenID:        testRecipe.KitchenID,
			recipeID:         testRecipe.RecipeID,
			expectedCode:     http.StatusOK,
			expectedResponse: testRecipe,
		},
		{
			name:                  "recipe not found",
			kitchenID:             testRecipe.KitchenID,
			recipeID:              "rcp_test",
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"recipe not found"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, body, err := request(http.MethodGet, "/v1/kitchen/"+tc.kitchenID+"/recipes/"+tc.recipeID, nil)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, status)

			// Check errors
			if tc.expectedErrorResponse != "" {
				assert.JSONEq(t, tc.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual recipes.Recipe
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// Hande time
			actual.CreatedAt = time.Time{}
			actual.UpdatedAt = time.Time{}
			actual.DeletedAt = null.Time{}

			assert.Equal(t, tc.expectedResponse, actual)
		})
	}
}

func TestDeleteRecipe(t *testing.T) {
	testRecipe := fixtures.GetTestRecipe()

	testCases := []struct {
		name                  string
		recipeID              string
		kitchenID             string
		expectedCode          int
		expectedErrorResponse string
	}{
		{
			name:         "successfully delete recipe",
			recipeID:     testRecipe.RecipeID,
			kitchenID:    testRecipe.KitchenID,
			expectedCode: http.StatusOK,
		},
		{
			name:                  "recipe not found",
			recipeID:              "rcp_test",
			kitchenID:             testRecipe.KitchenID,
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"recipe not found"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, body, err := request(http.MethodDelete, "/v1/kitchen/"+tc.kitchenID+"/recipes/"+tc.recipeID, nil)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, status)

			// Check errors
			if tc.expectedErrorResponse != "" {
				assert.JSONEq(t, tc.expectedErrorResponse, body.String())
				return
			}
		})
	}
}
