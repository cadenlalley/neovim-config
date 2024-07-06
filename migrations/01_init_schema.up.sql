-- =================
-- Schema 1.0
-- Establish the initial schema of the database.
-- =================

-- Accounts
CREATE TABLE IF NOT EXISTS accounts (
  account_id VARCHAR(31) PRIMARY KEY,
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
  kitchen_id VARCHAR(31) PRIMARY KEY,
  account_id VARCHAR(31) NOT NULL,
  kitchen_name VARCHAR(255) NOT NULL,
  bio VARCHAR(255),
  handle VARCHAR(30) UNIQUE NOT NULL,
  avatar VARCHAR(255),
  cover VARCHAR(255),
  is_private BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);
