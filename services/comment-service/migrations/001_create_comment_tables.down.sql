-- Drop triggers
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS comment_mentions;
DROP TABLE IF EXISTS comment_reactions;
DROP TABLE IF EXISTS comments;
