-- Drop triggers
DROP TRIGGER IF EXISTS update_sprints_updated_at ON sprints;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS sprint_issues;
DROP TABLE IF EXISTS sprints;
