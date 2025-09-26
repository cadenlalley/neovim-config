-- =================
-- Schema 1.9
-- Add meal_plan_category_order table for custom category ordering
-- =================

CREATE TABLE IF NOT EXISTS meal_plan_category_order (
    meal_plan_id CHAR(31) NOT NULL,
    category_order JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (meal_plan_id),
    FOREIGN KEY (meal_plan_id) REFERENCES meal_plans(meal_plan_id) ON DELETE CASCADE
);

-- Insert default category order for all existing meal plans
INSERT INTO meal_plan_category_order (meal_plan_id, category_order)
SELECT meal_plan_id, '["produce", "dairy", "meat", "fish", "snacks", "canned_goods", "breads_and_bakery", "dry_and_baking_goods", "frozen", "uncategorized"]'
FROM meal_plans
WHERE meal_plan_id NOT IN (SELECT meal_plan_id FROM meal_plan_category_order);