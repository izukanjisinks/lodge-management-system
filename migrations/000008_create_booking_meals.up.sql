-- Join table linking bookings to their selected meal add-ons
CREATE TABLE IF NOT EXISTS booking_meals (
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    meal_id    UUID NOT NULL REFERENCES meals(id) ON DELETE RESTRICT,
    PRIMARY KEY (booking_id, meal_id)
);

CREATE INDEX idx_booking_meals_booking_id ON booking_meals(booking_id);
CREATE INDEX idx_booking_meals_meal_id ON booking_meals(meal_id);
