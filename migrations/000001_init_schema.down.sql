-- Drop trigger
DROP TRIGGER IF EXISTS update_delivery_assignments_updated_at ON delivery_assignments;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_delivery_assignments_deleted_at;
DROP INDEX IF EXISTS idx_delivery_assignments_created_at;
DROP INDEX IF EXISTS idx_delivery_assignments_scheduled_pickup;
DROP INDEX IF EXISTS idx_delivery_assignments_status;
DROP INDEX IF EXISTS idx_delivery_assignments_driver_id;
DROP INDEX IF EXISTS idx_delivery_assignments_order_id;

-- Drop table
DROP TABLE IF EXISTS delivery_assignments;
