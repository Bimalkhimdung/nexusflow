DROP TRIGGER IF EXISTS update_teams_updated_at ON teams;
DROP TRIGGER IF EXISTS update_org_members_updated_at ON org_members;
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;

DROP TABLE IF EXISTS invites;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS org_members;
DROP TABLE IF EXISTS organizations;
