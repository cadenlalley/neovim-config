-- =================
-- Schema 1.9
-- Add recipe index
-- =================

ALTER TABLE recipes ADD FULLTEXT ft_name (recipe_name);
ALTER TABLE recipes ADD FULLTEXT ft_summary (summary);
ALTER TABLE recipes ADD FULLTEXT ft_name_summary_composite (recipe_name, summary);
