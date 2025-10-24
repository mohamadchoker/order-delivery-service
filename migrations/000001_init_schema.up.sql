-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create delivery_assignments table
CREATE TABLE IF NOT EXISTS delivery_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id VARCHAR(100) NOT NULL,
    driver_id VARCHAR(100),
    status VARCHAR(50) NOT NULL,
    pickup_address JSONB NOT NULL,
    delivery_address JSONB NOT NULL,
    scheduled_pickup_time TIMESTAMP NOT NULL,
    estimated_delivery_time TIMESTAMP NOT NULL,
    actual_pickup_time TIMESTAMP,
    actual_delivery_time TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_delivery_assignments_order_id ON delivery_assignments(order_id);
CREATE INDEX idx_delivery_assignments_driver_id ON delivery_assignments(driver_id);
CREATE INDEX idx_delivery_assignments_status ON delivery_assignments(status);
CREATE INDEX idx_delivery_assignments_scheduled_pickup ON delivery_assignments(scheduled_pickup_time);
CREATE INDEX idx_delivery_assignments_created_at ON delivery_assignments(created_at);
CREATE INDEX idx_delivery_assignments_deleted_at ON delivery_assignments(deleted_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_delivery_assignments_updated_at
    BEFORE UPDATE ON delivery_assignments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE delivery_assignments IS 'Stores order delivery assignments';
COMMENT ON COLUMN delivery_assignments.id IS 'Unique identifier for the delivery assignment';
COMMENT ON COLUMN delivery_assignments.order_id IS 'Reference to the order being delivered';
COMMENT ON COLUMN delivery_assignments.driver_id IS 'Reference to the assigned driver';
COMMENT ON COLUMN delivery_assignments.status IS 'Current status of the delivery';
COMMENT ON COLUMN delivery_assignments.pickup_address IS 'JSONB containing pickup address details';
COMMENT ON COLUMN delivery_assignments.delivery_address IS 'JSONB containing delivery address details';
