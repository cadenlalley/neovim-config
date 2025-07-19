-- =================
-- Schema 1.7
-- Add meal plan tables
-- =================

CREATE TABLE IF NOT EXISTS meal_plans (
    meal_plan_id CHAR(31) NOT NULL PRIMARY KEY,
    account_id CHAR(31) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE TABLE IF NOT EXISTS meal_plan_recipes (
	meal_plan_recipe_id INT NOT NULL AUTO_INCREMENT,
    meal_plan_id CHAR(31) NOT NULL,
    recipe_id CHAR(31) NOT NULL,
    day_number TINYINT NOT NULL,
    serving_size TINYINT NOT NULL,
    FOREIGN KEY (meal_plan_id) REFERENCES meal_plans(meal_plan_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
    UNIQUE (meal_plan_recipe_id, meal_plan_id, recipe_id, day_number)
);
