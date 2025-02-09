-- =================
-- Schema 1.2
-- Add kitchen follower management.
-- =================

-- Kitchen Followers
CREATE TABLE IF NOT EXISTS kitchen_followers (
  kitchen_id CHAR(31) NOT NULL,
  followed_kitchen_id CHAR(31) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (kitchen_id, followed_kitchen_id),
  FOREIGN KEY (kitchen_id) REFERENCES kitchens(kitchen_id),
  FOREIGN KEY (followed_kitchen_id) REFERENCES kitchens(kitchen_id)
);
