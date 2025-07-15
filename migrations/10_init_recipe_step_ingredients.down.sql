-- =================
-- Schema 1.10
-- Remove recipe step ingredients
-- =================

ALTER TABLE recipe_steps DROP COLUMN ingredient_ids;
