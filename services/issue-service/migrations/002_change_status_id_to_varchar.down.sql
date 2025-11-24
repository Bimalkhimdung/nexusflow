ALTER TABLE issues ALTER COLUMN status_id TYPE UUID USING status_id::uuid;
