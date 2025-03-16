-- =================
-- Schema 1.6
-- Add recipe reviews
-- =================

-- Recipe Reviews
CREATE TABLE IF NOT EXISTS recipe_reviews (
  review_id CHAR(31) PRIMARY KEY,
  recipe_id CHAR(31) NOT NULL,
  reviewer_kitchen_id CHAR(31) NOT NULL,
  review_description VARCHAR(2000),
  rating TINYINT NOT NULL,
  media_path VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
  FOREIGN KEY (reviewer_kitchen_id) REFERENCES kitchens(kitchen_id),
  UNIQUE(recipe_id, reviewer_kitchen_id)
);

-- Recipe Review Likes
CREATE TABLE IF NOT EXISTS recipe_review_likes (
  review_id CHAR(31) NOT NULL,
  kitchen_id CHAR(31) NOT NULL,
  PRIMARY KEY (review_id, kitchen_id),
  FOREIGN KEY (review_id) REFERENCES recipe_reviews(review_id),
  FOREIGN KEY (kitchen_id) REFERENCES kitchens(kitchen_id)
);
