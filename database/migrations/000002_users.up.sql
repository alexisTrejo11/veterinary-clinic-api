-- 000002_users.up.sql
-- Create users table

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    -- Authentication Information
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(255) UNIQUE,
    hashed_password VARCHAR(255),
    role VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,

    -- Personal Information
    name VARCHAR(255) NOT NULL,
    gender VARCHAR(50),
    photo_url TEXT DEFAULT 'https://example.com/default-profile.png',
    bio TEXT DEFAULT '',
    date_of_birth DATE,

    -- OAuth2 Information
    oauth_provider VARCHAR(50) DEFAULT 'local' NOT NULL,
    oauth_provider_id VARCHAR(255),
    oauth_access_token TEXT,
    oauth_refresh_token TEXT,
    oauth_token_expiry TIMESTAMP WITH TIME ZONE,
    email_verified BOOLEAN DEFAULT FALSE,
    
    -- 2FA (Two-Factor Authentication)
    two_fa_method VARCHAR(20) DEFAULT 'none' NOT NULL,
    two_fa_secret VARCHAR(255),
    two_fa_enabled BOOLEAN DEFAULT FALSE,
    two_fa_enabled_at TIMESTAMP WITH TIME ZONE,
    two_fa_backup_codes TEXT[], -- Array of hashed backup codes
    two_fa_backup_codes_generated_at TIMESTAMP WITH TIME ZONE,
    last_2fa_code_used_at TIMESTAMP WITH TIME ZONE,

    -- Activity Tracking
    last_login TIMESTAMP WITH TIME ZONE,
    login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT unique_oauth_provider_id UNIQUE(oauth_provider, oauth_provider_id),
    CONSTRAINT check_oauth_provider CHECK (oauth_provider IN ('google', 'facebook', 'apple', 'microsoft', 'github', 'local')),
    CONSTRAINT check_two_fa_method CHECK (two_fa_method IN ('totp', 'sms', 'email', 'none')),
    CONSTRAINT check_user_role CHECK (role IN ('veterinarian', 'receptionist', 'manager', 'customer', 'admin', 'superadmin')),
    CONSTRAINT check_user_status CHECK (status IN ('active', 'inactive', 'pending', 'banned', 'deleted'))
);

-- Indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_oauth_provider ON users(oauth_provider, oauth_provider_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_email_verified ON users(email_verified) WHERE email_verified = FALSE;
