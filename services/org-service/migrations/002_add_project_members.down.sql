-- Drop trigger
DROP TRIGGER IF EXISTS update_project_members_updated_at ON project_members;

-- Drop project_members table
DROP TABLE IF EXISTS project_members;

-- Remove added columns from org_members
ALTER TABLE org_members DROP COLUMN IF EXISTS invited_by;
ALTER TABLE org_members DROP COLUMN IF EXISTS invited_at;

-- Restore original role constraint (if needed)
ALTER TABLE org_members DROP CONSTRAINT IF EXISTS org_members_role_check;
