DROP TRIGGER IF EXISTS update_workflow_transitions_updated_at ON workflow_transitions;
DROP TRIGGER IF EXISTS update_workflow_statuses_updated_at ON workflow_statuses;
DROP TRIGGER IF EXISTS update_workflows_updated_at ON workflows;

DROP TABLE IF EXISTS workflow_transitions;
DROP TABLE IF EXISTS workflow_statuses;
DROP TABLE IF EXISTS workflows;
