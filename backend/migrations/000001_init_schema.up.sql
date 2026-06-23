-- 1. Initial Tables Setup for LifeMetric Multi-Tenant Ecosystem

CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    business_name VARCHAR(255) NOT NULL,
    plan_type VARCHAR(50) DEFAULT 'free' CHECK (plan_type IN ('free', 'premium')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE devices (
    device_uuid UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    room_label VARCHAR(100) NOT NULL,
    device_placement VARCHAR(50) DEFAULT 'ceiling' CHECK (device_placement IN ('ceiling', 'wall')),
    public_key TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'registered' CHECK (status IN ('registered', 'calibrating', 'active', 'offline')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    device_uuid UUID REFERENCES devices(device_uuid) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    encrypted_payload BYTEA NOT NULL, -- Encrypted using Merchant's Local RSA Public Key (One-Way Blind-Vault)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE webhooks (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    target_url TEXT NOT NULL,
    secret_key TEXT NOT NULL,
    event_types TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Activate Row-Level Security (RLS) to enforce perfect cross-facility tenant boundary isolation

ALTER TABLE devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE webhooks ENABLE ROW LEVEL SECURITY;

-- 3. Define the Multi-Tenant Security Policy
-- The database session context is populated by the Go API gateway on request connection setup

CREATE POLICY devices_isolation_policy ON devices
    FOR ALL USING (tenant_id = NULLIF(current_setting('app.current_tenant_id', true), '')::UUID);

CREATE POLICY transactions_isolation_policy ON transactions
    FOR ALL USING (tenant_id = NULLIF(current_setting('app.current_tenant_id', true), '')::UUID);

CREATE POLICY webhooks_isolation_policy ON webhooks
    FOR ALL USING (tenant_id = NULLIF(current_setting('app.current_tenant_id', true), '')::UUID);
