package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/kitchens-io/kitchens-api/internal/fixtures"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/stretchr/testify/assert"
)

func TestGetKitchen(t *testing.T) {
	// Setup expected output
	testKitchen := fixtures.GetTestKitchen()

	testCases := []struct {
		name                  string
		kitchenID             string
		expectedCode          int
		expectedResponse      kitchens.Kitchen
		expectedErrorResponse string
	}{
		{
			name:             "successfully get kitchen",
			kitchenID:        "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			expectedCode:     http.StatusOK,
			expectedResponse: testKitchen,
		},
		{
			name:                  "kitchen not found",
			kitchenID:             "ktc_test",
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"kitchen not found"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, res := getRequest("/v1/kitchen/" + tc.kitchenID)

			// Assert response
			assert.Equal(t, tc.expectedCode, status)

			// Check errors
			if tc.expectedErrorResponse != "" {
				assert.JSONEq(t, tc.expectedErrorResponse, res.String())
				return
			}

			// Actual
			var actual kitchens.Kitchen
			err := json.Unmarshal(res.Bytes(), &actual)
			assert.NoError(t, err)

			// Hande time
			actual.CreatedAt = time.Time{}
			actual.UpdatedAt = time.Time{}

			assert.Equal(t, tc.expectedResponse, actual)
		})
	}
}

func TestUpdateKitchen(t *testing.T) {
	// Setup expected output
	testKitchen := fixtures.GetTestKitchen()
	testKitchen.Name = "Slammin Sammies"

	testCases := []struct {
		name                  string
		kitchenID             string
		formData              map[string]string
		expectedCode          int
		expectedResponse      kitchens.Kitchen
		expectedErrorResponse string
	}{
		{
			name:      "successfully update kitchen",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			formData: map[string]string{
				"name": "Slammin Sammies",
			},
			expectedCode:     http.StatusOK,
			expectedResponse: testKitchen,
		},
		{
			name:                  "kitchen not found",
			kitchenID:             "ktc_test",
			formData:              map[string]string{},
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"kitchen not found"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, body, err := patchFormRequest("/v1/kitchen/"+tc.kitchenID, tc.formData)
			assert.NoError(t, err)

			// Assert response
			assert.Equal(t, tc.expectedCode, status)

			// Check errors
			if tc.expectedErrorResponse != "" {
				assert.JSONEq(t, tc.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual kitchens.Kitchen
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// Hande time
			actual.CreatedAt = time.Time{}
			actual.UpdatedAt = time.Time{}

			assert.Equal(t, tc.expectedResponse, actual)
		})
	}
}

func TestSearchKitchens(t *testing.T) {
	testKitchen := fixtures.GetTestKitchen()

	testCases := []struct {
		name                  string
		query                 string
		expectedCode          int
		expectedResponse      []kitchens.Kitchen
		expectedErrorResponse string
	}{
		{
			name:         "search kitchens by handle",
			query:        "sammycooks",
			expectedCode: http.StatusOK,
			expectedResponse: []kitchens.Kitchen{
				testKitchen,
			},
		},
		{
			name:         "search kitchens by owner name",
			query:        "smith",
			expectedCode: http.StatusOK,
			expectedResponse: []kitchens.Kitchen{
				testKitchen,
			},
		},
		{
			name:         "search kitchens by name",
			query:        "sam's kitchen",
			expectedCode: http.StatusOK,
			expectedResponse: []kitchens.Kitchen{
				testKitchen,
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
			status, body := getRequest("/v1/kitchens/search?q=" + url.QueryEscape(tc.query))

			// Assert response
			assert.Equal(t, tc.expectedCode, status)

			// Check errors
			if tc.expectedErrorResponse != "" {
				assert.JSONEq(t, tc.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual []kitchens.Kitchen
			err := json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// Hande time
			for i := range actual {
				actual[i].CreatedAt = time.Time{}
				actual[i].UpdatedAt = time.Time{}
			}

			assert.Equal(t, tc.expectedResponse, actual)
		})
	}
}
