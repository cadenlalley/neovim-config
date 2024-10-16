package media

// uploads/accounts/{accountID}/
func GetAccountMediaPath(accountID string) string {
	return "uploads/accounts/" + accountID + "/"
}

// uploads/kitchens/{kitchenID}/
func GetKitchenMediaPath(kitchenID string) string {
	return "uploads/kitchens/" + kitchenID + "/"
}

// uploads/recipes/{recipeID}/
func GetRecipeMediaPath(recipeID string) string {
	return "uploads/recipes/" + recipeID + "/"
}

// uploads/imports/{accountID}/
func GetImportMediaPath(accountID string) string {
	return "uploads/imports/" + accountID + "/"
}
