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
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Profile
CREATE TABLE IF NOT EXISTS profiles (
  profile_id VARCHAR(31) PRIMARY KEY,
  account_id VARCHAR(31) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  bio VARCHAR(255),
  handle VARCHAR(30) UNIQUE NOT NULL,
  avatar_photo VARCHAR(255),
  cover_photo VARCHAR(255),
  public BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);
