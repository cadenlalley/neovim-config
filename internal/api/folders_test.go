package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/kitchens-io/kitchens-api/internal/fixtures"
	"github.com/kitchens-io/kitchens-api/pkg/folders"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestKitchenFoldersAuthorization(t *testing.T) {
	testKitchenID := "ktc_2jEx1j3CVPIIAaOwGIORKqHfK89"
	testFolderID := "fld_test"

	testCases := []struct {
		name     string
		methods  []string
		endpoint string
	}{
		{
			name:     "prevents folder creation if not kitchen owner",
			methods:  []string{"POST"},
			endpoint: "/v1/kitchen/" + testKitchenID + "/folders",
		},
		{
			name:     "prevents folder modification if not kitchen owner",
			methods:  []string{"PUT", "DELETE"},
			endpoint: "/v1/kitchen/" + testKitchenID + "/folders/" + testFolderID,
		},
		{
			name:     "prevents folder recipe addition if not kitchen owner",
			methods:  []string{"POST"},
			endpoint: "/v1/kitchen/" + testKitchenID + "/folders/" + testFolderID + "/recipes/add",
		},
		{
			name:     "prevents folder recipe deletion if not kitchen owner",
			methods:  []string{"POST"},
			endpoint: "/v1/kitchen/" + testKitchenID + "/folders/" + testFolderID + "/recipes/delete",
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

func TestCreateKitchenFolder(t *testing.T) {
	tests := []struct {
		name                  string
		kitchenID             string
		formData              map[string]string
		expectedCode          int
		expectedResponse      folders.Folder
		expectedErrorResponse string
	}{
		{
			name:      "successfully create kitchen folder",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			formData: map[string]string{
				"name": "Test Folder",
			},
			expectedCode: http.StatusOK,
			expectedResponse: folders.Folder{
				FolderID:  "",
				KitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
				Name:      "Test Folder",
				Cover:     null.String{},
			},
		},
		{
			name:                  "kitchen not found",
			kitchenID:             "ktc_test",
			formData:              map[string]string{},
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"kitchen not found"}`,
		},
		{
			name:                  "missing name",
			kitchenID:             "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			formData:              map[string]string{},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"missing required field 'createkitchenfolderrequest.name'"}`,
		},
		{
			name:      "duplicate name",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			formData: map[string]string{
				"name": "Test Folder",
			},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"duplicate entry for folder"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, err := formRequest(http.MethodPost, "/v1/kitchen/"+tt.kitchenID+"/folders", tt.formData)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, status)

			// Check errors
			if tt.expectedErrorResponse != "" {
				assert.JSONEq(t, tt.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual folders.Folder
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// Handle non-deterministic fields
			actual.FolderID = ""
			actual.CreatedAt = time.Time{}

			assert.Equal(t, tt.expectedResponse, actual)
		})
	}
}

func TestUpdateKitchenFolder(t *testing.T) {
	tests := []struct {
		name                  string
		kitchenID             string
		folderID              string
		formData              map[string]string
		expectedCode          int
		expectedResponse      folders.Folder
		expectedErrorResponse string
	}{
		{
			name:      "successfully update kitchen folder",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:  "fld_2pPgQjn08dQzr5vjSk8WYSBTATo",
			formData: map[string]string{
				"name": "Breaky",
			},
			expectedCode: http.StatusOK,
			expectedResponse: folders.Folder{
				FolderID:  "fld_2pPgQjn08dQzr5vjSk8WYSBTATo",
				KitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
				Name:      "Breaky",
				Cover:     null.StringFrom("/uploads/folders/fld_2pPgQjn08dQzr5vjSk8WYSBTATo/2pR6BliLuDCHYJKq7Eqkb9l55bS.png"),
			},
		},
		{
			name:                  "folder not found",
			kitchenID:             "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:              "fld_test",
			formData:              map[string]string{},
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"folder not found"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, err := formRequest(http.MethodPut, "/v1/kitchen/"+tt.kitchenID+"/folders/"+tt.folderID, tt.formData)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, status)

			// Check errors
			if tt.expectedErrorResponse != "" {
				assert.JSONEq(t, tt.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual folders.Folder
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// Handle non-deterministic fields
			actual.CreatedAt = time.Time{}

			assert.Equal(t, tt.expectedResponse, actual)
		})
	}
}

func TestGetKitchenFolder(t *testing.T) {
	testFolder := fixtures.GetTestFolder()

	tests := []struct {
		name                  string
		kitchenID             string
		folderID              string
		expectedCode          int
		expectedResponse      folders.Folder
		expectedErrorResponse string
	}{
		{
			name:             "successfully get kitchen folder",
			kitchenID:        "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:         "fld_2pPgQjn08dQzr5vjSk8WYSBTATo",
			expectedCode:     http.StatusOK,
			expectedResponse: testFolder,
		},
		{
			name:                  "folder not found",
			kitchenID:             "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:              "fld_test",
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"folder not found"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, err := request(http.MethodGet, "/v1/kitchen/"+tt.kitchenID+"/folders/"+tt.folderID, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, status)

			// Check errors
			if tt.expectedErrorResponse != "" {
				assert.JSONEq(t, tt.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual folders.Folder
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// Handle non-deterministic fields
			actual.CreatedAt = time.Time{}
			actual.Recipes[0].CreatedAt = time.Time{}

			assert.Equal(t, tt.expectedResponse, actual)
		})
	}
}

func TestGetKitchenFolders(t *testing.T) {
	testFolders := fixtures.GetTestFolders()

	tests := []struct {
		name                  string
		kitchenID             string
		expectedCode          int
		expectedResponse      []folders.Folder
		expectedErrorResponse string
	}{
		{
			name:             "successfully get kitchen folders",
			kitchenID:        "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			expectedCode:     http.StatusOK,
			expectedResponse: testFolders,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, err := request(http.MethodGet, "/v1/kitchen/"+tt.kitchenID+"/folders", nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, status)

			// Check errors
			if tt.expectedErrorResponse != "" {
				assert.JSONEq(t, tt.expectedErrorResponse, body.String())
				return
			}

			// Actual
			var actual []folders.Folder
			err = json.Unmarshal(body.Bytes(), &actual)
			assert.NoError(t, err)

			// Handle non-deterministic fields
			for i := range actual {
				actual[i].CreatedAt = time.Time{}
			}

			assert.Equal(t, tt.expectedResponse, actual)
		})
	}
}

func TestDeleteKitchenFolder(t *testing.T) {
	tests := []struct {
		name                  string
		kitchenID             string
		folderID              string
		expectedCode          int
		expectedErrorResponse string
	}{
		{
			name:         "successfully delete kitchen folder",
			kitchenID:    "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:     "fld_2pPgQjn08dQzr5vjSk8WYSBTATo",
			expectedCode: http.StatusNoContent,
		},
		{
			name:                  "folder not found",
			kitchenID:             "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:              "fld_test",
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"folder not found"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, err := request(http.MethodDelete, "/v1/kitchen/"+tt.kitchenID+"/folders/"+tt.folderID, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, status)

			// Check errors
			if tt.expectedErrorResponse != "" {
				assert.JSONEq(t, tt.expectedErrorResponse, body.String())
				return
			}
		})
	}
}

func TestCreateKitchenFolderRecipes(t *testing.T) {
	tests := []struct {
		name                  string
		kitchenID             string
		folderID              string
		payload               CreateKitchenFolderRecipeRequest
		expectedCode          int
		expectedErrorResponse string
	}{
		{
			name:      "successfully create kitchen folder recipes",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:  "fld_2pPgQjn08dQzr5vjSk8WYSBTATo",
			payload: CreateKitchenFolderRecipeRequest{
				RecipeIDs: []string{"rcp_2oSUH6fs0iCWGNP1AF2XemKYClo"},
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:      "folder not found",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:  "fld_test",
			payload: CreateKitchenFolderRecipeRequest{
				RecipeIDs: []string{"rcp_2oSUH6fs0iCWGNP1AF2XemKYClo"},
			},
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"folder not found"}`,
		},
		{
			name:      "missing recipe IDs",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:  "fld_2pPgQjn08dQzr5vjSk8WYSBTATo",
			payload: CreateKitchenFolderRecipeRequest{
				RecipeIDs: []string{},
			},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"invalid value '[]' supplied for field 'recipeIds'"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, err := request(http.MethodPost, "/v1/kitchen/"+tt.kitchenID+"/folders/"+tt.folderID+"/recipes/add", tt.payload)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, status)

			// Check errors
			if tt.expectedErrorResponse != "" {
				assert.JSONEq(t, tt.expectedErrorResponse, body.String())
				return
			}
		})
	}
}

func TestDeleteKitchenFolderRecipes(t *testing.T) {
	tests := []struct {
		name                  string
		kitchenID             string
		folderID              string
		payload               DeleteKitchenFolderRecipesRequest
		expectedCode          int
		expectedErrorResponse string
	}{
		{
			name:      "successfully delete kitchen folder recipes",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:  "fld_2pPgQjn08dQzr5vjSk8WYSBTATo",
			payload: DeleteKitchenFolderRecipesRequest{
				RecipeIDs: []string{"rcp_2oSUH6fs0iCWGNP1AF2XemKYClo"},
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:      "folder not found",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:  "fld_test",
			payload: DeleteKitchenFolderRecipesRequest{
				RecipeIDs: []string{"rcp_2oSUH6fs0iCWGNP1AF2XemKYClo"},
			},
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: `{"message":"folder not found"}`,
		},
		{
			name:      "missing recipe IDs",
			kitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
			folderID:  "fld_2pPgQjn08dQzr5vjSk8WYSBTATo",
			payload: DeleteKitchenFolderRecipesRequest{
				RecipeIDs: []string{},
			},
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: `{"message":"invalid value '[]' supplied for field 'recipeIds'"}`,
		},
	}

	// Reset Database before testing.
	err := resetFixtures()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, err := request(http.MethodPost, "/v1/kitchen/"+tt.kitchenID+"/folders/"+tt.folderID+"/recipes/delete", tt.payload)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, status)

			// Check errors
			if tt.expectedErrorResponse != "" {
				assert.JSONEq(t, tt.expectedErrorResponse, body.String())
				return
			}
		})
	}
}
