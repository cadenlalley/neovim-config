-- =================
-- Schema 1.0
-- Add base user accounts for fixtures. The only one that can be logged into
-- is the test-service account. Other users are for show only.
-- =================

-- ACCOUNTS
-- =================
INSERT INTO accounts (account_id, user_id, email, first_name, last_name)
VALUES
  ('acc_2jEwcS7Rla6E5ik5ELa8uoULKOW', 'auth0|665e3646139d9f6300bad5e9', 'test-service@kitchens-app.com', 'Sam', 'Smith'),
  ('acc_2jEx1hZPbnNEoZRmkWqP2BDBB87', 'auth0|000000000000000000000001', 'test-mack@kitchens-app.com', 'Mack', 'Campbell'),
  ('acc_2jEx1fZrWeWKxxAciNcc1ng3fq5', 'auth0|000000000000000000000002', 'test-bill@kitchens-app.com', 'Bill', 'Williams');

-- KITCHENS
-- =================
INSERT INTO kitchens (kitchen_id, account_id, kitchen_name, bio, handle)
VALUES
  ('ktc_2jEx1e1esA5292rBisRGuJwXc14', 'acc_2jEwcS7Rla6E5ik5ELa8uoULKOW', "Sam's Kitchen", NULL, 'sammycooks'),
  ('ktc_2jEx1eCS13KMS8udlPoK12e5KPW', 'acc_2jEx1hZPbnNEoZRmkWqP2BDBB87', 'The Campbell Kitchen', "The Campbell's ladle out delicious delights with soup-erb flavor", 'Campbell_Soup'),
  ('ktc_2jEx1j3CVPIIAaOwGIORKqHfK89', 'acc_2jEx1fZrWeWKxxAciNcc1ng3fq5', 'Bill in the Kitchen', NULL, 'bbq_bill');

-- RECIPES
-- =================
INSERT INTO recipes (recipe_id, kitchen_id, recipe_name, summary, prep_time, cook_time, servings)
VALUES
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Homemade pumpkin pie', "With a combination of heavy cream and whole milk, this pumpkin pie has the creamiest filling, with warm spices and lovely flavor. It's baked in a flaky, buttery single crust.", 40, 60, 12),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 'ktc_2jEx1eCS13KMS8udlPoK12e5KPW', 'Spaghetti Bolognese', 'A classic Italian pasta dish made with a rich and savory meat sauce.', 20, 60, 4),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 'ktc_2jEx1j3CVPIIAaOwGIORKqHfK89', 'Chocolate Chip Cookies', 'Classic chocolate chip cookies with a soft and chewy texture.', 15, 12, 24);

INSERT INTO recipe_steps(recipe_id, step_id, instruction, group_name)
VALUES
  -- Pumpkin Pie
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 1, 'Cut the butter into slices (8-10 slices per stick). Put the butter in a bowl and place in the freezer. Fill a medium-sized measuring cup up with water, and add plenty of ice. Let both the butter and the ice sit for 5-10 minutes.', 'Pie crust'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 2, 'In the bowl of a standing mixer fitted with a paddle attachment, combine the flour, sugar, and salt. Add half of the chilled butter and mix on low, until the butter is just starting to break down, about a minute. Add the rest of the butter and continue mixing, until the butter is broken down and in various sizes. Slowly add the water, a few tablespoons at a time, and mix until the dough starts to come together but still is quite shaggy.', 'Pie crust'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 3, 'Dump the dough out on your work surface and flatten it slightly into a square. Gently fold the dough over onto itself and flatten again. Repeat this process 3 or 4 more times, until all the loose pieces are worked into the dough. Flatten the dough one last time into a circle, and wrap in plastic wrap. Refrigerate for 30 minutes (and up to 2 days) before using.', 'Pie crust'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 4, 'Adjust oven rack to lowest position, place rimmed baking sheet on rack, and heat oven to 400°F. Remove dough from refrigerator and roll out on generously floured (up to 1/4 cup) work surface to 12-inch circle about 1/8 inch thick. Roll dough loosely around rolling pin and unroll into pie plate, leaving at least 1-inch overhang on each side. Ease dough into plate by gently lifting edge of dough with one hand while pressing into plate bottom with the other.', 'Pie crust'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 5, 'Preheat oven to 400F. While the pie shell is baking, whisk cream, milk, eggs, yolks, and vanilla together in medium bowl. Combine pumpkin puree, sugars, maple syrup, cinnamon, ginger, nutmeg, and salt in large heavy-bottomed saucepan; bring to sputtering simmer over medium heat, 5 to 7 minutes. Continue to simmer pumpkin mixture, stirring constantly until thick and shiny, 10 to 15 minutes.', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 6, 'Remove pan from heat and stir in the black strap rum if using. Whisk in cream mixture until fully incorporated. Strain the mixture through fine-mesh strainer set over a medium bowl, using a spatula to press the solids through the strainer. Re-whisk the mixture and transfer to warm pre-baked pie shell. Return the pie plate with baking sheet to the oven and bake pie for 10 minutes. Reduce the heat to 300°F and continue baking until the edges of the pie are set and slightly puffed, and the center jiggles only slightly, 27 to 35 minutes longer. Transfer the pie to wire rack and cool to room temperature, 2 to 3 hours. Cut into wedges and serve with whipped cream.', 'Pie filling'),

  -- Spaghetti
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 1, 'Heat a large skillet over medium-high heat and add ground beef.', 'Cooking'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 2, 'Add chopped onion and minced garlic to the skillet and sauté until fragrant.', 'Cooking'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 3, 'Pour in the tomato sauce and let it simmer for 30 minutes.', 'Sauce'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 4, 'Cook the spaghetti according to the package instructions.', 'Pasta'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 5, 'Serve the Bolognese sauce over the cooked spaghetti.', 'Serving'),

  -- Chocolate Chip Cookies
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 1, 'Preheat the oven to 350°F (175°C).', 'Preparation'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 2, 'In a large bowl, cream together the butter, white sugar, and brown sugar until smooth.', 'Mixing'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 3, 'Beat in the eggs one at a time, then stir in the vanilla.', 'Mixing'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 4, 'Dissolve baking soda in hot water and add to the batter with salt.', 'Mixing'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 5, 'Stir in flour and chocolate chips.', 'Mixing'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 6, 'Drop large spoonfuls of dough onto ungreased pans.', 'Preparation'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 7, 'Bake for about 10 minutes, or until edges are nicely browned.', 'Baking');

INSERT INTO recipe_images(recipe_id, step_id, image_url)
VALUES
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 3, '/path-to-step-3-1.jpg'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 6, '/path-to-step-6-1.jpg'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 6, '/path-to-step-6-2.jpg');

INSERT INTO recipe_notes(recipe_id, step_id, note)
VALUES
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 2, 'If the dough is not coming together, add more water 1 tablespoon at a time until it does.');

INSERT INTO recipe_ingredients(recipe_id, ingredient_id, ingredient_name, quantity, unit, group_name)
VALUES
  -- Pumpkin Pie
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 1, 'unsalted butter, cold', 9.00, 'tbsp', 'Pie crust'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 2, 'all-purpose flour', 1.25, 'cups', 'Pie crust'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 3, 'heavy cream', 1.00, 'cup', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 4, 'whole milk', .5, 'cup', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 5, 'large eggs plus 2 large yolks', 3.00, NULL, 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 6, 'vanilla extract', 1.00, 'tsp', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 7, 'pumpkin puree', 1.00, '15 oz can', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 8, 'brown sugar', .5, 'cup', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 9, 'maple syrup', .25, 'cup', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 10, 'ground cinnamon', .75, 'tsp', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 11, 'ground ginger', .5, 'tsp', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 12, 'nutmeg', .25, 'tsp', 'Pie filling'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 13, 'salt', .75, 'tsp', 'Pie filling'),

  -- Spaghetti
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 1, 'Ground beef', 0.5, 'pound', 'Meat'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 2, 'Spaghetti', 1, 'box', 'Pasta'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 3, 'Tomato sauce', 24, 'ounce', 'Sauce'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 4, 'Garlic', 3, 'clove', 'Seasoning'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 5, 'Onion', 1, 'piece', 'Vegetable'),

  -- Chocolate Chip Cookies
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 1, 'Butter', 1, 'cup', 'Dairy'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 2, 'White sugar', 1, 'cup', 'Sweetener'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 3, 'Brown sugar', 1, 'cup', 'Sweetener'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 4, 'Eggs', 2, 'piece', 'Eggs'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 5, 'Vanilla extract', 2, 'tsp', 'Flavoring'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 6, 'Baking soda', 1, 'tsp', 'Leavening'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 7, 'Hot water', 2, 'tsp', 'Liquid'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 8, 'Salt', 0.5, 'tsp', 'Seasoning'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 9, 'All-purpose flour', 3, 'cup', 'Dry Ingredient'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 10, 'Chocolate chips', 2, 'cup', 'Add-in');

-- FOLDERS
-- =================
INSERT INTO folders (folder_id, kitchen_id, folder_name, created_at)
VALUES
  ('fld_2pPgQjn08dQzr5vjSk8WYSBTATo', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Breakfast', CURRENT_TIMESTAMP),
  ('fld_2jEx1eCS13KMS8udlPoK12e5KPW', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Healthy Lunch', CURRENT_TIMESTAMP),
  ('fld_2jEx1j3CVPIIAaOwGIORKqHfK89', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Mediteranean', CURRENT_TIMESTAMP);

INSERT INTO folder_recipes (folder_id, recipe_id, created_at)
VALUES
  ('fld_2pPgQjn08dQzr5vjSk8WYSBTATo', 'rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', CURRENT_TIMESTAMP);
