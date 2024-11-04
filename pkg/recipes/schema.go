package recipes

import "encoding/json"

var JsonSchema json.RawMessage = json.RawMessage(`{
		"type": "json_schema",
		"json_schema": {
				"name": "url_recipe_response",
				"schema": {
						"type": "object",
						"properties": {
								"name": {"type": "string"},
								"summary": {"type": "string"},
								"prepTime": {"type": "integer", "description": "the preparation time in minutes"},
								"cookTime": {"type": "integer", "description": "the cook time in minutes"},
								"servings": {"type": "integer", "description": "the number of servings"},
								"ingredients": {
										"type": "array",
										"items": {
												"type": "object",
												"properties": {
														"ingredientId": {"type": "integer"},
														"name": {"type": "string"},
														"quantity": {"type": "number"},
														"unit": {
															"type": ["string", "null"],
															"description": "optional unit of measurement, if a unit doesn't make sense for the ingredient set it to n/a",
															"enum": ["bag","bottle","box","can","clove","cup","dash","drop","gallon","gram","jar","kilogram","liter","milliliter","ounce","packet","piece","pint","pinch","pound","quart","slice","stick","tbsp","tsp", "n/a"]
														},
														"group": {"type": ["string", "null"]}
												},
												"required": ["ingredientId","name","quantity","group","unit"],
												"additionalProperties": false
										}
								},
								"steps": {
										"type": "array",
										"items": {
												"type": "object",
												"properties": {
														"stepId": {"type": "integer"},
														"instruction": {"type": "string"},
														"note": {"type": ["string", "null"]},
														"group": {"type": ["string", "null"]}
												},
												"required": ["stepId","instruction","note","group"],
												"additionalProperties": false
										}
								}
						},
						"required": ["name","summary","prepTime","cookTime","servings","ingredients","steps"],
						"additionalProperties": false
				},
				"strict": true
		}
}`)
