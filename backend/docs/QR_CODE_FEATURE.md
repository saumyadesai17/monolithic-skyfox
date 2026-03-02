# QR Code Generation Feature

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Database Schema](#database-schema)
4. [File Structure](#file-structure)
5. [Layer-by-Layer Walkthrough](#layer-by-layer-walkthrough)
   - [Request DTO](#1-request-dto)
   - [Response DTO](#2-response-dto)
   - [Models](#3-models)
   - [Repository](#4-repository)
   - [Service — QR Generation](#5-service--qr-generation)
   - [Controller](#6-controller)
   - [Router Wiring](#7-router-wiring)
   - [Configuration](#8-configuration)
6. [QR Payload — Fields & Signing](#qr-payload--fields--signing)
7. [Offline Verification Without a DB Call](#offline-verification-without-a-db-call)
8. [Caching Strategy](#caching-strategy)
9. [Error Handling](#error-handling)
10. [Testing](#testing)
11. [API Reference](#api-reference)
12. [Security Considerations](#security-considerations)
13. [Dependency](#dependency)

---

## Overview

The QR Code feature generates a unique, tamper-evident QR code for each confirmed booking. The QR code:

- Encodes a **self-contained signed payload** with everything a scanner needs to verify and display booking details — booking ID, show ID, customer ID, movie, theatre, show start/end time, and reserved seat IDs — without making a database call.
- Is signed with **HMAC-SHA256** using a server-side secret so any tampered payload is immediately detectable.
- Is rendered as a **256×256 PNG** and returned as a **base64 data URL** that can be embedded in an `<img>` tag or email directly.
- Is generated **on demand** via a GET endpoint and **cached** in the `bookings.qr_code_url` column — the first call generates and stores it; every subsequent call returns the stored copy instantly.

---

## Architecture

```
Client
  │
  │  GET /api/v1/bookings/:bookingId/qr
  ▼
┌──────────────────────────┐
│      QRController        │  ← ShouldBindUri validates UUID path param
└──────────┬───────────────┘
           │ GenerateQR(ctx, bookingID)
           ▼
┌──────────────────────────┐
│       QRService          │  ← cache check → enrich → sign → encode PNG → store
└──────────┬───────────────┘
           │
    ┌──────┴────────────────────────────────────┐
    │           QRBookingRepository             │
    │  FindBookingByID   ──► bookings table     │
    │  FindShowByID      ──► show table         │
    │  FindSeatsByBookingID ► booking_seats     │
    │  UpdateQRCodeURL   ──► bookings table     │
    └──────────────────────────────────────────┘
           │
           ▼
      PostgreSQL
```

Every layer talks **only to its direct neighbour via an interface**, keeping each layer independently testable and replaceable. The service makes **3 sequential DB reads** (booking → show → seats) on first generation, then **1 write** to cache the result.

---

## Database Schema

The feature reads from three tables across three migrations:

### `bookings` table — migration 018

| Column          | Type          | Notes                                              |
|-----------------|---------------|----------------------------------------------------|
| `id`            | UUID (PK)     | Booking identifier                                 |
| `show_id`       | UUID (FK)     | Links to the `show` table                          |
| `customer_id`   | UUID (FK)     | Links to the `users` table                         |
| `booking_status`| VARCHAR(20)   | RESERVED / CONFIRMED / CANCELLED / EXPIRED         |
| `total_amount`  | DECIMAL(10,2) |                                                    |
| `qr_code_url`   | TEXT          | **NULL** until first QR request; then a data URL   |
| `expires_at`    | TIMESTAMP     |                                                    |

### `show` table — migration 015 (read-only for QR)

| Column             | Type        | Notes                           |
|--------------------|-------------|---------------------------------|
| `id`               | UUID (PK)   | Show identifier                 |
| `movie_imdb_id`    | VARCHAR(20) | IMDB identifier encoded in QR   |
| `theatre_id`       | UUID (FK)   | Theatre UUID encoded in QR      |
| `screen_id`        | UUID (FK)   |                                 |
| `start_time`       | TIMESTAMP   | Show start — encoded in QR      |
| `end_time`         | TIMESTAMP   | Show end — encoded in QR        |

### `booking_seats` table — migration 019 (read-only for QR)

| Column       | Type      | Notes                           |
|--------------|-----------|---------------------------------|
| `id`         | UUID (PK) |                                 |
| `booking_id` | UUID (FK) | Filter key to find seats        |
| `seat_id`    | UUID (FK) | Seat UUID encoded in QR payload |
| `show_id`    | UUID (FK) |                                 |
| `price`      | DECIMAL   |                                 |

A `NULL` / empty `qr_code_url` means the QR has not been generated yet. A non-empty value is the cached result returned directly on subsequent requests.

---

## File Structure

```
backend/
├── bookings/
│   ├── controller/
│   │   ├── qr_controller.go          # HTTP handler
│   │   └── qr_controller_test.go     # Controller unit tests
│   ├── dto/
│   │   ├── request/
│   │   │   └── qr_request.go         # Path-param DTO (UUID validation)
│   │   └── response/
│   │       └── qr_response.go        # Response DTO
│   ├── model/
│   │   └── booking_new.go            # BookingRecord, ShowRecord, BookingSeat GORM models
│   ├── repository/
│   │   └── qr_repository.go          # DB access (FindBookingByID, FindShowByID, FindSeatsByBookingID, UpdateQRCodeURL)
│   └── service/
│       ├── qr_service.go             # Core business logic + HMAC signing
│       ├── qrService_test.go         # Service unit tests
│       └── mocks/
│           ├── mock_QRBookingRepository.go  # Auto-generated by mockery
│           └── mock_QRService.go            # Auto-generated by mockery
├── config/
│   ├── config.go                     # QRConfig struct
│   ├── config.yml                    # QR_SECRET env var mapping
│   └── config-local.yml              # Local dev secret
└── app/
    └── server/
        └── start.go                  # Dependency wiring + route registration
```

---

## Layer-by-Layer Walkthrough

### 1. Request DTO

**File:** `bookings/dto/request/qr_request.go`

```go
type QRCodeRequest struct {
    BookingID string `uri:"bookingId" binding:"required,uuid"`
}
```

- Uses Gin's `uri` tag to bind the `:bookingId` path parameter.
- The `uuid` binding rule rejects anything that is not a valid UUID v4 **before** the service is called.
- This means a malformed ID never reaches the database.

---

### 2. Response DTO

**File:** `bookings/dto/response/qr_response.go`

```go
type QRCodeResponse struct {
    QRCodeURL string `json:"qrCodeUrl" example:"data:image/png;base64,iVBORw0KGgo="`
}

func NewQRCodeResponse(qrCodeURL string) *QRCodeResponse {
    return &QRCodeResponse{QRCodeURL: qrCodeURL}
}
```

- Wraps the data URL in a typed struct rather than a raw `gin.H` map.
- The `example` tag is picked up by swaggo for Swagger UI documentation.

---

### 3. Models

**File:** `bookings/model/booking_new.go`

Three models are defined in this file. They are intentionally **separate from the legacy structs** so the QR feature can be developed and tested independently.

#### `BookingRecord` — `bookings` table (migration 018)

```go
type BookingRecord struct {
    Id            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    ShowId        string    `gorm:"type:uuid;not null"`
    CustomerId    string    `gorm:"type:uuid;not null"`
    BookingStatus string    `gorm:"default:RESERVED"`
    PaymentMode   string
    QRCodeURL     string    `gorm:"column:qr_code_url"`
    TotalAmount   float64   `gorm:"not null"`
    ExpiresAt     time.Time
    CreatedAt     time.Time
}
func (BookingRecord) TableName() string { return "bookings" }
```

#### `ShowRecord` — `show` table (migration 015)

```go
type ShowRecord struct {
    Id          string    `gorm:"primaryKey;type:uuid"`
    MovieImdbId string    `gorm:"column:movie_imdb_id"`
    ScreenId    string    `gorm:"column:screen_id;type:uuid"`
    TheatreId   string    `gorm:"column:theatre_id;type:uuid"`
    StartTime   time.Time `gorm:"column:start_time"`
    EndTime     time.Time `gorm:"column:end_time"`
    Status      string
}
func (ShowRecord) TableName() string { return "show" }
```

Used **read-only** — QR service fetches it to embed `movie_imdb_id`, `theatre_id`, `start_time`, and `end_time` into the payload.

#### `BookingSeat` — `booking_seats` table (migration 019)

```go
type BookingSeat struct {
    Id        string  `gorm:"primaryKey;type:uuid"`
    BookingId string  `gorm:"column:booking_id;type:uuid"`
    SeatId    string  `gorm:"column:seat_id;type:uuid"`
    ShowId    string  `gorm:"column:show_id;type:uuid"`
    Price     float64 `gorm:"column:price"`
}
func (BookingSeat) TableName() string { return "booking_seats" }
```

Used **read-only** — the service queries this table to collect all `seat_id` values for a booking and embeds them as a `[]string` in the QR payload.

---

### 4. Repository

**File:** `bookings/repository/qr_repository.go`

```go
type QRBookingRepository interface {
    FindBookingByID(ctx context.Context, id string) (*model.BookingRecord, error)
    FindShowByID(ctx context.Context, id string) (*model.ShowRecord, error)
    FindSeatsByBookingID(ctx context.Context, bookingID string) ([]string, error)
    UpdateQRCodeURL(ctx context.Context, id string, qrURL string) error
}
```

The interface is **deliberately narrow** — it exposes only the four operations the QR feature needs, keeping it independently testable without coupling to the broader `BookingRepository`.

#### `FindBookingByID`

```go
result := db.Where("id = ?", id).First(&booking)
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    return nil, nil   // caller treats nil → 404
}
```

- Returns `(nil, nil)` when no row exists. The service converts this to a `404 BookingNotFound` error.
- Returns `(nil, AppError{500})` for any genuine DB error.

#### `FindShowByID`

```go
result := db.Where("id = ?", id).First(&show)
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    return nil, nil   // caller treats nil → 404
}
```

- Fetches `movie_imdb_id`, `theatre_id`, `start_time`, `end_time` from the `show` table.
- Same `nil, nil` contract as `FindBookingByID`.

#### `FindSeatsByBookingID`

```go
db.Where("booking_id = ?", bookingID).Find(&seats)
// maps []BookingSeat → []string of seat_id values
```

- Returns a `[]string` of seat UUIDs for the booking.
- Returns an empty slice (not an error) if no seats are found — seats are optional for the QR to still be valid.

#### `UpdateQRCodeURL`

```go
db.Model(&model.BookingRecord{}).Where("id = ?", id).Update("qr_code_url", qrURL)
```

- Executes a targeted `UPDATE bookings SET qr_code_url = ? WHERE id = ?`.
- Does not touch any other column.

---

### 5. Service — QR Generation

**File:** `bookings/service/qr_service.go`

This is the heart of the feature. `GenerateQR` follows this exact flow:

```
1. FindBookingByID(bookingID)
        ├── DB error?            → return "", 500
        ├── nil?                 → return "", 404 BookingNotFound
        └── QRCodeURL != ""?     → return cached URL, nil  ← CACHE HIT (no further DB calls)
        │
2. FindShowByID(booking.ShowId)
        ├── DB error?            → return "", 500
        └── nil?                 → return "", 404 ShowNotFound
        │
3. FindSeatsByBookingID(bookingID)
        └── DB error?            → return "", 500
        │
4. buildSignedPayload(booking, show, seatIDs)
        │  assemble qrPayload{bid, sid, cid, mid, tid, sst, set, seats, iat}
        │  json.Marshal → base64.StdEncode
        │  HMAC-SHA256(secret, base64) → hex
        │  result: "<base64>.<hexsig>"
        │
5. qrcode.Encode(signedString, Medium, 256) → []byte PNG
        │
6. "data:image/png;base64," + base64(pngBytes) → dataURL
        │
7. UpdateQRCodeURL(bookingID, dataURL)
        │
8. return dataURL, nil
```

#### The `qrPayload` struct

```go
type qrPayload struct {
    BookingID  string   `json:"bid"`   // booking UUID
    ShowID     string   `json:"sid"`   // show UUID
    CustomerID string   `json:"cid"`   // customer UUID
    MovieID    string   `json:"mid"`   // movie IMDB ID (e.g. tt1234567)
    TheatreID  string   `json:"tid"`   // theatre UUID
    ShowStart  string   `json:"sst"`   // show start time RFC3339
    ShowEnd    string   `json:"set"`   // show end time RFC3339
    Seats      []string `json:"seats"` // seat UUIDs reserved for this booking
    IssuedAt   int64    `json:"iat"`   // Unix timestamp of QR generation
}
```

Fields are deliberately abbreviated to produce a shorter JSON string, reducing QR code visual density and making it easier to scan reliably. Together they make the payload **fully self-contained** — a scanner can display booking details and enforce a time window without a database lookup.

---

### 6. Controller

**File:** `bookings/controller/qr_controller.go`

```go
func (qc *QRController) GetQRCode(c *gin.Context) {
    var req request.QRCodeRequest
    if err := c.ShouldBindUri(&req); err != nil {
        c.JSON(400, ae.BadRequestError("InvalidBookingId", "bookingId must be a valid UUID", err))
        return
    }

    dataURL, err := qc.qrService.GenerateQR(c.Request.Context(), req.BookingID)
    if err != nil {
        appErr := err.(*ae.AppError)
        c.JSON(appErr.HTTPCode(), appErr)
        return
    }

    c.JSON(200, response.NewQRCodeResponse(dataURL))
}
```

- `ShouldBindUri` runs the `uuid` validator **before** anything else. An invalid ID is rejected immediately with `400`.
- The controller holds a local `QRService` interface (not the concrete service type), so it is fully mockable in tests without importing the service package.
- `*ae.AppError` preserves the exact HTTP status the service decided (404 or 500) — the controller never hardcodes error codes.

---

### 7. Router Wiring

**File:** `app/server/start.go`

```go
// Repositories
qrBookingRepository := repository.NewQRBookingRepository(db)

// Services
qrService := service.NewQRService(qrBookingRepository, cfg.QR.Secret)

// Controllers
qrController := controller.NewQRController(qrService)

// Routes — on the noAuthRouter (no authentication required)
v1 := noAuthRouter.Group("/api/v1")
{
    bookingsV1 := v1.Group("/bookings")
    {
        bookingsV1.GET("/:bookingId/qr", qrController.GetQRCode)
    }
}
```

The route is on `noAuthRouter` so clients can fetch a booking QR without needing an auth token (e.g., for display at a ticket kiosk or email link).

---

### 8. Configuration

**File:** `config/config.go`

```go
type QRConfig struct {
    Secret string `yaml:"secret"`
}

type AppConfig struct {
    // ...existing fields...
    QR QRConfig
}
```

**`config/config.yml`** (production):
```yaml
QR:
  secret: ${QR_SECRET}
```

**`config/config-local.yml`** (local dev):
```yaml
QR:
  secret: local-dev-qr-secret-change-in-prod
```

- In production the secret is injected via the `QR_SECRET` environment variable.
- The secret must be treated like a private key — rotating it invalidates all previously issued QR codes.

---

## QR Payload — Fields & Signing

### Wire format

The QR image does **not** contain a URL. It encodes a signed string:

```
<base64_payload>.<hex_signature>
```

### Payload fields

| Short key | Full meaning      | Value example                          |
|-----------|-------------------|----------------------------------------|
| `bid`     | Booking UUID      | `550e8400-e29b-41d4-a716-446655440000` |
| `sid`     | Show UUID         | `a3f8...`                              |
| `cid`     | Customer UUID     | `77b1...`                              |
| `mid`     | Movie IMDB ID     | `tt9362722`                            |
| `tid`     | Theatre UUID      | `c9d1...`                              |
| `sst`     | Show start (RFC3339) | `2026-03-10T18:00:00Z`              |
| `set`     | Show end (RFC3339)   | `2026-03-10T20:30:00Z`              |
| `seats`   | Seat UUIDs        | `["f1a2...", "b3c4..."]`               |
| `iat`     | Issued-at (Unix)  | `1741201200`                           |

### Step-by-step construction

```
Step 1 — Assemble qrPayload struct and marshal to JSON
    {
      "bid": "550e8400-...",
      "sid": "a3f8...",
      "cid": "77b1...",
      "mid": "tt9362722",
      "tid": "c9d1...",
      "sst": "2026-03-10T18:00:00Z",
      "set": "2026-03-10T20:30:00Z",
      "seats": ["f1a2...", "b3c4..."],
      "iat": 1741201200
    }

Step 2 — Base64-encode the JSON (StdEncoding)
    eyJiaWQiOiI1NTBlODQwMC0uLi4iLCJtaWQiOiJ0dDkzNjI3MjIifQ==

Step 3 — HMAC-SHA256 sign the base64 string
    mac = HMAC-SHA256(QR_SECRET, base64_string)
    sig = hex.EncodeToString(mac)   →  "3a9f...c21b"

Step 4 — Concatenate with a dot separator
    "eyJia...fQ==.3a9f...c21b"

Step 5 — Render as QR PNG
    qrcode.Encode(signedString, qrcode.Medium, 256)  →  []byte (PNG)

Step 6 — Wrap as browser-ready data URL
    "data:image/png;base64,iVBORw0KGgo..."
```

## Offline Verification Without a DB Call

A scanner (mobile app, kiosk, or verification endpoint) can fully verify and display the QR content using **only the server secret** — no database query required:

```
1. Scan QR → receive string  "<base64>.<hexsig>"

2. Split on "."
       base64_payload = parts[0]
       received_sig   = parts[1]

3. Recompute HMAC
       expected_sig = hex( HMAC-SHA256(QR_SECRET, base64_payload) )

4. Constant-time compare
       hmac.Equal(expected_sig_bytes, received_sig_bytes)
       → false  ⇒ REJECT — payload was tampered with or wrong secret

5. Decode and parse
       json = base64.StdDecode(base64_payload)
       payload = json.Unmarshal(json)

6. Enforce time window (optional)
       if show.EndTime (payload.set) is in the past → REJECT as expired

7. Display to venue staff
       Booking: payload.bid
       Movie:   payload.mid  (look up title from IMDB if needed)
       Theatre: payload.tid
       Show:    payload.sst → payload.set
       Seats:   payload.seats
       Customer: payload.cid
```

If the HMAC comparison fails the QR must be rejected — it was either tampered with or produced by a different signing key.

---

## Caching Strategy

The feature uses **write-through database caching**: the generated data URL is stored in `bookings.qr_code_url` on first generation and returned as-is on subsequent requests.

| Request | `qr_code_url` in DB | Action |
|---------|-------------------|--------|
| First   | NULL / empty      | Generate → store → return |
| Second+ | non-empty         | Return stored value immediately |

**Trade-offs:**
- No in-memory or Redis cache needed — the database acts as the persistent store.
- A rotate of `QR_SECRET` will not automatically regenerate stored QRs. If the secret is rotated, stored `qr_code_url` values should be cleared (`UPDATE bookings SET qr_code_url = NULL`) so new signed codes are generated on next request.

---

## Error Handling

| Scenario | HTTP Status | Error Code |
|---|---|---|
| `bookingId` path param is not a valid UUID | `400 Bad Request` | `InvalidBookingId` |
| Booking does not exist in the DB | `404 Not Found` | `BookingNotFound` |
| GORM error on `FindBookingByID` | `500 Internal Server Error` | `InternalServerError` |
| Show linked to booking does not exist | `404 Not Found` | `ShowNotFound` |
| GORM error on `FindShowByID` | `500 Internal Server Error` | `InternalServerError` |
| GORM error on `FindSeatsByBookingID` | `500 Internal Server Error` | `InternalServerError` |
| JSON marshal / HMAC failure in payload builder | `500 Internal Server Error` | `InternalServerError` |
| `qrcode.Encode` PNG generation failure | `500 Internal Server Error` | `InternalServerError` |
| GORM error on `UpdateQRCodeURL` | `500 Internal Server Error` | `InternalServerError` |

All errors are returned as `AppError` JSON:

```json
{
  "Code": "BookingNotFound",
  "Message": "booking not found",
  "httpCode": 404
}
```

---

## Testing

### Service Tests — `bookings/service/qrService_test.go`

Tests target `qrService` directly via a `MockQRBookingRepository` (generated by mockery). Test fixtures:
- `validBooking` — a `BookingRecord` with no `QRCodeURL` set
- `validShow` — a `ShowRecord` with `StartTime` / `EndTime` / `MovieImdbId` / `TheatreId`
- `validSeatIDs` — `[]string{"seat-uuid-0001", "seat-uuid-0002"}`

| Test Case | Mocked calls | What is asserted |
|---|---|---|
| Should generate and store QR when booking exists with show and seat data | `FindBookingByID` × 1, `FindShowByID` × 1, `FindSeatsByBookingID` × 1, `UpdateQRCodeURL` × 1 | Returned string starts with `data:image/png;base64,` |
| Should return cached QR without any further DB calls when qr_code_url is already set | `FindBookingByID` × 1 only | No other repo methods called; returned URL equals cached value |
| Should return 404 when booking is not found | `FindBookingByID` → `(nil, nil)` | `*AppError` HTTP 404, code `BookingNotFound` |
| Should return 500 when FindBookingByID returns a database error | `FindBookingByID` → error | `*AppError` HTTP 500 |
| Should return 404 when show is not found for the booking | `FindShowByID` → `(nil, nil)` | `*AppError` HTTP 404, code `ShowNotFound` |
| Should return 500 when FindShowByID returns a database error | `FindShowByID` → error | `*AppError` HTTP 500 |
| Should return 500 when FindSeatsByBookingID returns a database error | `FindSeatsByBookingID` → error | `*AppError` HTTP 500 |
| Should return 500 when UpdateQRCodeURL returns a database error | `UpdateQRCodeURL` → error | `*AppError` HTTP 500 |

Run:
```bash
go test ./bookings/service/... -v -run TestQRService
```

### Controller Tests — `bookings/controller/qr_controller_test.go`

Tests use `httptest.NewRecorder` + a Gin engine wired with `MockQRService`. No real HTTP server or database is required.

| Test Case | What is asserted |
|---|---|
| Should return 400 when bookingId is not a valid UUID | HTTP 400; service method **never called** |
| Should return 200 with qrCodeUrl when booking exists | HTTP 200; response body deserialises to `QRCodeResponse{QRCodeURL: testQRDataURL}` |
| Should return 404 when booking is not found | HTTP 404; response body contains a non-empty `Code` field |
| Should return 500 when QR service returns an internal error | HTTP 500; response body contains a non-empty `Code` field |

Run:
```bash
go test ./bookings/controller/... -v -run TestQRController
```

### Running all QR tests together

```bash
go test ./bookings/... -v -run "TestQR"
```

### Mock Generation

Mocks are generated by [mockery v3](https://github.com/vektra/mockery) using `.mockery.yml` at the repo root:

```bash
cd backend
mockery
```

Generated files (do not edit manually):
- `bookings/service/mocks/mock_QRBookingRepository.go`
- `bookings/service/mocks/mock_QRService.go`

---

## API Reference

### `GET /api/v1/bookings/:bookingId/qr`

Returns the QR code image for a booking as a base64-encoded PNG data URL.

#### Path Parameters

| Parameter  | Type   | Required | Validation | Description    |
|------------|--------|----------|------------|----------------|
| bookingId  | string | Yes      | UUID v4    | The booking ID |

#### Success Response — `200 OK`

```json
{
  "qrCodeUrl": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEA..."
}
```

The `qrCodeUrl` value is a complete [data URL](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URLs) and can be set directly as the `src` of an `<img>` tag in any web or email client.

#### Error Responses

**`400 Bad Request`** — bookingId is not a valid UUID:
```json
{ "Code": "InvalidBookingId", "Message": "bookingId must be a valid UUID" }
```

**`404 Not Found`** — no booking with that ID:
```json
{ "Code": "BookingNotFound", "Message": "booking not found" }
```

**`500 Internal Server Error`**:
```json
{ "Code": "InternalServerError", "Message": "failed to generate QR code" }
```

#### Authentication

This endpoint is on the **unauthenticated router** — no `Authorization` header is required. This enables use cases like:
- Embedding a QR link in a booking confirmation email.
- Displaying QR codes on a self-service kiosk without requiring user login.

---

## Security Considerations

| Concern | Mitigation |
|---|---|
| QR forgery | HMAC-SHA256 signature — any single-byte change to the payload produces a completely different signature, making tampering immediately detectable |
| Secret exposure | Secret is never hardcoded; injected at runtime via `QR_SECRET` environment variable |
| Secret rotation | Changing `QR_SECRET` invalidates all existing QRs; run `UPDATE bookings SET qr_code_url = NULL` to force regeneration on next request |
| Expired show admission | `sst`/`set` (show start/end) are embedded in the payload — scanners can enforce a time window without a DB call by comparing `set` against `time.Now()` |
| Replay attacks | `iat` (issued-at Unix timestamp) allows scanners to reject codes that are older than a configured maximum age |
| UUID injection | Gin's `binding:"uuid"` validates the path parameter before any DB query is executed |
| Seat/customer fraud | All identifying UUIDs (`bid`, `cid`, `tid`, `seats`) are included in the signed payload — forging any of them invalidates the HMAC |
| Booking enumeration | The endpoint returns `404 BookingNotFound` regardless of whether the booking doesn't exist or belongs to a different customer — if per-user auth is added later, move the endpoint to `authRouter` |

---

## Dependency

| Package | Version | Purpose |
|---|---|---|
| `github.com/skip2/go-qrcode` | `v0.0.0-20200617195104-da1b6568686e` | PNG QR code rendering |

Standard library packages used: `crypto/hmac`, `crypto/sha256`, `encoding/base64`, `encoding/hex`, `encoding/json`.
