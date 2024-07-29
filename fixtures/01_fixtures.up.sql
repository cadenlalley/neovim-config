-- =================
-- Schema 1.0
-- Add base user accounts for fixtures. The only one that can be logged into
-- is the test-service account. Other users are for show only.
-- =================
INSERT INTO accounts (account_id, user_id, email, first_name, last_name)
VALUES
  ('acc_2jEwcS7Rla6E5ik5ELa8uoULKOW', 'auth0|665e3646139d9f6300bad5e9', 'test-service@kitchens-app.com', 'Sam', 'Smith'),
  ('acc_2jEx1hZPbnNEoZRmkWqP2BDBB87', 'auth0|000000000000000000000001', 'test-mack@kitchens-app.com', 'Mack', 'Campbell'),
  ('acc_2jEx1fZrWeWKxxAciNcc1ng3fq5', 'auth0|000000000000000000000002', 'test-bill@kitchens-app.com', 'Bill', 'Williams');

INSERT INTO kitchens (kitchen_id, account_id, kitchen_name, bio, handle)
VALUES
  ('ktc_2jEx1e1esA5292rBisRGuJwXc14', 'acc_2jEwcS7Rla6E5ik5ELa8uoULKOW', "Sam's Kitchen", NULL, 'sammycooks'),
  ('ktc_2jEx1eCS13KMS8udlPoK12e5KPW', 'acc_2jEx1hZPbnNEoZRmkWqP2BDBB87', 'The Campbell Kitchen', "The Campbell's ladle out delicious delights with soup-erb flavor", 'Campbell_Soup'),
  ('ktc_2jEx1j3CVPIIAaOwGIORKqHfK89', 'acc_2jEx1fZrWeWKxxAciNcc1ng3fq5', 'Bill in the Kitchen', NULL, 'bbq_bill');

INSERT INTO recipes (recipe_id, kitchen_id, recipe_name, summary, prep_time, cook_time, servings)
VALUES
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 'ktc_2jEx1e1esA5292rBisRGuJwXc14', 'Homemade pumpkin pie', "With a combination of heavy cream and whole milk, this pumpkin pie has the creamiest filling, with warm spices and lovely flavor. It's baked in a flaky, buttery single crust.", 40, 60, 12);

INSERT INTO recipe_steps(recipe_id, step_id, instruction)
VALUES
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 1, 'Cut the butter into slices (8-10 slices per stick). Put the butter in a bowl and place in the freezer. Fill a medium-sized measuring cup up with water, and add plenty of ice. Let both the butter and the ice sit for 5-10 minutes.'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 2, 'In the bowl of a standing mixer fitted with a paddle attachment, combine the flour, sugar, and salt. Add half of the chilled butter and mix on low, until the butter is just starting to break down, about a minute. Add the rest of the butter and continue mixing, until the butter is broken down and in various sizes. Slowly add the water, a few tablespoons at a time, and mix until the dough starts to come together but still is quite shaggy.'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 3, 'Dump the dough out on your work surface and flatten it slightly into a square. Gently fold the dough over onto itself and flatten again. Repeat this process 3 or 4 more times, until all the loose pieces are worked into the dough. Flatten the dough one last time into a circle, and wrap in plastic wrap. Refrigerate for 30 minutes (and up to 2 days) before using.'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 4, 'Adjust oven rack to lowest position, place rimmed baking sheet on rack, and heat oven to 400°F. Remove dough from refrigerator and roll out on generously floured (up to 1/4 cup) work surface to 12-inch circle about 1/8 inch thick. Roll dough loosely around rolling pin and unroll into pie plate, leaving at least 1-inch overhang on each side. Ease dough into plate by gently lifting edge of dough with one hand while pressing into plate bottom with the other.'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 5, 'Preheat oven to 400F. While the pie shell is baking, whisk cream, milk, eggs, yolks, and vanilla together in medium bowl. Combine pumpkin puree, sugars, maple syrup, cinnamon, ginger, nutmeg, and salt in large heavy-bottomed saucepan; bring to sputtering simmer over medium heat, 5 to 7 minutes. Continue to simmer pumpkin mixture, stirring constantly until thick and shiny, 10 to 15 minutes.'),
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 6, 'Remove pan from heat and stir in the black strap rum if using. Whisk in cream mixture until fully incorporated. Strain the mixture through fine-mesh strainer set over a medium bowl, using a spatula to press the solids through the strainer. Re-whisk the mixture and transfer to warm pre-baked pie shell. Return the pie plate with baking sheet to the oven and bake pie for 10 minutes. Reduce the heat to 300°F and continue baking until the edges of the pie are set and slightly puffed, and the center jiggles only slightly, 27 to 35 minutes longer. Transfer the pie to wire rack and cool to room temperature, 2 to 3 hours. Cut into wedges and serve with whipped cream.');

  INSERT INTO recipe_images(recipe_id, step_id, image_url)
  VALUES
    ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 3, '/path-to-step-3-1.jpg'),
    ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 6, '/path-to-step-6-1.jpg'),
    ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 6, '/path-to-step-6-2.jpg');

INSERT INTO recipe_notes(recipe_id, step_id, note)
VALUES
  ('rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T', 2, 'If the dough is not coming together, add more water 1 tablespoon at a time until it does.');