-- =================
-- Schema 1.0
-- Remove all data from the database.
-- =================

DELETE FROM recipe_tags WHERE 1=1;
DELETE FROM recipe_review_likes WHERE 1=1;
DELETE FROM recipe_reviews WHERE 1=1;
DELETE FROM recipes_saved WHERE 1=1;
DELETE FROM kitchen_followers WHERE 1=1;
DELETE FROM folder_recipes WHERE 1=1;
DELETE FROM folders WHERE 1=1;
DELETE FROM recipe_ingredients WHERE 1=1;
DELETE FROM recipe_notes WHERE 1=1;
DELETE FROM recipe_images WHERE 1=1;
DELETE FROM recipe_steps WHERE 1=1;
DELETE FROM tags WHERE 1=1;
DELETE FROM recipes WHERE 1=1;
DELETE FROM kitchens WHERE 1=1;
DELETE FROM accounts WHERE 1=1;