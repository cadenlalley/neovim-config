-- =================
-- Schema 1.1
-- Add ingredient management.
-- =================

-- Recipe Ingredients
CREATE TABLE IF NOT EXISTS recipe_ingredients (
  recipe_id CHAR(31),
  ingredient_id TINYINT NOT NULL,
  ingredient_name VARCHAR(255) NOT NULL,
  quantity DECIMAL(4,2) NOT NULL,
  unit VARCHAR(255) NOT NULL,
  FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id)
);