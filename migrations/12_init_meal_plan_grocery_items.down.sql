-- =================
-- Schema 1.8
-- Add meal plan grocery list items table
-- =================

DROP TABLE IF EXISTS meal_plan_grocery_list_items;

ALTER TABLE meal_plans
DROP COLUMN grocery_list_is_dirty;

