-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Project Counters table (for generating sequential issue keys)
CREATE TABLE IF NOT EXISTS project_counters (
    project_id UUID PRIMARY KEY,
    next_issue_number BIGINT NOT NULL DEFAULT 1,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Issues table
CREATE TABLE IF NOT EXISTS issues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL, -- Reference to projects table in project-service
    key VARCHAR(50) NOT NULL, -- e.g., "PROJ-123"
    summary VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'task', -- epic, story, task, sub_task, bug
    priority VARCHAR(50) NOT NULL DEFAULT 'medium', -- lowest, low, medium, high, highest
    status_id UUID, -- Reference to workflow status (future)
    assignee_id UUID, -- Reference to users table
    reporter_id UUID, -- Reference to users table
    parent_id UUID REFERENCES issues(id), -- For sub-tasks
    sprint_id UUID, -- Reference to sprint (future)
    story_points INTEGER,
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    version BIGINT DEFAULT 1,
    UNIQUE(project_id, key)
);

CREATE INDEX idx_issues_project_id ON issues(project_id);
CREATE INDEX idx_issues_key ON issues(key);
CREATE INDEX idx_issues_assignee_id ON issues(assignee_id);
CREATE INDEX idx_issues_reporter_id ON issues(reporter_id);
CREATE INDEX idx_issues_parent_id ON issues(parent_id);
CREATE INDEX idx_issues_deleted_at ON issues(deleted_at);

-- Custom Fields table
CREATE TABLE IF NOT EXISTS custom_fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- text, number, date, select, etc.
    required BOOLEAN DEFAULT FALSE,
    default_value JSONB,
    options JSONB, -- For select/multi-select
    config JSONB, -- Type-specific config
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, name)
);

CREATE INDEX idx_custom_fields_project_id ON custom_fields(project_id);

-- Issue Custom Values table
CREATE TABLE IF NOT EXISTS issue_custom_values (
    issue_id UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    field_id UUID NOT NULL REFERENCES custom_fields(id) ON DELETE CASCADE,
    value JSONB,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (issue_id, field_id)
);

-- Issue Links table
CREATE TABLE IF NOT EXISTS issue_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_issue_id UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    target_issue_id UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- blocks, relates_to, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_issue_id, target_issue_id, type)
);

CREATE INDEX idx_issue_links_source ON issue_links(source_issue_id);
CREATE INDEX idx_issue_links_target ON issue_links(target_issue_id);

-- Issue Watchers table
CREATE TABLE IF NOT EXISTS issue_watchers (
    issue_id UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (issue_id, user_id)
);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_issues_updated_at
    BEFORE UPDATE ON issues
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_custom_fields_updated_at
    BEFORE UPDATE ON custom_fields
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_issue_custom_values_updated_at
    BEFORE UPDATE ON issue_custom_values
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
