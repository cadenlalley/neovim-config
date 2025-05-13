-- =================
-- Schema 1.7
-- Add recipe tags
-- =================

CREATE TABLE IF NOT EXISTS tags (
  tag_id MEDIUMINT PRIMARY KEY AUTO_INCREMENT,
  tag_type VARCHAR(10) NOT NULL,
  tag_value VARCHAR(36) NOT NULL,
  FULLTEXT(tag_value),
  UNIQUE(tag_type, tag_value)
);

CREATE TABLE IF NOT EXISTS recipe_tags (
  recipe_id CHAR(31) NOT NULL,
  tag_id MEDIUMINT NOT NULL,
  FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
  FOREIGN KEY (tag_id) REFERENCES tags(tag_id),
  UNIQUE(recipe_id, tag_id)
);
