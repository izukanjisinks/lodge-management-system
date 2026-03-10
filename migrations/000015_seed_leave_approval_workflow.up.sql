-- Seed data for Leave Approval Workflow Template
-- This creates a standard workflow: Submit → HR Review → Manager Approval → Completed

DO $$
DECLARE
    v_workflow_id UUID;
    v_super_admin_id UUID;
    v_step_submit_id UUID;
    v_step_hr_review_id UUID;
    v_step_manager_approval_id UUID;
    v_step_completed_id UUID;
BEGIN
    -- Get the super admin user ID, inserting a seed user if it doesn't exist yet.
    -- The app will overwrite the password hash on startup via its own seeding logic.
    SELECT user_id INTO v_super_admin_id FROM users WHERE email = 'admin@hr-system.com' LIMIT 1;

    IF v_super_admin_id IS NULL THEN
        INSERT INTO users (email, password, role_id, is_active)
        VALUES (
            'admin@hr-system.com',
            '$2a$10$placeholder.hash.will.be.overwritten.by.app.seed',
            (SELECT role_id FROM roles WHERE name = 'super_admin' LIMIT 1),
            true
        )
        RETURNING user_id INTO v_super_admin_id;
        RAISE NOTICE 'Created seed admin user with ID: %', v_super_admin_id;
    END IF;

    -- Create the Leave Approval Workflow template
    INSERT INTO workflows (id, name, description, is_active, created_by)
    VALUES (
        gen_random_uuid(),
        'Leave Request Approval',
        'Standard workflow for employee leave request approvals. Requires HR review and manager approval.',
        true,
        v_super_admin_id
    )
    RETURNING id INTO v_workflow_id;

    RAISE NOTICE 'Created workflow with ID: %', v_workflow_id;

    -- Step 1: Submit (Initial step - Employee submits)
    INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles)
    VALUES (
        gen_random_uuid(),
        v_workflow_id,
        'Submit',
        1,
        true,
        false,
        '["EMPLOYEE", "HR_MANAGER", "DEPARTMENT_HEAD", "SUPER_ADMIN"]'::jsonb
    )
    RETURNING id INTO v_step_submit_id;

    -- Step 2: HR Review
    INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles, min_approvals)
    VALUES (
        gen_random_uuid(),
        v_workflow_id,
        'HR Review',
        2,
        false,
        false,
        '["HR_MANAGER", "SUPER_ADMIN"]'::jsonb,
        1
    )
    RETURNING id INTO v_step_hr_review_id;

    -- Step 3: Manager Approval
    INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles, min_approvals)
    VALUES (
        gen_random_uuid(),
        v_workflow_id,
        'Manager Approval',
        3,
        false,
        false,
        '["DEPARTMENT_HEAD", "SUPER_ADMIN"]'::jsonb,
        1
    )
    RETURNING id INTO v_step_manager_approval_id;

    -- Step 4: Completed (Final step)
    INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles)
    VALUES (
        gen_random_uuid(),
        v_workflow_id,
        'Completed',
        4,
        false,
        true,
        '[]'::jsonb
    )
    RETURNING id INTO v_step_completed_id;

    -- Create transitions

    -- Submit → HR Review (when employee submits)
    INSERT INTO workflow_transitions (workflow_id, from_step_id, to_step_id, action_name, condition_type)
    VALUES (v_workflow_id, v_step_submit_id, v_step_hr_review_id, 'submit', 'any_user');

    -- HR Review → Manager Approval (when HR approves)
    INSERT INTO workflow_transitions (workflow_id, from_step_id, to_step_id, action_name, condition_type)
    VALUES (v_workflow_id, v_step_hr_review_id, v_step_manager_approval_id, 'approve', 'role_check');

    -- HR Review → Submit (when HR rejects)
    INSERT INTO workflow_transitions (workflow_id, from_step_id, to_step_id, action_name, condition_type)
    VALUES (v_workflow_id, v_step_hr_review_id, v_step_submit_id, 'reject', 'role_check');

    -- Manager Approval → Completed (when manager approves)
    INSERT INTO workflow_transitions (workflow_id, from_step_id, to_step_id, action_name, condition_type)
    VALUES (v_workflow_id, v_step_manager_approval_id, v_step_completed_id, 'approve', 'role_check');

    -- Manager Approval → HR Review (when manager rejects)
    INSERT INTO workflow_transitions (workflow_id, from_step_id, to_step_id, action_name, condition_type)
    VALUES (v_workflow_id, v_step_manager_approval_id, v_step_hr_review_id, 'reject', 'role_check');

    RAISE NOTICE 'Leave Approval Workflow template created successfully!';
    RAISE NOTICE 'Workflow ID: %', v_workflow_id;
    RAISE NOTICE 'Submit Step ID: %', v_step_submit_id;
    RAISE NOTICE 'HR Review Step ID: %', v_step_hr_review_id;
    RAISE NOTICE 'Manager Approval Step ID: %', v_step_manager_approval_id;
    RAISE NOTICE 'Completed Step ID: %', v_step_completed_id;

END $$;
