-- PostgreSQL Schema

-- I don't want business logic within the default postgresql schema
CREATE SCHEMA IF NOT EXISTS censys;
SET search_path TO censys, public;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA public;

-- if not exists doesn't work for create types
DO $$ BEGIN
    CREATE TYPE risk_level AS ENUM ('High', 'Medium', 'Low');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address VARCHAR(45) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    risk_level risk_level NOT NULL DEFAULT 'Low',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_ip_hostname UNIQUE(ip_address, hostname)
);

CREATE TABLE IF NOT EXISTS ports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    port_number INTEGER NOT NULL CHECK (port_number >= 1 AND port_number <= 65535),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_asset_port UNIQUE(asset_id, port_number)
);

CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    tag_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_asset_tag UNIQUE(asset_id, tag_name)
);

-- Performance
CREATE INDEX IF NOT EXISTS idx_assets_risk_level ON assets(risk_level);
CREATE INDEX IF NOT EXISTS idx_assets_ip_address ON assets(ip_address);
CREATE INDEX IF NOT EXISTS idx_assets_hostname ON assets(hostname);
CREATE INDEX IF NOT EXISTS idx_ports_asset_id ON ports(asset_id);
CREATE INDEX IF NOT EXISTS idx_tags_asset_id ON tags(asset_id);
CREATE INDEX IF NOT EXISTS idx_ports_port_number ON ports(port_number);
CREATE INDEX IF NOT EXISTS idx_tags_tag_name ON tags(tag_name);