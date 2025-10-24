# API Documentation

## Order Delivery Service gRPC API

This document describes the gRPC API endpoints for the Order Delivery Service.

### Service Definition

```protobuf
service DeliveryService {
  rpc CreateDeliveryAssignment(CreateDeliveryAssignmentRequest) returns (DeliveryAssignment);
  rpc GetDeliveryAssignment(GetDeliveryAssignmentRequest) returns (DeliveryAssignment);
  rpc UpdateDeliveryStatus(UpdateDeliveryStatusRequest) returns (DeliveryAssignment);
  rpc ListDeliveryAssignments(ListDeliveryAssignmentsRequest) returns (ListDeliveryAssignmentsResponse);
  rpc AssignDriver(AssignDriverRequest) returns (DeliveryAssignment);
  rpc GetDeliveryMetrics(GetDeliveryMetricsRequest) returns (DeliveryMetrics);
}
```

## Endpoints

### CreateDeliveryAssignment

Creates a new delivery assignment for an order.

**Request:**
```protobuf
message CreateDeliveryAssignmentRequest {
  string order_id = 1;                              // Required
  Address pickup_address = 2;                        // Required
  Address delivery_address = 3;                      // Required
  google.protobuf.Timestamp scheduled_pickup_time = 4;   // Required
  google.protobuf.Timestamp estimated_delivery_time = 5; // Required
  string notes = 6;                                  // Optional
}
```

**Response:**
```protobuf
message DeliveryAssignment {
  string id = 1;
  string order_id = 2;
  string driver_id = 3;
  DeliveryStatus status = 4;
  // ... more fields
}
```

**Example (grpcurl):**
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
  "notes": "Fragile package"
}' localhost:50051 delivery.DeliveryService/CreateDeliveryAssignment
```

### GetDeliveryAssignment

Retrieves a delivery assignment by ID.

**Request:**
```protobuf
message GetDeliveryAssignmentRequest {
  string id = 1;  // UUID format required
}
```

**Response:**
```protobuf
message DeliveryAssignment { /* ... */ }
```

**Example:**
```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:50051 delivery.DeliveryService/GetDeliveryAssignment
```

### UpdateDeliveryStatus

Updates the status of a delivery assignment.

**Request:**
```protobuf
message UpdateDeliveryStatusRequest {
  string id = 1;              // UUID format required
  DeliveryStatus status = 2;  // Required
  string notes = 3;           // Optional
}
```

**Valid Status Transitions:**
- PENDING → ASSIGNED, CANCELLED
- ASSIGNED → PICKED_UP, CANCELLED
- PICKED_UP → IN_TRANSIT, FAILED
- IN_TRANSIT → DELIVERED, FAILED
- DELIVERED → (final state)
- FAILED → (final state)
- CANCELLED → (final state)

**Example:**
```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DELIVERY_STATUS_PICKED_UP",
  "notes": "Package collected from sender"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus
```

### ListDeliveryAssignments

Lists delivery assignments with pagination and filtering.

**Request:**
```protobuf
message ListDeliveryAssignmentsRequest {
  int32 page = 1;              // Default: 1
  int32 page_size = 2;         // Default: 20, Max: 100
  DeliveryStatus status = 3;   // Optional filter
  string driver_id = 4;        // Optional filter
}
```

**Response:**
```protobuf
message ListDeliveryAssignmentsResponse {
  repeated DeliveryAssignment assignments = 1;
  int32 total_count = 2;
  int32 page = 3;
  int32 page_size = 4;
}
```

**Example:**
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10,
  "status": "DELIVERY_STATUS_IN_TRANSIT"
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

### AssignDriver

Assigns a driver to a delivery assignment.

**Request:**
```protobuf
message AssignDriverRequest {
  string id = 1;        // Assignment UUID
  string driver_id = 2; // Driver identifier
}
```

**Example:**
```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "driver_id": "DRIVER-789"
}' localhost:50051 delivery.DeliveryService/AssignDriver
```

### GetDeliveryMetrics

Retrieves aggregated delivery metrics for a time range.

**Request:**
```protobuf
message GetDeliveryMetricsRequest {
  google.protobuf.Timestamp start_time = 1;  // Required
  google.protobuf.Timestamp end_time = 2;    // Required
  string driver_id = 3;                       // Optional
}
```

**Response:**
```protobuf
message DeliveryMetrics {
  int32 total_deliveries = 1;
  int32 completed_deliveries = 2;
  int32 failed_deliveries = 3;
  int32 cancelled_deliveries = 4;
  double average_delivery_time_minutes = 5;
  double on_time_delivery_rate = 6;
}
```

**Example:**
```bash
grpcurl -plaintext -d '{
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-01-31T23:59:59Z"
}' localhost:50051 delivery.DeliveryService/GetDeliveryMetrics
```

## Status Codes

The service uses standard gRPC status codes:

- `OK` - Success
- `INVALID_ARGUMENT` - Invalid input (e.g., malformed UUID, missing required fields)
- `NOT_FOUND` - Resource not found
- `FAILED_PRECONDITION` - Invalid state transition
- `ALREADY_EXISTS` - Resource already exists
- `INTERNAL` - Internal server error

## Data Types

### DeliveryStatus Enum

```protobuf
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

### Address Message

```protobuf
message Address {
  string street = 1;
  string city = 2;
  string state = 3;
  string postal_code = 4;
  string country = 5;
  double latitude = 6;
  double longitude = 7;
}
```

## Testing with grpcurl

Install grpcurl:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

List available services:
```bash
grpcurl -plaintext localhost:50051 list
```

Describe a service:
```bash
grpcurl -plaintext localhost:50051 describe delivery.DeliveryService
```

## Health Check

The service implements the standard gRPC health checking protocol:

```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

## Error Handling

All endpoints return appropriate gRPC status codes with descriptive error messages. Example error response:

```
ERROR:
  Code: InvalidArgument
  Message: order_id is required
```
