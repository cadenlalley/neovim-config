-- =================
-- Schema 1.1
-- Add ingredient management and groups.
-- =================

-- Recipe Ingredients
CREATE TABLE IF NOT EXISTS recipe_ingredients (
  recipe_id CHAR(31),
  ingredient_id TINYINT NOT NULL,
  ingredient_name VARCHAR(255) NOT NULL,
  quantity DECIMAL(6,3),
  unit VARCHAR(255),
  group_name VARCHAR(255),
  FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id)
);

ALTER TABLE recipe_steps ADD COLUMN group_name VARCHAR(255);