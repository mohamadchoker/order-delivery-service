# API Examples - gRPC Service

## Protocol Buffer Enum Values

**IMPORTANT**: Proto enums use prefixed names. When using grpcurl, you must use the full enum name:

### DeliveryStatus Enum Values

| Domain (Internal) | Proto (gRPC/API) | Description |
|------------------|------------------|-------------|
| `PENDING` | `DELIVERY_STATUS_PENDING` | Delivery created, awaiting assignment |
| `ASSIGNED` | `DELIVERY_STATUS_ASSIGNED` | Driver assigned |
| `PICKED_UP` | `DELIVERY_STATUS_PICKED_UP` | Package picked up from sender |
| `IN_TRANSIT` | `DELIVERY_STATUS_IN_TRANSIT` | Package in transit |
| `DELIVERED` | `DELIVERY_STATUS_DELIVERED` | Package delivered successfully |
| `FAILED` | `DELIVERY_STATUS_FAILED` | Delivery failed |
| `CANCELLED` | `DELIVERY_STATUS_CANCELLED` | Delivery cancelled |

**Why the prefix?** Protocol Buffers best practice to avoid naming conflicts across enums.

---

## 🚀 Complete API Examples

### 1. Health Check

```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

**Response:**
```json
{
  "status": "SERVING"
}
```

---

### 2. List All Services

```bash
grpcurl -plaintext localhost:50051 list
```

**Response:**
```
delivery.DeliveryService
grpc.health.v1.Health
```

---

### 3. Describe Service

```bash
grpcurl -plaintext localhost:50051 describe delivery.DeliveryService
```

---

### 4. Create Delivery Assignment

```bash
grpcurl -plaintext -d '{
  "order_id": "ORDER-12345",
  "pickup_address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "USA",
    "latitude": 40.7128,
    "longitude": -74.0060
  },
  "delivery_address": {
    "street": "456 Oak Ave",
    "city": "Boston",
    "state": "MA",
    "postal_code": "02101",
    "country": "USA",
    "latitude": 42.3601,
    "longitude": -71.0589
  },
  "scheduled_pickup_time": "2024-01-15T10:00:00Z",
  "estimated_delivery_time": "2024-01-15T14:00:00Z",
  "notes": "Handle with care - fragile items"
}' localhost:50051 delivery.DeliveryService/CreateDeliveryAssignment
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "orderId": "ORDER-12345",
  "status": "DELIVERY_STATUS_PENDING",
  "pickupAddress": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "postalCode": "10001",
    "country": "USA",
    "latitude": 40.7128,
    "longitude": -74.006
  },
  "deliveryAddress": {
    "street": "456 Oak Ave",
    "city": "Boston",
    "state": "MA",
    "postalCode": "02101",
    "country": "USA",
    "latitude": 42.3601,
    "longitude": -71.0589
  },
  "scheduledPickupTime": "2024-01-15T10:00:00Z",
  "estimatedDeliveryTime": "2024-01-15T14:00:00Z",
  "notes": "Handle with care - fragile items",
  "createdAt": "2024-01-14T12:00:00Z",
  "updatedAt": "2024-01-14T12:00:00Z"
}
```

---

### 5. Get Delivery Assignment

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:50051 delivery.DeliveryService/GetDeliveryAssignment
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "orderId": "ORDER-12345",
  "status": "DELIVERY_STATUS_PENDING",
  "pickupAddress": { ... },
  "deliveryAddress": { ... },
  "scheduledPickupTime": "2024-01-15T10:00:00Z",
  "estimatedDeliveryTime": "2024-01-15T14:00:00Z",
  "createdAt": "2024-01-14T12:00:00Z",
  "updatedAt": "2024-01-14T12:00:00Z"
}
```

---

### 6. Assign Driver

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "driver_id": "DRIVER-123"
}' localhost:50051 delivery.DeliveryService/AssignDriver
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "orderId": "ORDER-12345",
  "driverId": "DRIVER-123",
  "status": "DELIVERY_STATUS_ASSIGNED",
  "pickupAddress": { ... },
  "deliveryAddress": { ... },
  "scheduledPickupTime": "2024-01-15T10:00:00Z",
  "estimatedDeliveryTime": "2024-01-15T14:00:00Z",
  "createdAt": "2024-01-14T12:00:00Z",
  "updatedAt": "2024-01-14T12:05:00Z"
}
```

---

### 7. Update Status to PICKED_UP

**IMPORTANT**: Use full enum name `DELIVERY_STATUS_PICKED_UP`, not just `PICKED_UP`

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DELIVERY_STATUS_PICKED_UP",
  "notes": "Package collected from sender at 10:05 AM"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "orderId": "ORDER-12345",
  "driverId": "DRIVER-123",
  "status": "DELIVERY_STATUS_PICKED_UP",
  "actualPickupTime": "2024-01-15T10:05:00Z",
  "notes": "Package collected from sender at 10:05 AM",
  "createdAt": "2024-01-14T12:00:00Z",
  "updatedAt": "2024-01-15T10:05:00Z"
}
```

---

### 8. Update Status to IN_TRANSIT

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DELIVERY_STATUS_IN_TRANSIT",
  "notes": "Package on the way to destination"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus
```

---

### 9. Update Status to DELIVERED

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DELIVERY_STATUS_DELIVERED",
  "notes": "Package delivered successfully. Signed by: John Doe"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "orderId": "ORDER-12345",
  "driverId": "DRIVER-123",
  "status": "DELIVERY_STATUS_DELIVERED",
  "actualPickupTime": "2024-01-15T10:05:00Z",
  "actualDeliveryTime": "2024-01-15T13:45:00Z",
  "notes": "Package delivered successfully. Signed by: John Doe",
  "createdAt": "2024-01-14T12:00:00Z",
  "updatedAt": "2024-01-15T13:45:00Z"
}
```

---

### 10. List All Deliveries (No Filter)

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

**Response:**
```json
{
  "assignments": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "orderId": "ORDER-12345",
      "status": "DELIVERY_STATUS_DELIVERED",
      ...
    },
    {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "orderId": "ORDER-12346",
      "status": "DELIVERY_STATUS_IN_TRANSIT",
      ...
    }
  ],
  "totalCount": 2,
  "page": 1,
  "pageSize": 20
}
```

---

### 11. List PENDING Deliveries

**CORRECT** - Use full enum name:
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20,
  "status": "DELIVERY_STATUS_PENDING"
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

**INCORRECT** - This will fail:
```bash
# ❌ ERROR: enum "delivery.DeliveryStatus" does not have value named "PENDING"
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20,
  "status": "PENDING"
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

---

### 12. List ASSIGNED Deliveries

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20,
  "status": "DELIVERY_STATUS_ASSIGNED"
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

---

### 13. List Deliveries by Driver

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20,
  "driver_id": "DRIVER-123"
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

---

### 14. List PENDING Deliveries for Specific Driver

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20,
  "status": "DELIVERY_STATUS_PENDING",
  "driver_id": "DRIVER-123"
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

---

### 15. Get Delivery Metrics (All Deliveries)

```bash
grpcurl -plaintext -d '{
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-01-31T23:59:59Z"
}' localhost:50051 delivery.DeliveryService/GetDeliveryMetrics
```

**Response:**
```json
{
  "totalDeliveries": 150,
  "completedDeliveries": 120,
  "failedDeliveries": 10,
  "cancelledDeliveries": 20,
  "averageDeliveryTimeMinutes": 185.5,
  "onTimeDeliveryRate": 0.92
}
```

---

### 16. Get Metrics for Specific Driver

```bash
grpcurl -plaintext -d '{
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-01-31T23:59:59Z",
  "driver_id": "DRIVER-123"
}' localhost:50051 delivery.DeliveryService/GetDeliveryMetrics
```

---

## 🔄 Status Transition Examples

### Valid Status Flow

```bash
# 1. Create delivery (PENDING)
grpcurl -plaintext -d '{...}' \
  localhost:50051 delivery.DeliveryService/CreateDeliveryAssignment

# 2. Assign driver (PENDING → ASSIGNED)
grpcurl -plaintext -d '{
  "id": "...",
  "driver_id": "DRIVER-123"
}' localhost:50051 delivery.DeliveryService/AssignDriver

# 3. Pick up package (ASSIGNED → PICKED_UP)
grpcurl -plaintext -d '{
  "id": "...",
  "status": "DELIVERY_STATUS_PICKED_UP"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus

# 4. Start transit (PICKED_UP → IN_TRANSIT)
grpcurl -plaintext -d '{
  "id": "...",
  "status": "DELIVERY_STATUS_IN_TRANSIT"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus

# 5. Complete delivery (IN_TRANSIT → DELIVERED)
grpcurl -plaintext -d '{
  "id": "...",
  "status": "DELIVERY_STATUS_DELIVERED"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus
```

### Invalid Status Transition (Will Fail)

```bash
# ❌ Cannot go from PENDING directly to DELIVERED
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DELIVERY_STATUS_DELIVERED"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus

# Error: invalid status transition
```

---

## 📊 Status Transition Diagram

```
PENDING → ASSIGNED → PICKED_UP → IN_TRANSIT → DELIVERED
   ↓         ↓          ↓            ↓
CANCELLED  CANCELLED  FAILED      FAILED
```

**Valid Transitions:**
- `PENDING` → `ASSIGNED`, `CANCELLED`
- `ASSIGNED` → `PICKED_UP`, `CANCELLED`
- `PICKED_UP` → `IN_TRANSIT`, `FAILED`
- `IN_TRANSIT` → `DELIVERED`, `FAILED`

**Terminal States** (no further transitions):
- `DELIVERED`
- `FAILED`
- `CANCELLED`

---

## 🐛 Common Errors

### Error: Enum value not found

```bash
# ❌ WRONG
"status": "PENDING"

# ✅ CORRECT
"status": "DELIVERY_STATUS_PENDING"
```

### Error: Invalid transition

```bash
# ❌ WRONG - Cannot skip states
PENDING → DELIVERED

# ✅ CORRECT - Follow valid transitions
PENDING → ASSIGNED → PICKED_UP → IN_TRANSIT → DELIVERED
```

### Error: Invalid timestamp format

```bash
# ❌ WRONG
"scheduled_pickup_time": "2024-01-15 10:00:00"

# ✅ CORRECT (RFC 3339 / ISO 8601)
"scheduled_pickup_time": "2024-01-15T10:00:00Z"
```

---

## 📝 Quick Reference

### All Status Values

```json
{
  "status": "DELIVERY_STATUS_UNSPECIFIED",  // 0 - Never use explicitly
  "status": "DELIVERY_STATUS_PENDING",      // 1
  "status": "DELIVERY_STATUS_ASSIGNED",     // 2
  "status": "DELIVERY_STATUS_PICKED_UP",    // 3
  "status": "DELIVERY_STATUS_IN_TRANSIT",   // 4
  "status": "DELIVERY_STATUS_DELIVERED",    // 5
  "status": "DELIVERY_STATUS_FAILED",       // 6
  "status": "DELIVERY_STATUS_CANCELLED"     // 7
}
```

### Field Naming

Protocol Buffers uses camelCase in JSON:

| Proto Field | JSON Field |
|------------|------------|
| `order_id` | `orderId` |
| `driver_id` | `driverId` |
| `pickup_address` | `pickupAddress` |
| `scheduled_pickup_time` | `scheduledPickupTime` |
| `created_at` | `createdAt` |

---

## 🔗 Resources

- **Proto Definition**: `proto/delivery.proto`
- **Service Implementation**: `internal/transport/grpc/handler.go`
- **Domain Logic**: `internal/domain/delivery.go`
- **gRPC Documentation**: https://grpc.io/docs/
- **grpcurl GitHub**: https://github.com/fullstorydev/grpcurl

---

## 💡 Tips

1. **Always use full enum names** with `DELIVERY_STATUS_` prefix
2. **Use RFC 3339 timestamps** (e.g., `2024-01-15T10:00:00Z`)
3. **Check valid transitions** before updating status
4. **Use pagination** for list operations
5. **Include notes** for audit trail

---

**Pro Tip**: Use grpcurl's describe command to see field types:

```bash
grpcurl -plaintext localhost:50051 describe delivery.DeliveryStatus
```

Output:
```
delivery.DeliveryStatus is an enum:
enum DeliveryStatus {
  DELIVERY_STATUS_UNSPECIFIED = 0;
  DELIVERY_STATUS_PENDING = 1;
  DELIVERY_STATUS_ASSIGNED = 2;
  DELIVERY_STATUS_PICKED_UP = 3;
  DELIVERY_STATUS_IN_TRANSIT = 4;
  DELIVERY_STATUS_DELIVERED = 5;
  DELIVERY_STATUS_FAILED = 6;
  DELIVERY_STATUS_CANCELLED = 7;
}
```
