-- =================
-- Schema 1.9
-- Add ingredients to the steps table
-- =================

ALTER TABLE recipe_steps ADD COLUMN ingredient_ids JSON NOT NULL DEFAULT ('[]');