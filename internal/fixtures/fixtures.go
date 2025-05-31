package fixtures

import (
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/kitchens-io/kitchens-api/pkg/ptr"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
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

// Return the base level test recipe default.
func GetTestRecipe() recipes.Recipe {
	return recipes.Recipe{
		RecipeID:  "rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T",
		KitchenID: "ktc_2jEx1e1esA5292rBisRGuJwXc14",
		Name:      "Homemade pumpkin pie",
		Summary:   null.NewString("With a combination of heavy cream and whole milk, this pumpkin pie has the creamiest filling, with warm spices and lovely flavor. It's baked in a flaky, buttery single crust.", true),
		PrepTime:  ptr.Int(40),
		CookTime:  ptr.Int(60),
		Servings:  ptr.Int(12),
		Cover:     null.NewString("uploads/recipes/rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T/2pR9B2cIFxj82GDTDB44lpMzYHu.png", true),
		Source:    null.String{},

		// Reviews
		ReviewCount:  4,
		ReviewRating: 3.5,

		// Recipe Data
		SourceDomain: null.String{},
		Ingredients: []recipes.RecipeIngredient{
			{
				IngredientID: 1,
				Name:         "unsalted butter, cold",
				Quantity:     null.NewFloat(9.00, true),
				Unit:         null.NewString("tbsp", true),
				Group:        null.NewString("Pie crust", true),
			},
			{
				IngredientID: 2,
				Name:         "all-purpose flour",
				Quantity:     null.NewFloat(1.25, true),
				Unit:         null.NewString("cups", true),
				Group:        null.NewString("Pie crust", true),
			},
			{
				IngredientID: 3,
				Name:         "heavy cream",
				Quantity:     null.NewFloat(1.00, true),
				Unit:         null.NewString("cup", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 4,
				Name:         "whole milk",
				Quantity:     null.NewFloat(0.5, true),
				Unit:         null.NewString("cup", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 5,
				Name:         "large eggs plus 2 large yolks",
				Quantity:     null.NewFloat(3.00, true),
				Unit:         null.String{},
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 6,
				Name:         "vanilla extract",
				Quantity:     null.NewFloat(1.00, true),
				Unit:         null.NewString("tsp", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 7,
				Name:         "pumpkin puree",
				Quantity:     null.NewFloat(1.00, true),
				Unit:         null.NewString("15 oz can", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 8,
				Name:         "brown sugar",
				Quantity:     null.NewFloat(0.5, true),
				Unit:         null.NewString("cup", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 9,
				Name:         "maple syrup",
				Quantity:     null.NewFloat(0.25, true),
				Unit:         null.NewString("cup", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 10,
				Name:         "ground cinnamon",
				Quantity:     null.NewFloat(0.75, true),
				Unit:         null.NewString("tsp", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 11,
				Name:         "ground ginger",
				Quantity:     null.NewFloat(0.5, true),
				Unit:         null.NewString("tsp", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 12,
				Name:         "nutmeg",
				Quantity:     null.NewFloat(0.25, true),
				Unit:         null.NewString("tsp", true),
				Group:        null.NewString("Pie filling", true),
			},
			{
				IngredientID: 13,
				Name:         "salt",
				Quantity:     null.NewFloat(0.75, true),
				Unit:         null.NewString("tsp", true),
				Group:        null.NewString("Pie filling", true),
			},
		},
		Steps: []recipes.RecipeStep{
			{
				StepID:      1,
				Instruction: "Cut the butter into slices (8-10 slices per stick). Put the butter in a bowl and place in the freezer. Fill a medium-sized measuring cup up with water, and add plenty of ice. Let both the butter and the ice sit for 5-10 minutes.",
				Group:       null.NewString("Pie crust", true),
			},
			{
				StepID:      2,
				Instruction: "In the bowl of a standing mixer fitted with a paddle attachment, combine the flour, sugar, and salt. Add half of the chilled butter and mix on low, until the butter is just starting to break down, about a minute. Add the rest of the butter and continue mixing, until the butter is broken down and in various sizes. Slowly add the water, a few tablespoons at a time, and mix until the dough starts to come together but still is quite shaggy.",
				Group:       null.NewString("Pie crust", true),
				Note:        "If the dough is not coming together, add more water 1 tablespoon at a time until it does.",
			},
			{
				StepID:      3,
				Instruction: "Dump the dough out on your work surface and flatten it slightly into a square. Gently fold the dough over onto itself and flatten again. Repeat this process 3 or 4 more times, until all the loose pieces are worked into the dough. Flatten the dough one last time into a circle, and wrap in plastic wrap. Refrigerate for 30 minutes (and up to 2 days) before using.",
				Group:       null.NewString("Pie crust", true),
				Images: []string{
					"/path-to-step-3-1.jpg",
				},
			},
			{
				StepID:      4,
				Instruction: "Adjust oven rack to lowest position, place rimmed baking sheet on rack, and heat oven to 400°F. Remove dough from refrigerator and roll out on generously floured (up to 1/4 cup) work surface to 12-inch circle about 1/8 inch thick. Roll dough loosely around rolling pin and unroll into pie plate, leaving at least 1-inch overhang on each side. Ease dough into plate by gently lifting edge of dough with one hand while pressing into plate bottom with the other.",
				Group:       null.NewString("Pie crust", true),
			},
			{
				StepID:      5,
				Instruction: "Preheat oven to 400F. While the pie shell is baking, whisk cream, milk, eggs, yolks, and vanilla together in medium bowl. Combine pumpkin puree, sugars, maple syrup, cinnamon, ginger, nutmeg, and salt in large heavy-bottomed saucepan; bring to sputtering simmer over medium heat, 5 to 7 minutes. Continue to simmer pumpkin mixture, stirring constantly until thick and shiny, 10 to 15 minutes.",
				Group:       null.NewString("Pie filling", true),
			},
			{
				StepID:      6,
				Instruction: "Remove pan from heat and stir in the black strap rum if using. Whisk in cream mixture until fully incorporated. Strain the mixture through fine-mesh strainer set over a medium bowl, using a spatula to press the solids through the strainer. Re-whisk the mixture and transfer to warm pre-baked pie shell. Return the pie plate with baking sheet to the oven and bake pie for 10 minutes. Reduce the heat to 300°F and continue baking until the edges of the pie are set and slightly puffed, and the center jiggles only slightly, 27 to 35 minutes longer. Transfer the pie to wire rack and cool to room temperature, 2 to 3 hours. Cut into wedges and serve with whipped cream.",
				Group:       null.NewString("Pie filling", true),
				Images: []string{
					"/path-to-step-6-1.jpg",
					"/path-to-step-6-2.jpg",
				},
			},
		},
		ShareURL: "/2jbgfAMKOCnKrWQroRBkXPIRI6T/homemade-pumpkin-pie",
	}
}
