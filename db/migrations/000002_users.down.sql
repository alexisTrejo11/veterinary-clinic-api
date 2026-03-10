-- 000002_users.down.sql
-- Drop users table and indexes

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_phone_number;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_oauth_provider;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_email_verified;

DROP TABLE IF EXISTS users CASCADE;
