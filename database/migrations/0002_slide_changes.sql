ALTER TABLE slides DROP COLUMN slide_number, DROP COLUMN metadata;
ALTER TABLE presentations ADD COLUMN slide_order uuid[] NOT NULL DEFAULT '{}';
