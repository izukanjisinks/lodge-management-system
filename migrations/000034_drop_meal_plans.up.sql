-- Phase 2: Remove the old meal plan system now that menus/orders are in place.
-- Run this only after menu data has been seeded in the new tables.

DROP TABLE IF EXISTS booking_meal_plans;
DROP TABLE IF EXISTS meal_plans;
