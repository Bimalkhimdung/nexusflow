-- Change status_id from UUID to VARCHAR to support string statuses for now
ALTER TABLE issues ALTER COLUMN status_id TYPE VARCHAR(50);
