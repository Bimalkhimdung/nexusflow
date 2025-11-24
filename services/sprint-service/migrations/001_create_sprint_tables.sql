-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Sprints table
CREATE TABLE IF NOT EXISTS sprints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    name TEXT NOT NULL,
    goal TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'planned',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sprints_project_id ON sprints(project_id);
CREATE INDEX IF NOT EXISTS idx_sprints_status ON sprints(status);

-- Sprint Issues table
CREATE TABLE IF NOT EXISTS sprint_issues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sprint_id UUID NOT NULL REFERENCES sprints(id) ON DELETE CASCADE,
    issue_id UUID NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(sprint_id, issue_id)
);

CREATE INDEX IF NOT EXISTS idx_sprint_issues_sprint_id ON sprint_issues(sprint_id);
CREATE INDEX IF NOT EXISTS idx_sprint_issues_issue_id ON sprint_issues(issue_id);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_sprints_updated_at BEFORE UPDATE ON sprints
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
