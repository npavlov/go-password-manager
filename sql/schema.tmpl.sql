-- Enable the pgcrypto extension to use gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create users table
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password TEXT NOT NULL
);

-- Create orders table
CREATE TABLE passwords (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                        name VARCHAR(255) NOT NULL,
                        login VARCHAR(255) NOT NULL,
                        password VARCHAR(255) NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        CONSTRAINT unique_user_password_name UNIQUE (user_id, name)
);

-- Add trigger function to update the updated_at column
CREATE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for orders table to automatically update updated_at
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON passwords
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

