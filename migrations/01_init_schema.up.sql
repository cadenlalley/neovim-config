-- =================
-- Schema 1.0
-- Establish the initial schema of the database.
-- =================

-- Accounts
CREATE TABLE IF NOT EXISTS accounts (
  account_id CHAR(31) PRIMARY KEY,
  user_id VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(320) UNIQUE NOT NULL,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  verified BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Kitchens
CREATE TABLE IF NOT EXISTS kitchens (
  kitchen_id CHAR(31) PRIMARY KEY,
  account_id CHAR(31) NOT NULL,
  kitchen_name VARCHAR(255) NOT NULL,
  bio TINYTEXT,
  handle VARCHAR(30) UNIQUE NOT NULL,
  avatar VARCHAR(255),
  cover VARCHAR(255),
  is_private BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- Recipes
CREATE TABLE IF NOT EXISTS recipes (
  recipe_id CHAR(31) PRIMARY KEY,
  kitchen_id CHAR(31) NOT NULL,
  recipe_name VARCHAR(255) NOT NULL,
  summary TEXT,
  prep_time SMALLINT NOT NULL,
  cook_time SMALLINT NOT NULL,
  servings TINYINT NOT NULL,
  cover VARCHAR(255),
  source VARCHAR(2048),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  FOREIGN KEY (kitchen_id) REFERENCES kitchens(kitchen_id)
);

-- Recipe Steps
CREATE TABLE IF NOT EXISTS recipe_steps (
  recipe_id CHAR(31),
  step_id TINYINT,
  instruction TEXT NOT NULL,
  FOREIGN KEY (recipe_id) REFERENCES recipes (recipe_id),
  UNIQUE (recipe_id, step_id)
);

CREATE TABLE IF NOT EXISTS recipe_images (
  recipe_id CHAR(31),
  step_id TINYINT,
  image_url VARCHAR(255),
  FOREIGN KEY (recipe_id, step_id) REFERENCES recipe_steps (recipe_id, step_id)
);

CREATE TABLE IF NOT EXISTS recipe_notes (
  recipe_id CHAR(31),
  step_id TINYINT,
  note TEXT NOT NULL,
  FOREIGN KEY (recipe_id, step_id) REFERENCES recipe_steps (recipe_id, step_id)
);
