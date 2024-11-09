-- =================
-- Schema 1.2
-- Add folder management.
-- =================

-- Folders
CREATE TABLE IF NOT EXISTS folders (
  folder_id CHAR(31) PRIMARY KEY,
  kitchen_id CHAR(31),
  folder_name VARCHAR(255) NOT NULL,
  cover VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (kitchen_id) REFERENCES kitchens(kitchen_id),
  UNIQUE(folder_name, kitchen_id)
);

CREATE TABLE IF NOT EXISTS folder_recipes (
  folder_id CHAR(31) PRIMARY KEY,
  recipe_id CHAR(31),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (folder_id) REFERENCES folders(folder_id),
  FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
  UNIQUE(folder_id, recipe_id)
);
