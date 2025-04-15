package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestGetIAM(t *testing.T) {
	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/v1/iam", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjNJS090Z0ZtWi02ejYzMjdFWWFYdiJ9.eyJuaWNrbmFtZSI6InRlc3Qtc2VydmljZSIsIm5hbWUiOiJ0ZXN0LXNlcnZpY2VAa2l0Y2hlbnMtYXBwLmNvbSIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yY2M3ODg2Y2Y1NTUwMDliZTNmNjE3ZjM5NDMyNjAxMj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRnRlLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDI1LTA0LTEyVDE0OjQ2OjI1LjY2MVoiLCJlbWFpbCI6InRlc3Qtc2VydmljZUBraXRjaGVucy1hcHAuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vZGV2LWI3bmRvN2l5MWNhbGl5bzQudXMuYXV0aDAuY29tLyIsImF1ZCI6IkZxYjNZcG85TFRSOThOb0M5QktZcE5hT2kyZEN0MThyIiwic3ViIjoiYXV0aDB8NjY1ZTM2NDYxMzlkOWY2MzAwYmFkNWU5IiwiaWF0IjoxNzQ0NDY5MTg3LCJleHAiOjE3NDUwNzM5ODcsInNpZCI6InlfSjkzWkFqQ2Z1UnhUYWZXX3NWcGRHbmplemd1MmpsIiwibm9uY2UiOiJYMGxIY2xsVmEyUXlkbVJEWlZkR01tZ3hOMFJLZGpkWE4wVm1hazkxT0dWUGMxUTNORE5FYkY5elNBPT0ifQ.n1G3wQt3SlzxP111TgZlASRN978vImYqzpCL3kDt32e3NCLMznT4y3juy6ujduz0iB42aCCy3ffQYEn2UyKWJMtzzNXKZOOUBNLQ1iIUMRqwVOxPZr1UgBjxrRYOYrecVn3MSCLzPcNPMLT84LK44b4I16yxJYQA8jLdAH6GM9V7cIHs00jjrLa95psUATSwlgnjm76diGHOO8JPzFUNwl28hfsyv4C_pTipfnZECgaxIIGQbaTcWEKEtWzzj-rFGZVX1Syt_VJA6Sk2Cwr0PWcvuQdLBh-vrR2TBPkPwvuly1ZQO_IfyAe7XkrzYzH5tUUZksQkeRl7Ajpd9hBm3Q")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the handler directly
	testApp.API.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Actual
	var actual GetIAMResponse
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(t, err)

	expected := GetIAMResponse{
		Account: accounts.Account{
			AccountID: "acc_2jEwcS7Rla6E5ik5ELa8uoULKOW",
			Email:     "test-service@kitchens-app.com",
			FirstName: "Sam",
			LastName:  "Smith",
			UserID:    "auth0|665e3646139d9f6300bad5e9",
			Verified:  false,
		},
		Kitchens: []kitchens.Kitchen{
			{
				AccountID: "acc_2jEwcS7Rla6E5ik5ELa8uoULKOW",
				KitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
				Name:      "Sam's Kitchen",
				Owner:     "Sam Smith",
				Handle:    "sammycooks",
				Private:   false,
				Avatar:    null.NewString("uploads/kitchens/ktc_2jEx1e1esA5292rBisRGuJwXc14/2pR9Azi4FC8IkJO3J6ZvjUqgz5Z.png", true),
			},
		},
	}

	// Fix time causing failed tests.
	expected.Account.CreatedAt = actual.Account.CreatedAt
	expected.Kitchens[0].CreatedAt = actual.Kitchens[0].CreatedAt
	expected.Kitchens[0].UpdatedAt = actual.Kitchens[0].UpdatedAt

	assert.Equal(t, expected, actual)
}
