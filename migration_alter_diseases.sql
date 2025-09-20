-- Migration to alter diseases table: solution to string, keep image_link as array
-- This migration converts only the solution column from text[] to text
-- image_link remains as text[] array

-- First, create a backup of the existing data
CREATE TABLE diseases_backup AS SELECT * FROM diseases;

-- Alter the solution column from text[] to text
-- Convert array data to comma-separated string
ALTER TABLE diseases 
ALTER COLUMN solution TYPE text 
USING array_to_string(solution, '; ');

-- image_link column remains as text[] - no changes needed

-- Optional: Drop the backup table after confirming everything works
-- DROP TABLE diseases_backup;
