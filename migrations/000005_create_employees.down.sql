ALTER TABLE departments DROP CONSTRAINT IF EXISTS fk_departments_manager;
DROP TABLE IF EXISTS employees;
DROP TYPE IF EXISTS employment_status;
DROP TYPE IF EXISTS employment_type;
