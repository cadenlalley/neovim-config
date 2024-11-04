package recipes

import "encoding/json"

var Sample json.RawMessage = json.RawMessage(`{
  "name": "Carrots and Ranch",
  "summary": "The best damn carrots and ranch you've ever had",
  "cover": null,
  "prepTime": 1,
  "cookTime": 0,
  "servings": 10,
  "source": null,
  "ingredients": [
    {
      "ingredientId": 1,
      "name": "Carrots",
      "quantity": 1,
      "unit": "pound",
      "group": "Veggies"
    },
    {
      "ingredientId": 2,
      "name": "Ranch",
      "quantity": 4,
      "unit": "tbsp",
      "group": "Dips"
    }
  ],
  "steps": [
    {
      "stepId": 1,
      "instruction": "unseal bag of carrots onto plate",
      "images": null,
      "note": "use scissors for easy access",
      "group": "Preparation"
    }, {
      "stepId": 2,
      "instruction": "unseal ranch and pour into bowl",
      "images": null,
      "note": null,
      "group": "Preparation"
    }, {
      "stepId": 3,
      "instruction": "arrange carrots and ranch together pleasingly",
      "images": null,
      "note": "throw some microgreens on the plate to be fancy",
      "group": "Presentation"
    }
  ]
}`)
