# gRPC-Gateway Integration Guide

## Overview

This service uses [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to automatically expose REST/HTTP APIs alongside gRPC endpoints. This provides dual protocol support from a single implementation.

## Architecture

```
┌─────────────┐         ┌──────────────┐
│ gRPC Client │         │ HTTP Client  │
└──────┬──────┘         └──────┬───────┘
       │                       │
       │ gRPC (Port 50051)     │ HTTP (Port 8080)
       │                       │
       ▼                       ▼
┌──────────────────────────────────────┐
│          gRPC Server                 │
│  ┌────────────┐   ┌───────────────┐ │
│  │  Handlers  │◄──┤ HTTP Gateway  │ │
│  └──────┬─────┘   │  (Reverse     │ │
│         │         │   Proxy)      │ │
│         ▼         └───────────────┘ │
│  ┌─────────────┐                    │
│  │   Service   │                    │
│  └──────┬──────┘                    │
│         ▼                            │
│  ┌─────────────┐                    │
│  │ Repository  │                    │
│  └─────────────┘                    │
└──────────────────────────────────────┘
```

## How It Works

1. **Proto Definitions**: Service methods are defined once in `proto/delivery.proto` with HTTP annotations
2. **Code Generation**: `make proto` generates:
   - `delivery.pb.go` - Protobuf message definitions
   - `delivery_grpc.pb.go` - gRPC server/client code
   - `delivery.pb.gw.go` - HTTP gateway reverse proxy
   - `api.swagger.json` - OpenAPI/Swagger specification
3. **Runtime**: Gateway translates HTTP/JSON ↔ gRPC/Protobuf automatically
4. **Single Implementation**: Same business logic serves both protocols

## Benefits

✅ **Single Source of Truth**: One proto definition for both gRPC and REST
✅ **No Code Duplication**: Same handlers serve both protocols
✅ **Type Safety**: Full protobuf validation on REST requests
✅ **Auto-Documentation**: OpenAPI spec generated automatically
✅ **Performance**: REST clients get gRPC-level performance internally
✅ **Gradual Migration**: Supports both protocols during transitions
✅ **Standards Compliant**: Follows Google API design guidelines

## REST API Endpoints

All gRPC methods are exposed as REST endpoints:

| Method | gRPC | REST | Description |
|--------|------|------|-------------|
| CreateDeliveryAssignment | `CreateDeliveryAssignment` | `POST /v1/deliveries` | Create new delivery |
| GetDeliveryAssignment | `GetDeliveryAssignment` | `GET /v1/deliveries/{id}` | Get delivery by ID |
| UpdateDeliveryStatus | `UpdateDeliveryStatus` | `PATCH /v1/deliveries/{id}/status` | Update status |
| ListDeliveryAssignments | `ListDeliveryAssignments` | `GET /v1/deliveries` | List with filters |
| AssignDriver | `AssignDriver` | `POST /v1/deliveries/{id}/assign-driver` | Assign driver |
| GetDeliveryMetrics | `GetDeliveryMetrics` | `GET /v1/deliveries/metrics` | Get metrics |

## Adding New Endpoints

To add a new endpoint that supports both gRPC and REST:

### 1. Define in Proto

```protobuf
// In proto/delivery.proto
rpc GetDeliveryStats(GetDeliveryStatsRequest) returns (DeliveryStats) {
  option (google.api.http) = {
    get: "/v1/deliveries/stats"
  };
}

message GetDeliveryStatsRequest {
  google.protobuf.Timestamp start_date = 1;
  google.protobuf.Timestamp end_date = 2;
}

message DeliveryStats {
  int32 total = 1;
  int32 completed = 2;
  double success_rate = 3;
}
```

### 2. Generate Code

```bash
make proto
```

This automatically creates both gRPC and REST endpoints.

### 3. Implement Handler

```go
// In internal/transport/grpc/delivery_handler.go
func (h *Handler) GetDeliveryStats(ctx context.Context, req *pb.GetDeliveryStatsRequest) (*pb.DeliveryStats, error) {
    // Implementation
    stats, err := h.service.GetDeliveryStats(ctx, req.StartDate, req.EndDate)
    if err != nil {
        return nil, err
    }

    return &pb.DeliveryStats{
        Total:       stats.Total,
        Completed:   stats.Completed,
        SuccessRate: stats.SuccessRate,
    }, nil
}
```

### 4. Test Both Protocols

**gRPC:**
```bash
grpcurl -plaintext -d '{"start_date": "2024-01-01T00:00:00Z", "end_date": "2024-01-31T23:59:59Z"}' \
  localhost:50051 delivery.DeliveryService/GetDeliveryStats
```

**REST:**
```bash
curl "http://localhost:8080/v1/deliveries/stats?start_date=2024-01-01T00:00:00Z&end_date=2024-01-31T23:59:59Z"
```

Both return identical data (JSON for REST, protobuf for gRPC internally).

## HTTP Annotation Options

### Query Parameters (GET)

```protobuf
rpc ListItems(ListRequest) returns (ListResponse) {
  option (google.api.http) = {
    get: "/v1/items"  // ?page=1&page_size=20
  };
}
```

### Path Parameters

```protobuf
rpc GetItem(GetRequest) returns (Item) {
  option (google.api.http) = {
    get: "/v1/items/{id}"  // /v1/items/123
  };
}
```

### Request Body (POST/PUT/PATCH)

```protobuf
rpc CreateItem(CreateRequest) returns (Item) {
  option (google.api.http) = {
    post: "/v1/items"
    body: "*"  // Entire request is the body
  };
}
```

### Partial Body

```protobuf
rpc UpdateItem(UpdateRequest) returns (Item) {
  option (google.api.http) = {
    patch: "/v1/items/{id}"
    body: "item"  // Only the 'item' field is in body
  };
}
```

### Multiple HTTP Methods

```protobuf
rpc SearchItems(SearchRequest) returns (SearchResponse) {
  option (google.api.http) = {
    get: "/v1/items/search"
    additional_bindings {
      post: "/v1/items/search"
      body: "*"
    }
  };
}
```

## OpenAPI/Swagger Documentation

The service automatically generates `api.swagger.json` in OpenAPI 2.0 format.

### Viewing the Documentation

**Option 1: Swagger UI**
```bash
docker run -p 8081:8080 -e SWAGGER_JSON=/api.swagger.json -v $(pwd):/usr/share/nginx/html/api swaggerapi/swagger-ui
# Open: http://localhost:8081
```

**Option 2: Swagger Editor**
```bash
docker run -p 8081:8080 swaggerapi/swagger-editor
# Upload api.swagger.json
```

**Option 3: Postman**
- Import `api.swagger.json` into Postman
- All endpoints configured automatically

### Customizing OpenAPI Output

Add metadata to proto file:

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "Order Delivery Service API";
    version: "1.0";
    description: "Enterprise delivery management system";
    contact: {
      name: "API Support";
      url: "https://github.com/company/order-delivery-service";
      email: "api-support@company.com";
    };
  };
  host: "api.example.com";
  schemes: HTTPS;
  schemes: HTTP;
  consumes: "application/json";
  produces: "application/json";
};
```

## Configuration

### Ports

Configure via environment variables:

```bash
PORT=50051        # gRPC server port
HTTP_PORT=8080    # HTTP gateway port
METRICS_PORT=9090 # Prometheus metrics
```

### Gateway Options

Customize in `cmd/server/main.go`:

```go
// Custom marshaler options
gwMux := runtime.NewServeMux(
    runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
        MarshalOptions: protojson.MarshalOptions{
            UseProtoNames:   true,  // Use proto field names (not camelCase)
            EmitUnpopulated: true,  // Include zero values
        },
    }),
    runtime.WithErrorHandler(customErrorHandler),
    runtime.WithMetadata(annotateContext),
)
```

## Performance Considerations

### Overhead
- Gateway adds ~1-2ms latency for HTTP→gRPC translation
- Negligible for most applications
- Internal services can still use gRPC directly for maximum performance

### Optimization
```go
// Use connection pooling for gRPC backend
opts := []grpc.DialOption{
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithDefaultCallOptions(
        grpc.MaxCallRecvMsgSize(10 * 1024 * 1024), // 10MB
    ),
}
```

## Troubleshooting

### Issue: Gateway not starting

**Check:**
```bash
# Verify gRPC server is running
grpcurl -plaintext localhost:50051 list

# Check gateway logs
grep "HTTP gateway" logs.txt
```

### Issue: 404 Not Found on REST endpoint

**Solution:**
1. Verify endpoint in `api.swagger.json`
2. Ensure `make proto` was run after proto changes
3. Check that gRPC method is implemented

### Issue: Request validation errors

**Debug:**
```bash
# Compare gRPC vs REST request
grpcurl -plaintext -d '{"order_id":"123"}' localhost:50051 delivery.DeliveryService/GetDeliveryAssignment
curl http://localhost:8080/v1/deliveries/123
```

### Issue: IDE shows proto import errors

See [IDE Setup Guide](IDE_SETUP.md) for GoLand/IntelliJ configuration.

## Testing

### Integration Tests

Test both protocols:

```go
func TestCreateDelivery_BothProtocols(t *testing.T) {
    // Start test server
    grpcAddr := startGRPCServer(t)
    httpAddr := startHTTPGateway(t, grpcAddr)

    // Test gRPC
    grpcResp := createViaGRPC(t, grpcAddr, testData)

    // Test HTTP
    httpResp := createViaHTTP(t, httpAddr, testData)

    // Should be identical
    assert.Equal(t, grpcResp.Id, httpResp.Id)
}
```

### Load Testing

```bash
# gRPC
ghz --insecure --proto proto/delivery.proto \
    --call delivery.DeliveryService/GetDeliveryAssignment \
    -d '{"id":"123"}' \
    -n 10000 -c 100 \
    localhost:50051

# HTTP
hey -n 10000 -c 100 http://localhost:8080/v1/deliveries/123
```

## References

- [gRPC-Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Google API Design Guide](https://cloud.google.com/apis/design)
- [HTTP to gRPC Transcoding](https://cloud.google.com/endpoints/docs/grpc/transcoding)
- [OpenAPI Specification](https://swagger.io/specification/)
