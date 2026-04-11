-- =============================================================================
-- Workflow seed data for Lodge Management System
-- Seeds: 1 workflow (Booking Approval), 2 steps, 1 transition,
--        3 workflow instances (tied to bookings from migration 000020),
--        3 assigned tasks, workflow history entries
--
-- Depends on: migration 000012 (workflow tables), 000014 (dev users),
--             000020 (demo bookings)
-- Idempotent: all inserts use ON CONFLICT DO NOTHING
-- =============================================================================

DO $$
DECLARE
    -- Users
    v_admin_id        UUID;
    v_manager_id      UUID;
    v_receptionist_id UUID;

    -- Workflow
    v_workflow_id UUID;

    -- Steps
    v_step_submission UUID;
    v_step_approval   UUID;

    -- Transition
    v_transition_id UUID;

    -- Bookings (from migration 000020)
    v_booking_confirmed  UUID;  -- Peter / Suite 301 / confirmed
    v_booking_pending    UUID;  -- Acme Corp / Conference A / pending
    v_booking_checked_in UUID;  -- Sarah / Room 201 / checked_in

    -- Workflow instances
    v_instance_1 UUID;
    v_instance_2 UUID;
    v_instance_3 UUID;

BEGIN

-- ---------------------------------------------------------------------------
-- Resolve user IDs
-- ---------------------------------------------------------------------------
SELECT user_id INTO v_admin_id        FROM users WHERE email = 'admin@lodge.dev'        LIMIT 1;
SELECT user_id INTO v_manager_id      FROM users WHERE email = 'manager@lodge.dev'      LIMIT 1;
SELECT user_id INTO v_receptionist_id FROM users WHERE email = 'receptionist@lodge.dev' LIMIT 1;

IF v_admin_id IS NULL THEN
    RAISE NOTICE 'Dev users not found — skipping workflow seed data (migration 000021)';
    RETURN;
END IF;

-- ---------------------------------------------------------------------------
-- Resolve booking IDs from migration 000020
-- ---------------------------------------------------------------------------
-- confirmed booking: Peter Okafor / Suite 301
SELECT b.id INTO v_booking_confirmed
FROM bookings b
JOIN individual_profiles ip ON b.client_id = ip.id
WHERE ip.email = 'peter.okafor@gmail.com' AND b.status = 'confirmed'
LIMIT 1;

-- pending booking: Acme Corp / Conference A
SELECT b.id INTO v_booking_pending
FROM bookings b
JOIN corporate_profiles cp ON b.client_id = cp.id
WHERE cp.email = 'bookings@acme.co.zm' AND b.status = 'pending'
LIMIT 1;

-- checked_in booking: Sarah Banda / Room 201
SELECT b.id INTO v_booking_checked_in
FROM bookings b
JOIN individual_profiles ip ON b.client_id = ip.id
WHERE ip.email = 'sarah.banda@gmail.com' AND b.status = 'checked_in'
LIMIT 1;

-- ---------------------------------------------------------------------------
-- Workflow definition: Booking Approval Workflow
-- ---------------------------------------------------------------------------
INSERT INTO workflows (id, name, description, is_active, created_by)
VALUES (
    gen_random_uuid(),
    'Booking Approval Workflow',
    'Standard booking approval process for all guest reservations.',
    TRUE,
    v_admin_id
)
ON CONFLICT DO NOTHING;

SELECT id INTO v_workflow_id
FROM workflows
WHERE name = 'Booking Approval Workflow'
LIMIT 1;

-- ---------------------------------------------------------------------------
-- Steps
-- ---------------------------------------------------------------------------
INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles, requires_all_approvers, min_approvals)
VALUES
    (gen_random_uuid(), v_workflow_id, 'Submission', 1, TRUE,  FALSE, '["receptionist","admin"]'::jsonb, FALSE, 1),
    (gen_random_uuid(), v_workflow_id, 'Approval',   2, FALSE, TRUE,  '["manager","admin"]'::jsonb,      FALSE, 1)
ON CONFLICT (workflow_id, step_order) DO NOTHING;

SELECT id INTO v_step_submission FROM workflow_steps WHERE workflow_id = v_workflow_id AND step_order = 1 LIMIT 1;
SELECT id INTO v_step_approval   FROM workflow_steps WHERE workflow_id = v_workflow_id AND step_order = 2 LIMIT 1;

-- ---------------------------------------------------------------------------
-- Transition: Submission → Approval
-- ---------------------------------------------------------------------------
INSERT INTO workflow_transitions (id, workflow_id, from_step_id, to_step_id, action_name, condition_type, condition_value)
VALUES (
    gen_random_uuid(),
    v_workflow_id,
    v_step_submission,
    v_step_approval,
    'Submit for Approval',
    'always',
    ''
)
ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------------------
-- Workflow instances (one per demo booking that needs approval)
-- ---------------------------------------------------------------------------

-- Instance 1: Peter's confirmed booking — in_progress (moved past submission)
v_instance_1 := gen_random_uuid();
INSERT INTO workflow_instances (id, workflow_id, current_step_id, status, task_details, created_by, due_date)
SELECT
    v_instance_1,
    v_workflow_id,
    v_step_approval,
    'in_progress',
    jsonb_build_object(
        'booking_id',       v_booking_confirmed,
        'room_name',        'Suite 301',
        'client_name',      'Peter Okafor',
        'client_type',      'individual',
        'check_in',         CURRENT_DATE + 5,
        'check_out',        CURRENT_DATE + 10,
        'guests',           2,
        'meal_plan_name',   'Bed & Breakfast',
        'total_amount',     3306.00,
        'special_requests', 'Late check-in, arriving after 10 PM'
    ),
    v_receptionist_id,
    CURRENT_DATE + 3
WHERE v_booking_confirmed IS NOT NULL
ON CONFLICT DO NOTHING;

-- Instance 2: Acme Corp pending booking — pending (at submission step)
v_instance_2 := gen_random_uuid();
INSERT INTO workflow_instances (id, workflow_id, current_step_id, status, task_details, created_by, due_date)
SELECT
    v_instance_2,
    v_workflow_id,
    v_step_submission,
    'pending',
    jsonb_build_object(
        'booking_id',       v_booking_pending,
        'room_name',        'Conference A',
        'client_name',      'Acme Zambia Ltd',
        'client_type',      'corporate',
        'check_in',         CURRENT_DATE + 2,
        'check_out',        CURRENT_DATE + 3,
        'guests',           15,
        'meal_plan_name',   NULL,
        'total_amount',     800.00,
        'special_requests', 'AV setup required, catering for 15'
    ),
    v_receptionist_id,
    CURRENT_DATE + 1
WHERE v_booking_pending IS NOT NULL
ON CONFLICT DO NOTHING;

-- Instance 3: Sarah's checked_in booking — completed (fully approved)
v_instance_3 := gen_random_uuid();
INSERT INTO workflow_instances (id, workflow_id, current_step_id, status, task_details, created_by, due_date, completed_at)
SELECT
    v_instance_3,
    v_workflow_id,
    v_step_approval,
    'completed',
    jsonb_build_object(
        'booking_id',       v_booking_checked_in,
        'room_name',        'Room 201',
        'client_name',      'Sarah Banda',
        'client_type',      'individual',
        'check_in',         CURRENT_DATE - 1,
        'check_out',        CURRENT_DATE + 3,
        'guests',           2,
        'meal_plan_name',   'Half Board',
        'total_amount',     1670.40,
        'special_requests', 'Extra pillows requested'
    ),
    v_receptionist_id,
    CURRENT_DATE - 1,
    CURRENT_DATE - 1
WHERE v_booking_checked_in IS NOT NULL
ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------------------
-- Assigned tasks
-- ---------------------------------------------------------------------------

-- Task 1: Instance 1 (Peter) — approval step, assigned to manager, pending
INSERT INTO assigned_tasks (id, instance_id, step_id, step_name, assigned_to, assigned_by, status, due_date)
SELECT gen_random_uuid(), v_instance_1, v_step_approval, 'Approval', v_manager_id, v_receptionist_id, 'pending', CURRENT_DATE + 3
WHERE v_booking_confirmed IS NOT NULL
ON CONFLICT DO NOTHING;

-- Task 2: Instance 2 (Acme) — submission step, assigned to receptionist, pending
INSERT INTO assigned_tasks (id, instance_id, step_id, step_name, assigned_to, assigned_by, status, due_date)
SELECT gen_random_uuid(), v_instance_2, v_step_submission, 'Submission', v_receptionist_id, v_admin_id, 'pending', CURRENT_DATE + 1
WHERE v_booking_pending IS NOT NULL
ON CONFLICT DO NOTHING;

-- Task 3: Instance 3 (Sarah) — approval step, assigned to manager, completed
INSERT INTO assigned_tasks (id, instance_id, step_id, step_name, assigned_to, assigned_by, status, due_date, completed_at)
SELECT gen_random_uuid(), v_instance_3, v_step_approval, 'Approval', v_manager_id, v_receptionist_id, 'completed', CURRENT_DATE - 1, CURRENT_DATE - 1
WHERE v_booking_checked_in IS NOT NULL
ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------------------
-- Workflow history
-- ---------------------------------------------------------------------------

-- History: Instance 1 moved from Submission → Approval by receptionist
INSERT INTO workflow_history (id, instance_id, from_step_id, to_step_id, action_taken, performed_by, performed_by_name, comments)
SELECT gen_random_uuid(), v_instance_1, v_step_submission, v_step_approval,
       'Submit for Approval', v_receptionist_id, 'Receptionist', 'Booking submitted for manager approval'
WHERE v_booking_confirmed IS NOT NULL
ON CONFLICT DO NOTHING;

-- History: Instance 3 completed — approved by manager
INSERT INTO workflow_history (id, instance_id, from_step_id, to_step_id, action_taken, performed_by, performed_by_name, comments)
SELECT gen_random_uuid(), v_instance_3, v_step_submission, v_step_approval,
       'Submit for Approval', v_receptionist_id, 'Receptionist', NULL
WHERE v_booking_checked_in IS NOT NULL
ON CONFLICT DO NOTHING;

INSERT INTO workflow_history (id, instance_id, from_step_id, to_step_id, action_taken, performed_by, performed_by_name, comments)
SELECT gen_random_uuid(), v_instance_3, v_step_approval, v_step_approval,
       'Approve', v_manager_id, 'Manager', 'Booking approved. Guest already checked in.'
WHERE v_booking_checked_in IS NOT NULL
ON CONFLICT DO NOTHING;

RAISE NOTICE 'Workflow seed data inserted successfully';

END $$;
