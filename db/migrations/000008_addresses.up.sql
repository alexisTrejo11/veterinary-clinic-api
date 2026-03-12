CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    zip_code VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    building_type VARCHAR(255) NOT NULL,
    building_outer_number VARCHAR(255) NOT NULL,
    building_inner_number VARCHAR(255) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    CREATED_AT TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UPDATED_AT TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    DELETED_AT TIMESTAMP WITH TIME ZONE,
    VERSION INT NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_building_type CHECK (building_type IN ('house', 'apartment', 'office', 'other')),
    CONSTRAINT chk_country CHECK (country IN ('USA', 'Mexico', 'Canada')),
    CONSTRAINT chk_is_default CHECK (is_default IN (TRUE, FALSE)),
    CONSTRAINT chk_street CHECK (street IS NOT NULL),
    CONSTRAINT chk_city CHECK (city IS NOT NULL)
);


CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses (user_id);
CREATE INDEX IF NOT EXISTS idx_addresses_street ON addresses (street);
CREATE INDEX IF NOT EXISTS idx_addresses_city ON addresses (city);
CREATE INDEX IF NOT EXISTS idx_addresses_state ON addresses (state);
CREATE INDEX IF NOT EXISTS idx_addresses_zip_code ON addresses (zip_code);
CREATE INDEX IF NOT EXISTS idx_addresses_country ON addresses (country);
CREATE INDEX IF NOT EXISTS idx_addresses_building_type ON addresses (building_type);
CREATE INDEX IF NOT EXISTS idx_addresses_building_outer_number ON addresses (building_outer_number);
CREATE INDEX IF NOT EXISTS idx_addresses_building_inner_number ON addresses (building_inner_number);
CREATE INDEX IF NOT EXISTS idx_addresses_is_default ON addresses (is_default);