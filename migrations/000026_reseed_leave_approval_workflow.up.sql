-- Remove the workflow seeded by migration 15 and replace it with a corrected version.
-- The original had a wrong description ("equipment procurement requests") and
-- inconsistent step/role naming. This migration deletes it and reseeds cleanly.

DO $$
DECLARE
    v_workflow_id     UUID;
    v_super_admin_id  UUID;
    v_step_submit_id  UUID;
    v_step_review_id  UUID;
    v_step_approve_id UUID;
BEGIN
    -- Delete instances referencing the old workflow first (assigned_tasks and workflow_history
    -- cascade automatically from workflow_instances ON DELETE CASCADE)
    DELETE FROM workflow_instances
    WHERE workflow_id IN (
        SELECT id FROM workflows WHERE name IN ('Leave Request Approval', 'Leave Request approval')
    );

    -- Delete the old workflow (cascades to workflow_steps and workflow_transitions)
    DELETE FROM workflows WHERE name IN ('Leave Request Approval', 'Leave Request approval');

    -- Get the super admin user ID
    SELECT user_id INTO v_super_admin_id
    FROM users WHERE email = 'admin@hr-system.com' LIMIT 1;

    IF v_super_admin_id IS NULL THEN
        RAISE EXCEPTION 'Super admin user not found. Cannot seed workflow.';
    END IF;

    -- Create the Leave Approval Workflow
    INSERT INTO workflows (id, name, description, is_active, created_by)
    VALUES (
        gen_random_uuid(),
        'Leave Request Approval',
        'Standard workflow for processing employee leave requests. The employee submits the request, HR reviews it for policy compliance, and the direct manager gives final approval.',
        true,
        v_super_admin_id
    )
    RETURNING id INTO v_workflow_id;

    -- Step 1: Submission (initial) — employee submits the leave request
    INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles, min_approvals)
    VALUES (
        gen_random_uuid(),
        v_workflow_id,
        'Submission',
        1,
        true,
        false,
        '["employee"]'::jsonb,
        1
    )
    RETURNING id INTO v_step_submit_id;

    -- Step 2: Review — HR manager reviews for policy compliance
    INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles, min_approvals)
    VALUES (
        gen_random_uuid(),
        v_workflow_id,
        'Review',
        2,
        false,
        false,
        '["hr_manager"]'::jsonb,
        1
    )
    RETURNING id INTO v_step_review_id;

    -- Step 3: Approve (final) — direct manager gives final approval
    INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles, min_approvals)
    VALUES (
        gen_random_uuid(),
        v_workflow_id,
        'Approve',
        3,
        false,
        true,
        '["manager"]'::jsonb,
        1
    )
    RETURNING id INTO v_step_approve_id;

    -- Transition: Submission → Review (employee submits)
    INSERT INTO workflow_transitions (workflow_id, from_step_id, to_step_id, action_name, condition_type)
    VALUES (v_workflow_id, v_step_submit_id, v_step_review_id, 'submit', 'always');

    -- Transition: Review → Approve (HR forwards to manager for approval)
    INSERT INTO workflow_transitions (workflow_id, from_step_id, to_step_id, action_name, condition_type)
    VALUES (v_workflow_id, v_step_review_id, v_step_approve_id, 'review', 'always');

    RAISE NOTICE 'Leave Request Approval workflow reseeded with ID: %', v_workflow_id;
END $$;
