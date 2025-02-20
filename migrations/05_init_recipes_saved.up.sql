-- =================
-- Schema 1.4
-- Add saved recipe management
-- =================

CREATE TABLE IF NOT EXISTS recipes_saved (
  kitchen_id CHAR(31) NOT NULL,
  recipe_id CHAR(31) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (kitchen_id) REFERENCES kitchens(kitchen_id),
  FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
  UNIQUE(kitchen_id, recipe_id)
);
