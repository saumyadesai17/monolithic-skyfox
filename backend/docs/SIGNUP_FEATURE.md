# Customer Signup Feature — Complete Developer Guide

## Table of Contents
1. [Overview](#overview)
2. [File Structure](#file-structure)
3. [Layer-by-Layer Breakdown](#layer-by-layer-breakdown)
   - [Model](#1-model--bookingsmodeluser_accountgo)
   - [DTOs](#2-dtos-data-transfer-objects)
   - [Repository](#3-repository--bookingsrepositoryuser_account_repositorygo)
   - [Service](#4-service--bookingsserviceauth_servicego)
   - [Controller](#5-controller--bookingscontrollerauth_controllergo)
   - [Router Registration](#6-router-registration--appserverstartgo)
   - [Error System](#7-error-system)
   - [Validation System](#8-validation-system)
4. [Full Request Flow](#full-request-flow)
5. [Test File Explained](#test-file-explained)
6. [Do I Need Integration Tests?](#do-i-need-integration-tests)
7. [Running the Tests](#running-the-tests)
8. [API Reference](#api-reference)

---

## Overview

The signup feature allows a **new customer to register an account** using their name, phone number, password, and optionally their email. It is a public endpoint — no authentication token is required.

**Endpoint:** `POST /api/v1/auth/signup`  
**Auth required:** No  
**Success response:** `201 Created`

---

## File Structure

```
backend/
├── bookings/
│   ├── model/
│   │   └── user_account.go          # GORM model → maps to `users` table
│   ├── dto/
│   │   ├── request/
│   │   │   └── signup_request.go    # What the client sends
│   │   └── response/
│   │       └── signup_response.go   # What the API returns
│   ├── repository/
│   │   └── user_account_repository.go  # DB read/write for users table
│   ├── service/
│   │   ├── auth_service.go          # All business logic
│   │   └── authService_test.go      # Unit tests for the service
│   └── controller/
│       └── auth_controller.go       # HTTP handler
├── common/
│   └── middleware/
│       └── validator/
│           ├── dto_validator.go     # Custom validation rules (passwordStrength, phoneNumber)
│           └── dto_validation_handler.go  # Translates validation errors to AppError
├── error/
│   ├── app_error.go                 # AppError struct
│   └── error_constructor.go        # Factory functions (BadRequestError, ConflictError, etc.)
└── app/
    └── server/
        └── start.go                 # Wires everything together, registers the route
```

---

## Layer-by-Layer Breakdown

### 1. Model — `bookings/model/user_account.go`

```go
type UserAccount struct {
    Id              string    `json:"id"    gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Phone           string    `json:"phone" gorm:"type:numeric(10);not null;unique"`
    Email           string    `json:"email,omitempty" gorm:"type:varchar(150)"`
    Name            string    `json:"name"  gorm:"type:varchar(100);not null"`
    AvatarUrl       string    `json:"-"     gorm:"type:text"`
    PasswordHash    string    `json:"-"     gorm:"type:text"`
    IsPhoneVerified bool      `json:"isPhoneVerified" gorm:"default:false"`
    IsEmailVerified bool      `json:"isEmailVerified" gorm:"default:false"`
    CreatedAt       time.Time `json:"createdAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
}

func (UserAccount) TableName() string { return "users" }
```

**Key design decisions:**

| Tag | What it does |
|-----|-------------|
| `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` | PostgreSQL generates the UUID — Go never sets this manually |
| `gorm:"type:numeric(10);not null;unique"` | DB-level unique constraint on phone, on top of the service-level duplicate check |
| `json:"-"` on `PasswordHash` | This field is **never included in any JSON response**, regardless of what the caller does |
| `json:"-"` on `AvatarUrl` | Not exposed at signup stage |
| `TableName()` returning `"users"` | Overrides GORM's default which would be `user_accounts` |

---

### 2. DTOs (Data Transfer Objects)

DTOs are the boundary objects — they define the exact shape of data coming **in** and going **out** of the API. The internal model is never exposed directly.

#### Request — `bookings/dto/request/signup_request.go`

```go
type SignupRequest struct {
    Name     string `json:"name"     binding:"required,min=2,max=100"`
    Phone    string `json:"phone"    binding:"required,phoneNumber"`
    Password string `json:"password" binding:"required,passwordStrength"`
    Email    string `json:"email"    binding:"omitempty,email"`
}
```

The `binding` tags are processed by Gin when `ShouldBindJSON` is called. Each tag is a validation rule:

| Field | Rules | What fails |
|-------|-------|-----------|
| `name` | `required`, `min=2`, `max=100` | Missing, single char, or over 100 chars |
| `phone` | `required`, `phoneNumber` | Missing, not exactly 10 digits |
| `password` | `required`, `passwordStrength` | Missing, or weak (see validator below) |
| `email` | `omitempty`, `email` | Only validated if provided — must be valid format |

`omitempty` on email means the field is **optional**. If the client doesn't send it, no error is thrown.

#### Response — `bookings/dto/response/signup_response.go`

```go
type SignupResponse struct {
    Id    string `json:"id"`
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

func NewSignupResponse(id, name, phone string) *SignupResponse {
    return &SignupResponse{Id: id, Name: name, Phone: phone}
}
```

The response deliberately exposes **only three fields** — id, name, phone. No email, no hash, no internal DB fields. This is the principle of least exposure — only return what the client actually needs.

---

### 3. Repository — `bookings/repository/user_account_repository.go`

The repository is the **only layer that talks to the database**. Nothing above it knows about GORM, SQL, or connection handling.

```go
type UserAccountRepository interface {
    FindByPhone(ctx context.Context, phone string) (*model.UserAccount, error)
    CreateUser(ctx context.Context, user *model.UserAccount) error
}
```

An **interface** is defined here so that:
- The service layer depends on the interface, not the concrete struct
- Tests can inject a mock that satisfies this interface without touching a real DB

#### `FindByPhone`

```go
func (r *userAccountRepository) FindByPhone(ctx context.Context, phone string) (*model.UserAccount, error) {
    var user model.UserAccount
    db, cancel := r.WithContext(ctx)
    defer cancel()

    result := db.Where("phone = ?", phone).First(&user)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil   // ← not found = nil, nil (not an error)
        }
        return nil, ae.InternalServerError(...)
    }
    return &user, nil
}
```

**Critical design:** When a phone is not found, this returns `nil, nil` — not a 404 error. The caller (service) checks if the returned pointer is nil. This is intentional — "user not found" is not an error during signup, it is the **happy path**.

#### `CreateUser`

```go
func (r *userAccountRepository) CreateUser(ctx context.Context, user *model.UserAccount) error {
    db, cancel := r.WithContext(ctx)
    defer cancel()

    if result := db.Create(user); result.Error != nil {
        return ae.InternalServerError(...)
    }
    return nil
}
```

`db.Create(user)` runs an `INSERT` statement. GORM automatically populates `user.Id` with the UUID generated by PostgreSQL after the insert completes.

Both methods use `r.WithContext(ctx)` which attaches the HTTP request context — this means if the HTTP request is cancelled (client disconnects), the DB query is also cancelled.

---

### 4. Service — `bookings/service/auth_service.go`

The service owns all **business logic**. It does not know about HTTP, JSON, or Gin — it just receives plain Go types and returns plain Go types or errors.

```go
type AuthUserRepository interface {
    FindByPhone(ctx context.Context, phone string) (*model.UserAccount, error)
    CreateUser(ctx context.Context, user *model.UserAccount) error
}
```

Notice the service defines its **own** interface for the repository — `AuthUserRepository`. This is a subset of the full `UserAccountRepository`. This pattern is called **Interface Segregation** — the service only declares what it needs, making it easy to mock in tests.

The `Signup` method runs **five steps in order**:

#### Step 1 — Password Strength (Service-Level Defence)

```go
if !isPasswordStrong(req.Password) {
    return nil, ae.BadRequestError("WeakPassword", "...", ...)
}
```

`isPasswordStrong` checks:
- Length >= 8
- At least one uppercase letter (`unicode.IsUpper`)
- At least one digit (`unicode.IsDigit`)
- At least one special character (`unicode.IsPunct || unicode.IsSymbol`)

This same validation also exists at the HTTP binding layer (via `passwordStrength` tag). Having it in **both** places is intentional — see [Validation System](#8-validation-system) for why.

#### Step 2 — Duplicate Phone Check

```go
existing, err := s.userRepo.FindByPhone(ctx, req.Phone)
if err != nil {
    return nil, err                        // DB error → 500
}
if existing != nil {
    return nil, ae.ConflictError("DuplicatePhone", "...", ...)   // 409
}
```

Checked **before** hashing because bcrypt is expensive (~300ms at cost 12). No point hashing a password for a phone that already exists.

#### Step 3 — Bcrypt Hash

```go
hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
```

- Cost `12` — each increment of cost doubles the computation time. Cost 12 ≈ 200–400ms on modern hardware, which is acceptable for signup but expensive enough to make brute-force attacks slow.
- Bcrypt automatically generates a random salt per hash — two calls with the same password produce **different** hashes.
- The hash is one-way — it cannot be reversed to recover the password.

#### Step 4 — Persist the User

```go
user := &model.UserAccount{
    Name:         req.Name,
    Phone:        req.Phone,
    Email:        req.Email,
    PasswordHash: string(hash),
}
s.userRepo.CreateUser(ctx, user)
```

After `CreateUser`, GORM populates `user.Id` with the UUID assigned by PostgreSQL.

#### Step 5 — Scrub Hash Before Returning

```go
user.PasswordHash = ""
return user, nil
```

The hash is cleared from the struct before it travels back up the call stack. The model also has `json:"-"` on `PasswordHash` as a second layer of protection, but explicitly blanking it here ensures no code path above the service can accidentally expose it — even through reflection or logging.

---

### 5. Controller — `bookings/controller/auth_controller.go`

The controller is the **HTTP adapter**. It translates between HTTP concepts (request body, status codes, JSON) and service concepts (Go structs, errors).

```go
type AuthService interface {
    Signup(ctx context.Context, req request.SignupRequest) (*model.UserAccount, error)
}
```

Again, a local interface — the controller depends on the interface, not the concrete `authService`. This mirrors the same pattern in the service layer and enables controller-level testing without a real service.

#### The Handler

```go
func (ac *AuthController) Signup(c *gin.Context) {
    var req request.SignupRequest

    // Step 1: Bind and validate the JSON body
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Error("signup validation failed: %v", err)
        c.AbortWithStatusJSON(http.StatusBadRequest, validator.HandleStructValidationError(err))
        return
    }

    // Step 2: Call the service
    user, err := ac.authService.Signup(c.Request.Context(), req)
    if err != nil {
        appErr := err.(*ae.AppError)
        logger.Error("signup error: %v", appErr)
        c.AbortWithStatusJSON(appErr.HTTPCode(), appErr)
        return
    }

    // Step 3: Return the response
    c.IndentedJSON(http.StatusCreated, response.NewSignupResponse(user.Id, user.Name, user.Phone))
}
```

Key points:
- `c.ShouldBindJSON` — binds the body AND runs all `binding:` tag validators before the service is touched
- `c.Request.Context()` — passes the HTTP request context down to the service and DB, enabling cancellation
- `appErr.HTTPCode()` — the service returns `*AppError` which carries its own HTTP status. The controller just asks the error what status to use — it doesn't hardcode anything
- `c.AbortWithStatusJSON` — stops any further middleware from running in addition to writing the response

---

### 6. Router Registration — `app/server/start.go`

```go
// Route is on noAuthRouter — no JWT/session middleware
v1 := noAuthRouter.Group("/api/v1")
{
    auth := v1.Group("/auth")
    {
        auth.POST("/signup", authController.Signup)
    }
}
```

The route is registered on `noAuthRouter`, not `authRouter`. This is what makes it a **public endpoint** — requests to `/api/v1/auth/signup` bypass the authentication middleware entirely. Users can't have a token before they have an account.

The custom validator is also registered here at server startup:

```go
binding.Validator = new(validator.DtoValidator)
```

This replaces Gin's default validator globally, so all `ShouldBindJSON` calls in the entire application use the custom rules including `phoneNumber`, `passwordStrength`, and `maxSeats`.

---

### 7. Error System

All errors in this codebase are `*AppError`, defined in `error/app_error.go`:

```go
type AppError struct {
    error           // wrapped original error (for logging/debugging)
    httpCode int    // HTTP status to return
    Code     string // machine-readable error code (e.g. "DuplicatePhone")
    Message  string // human-readable message returned to the client
}
```

Factory functions in `error/error_constructor.go` create typed errors:

```go
BadRequestError(...)    → 400
ConflictError(...)      → 409
InternalServerError(...) → 500
NotFoundError(...)      → 404
```

The controller type-asserts the error to `*AppError` and calls `.HTTPCode()` to get the status code. This design means:
- Each layer signals **what went wrong semantically** (DuplicatePhone, WeakPassword)
- The controller translates that to HTTP without needing to know the details

---

### 8. Validation System

Validation lives in **two layers** — both are intentional:

#### Layer 1 — HTTP Binding (`common/middleware/validator/dto_validator.go`)

`DtoValidator` implements Gin's `binding.Validator` interface. It is initialised once via `sync.Once` and registers custom rules via `go-playground/validator`:

```go
d.validate.RegisterValidation("passwordStrength", validatePasswordStrength())
d.validate.RegisterValidation("phoneNumber", validatePhoneNumber())
```

When `ShouldBindJSON` runs, these custom rules execute alongside the built-in ones (`required`, `min`, `max`, `email`). A weak password is rejected **before the controller even calls the service**.

`dto_validation_handler.go` translates the raw `validator.ValidationErrors` into a structured `AppError` with a human-readable message. It uses `errors.As` (not a raw type assertion) so it handles non-field-validation errors (e.g. malformed JSON body) without panicking.

#### Layer 2 — Service (`isPasswordStrong` in `auth_service.go`)

This is the **defence in depth** layer. The service validates the password independently, because:
- The service can be called from sources other than HTTP (scripts, other services, tests, future gRPC handlers)
- The service must never trust that upstream validation already ran
- Password strength is a **business rule**, not just an HTTP concern

If only the HTTP layer validated, any non-HTTP caller could write a weak-password user directly to the database.

---

## Full Request Flow

```
POST /api/v1/auth/signup
         │
         ▼
  [noAuthRouter] — skips auth middleware, goes straight to handler
         │
         ▼
  AuthController.Signup()
         │
         ├── c.ShouldBindJSON(&req)
         │         │
         │     DtoValidator runs binding tags:
         │       - required, min=2, max=100    → name
         │       - required, phoneNumber       → phone (must be 10 digits)
         │       - required, passwordStrength  → password (upper + digit + special)
         │       - omitempty, email            → email (optional, format check)
         │         │
         │         └── FAIL → HandleStructValidationError → 400 Bad Request
         │
         ▼ PASS
  authService.Signup(ctx, req)
         │
         ├── isPasswordStrong(req.Password)
         │         └── FAIL → 400 WeakPassword
         │
         ├── userRepo.FindByPhone(ctx, req.Phone)
         │         ├── DB error  → 500 InternalServerError
         │         ├── found     → 409 DuplicatePhone
         │         └── nil       → continue
         │
         ├── bcrypt.GenerateFromPassword(password, cost=12)
         │         └── error     → 500 InternalServerError
         │
         ├── userRepo.CreateUser(ctx, &UserAccount{...hash...})
         │         └── DB error  → 500 InternalServerError
         │
         ├── user.PasswordHash = ""   ← scrub before returning
         │
         └── return *UserAccount, nil
                   │
                   ▼
  AuthController builds SignupResponse{id, name, phone}
         │
         ▼
  201 Created
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Alice",
    "phone": "9876543210"
  }
```

---

## Test File Explained

**File:** `bookings/service/authService_test.go`  
**Package:** `service` (white-box test — same package as the code)  
**Framework:** `testify/assert` + `testify/mock`  
**Mocks:** Auto-generated by `mockery` in `bookings/service/mocks/`

### Setup

```go
var validReq = request.SignupRequest{
    Name:     "Alice",
    Phone:    "9876543210",
    Password: "Str0ng@Pass",
    Email:    "alice@example.com",
}
```

A **package-level variable** holding a valid request used as the base for all test cases. Individual test cases override fields (e.g. `weakPassword: true` swaps the password to `"weak"`).

### Test Structure

```go
tests := []struct {
    name          string
    setupMock     func(repo *servicemocks.MockAuthUserRepository)
    wantHTTPCode  int
    wantErrCode   string
    wantUser      bool   // true = expect a non-nil user
    wantEmptyHash bool   // true = PasswordHash must be blank on returned user
    weakPassword  bool   // true = override password with a weak one
}{...}
```

This is the **table-driven test pattern** — each test case is a row in a table. All cases run through the same loop with the same assertions structure. This makes it easy to add new scenarios without duplicating test runner code.

### The Mock

```go
repo := servicemocks.NewMockAuthUserRepository(t)
```

`MockAuthUserRepository` is auto-generated by `mockery` from the `AuthUserRepository` interface in `auth_service.go`. It lives in `bookings/service/mocks/mock_AuthUserRepository.go`. **Never edit it manually** — always run `mockery` to regenerate.

Passing `t` to the constructor means testify will automatically fail the test if any expected mock call is not made (or an unexpected one is).

### Test Cases Explained

| Test Case | What it proves |
|-----------|---------------|
| `Should create user and not return password hash when phone is not registered` | Happy path — mock returns nil for FindByPhone (phone is free), CreateUser runs and sets Id via `.Run()`, returned user has no hash |
| `Should return 400 when password does not meet strength requirements` | `isPasswordStrong` fires first — no mock calls happen at all because the service returns before touching the repo |
| `Should return 409 Conflict when phone is already registered` | FindByPhone returns an existingUser — service returns ConflictError without calling CreateUser |
| `Should return 500 when FindByPhone returns a database error` | FindByPhone returns a DB error — service propagates it as 500 |
| `Should return 500 when CreateUser returns a database error` | FindByPhone returns nil (phone is free), but CreateUser fails — service returns 500 |
| `Should store password as bcrypt hash when user is created` | `mock.MatchedBy` inspects the UserAccount passed to CreateUser and verifies the hash is a valid bcrypt hash of the original password |

### The `.Run()` Pattern (Happy Path)

```go
repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.UserAccount) bool {
    return u.Phone == validReq.Phone && u.Name == validReq.Name
})).
    Run(func(args mock.Arguments) {
        u := args.Get(1).(*model.UserAccount)
        u.Id = "new-uuid-1234"   // simulate what PostgreSQL does
    }).
    Return(nil)
```

In real PostgreSQL, after `db.Create(user)`, GORM populates `user.Id` with the UUID from the DB. In tests there is no DB, so `.Run()` manually sets the Id on the struct passed to the mock — simulating the DB side effect so the test can assert `user.Id == "new-uuid-1234"`.

### The `mock.MatchedBy` Pattern (Hash Verification)

```go
repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.UserAccount) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(validReq.Password))
    return err == nil && u.PasswordHash != validReq.Password
}))
```

This does not just check that *something* was passed — it calls `bcrypt.CompareHashAndPassword` on the hash being saved to confirm:
1. The hash is a valid bcrypt hash of the original password
2. The hash is not the plain password

---

## Do I Need Integration Tests?

**Yes, you should write one.** Here is why and what it should cover.

### What Unit Tests Already Cover ✅

The service unit tests (`authService_test.go`) test the business logic in full isolation — every success case, every error branch, the password check, the hash verification, the duplicate check. That is complete.

### What Unit Tests Cannot Cover ❌

Unit tests mock the repository. They never touch a real database. So these things are **untested**:

- Does GORM actually map the struct to the right table (`users`)?
- Does the `phone` unique constraint in the DB actually reject a duplicate?
- Does `gen_random_uuid()` actually populate the `Id` field after insert?
- Does the full HTTP stack — routing → binding → controller → service → repo → DB — work end to end?

### What Integration Test to Write

Following the pattern in `integration_test/booking_controller_test.go`, the signup integration test should:

```go
// integration_test/signup_controller_test.go
func Test_WhenSignup_WithValidPayload_ItShouldReturn201(t *testing.T) { ... }
func Test_WhenSignup_WithDuplicatePhone_ItShouldReturn409(t *testing.T) { ... }
func Test_WhenSignup_WithWeakPassword_ItShouldReturn400(t *testing.T) { ... }
func Test_WhenSignup_WithMissingField_ItShouldReturn400(t *testing.T) { ... }
```

### How Integration Tests Work in This Project

- `TestMain` in `integration_test/setup_test.go` spins up a **real PostgreSQL instance in Docker** using `testcontainers-go`
- Each test gets a real DB connection, runs real migrations, and makes real HTTP requests through the full Gin engine
- After the test suite, the container is terminated

So the integration test proves the entire vertical slice — HTTP → validation → service → GORM → real PostgreSQL — works correctly in combination.

**Short answer: unit tests are complete for the service logic. Write one integration test to prove the full HTTP-to-DB stack works.**

---

## Running the Tests

### Unit Tests (no DB required)

```powershell
# Run just the auth service tests
go test ./bookings/service/... -v -run TestAuthService_Signup

# Run all service tests
go test ./bookings/service/... -v

# Run with coverage report
go test ./bookings/service/... -cover
```

### File-Mode (test specific files only)

```powershell
cd bookings/service
go test -v -cover auth_service.go authService_test.go
```

### Regenerate Mocks (after changing an interface)

```powershell
cd backend
mockery
```

This reads `.mockery.yml` and regenerates all mock files into `{package}/mocks/mock_{Interface}.go`.

---

## API Reference

### `POST /api/v1/auth/signup`

**Request Body:**

```json
{
  "name": "Alice",
  "phone": "9876543210",
  "password": "Str0ng@Pass",
  "email": "alice@example.com"
}
```

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `name` | string | Yes | 2–100 characters |
| `phone` | string | Yes | Exactly 10 digits |
| `password` | string | Yes | Min 8 chars, 1 uppercase, 1 digit, 1 special character |
| `email` | string | No | Valid email format if provided |

**Responses:**

| Status | Code | When |
|--------|------|------|
| `201 Created` | — | User created successfully |
| `400 Bad Request` | `ValidationFailed` | Missing field or format error |
| `400 Bad Request` | `WeakPassword` | Password doesn't meet strength rules |
| `409 Conflict` | `DuplicatePhone` | Phone number already registered |
| `500 Internal Server Error` | `InternalServerError` | Database error |

**201 Response Body:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Alice",
  "phone": "9876543210"
}
```

**400 Response Body:**

```json
{
  "Code": "WeakPassword",
  "Message": "password must be at least 8 characters and contain an uppercase letter, a digit, and a special character"
}
```

**409 Response Body:**

```json
{
  "Code": "DuplicatePhone",
  "Message": "an account with this phone number already exists"
}
```
