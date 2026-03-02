# Skyfox Backend — Complete Developer Guide

> This document explains **everything** you need to know about the backend before adding new features. Read it top to bottom once; afterwards use it as a reference.

---

## Table of Contents

1. [What Is This App?](#1-what-is-this-app)
2. [Tech Stack](#2-tech-stack)
3. [Project Layout](#3-project-layout)
4. [How the App Boots (Startup Flow)](#4-how-the-app-boots-startup-flow)
5. [Configuration](#5-configuration)
6. [Database — Connection & BaseDB](#6-database--connection--basedb)
7. [Schema Migrations](#7-schema-migrations)
8. [Domain Models](#8-domain-models)
9. [Repositories (Data Access Layer)](#9-repositories-data-access-layer)
10. [Services (Business Logic Layer)](#10-services-business-logic-layer)
11. [Controllers (HTTP Handler Layer)](#11-controllers-http-handler-layer)
12. [HTTP Routes & Endpoints](#12-http-routes--endpoints)
13. [Middleware](#13-middleware)
14. [DTOs — Request & Response Shapes](#14-dtos--request--response-shapes)
15. [Error Handling](#15-error-handling)
16. [External Dependency — Movie Service Gateway](#16-external-dependency--movie-service-gateway)
17. [Seeding the Database](#17-seeding-the-database)
18. [Logging](#18-logging)
19. [Testing Strategy](#19-testing-strategy)
20. [Docker & Local Development](#20-docker--local-development)
21. [Full Request Lifecycle — End to End](#21-full-request-lifecycle--end-to-end)
22. [Business Rules & Constants](#22-business-rules--constants)
23. [Dependency Graph (Summary)](#23-dependency-graph-summary)
24. [Where to Add New Features](#24-where-to-add-new-features)

---

## 1. What Is This App?

**Skyfox** is a **movie ticket booking system**. The backend is a REST API that lets:

- Users **log in** with Basic Auth.
- Authenticated users **browse shows** available on a given date.
- Authenticated users **book seats** for a show.
- Authenticated users **query revenue** earned on a given date.

The backend talks to a separate **Movie Service** (Ruby / Go micro-service in `movie-service/`) over HTTP to fetch movie metadata (title, plot, runtime). The backend itself manages shows, slots, bookings, customers, and users.

---

## 2. Tech Stack

| Concern | Library / Tool |
|---|---|
| Language | Go 1.22 |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) with `gorm.io/driver/postgres` |
| Database | PostgreSQL |
| Migrations | [golang-migrate/migrate](https://github.com/golang-migrate/migrate) |
| Config loading | [Viper](https://github.com/spf13/viper) |
| Logging | [Uber Zap](https://go.uber.org/zap) via `gin-contrib/zap` |
| Validation | `go-playground/validator/v10` |
| Swagger docs | `swaggo/swag` + `swaggo/gin-swagger` |
| Integration tests | `testcontainers-go` (real Postgres in a Docker container) |
| Unit test HTTP | `appleboy/gofight` |
| HTTP mocking | `gopkg.in/h2non/gock.v1` |

---

## 3. Project Layout

```
backend/
├── main.go                          ← Entry point
├── go.mod / go.sum                  ← Go module definition
├── config/
│   ├── config.go                    ← Config structs + Viper loader
│   ├── config.yml                   ← Production config template (env-var based)
│   └── config-local.yml             ← Local dev overrides
├── app/
│   └── server/
│       └── start.go                 ← Wires everything together, sets up router
├── bookings/
│   ├── constants/
│   │   └── app_constant.go          ← Route paths, seat limits
│   ├── controller/                  ← HTTP handlers (one file per domain)
│   ├── service/                     ← Business logic (one file per domain)
│   ├── repository/                  ← DB queries (one file per domain)
│   ├── model/                       ← GORM models (DB-mapped structs)
│   ├── dto/
│   │   ├── request/                 ← Incoming JSON shapes
│   │   └── response/                ← Outgoing JSON shapes
│   └── database/
│       ├── common/base_db.go        ← BaseDB wrapper (timeout, gorm instance)
│       ├── connection/connection.go ← Singleton Postgres connection
│       └── seed/dataseeder.go       ← Creates default users at startup
├── common/
│   ├── logger/logger.go             ← Zap logger initializer
│   └── middleware/
│       ├── cors/                    ← CORS middleware
│       ├── security/basic_auth.go   ← Basic Auth middleware
│       └── validator/               ← Custom validator + error formatter
├── error/
│   ├── app_error.go                 ← AppError type
│   └── error_constructor.go         ← NotFound, BadRequest, etc. helpers
├── migration/
│   ├── migration.go                 ← Standalone binary to run SQL migrations
│   └── scripts/                     ← Numbered SQL up/down files
├── movieservice/
│   └── movie_gateway/               ← HTTP client to the Movie Service
├── integration_test/                ← End-to-end tests using a real DB container
├── _mocks/repomocks/                ← Hand-written mocks for repositories
└── docs/                            ← Auto-generated Swagger JSON/YAML
```

---

## 4. How the App Boots (Startup Flow)

**`main.go`** is the sole entry point.

```
main()
 │
 ├── flag.Parse()                        ← accepts -configFile=/path/to/config.yml
 ├── config.LoadConfig(configFile)       ← reads YAML + expands env vars
 └── server.Init(cfg)                    ← does everything below
```

**`app/server/start.go` → `Init(cfg)`** does the following in order:

```
1.  logger.InitAppLogger(cfg.Logger)         ← sets up Zap logger

2.  connection.NewDBHandler(cfg.Database)    ← builds Postgres DSN
    handler.Instance()                       ← opens GORM connection (singleton)

3.  movieservice.NewMovieGateway(cfg)        ← creates HTTP client to Movie Service

4.  Instantiate Repositories
      bookingRepository  = repository.NewBookingRepository(db)
      showRepository     = repository.NewShowRepository(db)
      userRepository     = repository.NewUserRepository(db)
      customerRepository = repository.NewCustomerRepository(db)

5.  database.SeedDB(userRepository)          ← creates seed-user-1 & seed-user-2 if missing

6.  Instantiate Services
      bookingService  = service.NewBookingService(bookingRepo, showRepo)
      bookingService.SetCustomerRepository(customerRepo)
      showService     = service.NewShowService(showRepo, movieGateway)
      userService     = service.NewUserService(userRepo)
      revenueService  = service.NewRevenueService(bookingRepo, showRepo)

7.  Instantiate Controllers
      bookingController = controller.NewBookingController(bookingService)
      showController    = controller.NewShowController(showService)
      userController    = controller.NewUserController(userService)
      revenueController = controller.NewRevenueController(revenueService)

8.  setupApp(cfg)                            ← creates Gin engine, registers middleware
    routerGroupWithAuth(...)                 ← routes protected by Basic Auth
    routerGroupWithNoAuth(...)               ← public routes

9.  Register routes (see Section 12)

10. http.Server.ListenAndServe()             ← starts listening
```

---

## 5. Configuration

Config is loaded from a YAML file + environment variable expansion.

**`config/config.go`** defines:

```go
type AppConfig struct {
    Server       ServerConfig        // HTTP server settings
    Database     DbConfig            // Postgres settings
    MovieGateway MovieGatewayConfig  // URL of the movie microservice
    Logger       LoggerConfig        // Log level
}
```

**`config/config.yml`** (the real config template):

```yaml
Server:
  port: ${SERVER_PORT}       # e.g. 8080
  ReadTimeout: 5             # seconds
  WriteTimeout: 5
  GineMode: debug

Database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
  user: ${POSTGRES_USERNAME}
  password: ${POSTGRES_PASSWORD}
  MigrationPath: migration/scripts

Logger:
  Level: debug

MovieGateway:
  movieServiceHost: ${MOVIE_SERVICE_HOST}   # e.g. http://localhost:3000/
```

All values with `${VAR}` are **read from environment variables** at startup via `os.ExpandEnv`. This makes the app 12-factor compliant — no secrets baked in.

For local dev, `config-local.yml` overrides these values with hardcoded local addresses.

---

## 6. Database — Connection & BaseDB

### Connection (`bookings/database/connection/connection.go`)

- Uses a **`sync.Once`** to ensure only one GORM connection is ever opened (singleton pattern).
- Builds a Postgres DSN string: `host=... user=... password=... dbname=... port=...`
- Opens with `gorm.Open(postgres.Open(dsn), &gorm.Config{})`.

### BaseDB (`bookings/database/common/base_db.go`)

Every repository embeds `*common.BaseDB`. It provides two things:

1. **`WithContext(ctx)`** — wraps every DB query in a **5-second timeout** using `context.WithTimeout`. Always call `defer cancel()` after using this.
2. **`SqlDB()`** / **`GormDB()`** — escape hatches to the raw `*sql.DB` or `*gorm.DB` if needed.

**Pattern used in every repository:**
```go
db, cancel := repo.WithContext(ctx)
defer cancel()
db.Where(...).Find(&results)
```

---

## 7. Schema Migrations

> **IMPORTANT — Current State of the Project**
>
> The database schema has been **completely redesigned** (migrations 006–021 applied). The old tables have been dropped and replaced with a richer, multi-theatre model. **The Go source code (models, repositories, services, controllers) has NOT been updated yet** — that is the next phase of work. Sections 8–21 of this guide still describe the old code behaviour for reference. Use this section as the source of truth for what the database actually looks like right now.

Migrations live in **`migration/scripts/`** as plain numbered SQL files. Each migration has a matching `.down.sql` that reverses it.

`migration/migration.go` is a **separate `main` package** (not the API server). Run migrations like:

```bash
go run migration/migration.go -configFile=config/config-local.yml up    # apply all pending
go run migration/migration.go -configFile=config/config-local.yml down  # rollback one step
```

Or via the Makefile: `make migrate-up` / `make migrate-down`.

---

### Migration History

#### Phase 1 — Original Schema (001–005) — DROPPED

These created the initial simple schema and have since been rolled forward through Phase 2.

| # | Script | Old Table Created |
|---|---|---|
| 1 | `000001_create_slot_table` | `slot` |
| 2 | `000002_create_show_table` | `show` |
| 3 | `000003_create_customer_table` | `customer` |
| 4 | `000004_create_booking_table` | `booking` |
| 5 | `000005_create_user_table` | `usertable` |

#### Phase 2 — Teardown (006–010)

Dropped all old tables in reverse dependency order to make way for the new schema.

| # | Script | Drops |
|---|---|---|
| 6 | `000006_drop_booking_table` | `booking` |
| 7 | `000007_drop_show_table` | `show` |
| 8 | `000008_drop_customer_table` | `customer` |
| 9 | `000009_drop_slot_table` | `slot` |
| 10 | `000010_drop_user_table` | `usertable` |

#### Phase 3 — New Schema (011–021) — ACTIVE

| # | Script | New Table / Object |
|---|---|---|
| 11 | `000011_create_admins_table` | `admins` |
| 12 | `000012_create_users_table` | `users` |
| 13 | `000013_create_theatre_table` | `theatres` |
| 14 | `000014_create_screens_table` | `screens` |
| 15 | `000015_create_show_table` | `show` (redesigned) |
| 16 | `000016_create_seat_types_table` | `seat_types` |
| 17 | `000017_create_seat_table` | `seats` |
| 18 | `000018_create_bookings_table` | `bookings` (redesigned) |
| 19 | `000019_create_booking_seats_table` | `booking_seats` |
| 20 | `000020_create_seat_lock_table` | `seat_lock` |
| 21 | `000021_create_payments_table` | `payments` |

---

### New Schema in Detail (Phase 3)

#### `admins`
```sql
id            UUID PRIMARY KEY DEFAULT gen_random_uuid()
username      VARCHAR(255) NOT NULL UNIQUE
password_hash VARCHAR(255) NOT NULL          -- hashed, unlike the old plain-text usertable
role          admin_role NOT NULL DEFAULT 'admin'  -- ENUM: 'owner' | 'admin'
created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
```
Replaces the old `usertable`. Now uses UUIDs, hashed passwords, and role-based access.

#### `users`
```sql
id                UUID PRIMARY KEY DEFAULT gen_random_uuid()
phone             NUMERIC(10) NOT NULL UNIQUE    -- primary identifier (10-digit)
email             VARCHAR(150) UNIQUE
name              VARCHAR(100) NOT NULL
avatar_url        TEXT
password_hash     TEXT
counter_no        VARCHAR(50)                    -- for counter/box-office users
is_phone_verified BOOLEAN NOT NULL DEFAULT FALSE
is_email_verified BOOLEAN NOT NULL DEFAULT FALSE
created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
updated_at        TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
```
Book-buying customers. Identified by phone number. Supports both app and counter bookings.

#### `theatres`
```sql
id                UUID PRIMARY KEY DEFAULT gen_random_uuid()
name              VARCHAR(255) NOT NULL
location          TEXT NOT NULL
number_of_screens INT NOT NULL
```
Top-level venue. A theatre has one or more screens.

#### `screens`
```sql
id          UUID PRIMARY KEY DEFAULT gen_random_uuid()
theatre_id  UUID REFERENCES theatres(id)
name        VARCHAR(50) NOT NULL         -- e.g. "Screen 1", "IMAX"
total_seats INT NOT NULL
```
A physical auditorium inside a theatre. Seats are defined per screen.

#### `show` (redesigned)
```sql
id                   UUID PRIMARY KEY DEFAULT gen_random_uuid()
movie_imdb_id        VARCHAR(20) NOT NULL             -- IMDb ID, fetched from Movie Service
screen_id            UUID NOT NULL REFERENCES screens(id)
theatre_id           UUID NOT NULL REFERENCES theatres(id)
start_time           TIMESTAMP NOT NULL               -- full datetime, not just a slot reference
end_time             TIMESTAMP NOT NULL
status               VARCHAR(20) DEFAULT 'ACTIVE'
seat_layout_snapshot JSONB                            -- snapshot of seat layout at show creation
created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP
updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```
Key changes from old schema: linked to a specific **screen** (not just a slot), uses full **timestamps** instead of date+slot, has a **status** field, and stores a **JSONB snapshot** of the seat layout.

#### `seat_types`
```sql
id    UUID PRIMARY KEY DEFAULT gen_random_uuid()
name  VARCHAR(50) NOT NULL        -- e.g. "Regular", "Premium", "Recliner"
price DECIMAL(10,2) NOT NULL      -- base price for this category
```
Defines pricing tiers. Seats belong to a type.

#### `seats`
```sql
id           UUID PRIMARY KEY DEFAULT gen_random_uuid()
screen_id    UUID NOT NULL REFERENCES screens(id)
seat_type_id UUID NOT NULL REFERENCES seat_types(id)
seat_label   VARCHAR(20) NOT NULL             -- e.g. "A1", "B12"
status       VARCHAR(20) DEFAULT 'AVAILABLE'  -- ENUM: 'AVAILABLE' | 'RESERVED' | 'MAINTENANCE'
created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
UNIQUE (screen_id, seat_label)               -- no duplicate labels per screen
```
Individual physical seats. A seat's availability for a specific show is managed via `seat_lock` and `booking_seats`.

#### `bookings` (redesigned)
```sql
id             UUID PRIMARY KEY DEFAULT gen_random_uuid()
show_id        UUID REFERENCES show(id)
customer_id    UUID REFERENCES users(id)
booking_status VARCHAR(20) DEFAULT 'RESERVED'  -- ENUM: 'RESERVED' | 'CONFIRMED' | 'CANCELLED' | 'EXPIRED'
payment_mode   VARCHAR(20)                     -- ENUM: 'CASH' | 'ONLINE' | NULL
qr_code_url    TEXT                            -- for show entry
total_amount   DECIMAL(10,2) NOT NULL
expires_at     TIMESTAMP                       -- booking held until this time
created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```
Key changes: UUIDs throughout, proper status lifecycle, payment mode captured, QR code support, expiry time for unpaid bookings.

#### `booking_seats`
```sql
id         UUID PRIMARY KEY DEFAULT gen_random_uuid()
booking_id UUID REFERENCES bookings(id)
seat_id    UUID REFERENCES seats(id)
show_id    UUID REFERENCES show(id)
price      DECIMAL(10,2) NOT NULL    -- price locked at time of booking
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```
Junction table linking a booking to its specific seats. Storing `price` here means the price is locked at booking time even if `seat_types.price` changes later.

#### `seat_lock`
```sql
id         UUID PRIMARY KEY DEFAULT gen_random_uuid()
show_id    UUID NOT NULL REFERENCES show(id)
seat_id    UUID NOT NULL REFERENCES seats(id)
locked_by  UUID NOT NULL REFERENCES users(id)
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
expires_at TIMESTAMP NOT NULL
UNIQUE (show_id, seat_id)             -- one lock per seat per show at a time
```
Temporary lock held while a user is in the checkout flow. Prevents two users from booking the same seat simultaneously. Locks expire automatically (via `expires_at`).

#### `payments`
```sql
id             INT GENERATED BY DEFAULT AS IDENTITY (PK)
booking_id     UUID NOT NULL UNIQUE REFERENCES bookings(id)
payment_mode   VARCHAR(20) NOT NULL           -- 'cash' | 'online'
gateway        VARCHAR(50)                    -- e.g. 'razorpay', 'stripe'
transaction_id VARCHAR(255)
amount         DECIMAL(10,2) NOT NULL
status         VARCHAR(20)                    -- 'PENDING' | 'SUCCESS' | 'FAILED'
created_at     TIMESTAMP
paid_at        TIMESTAMP
```
Payment record per booking. Decoupled from `bookings` to avoid the circular dependency noted in the SQL comment.

---

### Entity Relationship (New Schema)

```
theatres
  └─< screens
        └─< seats >──── seat_types
        └─< show
              └─< bookings >──── users
              │     └─< booking_seats >── seats
              │     └── payments
              └─< seat_lock >─── users
                                seats
```

---

## 8. Domain Models

> **Note:** The models below reflect the **old schema** (original code). They have not been updated to match the new database schema yet. Treat this section as historical context while you build out the new feature code.

All models live in **`bookings/model/`** and are mapped to DB tables via GORM tags.

### `User`
```go
type User struct {
    Id       int    // PK
    Username string
    Password string // plain text
}
// Table: usertable
```

### `Slot`
```go
type Slot struct {
    Id        int
    Name      string  // "Morning", "Afternoon", "Evening"
    StartTime string
    EndTime   string
}
// Table: slot
```

### `Show`
```go
type Show struct {
    Id      int
    MovieId string  // IMDb ID — looked up from Movie Service
    Date    string
    Slot    Slot    // preloaded via GORM
    SlotId  int
    Cost    float64 // per seat
}
// Table: show
```

### `Customer`
```go
type Customer struct {
    Id          int
    Oid         string // external owner id, defaults to "system"
    Name        string // 2–15 chars, binding:"required,min=2,max=15"
    PhoneNumber string // exactly 10 digits, binding:"phoneNumber"
}
// Table: customer
```

### `Booking`
```go
type Booking struct {
    Id         int
    Date       string
    Show       Show     // preloaded via GORM
    ShowId     int
    Customer   Customer // preloaded via GORM
    CustomerId int
    NoOfSeats  int
    AmountPaid float64
}
// Table: booking
```

### `Movie` (NOT in DB)
```go
type Movie struct {
    MovieId  string // IMDb ID
    Name     string
    Duration string // e.g. "2h22m"
    Plot     string
}
```
`Movie` is fetched from the external Movie Service at runtime — it is **not stored** in the Postgres database.

---

## 9. Repositories (Data Access Layer)

Repositories are the only layer that talks to the database. Each one embeds `*common.BaseDB`.

### `bookingRepository`
| Method | SQL Operation |
|---|---|
| `Create(ctx, *Booking) error` | `INSERT INTO booking` (conflict → do nothing) |
| `BookedSeatsByShow(ctx, showId) int` | `SELECT SUM(no_of_seats) FROM booking WHERE show_id=?` |
| `BookingAmountByShows(ctx, []showIds) float64` | `SELECT SUM(amount_paid) FROM booking WHERE show_id IN (?)` |

### `showRepository`
| Method | SQL Operation |
|---|---|
| `GetAllShowsOn(ctx, date) ([]Show, error)` | `SELECT * FROM show WHERE date=? PRELOAD Slot` |
| `FindById(ctx, id) (Show, error)` | `SELECT * FROM show WHERE id=? PRELOAD Slot` (404 if not found) |

### `userRepository`
| Method | SQL Operation |
|---|---|
| `FindByUsername(ctx, username) (User, error)` | `SELECT * FROM usertable WHERE username=?` |
| `Create(ctx, *User) error` | `INSERT INTO usertable` |

### `customerRepository`
| Method | SQL Operation |
|---|---|
| `Create(ctx, *Customer) error` | `INSERT INTO customer` (auto-increment Id if 0; conflict → do nothing) |

> **Important**: `CustomerRepository` is defined as an **interface** in the repository package. Other repositories expose their structs directly. This is slightly inconsistent but not a bug.

---

## 10. Services (Business Logic Layer)

Services hold all business rules. They receive interfaces (not concrete types) for repositories and other services, making them easy to unit-test with mocks.

### `bookingService`

**`Book(ctx, BookingRequest) (*Booking, error)`** — the most complex service method:

```
1. Find the Show by showId           → 404 if not found
2. Validate NoOfSeats > 0
3. Check NoOfSeats ≤ MAX_NO_OF_SEATS_PER_BOOKING (15)
4. Calculate amountPaid = show.Cost * noOfSeats
5. Calculate availableSeats = TOTAL_NO_OF_SEATS(100) - alreadyBookedSeats
6. Ensure availableSeats ≥ seatsRequested   → 400 if not enough seats
7. Save the Customer to DB
8. Create the Booking record
9. Return the created Booking
```

### `showService`

| Method | What it does |
|---|---|
| `GetShows(ctx, date)` | Delegates to `showRepository.GetAllShowsOn`. |
| `GetMovieById(ctx, movieId)` | Calls `movieGateway.MovieById` (HTTP call to Movie Service). |

### `userService`

| Method | What it does |
|---|---|
| `UserDetails(ctx, username)` | Looks up user by username. Used both for login and for auth middleware. |

### `revenueService`

| Method | What it does |
|---|---|
| `RevenueOn(ctx, date)` | Gets all shows on that date, extracts their IDs, sums up `amount_paid` of all bookings for those shows. |

---

## 11. Controllers (HTTP Handler Layer)

Controllers translate HTTP ↔ Service. They:
1. Parse/validate the request (via `ShouldBindJSON` or query params).
2. Call the appropriate service method.
3. Map the result to a response DTO.
4. Write the HTTP response (`c.IndentedJSON` or `c.AbortWithStatusJSON`).

### `BookingController`
- **`CreateBooking(c *gin.Context)`** — `POST /bookings`
  - Binds JSON body to `BookingRequest`.
  - Calls `bookingService.Book(...)`.
  - Returns `BookingConfirmationResponse` with status **201 Created**.

### `showController`
- **`Shows(c *gin.Context)`** — `GET /shows?date=YYYY-MM-DD`
  - Reads `date` from query string.
  - Calls `showService.GetShows(...)` → list of shows.
  - For each show, calls `showService.GetMovieById(...)` (hits Movie Service).
  - Returns list of `ShowResponse`.

### `UserController`
- **`Login(c *gin.Context)`** — `GET /login` (with Basic Auth header)
  - Extracts username from the Basic Auth header (password was already validated by middleware).
  - Calls `userService.UserDetails(...)`.
  - Returns the username as a JSON string with status **200 OK**.

### `RevenueController`
- **`GetRevenue(c *gin.Context)`** — `GET /revenue?date=YYYY-MM-DD`
  - Reads `date` from query string.
  - Calls `revenueService.RevenueOn(...)`.
  - Returns a single float64 (total revenue amount).

---

## 12. HTTP Routes & Endpoints

```
GET  /swagger/*any          → Swagger UI (public, no auth)
GET  /login                 → UserController.Login     (Basic Auth required)
GET  /shows?date=...        → showController.Shows     (Basic Auth required)
POST /bookings              → BookingController.CreateBooking  (Basic Auth required)
GET  /revenue?date=...      → RevenueController.GetRevenue    (Basic Auth required)
```

Route constants are in `bookings/constants/app_constant.go`:
```go
RevenueEndPoint = "/revenue"
BookingEndPoint = "/bookings"
ShowEndPoint    = "/shows"
LoginEndPoint   = "/login"
```

The router is split into two groups:
- **`authRouter`** — passes through the `security.Authenticate` middleware. All business endpoints live here.
- **`noAuthRouter`** — no middleware. Only Swagger UI lives here.

---

## 13. Middleware

### Basic Auth — `common/middleware/security/basic_auth.go`

Function: **`Authenticate(userService) gin.HandlerFunc`**

Flow:
```
1. Extract username + password from Authorization header (Basic scheme).
2. Return 400 if header is missing or credentials are empty.
3. Look up user by username via userService.UserDetails().
4. Return 401 if user not found or password does not match.
5. Call c.Next() to pass control to the route handler.
```

> **Note**: Passwords are compared as **plain text**. There is no hashing.

### Input Validation — `common/middleware/validator/`

**`DtoValidator`** is a custom Gin binding validator. It is lazily initialised once (`sync.Once`) and registers two custom rules:

| Tag | Rule |
|---|---|
| `phoneNumber` | Field must be exactly 10 decimal digits (`\d{10}`). |
| `maxSeats` | Field must be ≤ `MAX_NO_OF_SEATS_PER_BOOKING` (15). |

**`HandleStructValidationError(err)`** — turns `validator.ValidationErrors` into a human-readable `AppError` (400 Bad Request) with a message like `"field 'NoOfSeats', condition: ..."`.

### CORS — `common/middleware/cors/`

Standard CORS middleware registered on the Gin engine (exact origins/methods defined in that file).

---

## 14. DTOs — Request & Response Shapes

### Request

**`BookingRequest`** (`bookings/dto/request/booking_request.go`):
```json
{
  "date":      "2025-07-15",       // required, format YYYY-MM-DD
  "showId":    3,                  // required, >= 1
  "customer": {
    "name":        "Alice",        // required, 2–15 chars
    "phoneNumber": "9876543210"    // required, exactly 10 digits
  },
  "noOfSeats": 2                   // required, 1–15 (custom maxSeats validator)
}
```

### Responses

**`BookingConfirmationResponse`**:
```json
{
  "id":           42,
  "customerName": "Alice",
  "showDate":     "2025-07-15",
  "startTime":    "09:00:00",
  "amountPaid":   300.00,
  "noOfSeats":    2
}
```

**`ShowResponse`** (one element in the array returned by `GET /shows`):
```json
{
  "movie": {
    "id":       "tt0111161",
    "name":     "The Shawshank Redemption",
    "duration": "2h22m0s",
    "plot":     "..."
  },
  "slot": {
    "id":        1,
    "name":      "Morning",
    "startTime": "09:00:00",
    "endTime":   "12:00:00"
  },
  "id":   5,
  "date": "2025-07-15",
  "cost": 150.00
}
```

---

## 15. Error Handling

All errors returned by services and repositories are of type `*AppError` (`error/app_error.go`).

```go
type AppError struct {
    error               // wrapped underlying Go error
    httpCode int        // HTTP status code
    Code     string     // machine-readable code, e.g. "ShowNotFound"
    Message  string     // human-readable message
}
```

**Constructors** (`error/error_constructor.go`):

| Function | HTTP Status |
|---|---|
| `NotFoundError(code, msg, err)` | 404 |
| `BadRequestError(code, msg, err)` | 400 |
| `UnProcessableError(code, msg, err)` | 422 |
| `InternalServerError(code, msg, err)` | 500 |
| `InvalidCredentialsError(code, msg, err)` | 401 |

**How controllers use it:**
```go
if responseError != nil {
    err := responseError.(*ae.AppError)   // type-assert to AppError
    c.AbortWithStatusJSON(err.HTTPCode(), err)
    return
}
```

**Error JSON response shape:**
```json
{
  "Code":    "ShowNotFound",
  "Message": "Show not found for id : 99"
}
```

---

## 16. External Dependency — Movie Service Gateway

The backend **does not store movie metadata**. Instead it calls the Movie Service microservice for each show.

**`movieservice/movie_gateway/movie_gateway.go`**:
- `MovieById(ctx, id)` → `GET {movieServiceHost}movies/{id}`
- Parses the response as `MovieServiceResponse` (IMDb ID, title, runtime, plot).
- Converts runtime string (e.g. `"142 min"`) into a Go `time.Duration` string (`"2h22m0s"`).
- Returns `*model.Movie`.

If the Movie Service is **unavailable**, the entire `GET /shows` response will fail with a 500 error because every show needs its movie data.

> **Tip for local dev**: make sure `docker-compose-local.yml` brings up the movie service, or point `MOVIE_SERVICE_HOST` to a running instance.

---

## 17. Seeding the Database

`bookings/database/seed/dataseeder.go` → `SeedDB(userRepo)` runs **at every startup**.

> **Note:** The seeder currently targets the old `usertable`. With the new schema it will need to target the `admins` table instead. The seeder has not been updated yet.

Current behaviour in the old code — creates these users if they don't exist:

| Username | Password |
|---|---|
| `seed-user-1` | `foobar` |
| `seed-user-2` | `foobar` |

With the new schema, authentication will need to seed the `admins` table with hashed passwords.

---

## 18. Logging

`common/logger/logger.go` wraps **Uber Zap**.

- Initialised once in `server.Init` via `logger.InitAppLogger(cfg.Logger)`.
- `cfg.Logger.Level` controls verbosity (`debug`, `info`, `warn`, `error`).
- Usage anywhere in the codebase: `logger.Error(...)`, `logger.Info(...)`, etc.
- Gin request logs are also handled by `gin-contrib/zap` middleware registered in `setupApp`.

---

## 19. Testing Strategy

The codebase has three levels of tests:

### Unit Tests (`bookings/service/*_test.go`, `bookings/repository/*_test.go`)

- Services are tested using **hand-written mocks** from `_mocks/repomocks/`.
- The mocks implement the same interfaces (e.g. `BookingRepository`, `ShowRepository`) that services depend on.
- Fast — no database, no network.

### Repository Tests (`bookings/repository/*_test.go`)

- Use **`testcontainers-go`** to spin up a real Postgres container.
- `setup_test.go` in the repository package handles container lifecycle.
- Tests run real SQL against a real DB schema.

### Integration Tests (`integration_test/`)

- Full HTTP-level tests: spin up a real Postgres container, run all migrations, seed data, start the Gin router, then fire real HTTP requests using `appleboy/gofight`.
- `setup_test.go` in `integration_test/` manages all of this.
- Files: `booking_controller_test.go`, `shows_controller_test.go`, `user_controller_test.go`, `revenue_controller_test.go`.

### Running Tests

```bash
# All unit + repository tests
go test ./...

# Integration tests (requires Docker)
go test ./integration_test/...

# Or via Makefile
make test
make integration-test
```

---

## 20. Docker & Local Development

### `docker-compose-local.yml`

Brings up:
- **Postgres** database
- **Movie Service** (the Ruby/Go microservice)
- **Backend** (the Go API)

### Starting locally

```bash
# From backend/
./run-backend-local.sh
# or
docker compose -f docker-compose-local.yml up
```

### Environment Variables needed

```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=skyfox
POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=password
MOVIE_SERVICE_HOST=http://localhost:3000/
```

### Swagger UI

Once running: `http://localhost:8080/swagger/index.html`

---

## 21. Full Request Lifecycle — End to End

### Example: `POST /bookings`

```
Client
  │
  │  POST /bookings
  │  Authorization: Basic c2VlZC11c2VyLTE6Zm9vYmFy
  │  Body: { "date": "2025-07-15", "showId": 3, "customer": {...}, "noOfSeats": 2 }
  ▼
Gin Router
  │
  ├── security.Authenticate middleware
  │     ├── Parse Basic Auth header → username="seed-user-1", password="foobar"
  │     ├── userService.UserDetails("seed-user-1")
  │     │     └── userRepository.FindByUsername(ctx, "seed-user-1")
  │     │           └── SELECT * FROM usertable WHERE username='seed-user-1'
  │     ├── Compare password → match ✓
  │     └── c.Next()
  │
  └── BookingController.CreateBooking
        ├── c.ShouldBindJSON(&BookingRequest)
        │     ├── validate date format
        │     ├── validate showId >= 1
        │     ├── validate customer.name & phoneNumber
        │     └── validate noOfSeats (maxSeats: 1–15)
        │
        └── bookingService.Book(ctx, bookingRequest)
              ├── showRepository.FindById(ctx, 3)
              │     └── SELECT * FROM show WHERE id=3 PRELOAD Slot
              ├── check noOfSeats ≤ 15 ✓
              ├── amountPaid = show.Cost * 2 = 150.00 * 2 = 300.00
              ├── bookingRepository.BookedSeatsByShow(ctx, 3)
              │     └── SELECT SUM(no_of_seats) FROM booking WHERE show_id=3 → 20
              ├── availableSeats = 100 - 20 = 80 ≥ 2 ✓
              ├── customerRepository.Create(ctx, &customer)
              │     └── INSERT INTO customer (...) ON CONFLICT DO NOTHING
              ├── bookingRepository.Create(ctx, &newBooking)
              │     └── INSERT INTO booking (...)
              └── return &Booking{...}
        │
        └── response.NewBookingConfirmationResponse(...)
              └── c.IndentedJSON(201, BookingConfirmationResponse{...})
```

---

## 22. Business Rules & Constants

> **Note:** The constants and rules below belong to the **old code** and will need to be rethought for the new schema. The new DB has per-seat granularity, a status lifecycle, and seat locks — the old blunt `TOTAL_NO_OF_SEATS - SUM(booked)` approach no longer applies.

**Old constants** (`bookings/constants/app_constant.go` — still in source, not yet updated):

| Constant | Value | Meaning |
|---|---|---|
| `TOTAL_NO_OF_SEATS` | `100` | Total seats available per show (hardcoded — obsolete) |
| `MAX_NO_OF_SEATS_PER_BOOKING` | `15` | Max seats one booking can request |

**Old derived rules (for reference only):**
- `availableSeats = 100 - SUM(no_of_seats for that show)`
- `amountPaid = show.Cost * noOfSeats` (single flat price)
- A show was uniquely identified by `(date, slot_id)`.

**New rules implied by the new schema:**
- Seat availability per show is determined by checking `seats` not in `booking_seats` for that show, and not currently locked in `seat_lock` with a non-expired `expires_at`.
- Total amount = SUM of `seat_types.price` for each selected seat (prices can differ per seat type — Regular / Premium / Recliner etc.).
- Booking lifecycle: `RESERVED` → `CONFIRMED` (payment done) or `EXPIRED` (payment not made before `expires_at`) or `CANCELLED`.
- A seat is held temporarily via `seat_lock` during the checkout window to prevent race conditions.

---

## 23. Dependency Graph (Summary)

```
main.go
  └── server.Init()
        ├── config.LoadConfig()
        ├── connection.NewDBHandler()  ──────────────────┐
        │     └── *common.BaseDB                         │ (embedded in all repos)
        ├── movieGateway                                  │
        │                                                 │
        ├── Repositories ◄────────────────────────────────┘
        │     ├── bookingRepository
        │     ├── showRepository
        │     ├── userRepository
        │     └── customerRepository
        │
        ├── Services (depend on Repositories via interfaces)
        │     ├── bookingService  ← bookingRepository, showRepository, customerRepository
        │     ├── showService     ← showRepository, movieGateway
        │     ├── userService     ← userRepository
        │     └── revenueService  ← bookingRepository, showRepository
        │
        └── Controllers (depend on Services via interfaces)
              ├── BookingController  ← bookingService
              ├── showController     ← showService
              ├── UserController     ← userService
              └── RevenueController  ← revenueService
```

Every dependency is injected through **interfaces** — this is why mocking works cleanly in tests.

---

## 24. Where to Add New Features

### Immediate Priority — Bring Code Up to the New Schema

Before building new features, the existing code needs to be migrated to match the new database. The recommended order:

1. **Update models** — replace old structs (`Slot`, `Show`, `Customer`, `Booking`, `User`) with new ones (`Admin`, `User`, `Theatre`, `Screen`, `Show`, `SeatType`, `Seat`, `Booking`, `BookingSeat`, `SeatLock`, `Payment`).
2. **Update repositories** — rewrite all DB queries to target new table names and UUID-based IDs.
3. **Update services** — revise business logic (e.g. seat availability is now per-seat, not a simple count; bookings have a status lifecycle).
4. **Update controllers & DTOs** — new request/response shapes.
5. **Update routes** — new endpoints will likely be needed (e.g. seat lock, payment confirmation).
6. **Update seeder** — seed `admins` table instead of `usertable`.
7. **Update mocks** — regenerate with `mockery` after interfaces change.

---

### Step-by-Step Recipe for Any New Feature (after code is caught up)

Example: "cancel a booking"

### Step 1 — Migration (if new DB columns/tables are needed)
Create `migration/scripts/000022_...up.sql` and `.down.sql`.

### Step 2 — Model
Update or add a struct in `bookings/model/`.

### Step 3 — Repository
Add the new method to the appropriate repository struct. If it's a new query type, add it to the repository file (e.g. `booking_repository.go`). Add the method signature to any interface that declares it.

### Step 4 — Service
Add the business logic method to the appropriate service. Declare what repository method you need as a new entry in the service's local interface (e.g. `BookingRepository` interface inside `booking_service.go`).

### Step 5 — DTO
Add a new request struct in `bookings/dto/request/` and/or a response struct in `bookings/dto/response/`.

### Step 6 — Controller
Add a new handler method to the appropriate controller struct, or create a new controller file.

### Step 7 — Route
Register the new route in `app/server/start.go` inside `Init()`.

### Step 8 — Mocks
Update or add a mock in `_mocks/repomocks/` to keep unit tests working.

### Step 9 — Tests
- Add a unit test in `bookings/service/`.
- Add an integration test in `integration_test/`.

### Step 10 — Swagger
Run `swag init` from the `backend/` directory to regenerate `docs/` after adding Swaggo annotations to the handler.

```bash
swag init -g main.go
```

---

*Generated from source code analysis of the `backend/` directory.*
