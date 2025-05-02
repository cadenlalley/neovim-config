package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kitchens-io/kitchens-api/internal/fixtures"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/stretchr/testify/assert"
)

func TestGetIAM(t *testing.T) {
	testAccount := fixtures.GetTestAccount()
	testKitchen := fixtures.GetTestKitchen()

	testCases := []struct {
		name                  string
		expectedCode          int
		expectedResponse      GetIAMResponse
		expectedErrorResponse string
	}{
		{
			name:         "successfully get iam",
			expectedCode: http.StatusOK,
			expectedResponse: GetIAMResponse{
				Account: testAccount,
				Kitchens: []kitchens.Kitchen{
					testKitchen,
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
			status, body := getRequest("/v1/iam")

			// Assert response
			assert.Equal(t, http.StatusOK, status)

			// Actual
			var actual GetIAMResponse
			err := json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// Fix time causing failed tests.
			tc.expectedResponse.Account.CreatedAt = actual.Account.CreatedAt
			tc.expectedResponse.Kitchens[0].CreatedAt = actual.Kitchens[0].CreatedAt
			tc.expectedResponse.Kitchens[0].UpdatedAt = actual.Kitchens[0].UpdatedAt

			assert.Equal(t, tc.expectedResponse, actual)
		})
	}
}
