CREATE TABLE IF NOT EXISTS roles (
    role_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO roles (name, description) VALUES
    ('admin',        'Full system access — manages users, rooms, bookings, and configuration'),
    ('manager',      'Oversees operations — approves bookings, views reports, manages rooms'),
    ('receptionist', 'Front-desk staff — handles bookings, clients, and invoices'),
    ('cleaner',      'Housekeeping staff — views assigned rooms and cleaning schedule')
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users
    ADD CONSTRAINT fk_users_role
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE SET NULL;
