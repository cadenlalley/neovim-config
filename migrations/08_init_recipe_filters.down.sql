-- =================
-- Schema 1.8
-- Add recipe filters
-- =================

ALTER TABLE recipes
  DROP COLUMN difficulty,
  DROP COLUMN course,
  DROP COLUMN class,
  DROP COLUMN cuisine;