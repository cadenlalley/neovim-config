package fixtures

import (
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"gopkg.in/guregu/null.v4"
)

// The fixtures package is meant to provide a simple way to access
// the data that has been defined in the default fixtures on spin up.

// Return the base level test account defaults.
func GetTestAccount() accounts.Account {
	return accounts.Account{
		AccountID: "acc_2jEwcS7Rla6E5ik5ELa8uoULKOW",
		Email:     "test-service@kitchens-app.com",
		FirstName: "Sam",
		LastName:  "Smith",
		UserID:    "auth0|665e3646139d9f6300bad5e9",
		Verified:  false,
	}
}

// Return the base level test kitchen default.
func GetTestKitchen() kitchens.Kitchen {
	return kitchens.Kitchen{
		KitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
		AccountID: "acc_2jEwcS7Rla6E5ik5ELa8uoULKOW",
		Owner:     "Sam Smith",
		Name:      "Sam's Kitchen",
		Bio:       null.String{},
		Handle:    "sammycooks",
		Avatar:    null.NewString("uploads/kitchens/ktc_2jEx1e1esA5292rBisRGuJwXc14/2pR9Azi4FC8IkJO3J6ZvjUqgz5Z.png", true),
		Cover:     null.String{},
		Private:   false,
	}
}
