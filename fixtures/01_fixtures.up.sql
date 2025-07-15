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
  ('acc_2jEx1fZrWeWKxxAciNcc1ng3fq5', 'auth0|000000000000000000000002', 'test-bill@kitchens-app.com', 'Bill', 'Williams'),
  ('acc_2qwGARyD8go0GdknbRK83NOob50', 'auth0|000000000000000000000003', 'test-sarah@kitchens-app.com', 'Sarah', 'Marry');

-- KITCHENS
-- =================
INSERT INTO kitchens (kitchen_id, account_id, kitchen_name, bio, handle, avatar, cover)
VALUES
  ('ktc_2jEx1e1esA5292rBisRGuJwXc14', 'acc_2jEwcS7Rla6E5ik5ELa8uoULKOW', "Sam's Kitchen", NULL, 'sammycooks', "uploads/kitchens/ktc_2jEx1e1esA5292rBisRGuJwXc14/2pR9Azi4FC8IkJO3J6ZvjUqgz5Z.png", NULL),
  ('ktc_2jEx1eCS13KMS8udlPoK12e5KPW', 'acc_2jEx1hZPbnNEoZRmkWqP2BDBB87', 'The Campbell Kitchen', "The Campbell's ladle out delicious delights with soup-erb flavor", 'Campbell_Soup', NULL, NULL),
  ('ktc_2jEx1j3CVPIIAaOwGIORKqHfK89', 'acc_2jEx1fZrWeWKxxAciNcc1ng3fq5', 'Bill in the Kitchen', NULL, 'bbq_bill', "uploads/kitchens/ktc_2jEx1j3CVPIIAaOwGIORKqHfK89/2qwGAZKuTbHO96b6VXZohC6VBes.jpg", NULL),
  ('ktc_2qwGAW3mHabA1ICXhRms3SJU11E', 'acc_2qwGARyD8go0GdknbRK83NOob50', "Sarah's Sauce House", NULL, 'saucy_sarah', NULL, NULL);

-- RECIPES
-- =================
INSERT INTO recipes (recipe_id, kitchen_id, recipe_name, summary, prep_time, cook_time, servings, difficulty, course, class, cuisine, cover, created_at)
VALUES
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Homemade pumpkin pie', "With a combination of heavy cream and whole milk, this pumpkin pie has the creamiest filling, with warm spices and a flavor to love. It's baked in a flaky, buttery single crust.", 40, 60, 12, 2, 'dessert', 'dessert', 'American', "uploads/recipes/rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T/2pR9B2cIFxj82GDTDB44lpMzYHu.png", CURRENT_TIMESTAMP),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 'ktc_2jEx1eCS13KMS8udlPoK12e5KPW', 'Spaghetti Bolognese', "A classic Italian pasta dish made with a rich and savory meat sauce that you'll love.", 20, 60, 4, 2, 'dinner', 'main', 'Italian', NULL, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 1 DAY)),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 'ktc_2jEx1j3CVPIIAaOwGIORKqHfK89', "Love's Chocolate Chip Cookies", 'Classic homemade chocolate chip cookies with a soft and chewy texture.', 15, 12, 24, 1, 'dessert', 'dessert', 'American', NULL, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 2 DAY));

INSERT INTO recipe_steps(recipe_id, step_id, instruction, group_name, ingredient_ids)
VALUES
  -- Pumpkin Pie
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 1, 'Cut the butter into slices (8-10 slices per stick). Put the butter in a bowl and place in the freezer. Fill a medium-sized measuring cup up with water, and add plenty of ice. Let both the butter and the ice sit for 5-10 minutes.', 'Pie crust', '[1]'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 2, 'In the bowl of a standing mixer fitted with a paddle attachment, combine the flour, sugar, and salt. Add half of the chilled butter and mix on low, until the butter is just starting to break down, about a minute. Add the rest of the butter and continue mixing, until the butter is broken down and in various sizes. Slowly add the water, a few tablespoons at a time, and mix until the dough starts to come together but still is quite shaggy.', 'Pie crust', '[1,2]'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 3, 'Dump the dough out on your work surface and flatten it slightly into a square. Gently fold the dough over onto itself and flatten again. Repeat this process 3 or 4 more times, until all the loose pieces are worked into the dough. Flatten the dough one last time into a circle, and wrap in plastic wrap. Refrigerate for 30 minutes (and up to 2 days) before using.', 'Pie crust', '[]'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 4, 'Adjust oven rack to lowest position, place rimmed baking sheet on rack, and heat oven to 400°F. Remove dough from refrigerator and roll out on generously floured (up to 1/4 cup) work surface to 12-inch circle about 1/8 inch thick. Roll dough loosely around rolling pin and unroll into pie plate, leaving at least 1-inch overhang on each side. Ease dough into plate by gently lifting edge of dough with one hand while pressing into plate bottom with the other.', 'Pie crust', '[]'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 5, 'Preheat oven to 400F. While the pie shell is baking, whisk cream, milk, eggs, yolks, and vanilla together in medium bowl. Combine pumpkin puree, sugars, maple syrup, cinnamon, ginger, nutmeg, and salt in large heavy-bottomed saucepan; bring to sputtering simmer over medium heat, 5 to 7 minutes. Continue to simmer pumpkin mixture, stirring constantly until thick and shiny, 10 to 15 minutes.', 'Pie filling', '[3,4,5,6,7,8,9,10,11,12,13]'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 6, 'Remove pan from heat and stir in the black strap rum if using. Whisk in cream mixture until fully incorporated. Strain the mixture through fine-mesh strainer set over a medium bowl, using a spatula to press the solids through the strainer. Re-whisk the mixture and transfer to warm pre-baked pie shell. Return the pie plate with baking sheet to the oven and bake pie for 10 minutes. Reduce the heat to 300°F and continue baking until the edges of the pie are set and slightly puffed, and the center jiggles only slightly, 27 to 35 minutes longer. Transfer the pie to wire rack and cool to room temperature, 2 to 3 hours. Cut into wedges and serve with whipped cream.', 'Pie filling', '[]'),

  -- Spaghetti
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 1, 'Heat a large skillet over medium-high heat and add ground beef.', 'Cooking', '[]'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 2, 'Add chopped onion and minced garlic to the skillet and sauté until fragrant.', 'Cooking', '[]'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 3, 'Pour in the tomato sauce and let it simmer for 30 minutes.', 'Sauce', '[]'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 4, 'Cook the spaghetti according to the package instructions.', 'Pasta', '[]'),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 5, 'Serve the Bolognese sauce over the cooked spaghetti.', 'Serving', '[]'),

  -- Chocolate Chip Cookies
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 1, 'Preheat the oven to 350°F (175°C).', 'Preparation', '[]'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 2, 'In a large bowl, cream together the butter, white sugar, and brown sugar until smooth.', 'Mixing', '[1,2,3]'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 3, 'Beat in the eggs one at a time, then stir in the vanilla.', 'Mixing', '[4,5]'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 4, 'Dissolve baking soda in hot water and add to the batter with salt.', 'Mixing', '[6,7,8]'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 5, 'Stir in flour and chocolate chips.', 'Mixing', '[9,10]'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 6, 'Drop large spoonfuls of dough onto ungreased pans.', 'Preparation', '[]'),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 7, 'Bake for about 10 minutes, or until edges are nicely browned.', 'Baking', '[]');

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
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 1, 'ground beef', 0.5, 'pound', NULL),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 2, 'spaghetti', 1, 'box', NULL),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 3, 'tomato sauce', 24, 'ounce', NULL),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 4, 'garlic', 3, 'clove', NULL),
  ('rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', 5, 'onion', 1, 'piece', NULL),

  -- Chocolate Chip Cookies
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 1, 'butter', 1, 'cup', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 2, 'white sugar', 1, 'cup', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 3, 'brown sugar', 1, 'cup', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 4, 'eggs', 2, 'piece', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 5, 'vanilla extract', 2, 'tsp', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 6, 'baking soda', 1, 'tsp', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 7, 'hot water', 2, 'tsp', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 8, 'salt', 0.5, 'tsp', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 9, 'all-purpose flour', 3, 'cup', NULL),
  ('rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', 10, 'chocolate chips', 2, 'cup', NULL);

-- FOLDERS
-- =================
INSERT INTO folders (folder_id, kitchen_id, folder_name, cover)
VALUES
  ('fld_2pPgQjn08dQzr5vjSk8WYSBTATo', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Breakfast', '/uploads/folders/fld_2pPgQjn08dQzr5vjSk8WYSBTATo/2pR6BliLuDCHYJKq7Eqkb9l55bS.png'),
  ('fld_2jEx1eCS13KMS8udlPoK12e5KPW', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Healthy Lunch', "uploads/folders/fld_2pPgQjn08dQzr5vjSk8WYSBTATo/2pR6Bouca6LsCl1fXHDK9p264bi.png"),
  ('fld_2jEx1j3CVPIIAaOwGIORKqHfK89', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Mediteranean', "uploads/folders/fld_2pPgQjn08dQzr5vjSk8WYSBTATo/2pR6BrWvEtBW5uBwoNc213qALBc.png");

INSERT INTO folder_recipes (folder_id, recipe_id)
VALUES
  ('fld_2pPgQjn08dQzr5vjSk8WYSBTATo', 'rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T'); -- homemade pumpkin pie -> breakfast

-- FOLLOWERS
-- =================
INSERT INTO kitchen_followers (kitchen_id, followed_kitchen_id)
VALUES
  ('ktc_2jEx1e1esA5292rBisRGuJwXc14', 'ktc_2jEx1eCS13KMS8udlPoK12e5KPW'), -- sammycooks -> Campbell_Soup
  ('ktc_2jEx1e1esA5292rBisRGuJwXc14', 'ktc_2jEx1j3CVPIIAaOwGIORKqHfK89'), -- sammycooks -> bbqbill
  ('ktc_2jEx1j3CVPIIAaOwGIORKqHfK89', 'ktc_2jEx1e1esA5292rBisRGuJwXc14'); -- bbqbill -> sammycooks

-- SAVED RECIPES
-- =================
INSERT INTO recipes_saved (kitchen_id, recipe_id)
VALUES
  ('ktc_2jEx1e1esA5292rBisRGuJwXc14', 'rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7'); -- sammycooks -> spaghetti bolognese;

-- RECIPE REVIEWS
-- =================
INSERT INTO recipe_reviews (review_id, recipe_id, reviewer_kitchen_id, rating, review_description)
VALUES
  ('rvw_2tvEpuxUBa47rgYIwxPoyMfPMl6', 'rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 5, "Incredible spaghetti recipe! The sauce is perfectly balanced with rich herbs, and the pasta is cooked to al dente perfection. A true homemade taste that brings back memories of my grandmother's cooking. Absolutely delicious and will definitely make again!"),
  ('rvw_2tvEpsFp6GXKthzifUVPDRObPLP', 'rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 'ktc_2jEx1eCS13KMS8udlPoK12e5KPW', 4, "Really good spaghetti that's quick and easy to make. My only minor critique is that I might add a bit more garlic next time for extra punch. Overall, a solid weeknight dinner that the whole family enjoyed."),
  ('rvw_2tvEppPKbjnGgYZHdud9OIxFmGB', 'rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 'ktc_2jEx1j3CVPIIAaOwGIORKqHfK89', 3, "It's a standard, no-frills approach that gets the job done. The sauce was okay, but could use some more depth. It's reliable but lacks that wow factor that would make me want to make it again and again."),
  ('rvw_2tvEpsnHblNIdAXshLNyWHpy4MS', 'rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 'ktc_2qwGAW3mHabA1ICXhRms3SJU11E', 2, "Disappointing spaghetti that missed the mark. The sauce was too watery, and the seasoning felt bland and uninspired. There are definitely better spaghetti recipes out there that I'd recommend over this one.");

-- RECIPE TAGS
-- =================
INSERT INTO tags (tag_id, tag_type, tag_value)
VALUES
  -- Generated from Pumpkin Pie
  (1, 'diet', 'vegetarian'),
  (2, 'ingredient', 'brown-sugar'),
  (3, 'ingredient', 'heavy-cream'),
  (4, 'ingredient', 'pumpkin-puree'),
  (5, 'keyword', 'baked-pie'),
  (6, 'keyword', 'creamy-filling'),
  (7, 'keyword', 'fall-dessert'),
  (8, 'keyword', 'homemade-pie'),
  (9, 'keyword', 'spice-flavors'),
  (10, 'keyword', 'thanksgiving-dessert'),
  -- Generated from Spaghetti
  (11, 'ingredient', 'ground-beef'),
  (12, 'ingredient', 'spaghetti'),
  (13, 'ingredient', 'tomato-sauce'),
  (14, 'keyword', 'classic-pasta-dish'),
  (15, 'keyword', 'comfort-food'),
  (16, 'keyword', 'homemade-bolognese'),
  (17, 'keyword', 'quick-recipe'),
  (18, 'keyword', 'savory-meat-sauce'),
  -- Generated from Chocolate Chip Cookies
  (19, 'ingredient', 'all-purpose-flour'),
  (20, 'ingredient', 'butter'),
  (21, 'ingredient', 'chocolate-chips'),
  (22, 'keyword', 'baking-from-scratch'),
  (23, 'keyword', 'classic-cookie'),
  (24, 'keyword', 'easy-recipe'),
  (25, 'keyword', 'kid-friendly'),
  (26, 'keyword', 'soft-and-chewy');

INSERT INTO recipe_tags (recipe_id, tag_id)
  SELECT 'rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', tag_id
    FROM tags WHERE (tag_id > 0 AND tag_id <= 10);

INSERT INTO recipe_tags (recipe_id, tag_id)
SELECT 'rcp_2oSUH8pCU0gfKbQPxg1JLD8DzJ7', tag_id
    FROM tags WHERE tag_id > 10 AND tag_id <= 18;

INSERT INTO recipe_tags (recipe_id, tag_id)
SELECT 'rcp_2oSUH6fs0iCWGNP1AF2XemKYClo', tag_id
    FROM tags WHERE tag_id > 18 AND tag_id <= 26;
