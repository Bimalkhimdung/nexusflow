-- Rollback authentication fields from users table
ALTER TABLE users 
DROP COLUMN IF EXISTS password_hash,
DROP COLUMN IF EXISTS email_verified,
DROP COLUMN IF EXISTS verification_token,
DROP COLUMN IF EXISTS reset_token,
DROP COLUMN IF EXISTS reset_token_expiry,
DROP COLUMN IF EXISTS oauth_provider,
DROP COLUMN IF EXISTS oauth_id,
DROP COLUMN IF EXISTS oauth_access_token;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_verification_token;
DROP INDEX IF EXISTS idx_users_reset_token;
DROP INDEX IF EXISTS idx_users_oauth;
DROP INDEX IF EXISTS idx_users_email_verified;
