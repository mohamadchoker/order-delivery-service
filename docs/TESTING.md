# Testing Guide

## Overview

This project follows a comprehensive testing strategy with unit tests, integration tests, and end-to-end testing capabilities.

## Test Structure

```
order-delivery-service/
├── internal/
│   ├── entity/
│   │   └── delivery_test.go       # Entity unit tests
│   └── usecase/
│       └── delivery_usecase_test.go  # Use case unit tests
└── tests/
    └── integration/
        └── delivery_integration_test.go  # Integration tests
```

## Running Tests

### All Tests

```bash
make test
```

or

```bash
go test -v ./...
```

### With Coverage

```bash
make test-coverage
```

This generates `coverage.html` that you can open in a browser.

### Specific Package

```bash
go test -v ./internal/entity/
go test -v ./internal/usecase/
```

### Integration Tests

```bash
make test-integration
```

or

```bash
go test -v -tags=integration ./tests/integration/...
```

### Race Condition Detection

```bash
go test -race ./...
```

### Verbose Output

```bash
go test -v ./...
```

## Writing Tests

### Unit Test Example

```go
func TestNewDeliveryAssignment(t *testing.T) {
    // Arrange
    orderID := "ORDER-123"
    pickupAddr := entity.Address{City: "New York"}
    deliveryAddr := entity.Address{City: "Boston"}
    scheduledTime := time.Now().Add(1 * time.Hour)
    estimatedTime := time.Now().Add(3 * time.Hour)

    // Act
    assignment := entity.NewDeliveryAssignment(
        orderID,
        pickupAddr,
        deliveryAddr,
        scheduledTime,
        estimatedTime,
        "Test notes",
    )

    // Assert
    assert.NotNil(t, assignment)
    assert.Equal(t, orderID, assignment.OrderID)
    assert.Equal(t, entity.DeliveryStatusPending, assignment.Status)
}
```

### Table-Driven Test Example

```go
func TestUpdateStatus(t *testing.T) {
    tests := []struct {
        name          string
        currentStatus entity.DeliveryStatus
        newStatus     entity.DeliveryStatus
        expectError   bool
    }{
        {
            name:          "pending to assigned",
            currentStatus: entity.DeliveryStatusPending,
            newStatus:     entity.DeliveryStatusAssigned,
            expectError:   false,
        },
        {
            name:          "invalid transition",
            currentStatus: entity.DeliveryStatusPending,
            newStatus:     entity.DeliveryStatusDelivered,
            expectError:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assignment := &entity.DeliveryAssignment{
                Status: tt.currentStatus,
            }

            err := assignment.UpdateStatus(tt.newStatus)

            if tt.expectError {
                assert.Error(t, err)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.newStatus, assignment.Status)
            }
        })
    }
}
```

### Mock Repository Example

```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.DeliveryAssignment, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.DeliveryAssignment), args.Error(1)
}

func TestGetDeliveryAssignment(t *testing.T) {
    mockRepo := new(MockRepository)
    uc := usecase.NewDeliveryUseCase(mockRepo, logger)

    ctx := context.Background()
    id := uuid.New()

    expected := &entity.DeliveryAssignment{
        ID:      id,
        OrderID: "ORDER-123",
    }

    mockRepo.On("GetByID", ctx, id).Return(expected, nil)

    result, err := uc.GetDeliveryAssignment(ctx, id)

    require.NoError(t, err)
    assert.Equal(t, expected, result)
    mockRepo.AssertExpectations(t)
}
```

## Test Coverage Goals

- **Unit Tests**: Aim for >80% coverage
- **Integration Tests**: Cover critical paths
- **Edge Cases**: Test boundary conditions
- **Error Paths**: Test all error scenarios

## Testing Best Practices

### 1. Arrange-Act-Assert Pattern

```go
func TestExample(t *testing.T) {
    // Arrange: Set up test data and dependencies
    input := "test"
    
    // Act: Execute the function under test
    result := Function(input)
    
    // Assert: Verify the results
    assert.Equal(t, "expected", result)
}
```

### 2. Use Table-Driven Tests

For testing multiple scenarios:

```go
tests := []struct {
    name     string
    input    string
    expected string
}{
    {"case1", "input1", "output1"},
    {"case2", "input2", "output2"},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := Function(tt.input)
        assert.Equal(t, tt.expected, result)
    })
}
```

### 3. Use Testify Assertions

```go
// Equality
assert.Equal(t, expected, actual)

// Nil checks
assert.Nil(t, err)
assert.NotNil(t, result)

// Error checks
assert.Error(t, err)
assert.NoError(t, err)

// Require (stops test on failure)
require.NoError(t, err)
```

### 4. Mock External Dependencies

Always mock:
- Database calls
- External API calls
- Time-dependent operations (use fixed times in tests)

### 5. Clean Up Resources

```go
func TestWithCleanup(t *testing.T) {
    resource := setupResource()
    defer resource.Close()
    
    // Test code...
}
```

## Integration Testing

### Setup

Integration tests require a PostgreSQL database:

```bash
docker-compose up -d postgres
```

### Example Integration Test

```go
//go:build integration
// +build integration

func TestIntegration_CreateAndRetrieve(t *testing.T) {
    // Setup database connection
    db, err := setupTestDB()
    require.NoError(t, err)
    defer cleanupTestDB(db)
    
    // Create repository and use case
    repo := repository.NewPostgresRepository(db)
    uc := usecase.NewDeliveryUseCase(repo, logger)
    
    // Create delivery
    assignment, err := uc.CreateDeliveryAssignment(ctx, input)
    require.NoError(t, err)
    
    // Retrieve and verify
    retrieved, err := uc.GetDeliveryAssignment(ctx, assignment.ID)
    require.NoError(t, err)
    assert.Equal(t, assignment.ID, retrieved.ID)
}
```

## Performance Testing

### Benchmark Tests

```go
func BenchmarkCreateDeliveryAssignment(b *testing.B) {
    uc := setupUseCase()
    input := createTestInput()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = uc.CreateDeliveryAssignment(context.Background(), input)
    }
}
```

Run benchmarks:

```bash
go test -bench=. -benchmem ./...
```

## Continuous Integration

Tests run automatically on:
- Pull requests
- Pushes to main/develop branches

See `.github/workflows/ci.yml` for CI configuration.

## Troubleshooting

### Tests Failing Locally

1. **Database connection issues**:
   ```bash
   docker-compose up -d postgres
   ```

2. **Missing dependencies**:
   ```bash
   go mod download
   ```

3. **Proto files not generated**:
   ```bash
   make proto
   ```

### Race Condition Warnings

If you see race condition warnings:
1. Review concurrent access to shared resources
2. Add proper synchronization (mutexes, channels)
3. Run with `-race` flag to debug

### Coverage Reports Not Generated

Ensure you're using the coverage commands:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Test Data Helpers

Create reusable test data:

```go
func createTestAddress() entity.Address {
    return entity.Address{
        Street:     "123 Test St",
        City:       "Test City",
        State:      "TS",
        PostalCode: "12345",
        Country:    "USA",
        Latitude:   40.7128,
        Longitude:  -74.0060,
    }
}

func createTestAssignment() *entity.DeliveryAssignment {
    return entity.NewDeliveryAssignment(
        "ORDER-TEST",
        createTestAddress(),
        createTestAddress(),
        time.Now().Add(1*time.Hour),
        time.Now().Add(3*time.Hour),
        "Test notes",
    )
}
```

## Next Steps

1. Write tests for new features
2. Maintain >80% coverage
3. Run tests before committing
4. Review test failures in CI
5. Add integration tests for critical flows
