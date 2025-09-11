-- =================
-- Schema 1.8
-- Add meal plan grocery list items table
-- =================

CREATE TABLE IF NOT EXISTS meal_plan_grocery_list_items (
	item_id INT AUTO_INCREMENT,
	alike_id INT NOT NULL,
	meal_plan_id CHAR(31) NOT NULL,
	recipe_id CHAR(31) NOT NULL,
	name CHAR(255) NOT NULL,
	quantity DECIMAL(10,2) NOT NULL,
	unit CHAR(255) NOT NULL,
	category CHAR(255) NOT NULL,
	is_user_created BOOL DEFAULT false,
	marked BOOL DEFAULT FALSE,
	FOREIGN KEY (meal_plan_id) REFERENCES meal_plans(meal_plan_id),
	UNIQUE (item_id, meal_plan_id)
);

ALTER TABLE meal_plans
ADD COLUMN grocery_list_is_dirty BOOL DEFAULT false;
