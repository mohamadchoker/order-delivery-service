# Naming Conventions Quick Reference

## Proto Files

| Type | Convention | Example | ❌ Wrong |
|------|-----------|---------|----------|
| **File name** | lowercase, snake_case | `driver.proto` | `driverService.proto` |
| **Package** | lowercase, singular | `package driver;` | `package drivers;` |
| **Service** | PascalCase + "Service" | `service DriverService` | `service Driver` |
| **Message** | PascalCase | `message Driver` | `message driver` |
| **Request** | PascalCase + "Request" | `CreateDriverRequest` | `DriverCreate` |
| **Response** | PascalCase + "Response" | `ListDriversResponse` | `DriverList` |
| **Enum** | PascalCase | `enum DriverStatus` | `enum driverStatus` |
| **Enum values** | UPPERCASE_SNAKE_CASE | `AVAILABLE`, `IN_TRANSIT` | `Available` |

## Go Files

### Domain Layer (`internal/domain/`)

| Type | Convention | Example |
|------|-----------|---------|
| **File name** | singular, lowercase | `driver.go` |
| **Entity struct** | PascalCase, singular | `type Driver struct` |
| **Status const** | PascalCase prefix | `DriverStatusAvailable` |
| **Constructor** | `New{Entity}` | `func NewDriver()` |
| **Methods** | PascalCase, verb first | `func (d *Driver) UpdateStatus()` |

### Repository Layer (`internal/repository/postgres/`)

| Type | Convention | Example |
|------|-----------|---------|
| **Interface file** | `{entity}_repository.go` | `driver_repository.go` |
| **Impl file** | `{entity}_repository_impl.go` | `driver_repository_impl.go` |
| **Interface name** | `{Entity}Repository` | `type DriverRepository interface` |
| **Impl struct** | `{entity}Repository` (private) | `type driverRepository struct` |
| **Constructor** | `New{Entity}Repository` | `func NewDriverRepository()` |
| **Methods** | CRUD verbs | `Create`, `GetByID`, `List`, `Update`, `Delete` |

### Database Model (`internal/repository/postgres/model/`)

| Type | Convention | Example |
|------|-----------|---------|
| **File name** | singular | `driver.go` |
| **Struct name** | PascalCase, singular | `type Driver struct` |
| **Table name** | plural, snake_case | `"drivers"` |
| **Converter to domain** | `ToDomain()` | `func (d *Driver) ToDomain()` |
| **Converter from domain** | `{Entity}FromDomain()` | `func DriverFromDomain()` |

### Service Layer (`internal/service/`)

| Type | Convention | Example |
|------|-----------|---------|
| **File name** | `{entity}_usecase.go` | `driver_usecase.go` |
| **Interface name** | `{Entity}UseCase` | `type DriverUseCase interface` |
| **Impl struct** | `{entity}UseCase` (private) | `type driverUseCase struct` |
| **Constructor** | `New{Entity}UseCase` | `func NewDriverUseCase()` |
| **Input struct** | `{Action}{Entity}Input` | `type CreateDriverInput struct` |
| **Methods** | Verb + Entity | `CreateDriver`, `GetDriver`, `ListDrivers` |

### Transport Layer (`internal/transport/grpc/`)

| Type | Convention | Example |
|------|-----------|---------|
| **Handler file** | `{entity}_handler.go` | `driver_handler.go` |
| **Converter file** | `{entity}_converter.go` | `driver_converter.go` |
| **Handler struct** | `{Entity}Handler` | `type DriverHandler struct` |
| **Constructor** | `New{Entity}Handler` | `func NewDriverHandler()` |
| **Converter funcs** | `{proto/domain}To{Domain/Proto}` | `protoToDriver()`, `driverToProto()` |

## Database

| Type | Convention | Example |
|------|-----------|---------|
| **Table name** | plural, snake_case | `drivers`, `delivery_assignments` |
| **Column name** | snake_case | `created_at`, `driver_id` |
| **Primary key** | `id` | `id UUID PRIMARY KEY` |
| **Foreign key** | `{table}_id` | `driver_id`, `order_id` |
| **Index** | `idx_{table}_{column}` | `idx_drivers_status` |
| **Migration** | `{seq}_{description}` | `000002_create_drivers_table` |

## Variables & Functions

| Type | Convention | Example |
|------|-----------|---------|
| **Public var/func** | PascalCase | `DriverRepository`, `NewDriver` |
| **Private var/func** | camelCase | `driverRepo`, `validateInput` |
| **Constants** | PascalCase with prefix | `DriverStatusAvailable` |
| **Package-level const** | PascalCase or UPPER | `DefaultPageSize` |
| **Receiver** | Single letter or short | `(d *Driver)`, `(uc *useCase)` |

## HTTP/gRPC Endpoints

| Type | Convention | Example |
|------|-----------|---------|
| **RPC method** | PascalCase, verb first | `CreateDriver`, `GetDriver`, `ListDrivers` |
| **Package** | lowercase, singular | `driver`, `delivery` |

## Common Verbs

| Action | Proto RPC | Go Method | SQL |
|--------|-----------|-----------|-----|
| Create new | `CreateDriver` | `CreateDriver` | `INSERT` |
| Get one | `GetDriver` | `GetDriver` or `GetByID` | `SELECT` |
| Get many | `ListDrivers` | `ListDrivers` or `List` | `SELECT` |
| Update | `UpdateDriver` | `UpdateDriver` or `Update` | `UPDATE` |
| Delete | `DeleteDriver` | `DeleteDriver` or `Delete` | `DELETE` |
| Update status | `UpdateDriverStatus` | `UpdateStatus` | `UPDATE` |

## Examples

### ✅ Good Naming

```go
// Domain
type Driver struct { ... }
const DriverStatusAvailable DriverStatus = "AVAILABLE"
func NewDriver() *Driver { ... }
func (d *Driver) UpdateStatus() { ... }

// Repository
type DriverRepository interface { ... }
type driverRepository struct { ... }
func NewDriverRepository() DriverRepository { ... }
func (r *driverRepository) Create(ctx context.Context, driver *domain.Driver) error { ... }

// Service
type DriverUseCase interface { ... }
type driverUseCase struct { ... }
type CreateDriverInput struct { ... }
func NewDriverUseCase() DriverUseCase { ... }
func (uc *driverUseCase) CreateDriver(ctx context.Context, input CreateDriverInput) (*domain.Driver, error) { ... }

// Transport
type DriverHandler struct { ... }
func NewDriverHandler() *DriverHandler { ... }
func driverToProto(d *domain.Driver) *pb.Driver { ... }
func protoStatusToDriverDomain(s pb.DriverStatus) domain.DriverStatus { ... }
```

### ❌ Bad Naming

```go
// Wrong!
type Drivers struct { ... }                    // Should be singular
const DRIVER_STATUS_AVAILABLE = "available"    // Use PascalCase
func CreateNewDriver() *Driver { ... }         // Too verbose
func (driver *Driver) Status() { ... }         // Unclear

type driverRepo interface { ... }              // Interface should be public
func NewRepo() DriverRepository { ... }        // Not specific enough
func (r *repo) GetDriver() { ... }             // Receiver name too generic

type driver_use_case struct { ... }            // Use camelCase for private
func (u *DriverUseCase) Create() { ... }       // Not specific enough

type driverHandler struct { ... }              // Handler should be public
func convertDriver() { ... }                   // Which direction?
```

## Proto-to-Go Mappings

| Proto Type | Go Type |
|------------|---------|
| `string` | `string` |
| `int32` | `int32` |
| `int64` | `int64` |
| `bool` | `bool` |
| `double` | `float64` |
| `float` | `float32` |
| `bytes` | `[]byte` |
| `repeated string` | `[]string` |
| `message Address` | `*Address` |
| `repeated message` | `[]*Message` |
| `google.protobuf.Timestamp` | `time.Time` (via `timestamppb`) |
| `enum Status` | Custom type (e.g., `type Status string`) |

## Acronyms

When using acronyms, use consistent casing:

| Context | Convention | Example |
|---------|-----------|---------|
| **Start of name** | All caps | `IDGenerator`, `URLBuilder` |
| **Middle/end** | All caps | `UserID`, `GetURL` |
| **Proto** | Uppercase | `ID`, `URL`, `API` |

Common acronyms: `ID`, `URL`, `HTTP`, `JSON`, `API`, `UUID`, `DB`, `SQL`

---

## Quick Decision Tree

**Naming a file?**
- Domain entity? → `{entity}.go` (singular)
- Repository? → `{entity}_repository.go`
- Service? → `{entity}_usecase.go`
- Handler? → `{entity}_handler.go`
- Converter? → `{entity}_converter.go`
- Proto? → `{service}.proto`

**Naming a struct?**
- Is it public? → PascalCase
- Is it private? → camelCase
- Is it a model? → Singular noun
- Is it a collection? → Plural (rare, usually use slice)

**Naming a function?**
- Constructor? → `New{Type}`
- Converter to domain? → `ToDomain()`
- Converter from domain? → `{Entity}FromDomain()`
- Proto to domain? → `protoTo{Entity}()`
- Domain to proto? → `{entity}ToProto()`

**Naming a method?**
- CRUD? → `Create`, `GetByID`, `List`, `Update`, `Delete`
- Business logic? → Verb first (e.g., `UpdateStatus`, `AssignDriver`)

---

Follow these conventions consistently and your code will be clean, predictable, and easy to navigate!
