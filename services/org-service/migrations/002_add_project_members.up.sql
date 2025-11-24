-- Add project_members table for project-level access control
CREATE TABLE IF NOT EXISTS project_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    added_by UUID,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, user_id)
);

CREATE INDEX idx_project_members_project ON project_members(project_id);
CREATE INDEX idx_project_members_user ON project_members(user_id);

-- Add invite tracking fields to org_members if they don't exist
ALTER TABLE org_members ADD COLUMN IF NOT EXISTS invited_by UUID;
ALTER TABLE org_members ADD COLUMN IF NOT EXISTS invited_at TIMESTAMP WITH TIME ZONE;

-- Update org_members role constraint to match our requirements
ALTER TABLE org_members DROP CONSTRAINT IF EXISTS org_members_role_check;
ALTER TABLE org_members ADD CONSTRAINT org_members_role_check CHECK (role IN ('member', 'admin'));

-- Add trigger for project_members updated_at
CREATE TRIGGER update_project_members_updated_at
    BEFORE UPDATE ON project_members
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
