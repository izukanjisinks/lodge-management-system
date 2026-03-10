-- Drop workflow tables in reverse order of creation
DROP TABLE IF EXISTS workflow_history;
DROP TABLE IF EXISTS assigned_tasks;
DROP TABLE IF EXISTS workflow_instances;
DROP TABLE IF EXISTS workflow_transitions;
DROP TABLE IF EXISTS workflow_steps;
DROP TABLE IF EXISTS workflows;
