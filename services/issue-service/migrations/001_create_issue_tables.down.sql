DROP TRIGGER IF EXISTS update_issue_custom_values_updated_at ON issue_custom_values;
DROP TRIGGER IF EXISTS update_custom_fields_updated_at ON custom_fields;
DROP TRIGGER IF EXISTS update_issues_updated_at ON issues;

DROP TABLE IF EXISTS issue_watchers;
DROP TABLE IF EXISTS issue_links;
DROP TABLE IF EXISTS issue_custom_values;
DROP TABLE IF EXISTS custom_fields;
DROP TABLE IF EXISTS issues;
DROP TABLE IF EXISTS project_counters;
