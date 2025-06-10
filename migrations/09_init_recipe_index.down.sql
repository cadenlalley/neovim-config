-- =================
-- Schema 1.9
-- Add recipe index
-- =================

ALTER TABLE recipes DROP INDEX ft_name;
ALTER TABLE recipes DROP INDEX ft_summary;
ALTER TABLE recipes DROP INDEX ft_name_summary_composite;
