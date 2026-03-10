-- Seed demo data: departments, positions, users, employees, leave types, leave balances
-- Admin user is intentionally excluded (handled by the API on startup via SeedSuperAdmin)

DO $$
DECLARE
    v_admin_id UUID;

    -- Role IDs
    v_role_super_admin  UUID;
    v_role_hr_manager   UUID;
    v_role_manager      UUID;
    v_role_employee     UUID;

    -- Department IDs
    v_dept_eng   UUID; v_dept_hr    UUID; v_dept_fin   UUID;
    v_dept_sal   UUID; v_dept_ops   UUID;
    v_dept_be    UUID; v_dept_fe    UUID; v_dept_do    UUID;
    v_dept_ta    UUID; v_dept_pay   UUID;
    v_dept_acc   UUID; v_dept_tre   UUID;
    v_dept_sls   UUID; v_dept_mkt   UUID;

    -- Position IDs
    v_pos_swe1   UUID; v_pos_swe2   UUID; v_pos_swesr  UUID; v_pos_swest  UUID;
    v_pos_fe1    UUID; v_pos_fesr   UUID; v_pos_do1    UUID; v_pos_engmgr UUID;
    v_pos_hrgen  UUID; v_pos_hrtas  UUID; v_pos_hrmgr  UUID; v_pos_paysp  UUID;
    v_pos_acc1   UUID; v_pos_acc2   UUID; v_pos_finmgr UUID;
    v_pos_salrep UUID; v_pos_salae  UUID; v_pos_mktsp  UUID; v_pos_salmgr UUID;
    v_pos_opsan  UUID; v_pos_opsmgr UUID;

    -- User IDs
    v_u_hrmgr    UUID; v_u_engmgr   UUID; v_u_finmgr   UUID;
    v_u_salmgr   UUID; v_u_opsmgr   UUID;
    v_u_alice    UUID; v_u_bob      UUID; v_u_carol    UUID;
    v_u_david    UUID; v_u_eve      UUID; v_u_frank    UUID;
    v_u_grace    UUID; v_u_henry    UUID; v_u_iris     UUID;
    v_u_jack     UUID;
    -- extra 10
    v_u_liam     UUID; v_u_mia      UUID; v_u_noah     UUID;
    v_u_sofia    UUID; v_u_ethan    UUID; v_u_chloe    UUID;
    v_u_lucas    UUID; v_u_amara    UUID; v_u_ryan     UUID;
    v_u_zoe      UUID;

    -- Employee IDs
    v_e_001 UUID; v_e_002 UUID; v_e_003 UUID; v_e_004 UUID; v_e_005 UUID;
    v_e_006 UUID; v_e_007 UUID; v_e_008 UUID; v_e_009 UUID; v_e_010 UUID;
    v_e_011 UUID; v_e_012 UUID; v_e_013 UUID; v_e_014 UUID; v_e_015 UUID;
    v_e_016 UUID; v_e_017 UUID; v_e_018 UUID; v_e_019 UUID; v_e_020 UUID;
    v_e_021 UUID; v_e_022 UUID; v_e_023 UUID; v_e_024 UUID; v_e_025 UUID;

    -- Leave type IDs
    v_lt_al UUID; v_lt_sl UUID; v_lt_pl UUID;
    v_lt_ul UUID; v_lt_cl UUID; v_lt_ml UUID;

BEGIN

-- ------------------------------------------------------------------ roles
SELECT role_id INTO v_role_super_admin FROM roles WHERE name = 'super_admin' LIMIT 1;
SELECT role_id INTO v_role_hr_manager  FROM roles WHERE name = 'hr_manager'  LIMIT 1;
SELECT role_id INTO v_role_manager     FROM roles WHERE name = 'manager'     LIMIT 1;
SELECT role_id INTO v_role_employee    FROM roles WHERE name = 'employee'    LIMIT 1;

IF v_role_super_admin IS NULL THEN
    RAISE EXCEPTION 'Roles not found — run API first to initialise roles';
END IF;

-- Get admin user for FK references (created by API startup or migration 15)
SELECT user_id INTO v_admin_id FROM users WHERE email = 'admin@hr-system.com' LIMIT 1;
IF v_admin_id IS NULL THEN
    RAISE EXCEPTION 'Admin user not found — ensure migration 15 has run';
END IF;

-- ------------------------------------------------------------------ departments
INSERT INTO departments (name, code, description) VALUES
    ('Engineering',       'ENG', 'Software engineering and product development'),
    ('Human Resources',   'HR',  'People operations and HR management'),
    ('Finance',           'FIN', 'Finance and accounting'),
    ('Sales & Marketing', 'SAL', 'Revenue generation and brand marketing'),
    ('Operations',        'OPS', 'Business operations and logistics')
ON CONFLICT (code) DO NOTHING;

SELECT id INTO v_dept_eng FROM departments WHERE code = 'ENG';
SELECT id INTO v_dept_hr  FROM departments WHERE code = 'HR';
SELECT id INTO v_dept_fin FROM departments WHERE code = 'FIN';
SELECT id INTO v_dept_sal FROM departments WHERE code = 'SAL';
SELECT id INTO v_dept_ops FROM departments WHERE code = 'OPS';

INSERT INTO departments (name, code, description, parent_department_id) VALUES
    ('Backend',              'ENG-BE',  'Backend services and APIs',        v_dept_eng),
    ('Frontend',             'ENG-FE',  'Web and mobile frontend',           v_dept_eng),
    ('DevOps',               'ENG-DO',  'Infrastructure and CI/CD',          v_dept_eng),
    ('Talent Acquisition',   'HR-TA',   'Recruiting and onboarding',         v_dept_hr),
    ('Payroll',              'HR-PAY',  'Payroll processing',                v_dept_hr),
    ('Accounting',           'FIN-ACC', 'Bookkeeping and accounting',        v_dept_fin),
    ('Treasury',             'FIN-TRE', 'Cash management and treasury',      v_dept_fin),
    ('Sales',                'SAL-SLS', 'Direct and indirect sales',         v_dept_sal),
    ('Marketing',            'SAL-MKT', 'Brand and digital marketing',       v_dept_sal)
ON CONFLICT (code) DO NOTHING;

SELECT id INTO v_dept_be  FROM departments WHERE code = 'ENG-BE';
SELECT id INTO v_dept_fe  FROM departments WHERE code = 'ENG-FE';
SELECT id INTO v_dept_do  FROM departments WHERE code = 'ENG-DO';
SELECT id INTO v_dept_ta  FROM departments WHERE code = 'HR-TA';
SELECT id INTO v_dept_pay FROM departments WHERE code = 'HR-PAY';
SELECT id INTO v_dept_acc FROM departments WHERE code = 'FIN-ACC';
SELECT id INTO v_dept_tre FROM departments WHERE code = 'FIN-TRE';
SELECT id INTO v_dept_sls FROM departments WHERE code = 'SAL-SLS';
SELECT id INTO v_dept_mkt FROM departments WHERE code = 'SAL-MKT';

-- ------------------------------------------------------------------ positions
INSERT INTO positions (title, code, department_id, grade_level, base_salary, housing_allowance, transport_allowance, medical_allowance) VALUES
    ('Software Engineer I',           'SWE-1',   v_dept_be,  'L1',  50000, 5000, 3000, 2000),
    ('Software Engineer II',          'SWE-2',   v_dept_be,  'L2',  70000, 7000, 3000, 2500),
    ('Senior Software Engineer',      'SWE-SR',  v_dept_be,  'L3',  95000, 9500, 4000, 3000),
    ('Staff Engineer',                'SWE-ST',  v_dept_be,  'L4', 130000,13000, 5000, 4000),
    ('Frontend Engineer',             'FE-1',    v_dept_fe,  'L2',  65000, 6500, 3000, 2500),
    ('Senior Frontend Engineer',      'FE-SR',   v_dept_fe,  'L3',  90000, 9000, 4000, 3000),
    ('DevOps Engineer',               'DO-1',    v_dept_do,  'L2',  70000, 7000, 3000, 2500),
    ('Engineering Manager',           'ENG-MGR', v_dept_eng, 'M1', 140000,14000, 6000, 5000),
    ('HR Generalist',                 'HR-GEN',  v_dept_ta,  'L2',  45000, 4500, 2500, 2000),
    ('Talent Acquisition Specialist', 'HR-TAS',  v_dept_ta,  'L2',  50000, 5000, 2500, 2000),
    ('HR Manager',                    'HR-MGR',  v_dept_hr,  'M1',  80000, 8000, 4000, 3000),
    ('Payroll Specialist',            'PAY-SP',  v_dept_pay, 'L2',  50000, 5000, 2500, 2000),
    ('Accountant',                    'FIN-ACC1',v_dept_acc, 'L2',  55000, 5500, 3000, 2000),
    ('Senior Accountant',             'FIN-ACC2',v_dept_acc, 'L3',  75000, 7500, 3500, 2500),
    ('Finance Manager',               'FIN-MGR', v_dept_fin, 'M1', 100000,10000, 5000, 3500),
    ('Sales Representative',          'SAL-REP', v_dept_sls, 'L1',  40000, 4000, 2500, 1500),
    ('Account Executive',             'SAL-AE',  v_dept_sls, 'L2',  60000, 6000, 3000, 2000),
    ('Marketing Specialist',          'MKT-SP',  v_dept_mkt, 'L2',  50000, 5000, 2500, 2000),
    ('Sales Manager',                 'SAL-MGR', v_dept_sal, 'M1',  90000, 9000, 4500, 3000),
    ('Operations Analyst',            'OPS-AN',  v_dept_ops, 'L2',  50000, 5000, 2500, 2000),
    ('Operations Manager',            'OPS-MGR', v_dept_ops, 'M1',  85000, 8500, 4000, 3000)
ON CONFLICT (code) DO NOTHING;

SELECT id INTO v_pos_swe1   FROM positions WHERE code = 'SWE-1';
SELECT id INTO v_pos_swe2   FROM positions WHERE code = 'SWE-2';
SELECT id INTO v_pos_swesr  FROM positions WHERE code = 'SWE-SR';
SELECT id INTO v_pos_swest  FROM positions WHERE code = 'SWE-ST';
SELECT id INTO v_pos_fe1    FROM positions WHERE code = 'FE-1';
SELECT id INTO v_pos_fesr   FROM positions WHERE code = 'FE-SR';
SELECT id INTO v_pos_do1    FROM positions WHERE code = 'DO-1';
SELECT id INTO v_pos_engmgr FROM positions WHERE code = 'ENG-MGR';
SELECT id INTO v_pos_hrgen  FROM positions WHERE code = 'HR-GEN';
SELECT id INTO v_pos_hrtas  FROM positions WHERE code = 'HR-TAS';
SELECT id INTO v_pos_hrmgr  FROM positions WHERE code = 'HR-MGR';
SELECT id INTO v_pos_paysp  FROM positions WHERE code = 'PAY-SP';
SELECT id INTO v_pos_acc1   FROM positions WHERE code = 'FIN-ACC1';
SELECT id INTO v_pos_acc2   FROM positions WHERE code = 'FIN-ACC2';
SELECT id INTO v_pos_finmgr FROM positions WHERE code = 'FIN-MGR';
SELECT id INTO v_pos_salrep FROM positions WHERE code = 'SAL-REP';
SELECT id INTO v_pos_salae  FROM positions WHERE code = 'SAL-AE';
SELECT id INTO v_pos_mktsp  FROM positions WHERE code = 'MKT-SP';
SELECT id INTO v_pos_salmgr FROM positions WHERE code = 'SAL-MGR';
SELECT id INTO v_pos_opsan  FROM positions WHERE code = 'OPS-AN';
SELECT id INTO v_pos_opsmgr FROM positions WHERE code = 'OPS-MGR';

-- ------------------------------------------------------------------ users
-- Passwords are bcrypt of the plaintext shown in comments (cost 10)
INSERT INTO users (email, password, role_id, is_active) VALUES
    ('hr.manager@hr-system.com',    '$2a$10$uHyip8mIKwtL3IjbFyKP4eZQm.ATOK66a1QWk/yEa9q4qoBIykNam', v_role_hr_manager, true),  -- HrManager@123
    ('eng.manager@hr-system.com',   '$2a$10$CG7B/dlg5E89hqYHkAMsNOOoOyHNMMxfIo8ILdqM.G21Kyx9p.cRi', v_role_manager,    true),  -- Manager@123
    ('fin.manager@hr-system.com',   '$2a$10$CG7B/dlg5E89hqYHkAMsNOOoOyHNMMxfIo8ILdqM.G21Kyx9p.cRi', v_role_manager,    true),
    ('sales.manager@hr-system.com', '$2a$10$CG7B/dlg5E89hqYHkAMsNOOoOyHNMMxfIo8ILdqM.G21Kyx9p.cRi', v_role_manager,    true),
    ('ops.manager@hr-system.com',   '$2a$10$CG7B/dlg5E89hqYHkAMsNOOoOyHNMMxfIo8ILdqM.G21Kyx9p.cRi', v_role_manager,    true),
    ('alice.smith@hr-system.com',   '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),  -- Employee@123
    ('bob.jones@hr-system.com',     '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('carol.white@hr-system.com',   '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('david.brown@hr-system.com',   '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('eve.davis@hr-system.com',     '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('frank.miller@hr-system.com',  '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('grace.wilson@hr-system.com',  '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('henry.moore@hr-system.com',   '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('iris.taylor@hr-system.com',   '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('jack.anderson@hr-system.com', '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('liam.chen@hr-system.com',     '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),  -- Employee@123
    ('mia.patel@hr-system.com',     '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('noah.garcia@hr-system.com',   '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('sofia.lee@hr-system.com',     '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('ethan.kim@hr-system.com',     '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('chloe.nguyen@hr-system.com',  '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('lucas.martin@hr-system.com',  '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('amara.obi@hr-system.com',     '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('ryan.scott@hr-system.com',    '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true),
    ('zoe.walker@hr-system.com',    '$2a$10$Rt69BXys5MUXnUj11Loh/.3/0c38bWMvynonW8K7ubac5i3FU6T/W', v_role_employee,   true)
ON CONFLICT (email) DO NOTHING;

SELECT user_id INTO v_u_hrmgr  FROM users WHERE email = 'hr.manager@hr-system.com';
SELECT user_id INTO v_u_engmgr FROM users WHERE email = 'eng.manager@hr-system.com';
SELECT user_id INTO v_u_finmgr FROM users WHERE email = 'fin.manager@hr-system.com';
SELECT user_id INTO v_u_salmgr FROM users WHERE email = 'sales.manager@hr-system.com';
SELECT user_id INTO v_u_opsmgr FROM users WHERE email = 'ops.manager@hr-system.com';
SELECT user_id INTO v_u_alice  FROM users WHERE email = 'alice.smith@hr-system.com';
SELECT user_id INTO v_u_bob    FROM users WHERE email = 'bob.jones@hr-system.com';
SELECT user_id INTO v_u_carol  FROM users WHERE email = 'carol.white@hr-system.com';
SELECT user_id INTO v_u_david  FROM users WHERE email = 'david.brown@hr-system.com';
SELECT user_id INTO v_u_eve    FROM users WHERE email = 'eve.davis@hr-system.com';
SELECT user_id INTO v_u_frank  FROM users WHERE email = 'frank.miller@hr-system.com';
SELECT user_id INTO v_u_grace  FROM users WHERE email = 'grace.wilson@hr-system.com';
SELECT user_id INTO v_u_henry  FROM users WHERE email = 'henry.moore@hr-system.com';
SELECT user_id INTO v_u_iris   FROM users WHERE email = 'iris.taylor@hr-system.com';
SELECT user_id INTO v_u_jack   FROM users WHERE email = 'jack.anderson@hr-system.com';
SELECT user_id INTO v_u_liam   FROM users WHERE email = 'liam.chen@hr-system.com';
SELECT user_id INTO v_u_mia    FROM users WHERE email = 'mia.patel@hr-system.com';
SELECT user_id INTO v_u_noah   FROM users WHERE email = 'noah.garcia@hr-system.com';
SELECT user_id INTO v_u_sofia  FROM users WHERE email = 'sofia.lee@hr-system.com';
SELECT user_id INTO v_u_ethan  FROM users WHERE email = 'ethan.kim@hr-system.com';
SELECT user_id INTO v_u_chloe  FROM users WHERE email = 'chloe.nguyen@hr-system.com';
SELECT user_id INTO v_u_lucas  FROM users WHERE email = 'lucas.martin@hr-system.com';
SELECT user_id INTO v_u_amara  FROM users WHERE email = 'amara.obi@hr-system.com';
SELECT user_id INTO v_u_ryan   FROM users WHERE email = 'ryan.scott@hr-system.com';
SELECT user_id INTO v_u_zoe    FROM users WHERE email = 'zoe.walker@hr-system.com';

-- ------------------------------------------------------------------ employees (managers first)
INSERT INTO employees (employee_number, user_id, first_name, last_name, email, phone, date_of_birth, gender, national_id, city, country, department_id, position_id, manager_id, hire_date, employment_type, employment_status)
VALUES
    ('EMP-001', v_u_hrmgr,  'Sarah',   'Connor',  'hr.manager.emp@hr-system.com',    '+12025550101', '1982-03-14', 'female', 'NAT-001', 'New York',      'USA', v_dept_hr,  v_pos_hrmgr,  NULL, '2025-12-01', 'full_time', 'active'),
    ('EMP-002', v_u_engmgr, 'James',   'Kirk',    'eng.manager.emp@hr-system.com',   '+12025550102', '1980-07-22', 'male',   'NAT-002', 'San Francisco', 'USA', v_dept_eng, v_pos_engmgr, NULL, '2025-12-08', 'full_time', 'active'),
    ('EMP-003', v_u_finmgr, 'Laura',   'Palmer',  'fin.manager.emp@hr-system.com',   '+12025550103', '1979-11-05', 'female', 'NAT-003', 'Chicago',       'USA', v_dept_fin, v_pos_finmgr, NULL, '2025-12-15', 'full_time', 'active'),
    ('EMP-004', v_u_salmgr, 'Tony',    'Stark',   'sales.manager.emp@hr-system.com', '+12025550104', '1984-05-29', 'male',   'NAT-004', 'Austin',        'USA', v_dept_sal, v_pos_salmgr, NULL, '2025-12-22', 'full_time', 'active'),
    ('EMP-005', v_u_opsmgr, 'Natasha', 'Romanova','ops.manager.emp@hr-system.com',   '+12025550105', '1986-09-18', 'female', 'NAT-005', 'Boston',        'USA', v_dept_ops, v_pos_opsmgr, NULL, '2026-01-05', 'full_time', 'active')
ON CONFLICT (employee_number) DO NOTHING;

SELECT id INTO v_e_001 FROM employees WHERE employee_number = 'EMP-001';
SELECT id INTO v_e_002 FROM employees WHERE employee_number = 'EMP-002';
SELECT id INTO v_e_003 FROM employees WHERE employee_number = 'EMP-003';
SELECT id INTO v_e_004 FROM employees WHERE employee_number = 'EMP-004';
SELECT id INTO v_e_005 FROM employees WHERE employee_number = 'EMP-005';

-- Regular employees
INSERT INTO employees (employee_number, user_id, first_name, last_name, email, phone, date_of_birth, gender, national_id, city, country, department_id, position_id, manager_id, hire_date, employment_type, employment_status)
VALUES
    ('EMP-006', v_u_alice, 'Alice',  'Smith',    'alice.smith.emp@hr-system.com',   '+12025550106', '1994-02-14', 'female', 'NAT-006', 'New York',      'USA', v_dept_be,  v_pos_swe2,  v_e_002, '2026-01-12', 'full_time', 'active'),
    ('EMP-007', v_u_bob,   'Bob',    'Jones',    'bob.jones.emp@hr-system.com',     '+12025550107', '1992-06-30', 'male',   'NAT-007', 'New York',      'USA', v_dept_be,  v_pos_swesr, v_e_002, '2026-01-19', 'full_time', 'active'),
    ('EMP-008', v_u_carol, 'Carol',  'White',    'carol.white.emp@hr-system.com',   '+12025550108', '1996-10-03', 'female', 'NAT-008', 'Remote',        'USA', v_dept_be,  v_pos_swe1,  v_e_002, '2026-01-26', 'full_time', 'active'),
    ('EMP-009', v_u_david, 'David',  'Brown',    'david.brown.emp@hr-system.com',   '+12025550109', '1993-04-20', 'male',   'NAT-009', 'San Francisco', 'USA', v_dept_fe,  v_pos_fesr,  v_e_002, '2026-02-02', 'full_time', 'active'),
    ('EMP-010', v_u_eve,   'Eve',    'Davis',    'eve.davis.emp@hr-system.com',     '+12025550110', '1997-08-08', 'female', 'NAT-010', 'Seattle',       'USA', v_dept_fe,  v_pos_fe1,   v_e_002, '2026-02-09', 'full_time', 'active'),
    ('EMP-011', v_u_frank, 'Frank',  'Miller',   'frank.miller.emp@hr-system.com',  '+12025550111', '1990-12-25', 'male',   'NAT-011', 'New York',      'USA', v_dept_ta,  v_pos_hrtas, v_e_001, '2025-12-29', 'full_time', 'active'),
    ('EMP-012', v_u_grace, 'Grace',  'Wilson',   'grace.wilson.emp@hr-system.com',  '+12025550112', '1995-03-17', 'female', 'NAT-012', 'Chicago',       'USA', v_dept_pay, v_pos_paysp, v_e_001, '2026-02-16', 'full_time', 'active'),
    ('EMP-013', v_u_henry, 'Henry',  'Moore',    'henry.moore.emp@hr-system.com',   '+12025550113', '1988-07-04', 'male',   'NAT-013', 'Chicago',       'USA', v_dept_acc, v_pos_acc2,  v_e_003, '2026-02-23', 'full_time', 'active'),
    ('EMP-014', v_u_iris,  'Iris',   'Taylor',   'iris.taylor.emp@hr-system.com',   '+12025550114', '1998-01-29', 'female', 'NAT-014', 'Austin',        'USA', v_dept_sls, v_pos_salae, v_e_004, '2026-03-02', 'full_time', 'active'),
    ('EMP-015', v_u_jack,  'Jack',   'Anderson', 'jack.anderson.emp@hr-system.com', '+12025550115', '1991-05-11', 'male',   'NAT-015', 'Dallas',        'USA', v_dept_mkt, v_pos_mktsp, v_e_004, '2026-03-09', 'full_time', 'on_leave')
ON CONFLICT (employee_number) DO NOTHING;

SELECT id INTO v_e_006 FROM employees WHERE employee_number = 'EMP-006';
SELECT id INTO v_e_007 FROM employees WHERE employee_number = 'EMP-007';
SELECT id INTO v_e_008 FROM employees WHERE employee_number = 'EMP-008';
SELECT id INTO v_e_009 FROM employees WHERE employee_number = 'EMP-009';
SELECT id INTO v_e_010 FROM employees WHERE employee_number = 'EMP-010';
SELECT id INTO v_e_011 FROM employees WHERE employee_number = 'EMP-011';
SELECT id INTO v_e_012 FROM employees WHERE employee_number = 'EMP-012';
SELECT id INTO v_e_013 FROM employees WHERE employee_number = 'EMP-013';
SELECT id INTO v_e_014 FROM employees WHERE employee_number = 'EMP-014';
SELECT id INTO v_e_015 FROM employees WHERE employee_number = 'EMP-015';

-- Extra 10 employees — hire dates clustered Dec 2025 – Feb 2026
INSERT INTO employees (employee_number, user_id, first_name, last_name, email, phone, date_of_birth, gender, national_id, city, country, department_id, position_id, manager_id, hire_date, employment_type, employment_status)
VALUES
    ('EMP-016', v_u_liam,  'Liam',  'Chen',    'liam.chen.emp@hr-system.com',    '+12025550116', '1995-04-10', 'male',   'NAT-016', 'New York',      'USA', v_dept_be,  v_pos_swe1,  v_e_002, '2025-12-02', 'full_time', 'active'),
    ('EMP-017', v_u_mia,   'Mia',   'Patel',   'mia.patel.emp@hr-system.com',    '+12025550117', '1997-07-23', 'female', 'NAT-017', 'Chicago',       'USA', v_dept_ta,  v_pos_hrgen, v_e_001, '2025-12-09', 'full_time', 'active'),
    ('EMP-018', v_u_noah,  'Noah',  'Garcia',  'noah.garcia.emp@hr-system.com',  '+12025550118', '1993-11-15', 'male',   'NAT-018', 'Austin',        'USA', v_dept_sls, v_pos_salrep,v_e_004, '2025-12-16', 'full_time', 'active'),
    ('EMP-019', v_u_sofia, 'Sofia', 'Lee',     'sofia.lee.emp@hr-system.com',    '+12025550119', '1996-02-28', 'female', 'NAT-019', 'San Francisco', 'USA', v_dept_fe,  v_pos_fe1,   v_e_002, '2025-12-23', 'full_time', 'active'),
    ('EMP-020', v_u_ethan, 'Ethan', 'Kim',     'ethan.kim.emp@hr-system.com',    '+12025550120', '1991-09-05', 'male',   'NAT-020', 'Seattle',       'USA', v_dept_do,  v_pos_do1,   v_e_002, '2025-12-30', 'full_time', 'active'),
    ('EMP-021', v_u_chloe, 'Chloe', 'Nguyen',  'chloe.nguyen.emp@hr-system.com', '+12025550121', '1998-05-17', 'female', 'NAT-021', 'Boston',        'USA', v_dept_acc, v_pos_acc1,  v_e_003, '2026-01-06', 'full_time', 'active'),
    ('EMP-022', v_u_lucas, 'Lucas', 'Martin',  'lucas.martin.emp@hr-system.com', '+12025550122', '1994-08-30', 'male',   'NAT-022', 'Dallas',        'USA', v_dept_mkt, v_pos_mktsp, v_e_004, '2026-01-13', 'full_time', 'active'),
    ('EMP-023', v_u_amara, 'Amara', 'Obi',     'amara.obi.emp@hr-system.com',    '+12025550123', '1992-12-04', 'female', 'NAT-023', 'New York',      'USA', v_dept_pay, v_pos_paysp, v_e_001, '2026-01-20', 'full_time', 'active'),
    ('EMP-024', v_u_ryan,  'Ryan',  'Scott',   'ryan.scott.emp@hr-system.com',   '+12025550124', '1990-03-19', 'male',   'NAT-024', 'Chicago',       'USA', v_dept_ops, v_pos_opsan, v_e_005, '2026-01-27', 'full_time', 'active'),
    ('EMP-025', v_u_zoe,   'Zoe',   'Walker',  'zoe.walker.emp@hr-system.com',   '+12025550125', '1999-06-11', 'female', 'NAT-025', 'Miami',         'USA', v_dept_be,  v_pos_swe2,  v_e_002, '2026-02-03', 'full_time', 'active')
ON CONFLICT (employee_number) DO NOTHING;

SELECT id INTO v_e_016 FROM employees WHERE employee_number = 'EMP-016';
SELECT id INTO v_e_017 FROM employees WHERE employee_number = 'EMP-017';
SELECT id INTO v_e_018 FROM employees WHERE employee_number = 'EMP-018';
SELECT id INTO v_e_019 FROM employees WHERE employee_number = 'EMP-019';
SELECT id INTO v_e_020 FROM employees WHERE employee_number = 'EMP-020';
SELECT id INTO v_e_021 FROM employees WHERE employee_number = 'EMP-021';
SELECT id INTO v_e_022 FROM employees WHERE employee_number = 'EMP-022';
SELECT id INTO v_e_023 FROM employees WHERE employee_number = 'EMP-023';
SELECT id INTO v_e_024 FROM employees WHERE employee_number = 'EMP-024';
SELECT id INTO v_e_025 FROM employees WHERE employee_number = 'EMP-025';

-- ------------------------------------------------------------------ emergency contacts
INSERT INTO emergency_contacts (employee_id, name, relationship, phone, email) VALUES
    (v_e_006, 'Michael Smith',    'Spouse',  '+12025550201', 'michael.smith@gmail.com'),
    (v_e_007, 'Linda Jones',      'Mother',  '+12025550202', 'linda.jones@gmail.com'),
    (v_e_008, 'Peter White',      'Father',  '+12025550203', 'peter.white@gmail.com'),
    (v_e_009, 'Emma Brown',       'Spouse',  '+12025550204', 'emma.brown@gmail.com'),
    (v_e_010, 'Chris Davis',      'Brother', '+12025550205', 'chris.davis@gmail.com'),
    (v_e_011, 'Anna Miller',      'Spouse',  '+12025550206', 'anna.miller@gmail.com'),
    (v_e_012, 'Tom Wilson',       'Father',  '+12025550207', 'tom.wilson@gmail.com'),
    (v_e_013, 'Sandra Moore',     'Spouse',  '+12025550208', 'sandra.moore@gmail.com'),
    (v_e_014, 'Kevin Taylor',     'Brother', '+12025550209', 'kevin.taylor@gmail.com'),
    (v_e_015, 'Olivia Anderson',  'Spouse',  '+12025550210', 'olivia.anderson@gmail.com'),
    (v_e_016, 'Wei Chen',         'Father',  '+12025550211', 'wei.chen@gmail.com'),
    (v_e_017, 'Raj Patel',        'Spouse',  '+12025550212', 'raj.patel@gmail.com'),
    (v_e_018, 'Maria Garcia',     'Mother',  '+12025550213', 'maria.garcia@gmail.com'),
    (v_e_019, 'James Lee',        'Brother', '+12025550214', 'james.lee@gmail.com'),
    (v_e_020, 'Jenny Kim',        'Spouse',  '+12025550215', 'jenny.kim@gmail.com'),
    (v_e_021, 'Tuan Nguyen',      'Father',  '+12025550216', 'tuan.nguyen@gmail.com'),
    (v_e_022, 'Claire Martin',    'Spouse',  '+12025550217', 'claire.martin@gmail.com'),
    (v_e_023, 'Chidi Obi',        'Brother', '+12025550218', 'chidi.obi@gmail.com'),
    (v_e_024, 'Emma Scott',       'Spouse',  '+12025550219', 'emma.scott@gmail.com'),
    (v_e_025, 'Daniel Walker',    'Father',  '+12025550220', 'daniel.walker@gmail.com');

-- ------------------------------------------------------------------ employee documents
INSERT INTO employee_documents (employee_id, document_type, title, file_url, file_name, file_size, mime_type, uploaded_by, is_verified, verified_by, verified_at) VALUES
    (v_e_006, 'contract',      'Employment Contract 2026', 'https://storage.hr-system.local/docs/contract_alice_2026.pdf',  'contract_alice_2026.pdf',  102400, 'application/pdf', v_admin_id, true,  v_admin_id, NOW() - INTERVAL '1 day'),
    (v_e_006, 'id_document',   'National ID',              'https://storage.hr-system.local/docs/id_alice.pdf',            'id_alice.pdf',             102400, 'application/pdf', v_admin_id, true,  v_admin_id, NOW() - INTERVAL '1 day'),
    (v_e_007, 'contract',      'Employment Contract 2026', 'https://storage.hr-system.local/docs/contract_bob_2026.pdf',   'contract_bob_2026.pdf',    102400, 'application/pdf', v_admin_id, true,  v_admin_id, NOW() - INTERVAL '1 day'),
    (v_e_007, 'certification', 'AWS Solutions Architect',  'https://storage.hr-system.local/docs/aws_cert_bob.pdf',        'aws_cert_bob.pdf',         102400, 'application/pdf', v_admin_id, true,  v_admin_id, NOW() - INTERVAL '1 day'),
    (v_e_008, 'contract',      'Employment Contract 2026', 'https://storage.hr-system.local/docs/contract_carol_2026.pdf', 'contract_carol_2026.pdf',  102400, 'application/pdf', v_admin_id, false, NULL,       NULL),
    (v_e_009, 'contract',      'Employment Contract 2026', 'https://storage.hr-system.local/docs/contract_david_2026.pdf', 'contract_david_2026.pdf',  102400, 'application/pdf', v_admin_id, true,  v_admin_id, NOW() - INTERVAL '1 day'),
    (v_e_010, 'contract',      'Employment Contract 2026', 'https://storage.hr-system.local/docs/contract_eve_2026.pdf',   'contract_eve_2026.pdf',    102400, 'application/pdf', v_admin_id, false, NULL,       NULL),
    (v_e_011, 'contract',      'Employment Contract 2025', 'https://storage.hr-system.local/docs/contract_frank_2025.pdf', 'contract_frank_2025.pdf',  102400, 'application/pdf', v_admin_id, true,  v_admin_id, NOW() - INTERVAL '1 day'),
    (v_e_013, 'contract',      'Employment Contract 2026', 'https://storage.hr-system.local/docs/contract_henry_2026.pdf', 'contract_henry_2026.pdf',  102400, 'application/pdf', v_admin_id, true,  v_admin_id, NOW() - INTERVAL '1 day');

-- ------------------------------------------------------------------ leave types
INSERT INTO leave_types (code, name, description, default_days_per_year, is_paid, is_carry_forward_allowed, max_carry_forward_days, requires_approval, requires_document, is_active)
VALUES
    ('AL', 'Annual Leave',       '', 24, true,  true,  5, true, false, true),
    ('SL', 'Sick Leave',         '', 15, true,  false, 0, true, true,  true),
    ('PL', 'Parental Leave',     '', 90, true,  false, 0, true, false, true),
    ('UL', 'Unpaid Leave',       '',  0, false, false, 0, true, false, true),
    ('CL', 'Compassionate Leave','',  5, true,  false, 0, true, false, true),
    ('ML', 'Marriage Leave',     '',  3, true,  false, 0, true, false, true)
ON CONFLICT (code) DO NOTHING;

SELECT id INTO v_lt_al FROM leave_types WHERE code = 'AL';
SELECT id INTO v_lt_sl FROM leave_types WHERE code = 'SL';
SELECT id INTO v_lt_pl FROM leave_types WHERE code = 'PL';
SELECT id INTO v_lt_ul FROM leave_types WHERE code = 'UL';
SELECT id INTO v_lt_cl FROM leave_types WHERE code = 'CL';
SELECT id INTO v_lt_ml FROM leave_types WHERE code = 'ML';

-- ------------------------------------------------------------------ leave balances (2026)
INSERT INTO leave_balances (employee_id, leave_type_id, year, total_entitled, carried_forward, earned_leave_days, used, pending)
SELECT e.id, lt.id, 2026, lt.default_days_per_year, 0, 0, 0, 0
FROM employees e
CROSS JOIN leave_types lt
WHERE e.employee_number IN ('EMP-001','EMP-002','EMP-003','EMP-004','EMP-005',
                             'EMP-006','EMP-007','EMP-008','EMP-009','EMP-010',
                             'EMP-011','EMP-012','EMP-013','EMP-014','EMP-015',
                             'EMP-016','EMP-017','EMP-018','EMP-019','EMP-020',
                             'EMP-021','EMP-022','EMP-023','EMP-024','EMP-025')
ON CONFLICT (employee_id, leave_type_id, year) DO NOTHING;

RAISE NOTICE 'Demo data seeded successfully';

END $$;
