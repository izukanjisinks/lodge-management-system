ALTER TABLE rooms DROP CONSTRAINT IF EXISTS uq_rooms_name_org;
ALTER TABLE roles DROP CONSTRAINT IF EXISTS uq_roles_name_org;
ALTER TABLE individual_profiles DROP CONSTRAINT IF EXISTS uq_individual_profiles_email_org;
ALTER TABLE corporate_profiles DROP CONSTRAINT IF EXISTS uq_corporate_profiles_email_org;
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS uq_bookings_booking_number_org;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS uq_orders_order_number_org;
