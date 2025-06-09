-- =================
-- Schema 1.8
-- Add recipe filters
-- =================

ALTER TABLE recipes
  ADD COLUMN difficulty TINYINT NOT NULL DEFAULT 0 AFTER servings,
  ADD COLUMN course VARCHAR(16) AFTER difficulty,
  ADD COLUMN class VARCHAR(16) AFTER course,
  ADD COLUMN cuisine VARCHAR(16) AFTER class;