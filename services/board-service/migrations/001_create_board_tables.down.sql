-- Drop triggers
DROP TRIGGER IF EXISTS update_boards_updated_at ON boards;
DROP TRIGGER IF EXISTS update_cards_updated_at ON cards;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS boards;
