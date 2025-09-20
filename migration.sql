-- Migration script to add role column to users table
-- Run this SQL script on your existing database to add the role column

-- Add role column to users table with default value 'user'
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';

-- Create an index on the role column for better query performance
CREATE INDEX idx_users_role ON users(role);

-- Optional: Create an admin user (replace with your desired admin credentials)
-- INSERT INTO users (id, email, username, password, full_name, role, created_at, updated_at) 
-- VALUES (
--     gen_random_uuid(),
--     'admin@plantheon.com',
--     'admin',
--     '$2a$14$example_hashed_password_here', -- Use bcrypt to hash your password
--     'System Administrator',
--     'admin',
--     NOW(),
--     NOW()
-- );

-- Verify the migration
-- SELECT id, email, username, role FROM users;
