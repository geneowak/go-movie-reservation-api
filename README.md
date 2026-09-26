# Cinehold

**A Go API for holding and booking cinema seats.**

Cinehold is a REST backend for a movie ticket reservation system: users browse what's
playing, hold a seat while they decide, and convert that hold into a booking. Admins
manage the catalogue — movies, locations, cinemas, seat maps and showtimes.

Go 1.27 · PostgreSQL 18+ · stdlib `net/http` · sqlc · goose · no web framework, no ORM

---

## Description

Cinehold implements the [Movie Reservation System](https://roadmap.sh/projects/movie-reservation-system)
spec end to end. The interesting part isn't the CRUD — it's what happens when two people
want the same seat at the same time.

- **Auth** — short-lived HS256 access tokens (1 hour) paired with long-lived opaque
  refresh tokens (7 days) that live in the database, so they can be revoked server-side.
  Passwords are hashed with Argon2id.
- **Catalogue management** — admins create movies, locations, cinemas and showtimes.
  A cinema's seating layout is stored as JSONB, so an auditorium can have an irregular
  shape without a schema change.
- **Hold, then book** — `POST /api/seats/reserve` puts a 10-minute hold on a seat while
  the user is still deciding. Nothing is sold until `POST /api/seats/book` confirms it.
- **Double-booking protection** — a `UNIQUE (show_time_id, seat_no)` constraint makes it
  structurally impossible for one seat to be held or booked twice, no matter how many API
  instances are running.
- **Abandoned hold recovery** — holds expire after 10 minutes and are reclaimed lazily on
  the next reservation attempt for that seat, so no background sweeper is required.

Everything is generated or type-checked: SQL is written by hand in `sql/queries` and turned
into typed Go by sqlc, and the schema is migrated by goose. There is no ORM and no
reflection-based query builder anywhere in the project.

---

## Motivation

**The concurrency problem is the heart of this project.** Every ticketing system
eventually hits the same wall: you cannot hold a database row lock for the ten minutes a
user spends deciding, and you cannot let two requests both believe they own seat `A:4`.
Cinehold resolves that by pushing the invariant down into the schema — a unique
constraint, a status column, a timestamp — instead of trusting every code path to remember
to check. The application-layer checks are there to return friendlier errors; the database
is what makes correctness guaranteed.

**The rest is deliberate engineering rather than convenience.** No web framework, no ORM,
no dependency injection container. Routing is `net/http.ServeMux` with Go 1.22+
method-and-pattern matching. Handlers depend on a `database.Querier` *interface* supplied
through a single `ApiConfig` struct, which is what makes them straightforward to test. SQL
is compiled by sqlc rather than assembled at runtime, so a malformed query fails at build
time instead of in production. Each of those decisions costs more on day one and pays for
itself every time the code is read or changed afterwards.

**The API surface is complete and coherent.** What remains is hardening — transactions,
pagination, observability — and that list is written down plainly in
[Retrospective & Future Work](#retrospective--future-work) rather than left implicit.

---

## Quick Start

### Prerequisites

- **Go 1.27+**
- **PostgreSQL 18+** — every primary key is generated with the built-in `uuidv7()`
  function, which landed in PostgreSQL 18
- **`sqlc`** — only needed if you change anything under `sql/`
- **`goose`** — only needed to run migrations

### 1. Get the code and configure it

```bash
git clone https://github.com/geneowak/go-movie-reservation-api.git
cd go-movie-reservation-api

cp .env.example .env
```

Then fill in `.env`:

```ini
DB_URL="postgres://postgres:postgres@localhost:5432/go_movie_reservation_api?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="<paste a long random string here>"
PORT="8080"
```

Generate a secret with:

```bash
openssl rand -base64 64
```

| Variable | Required | Notes |
| --- | --- | --- |
| `DB_URL` | yes | pgx connection string. The server exits immediately if it's missing. |
| `PLATFORM` | yes | Required by the config, though it isn't read anywhere yet. |
| `JWT_SECRET` | yes | HS256 signing key for access tokens. |
| `PORT` | no | Defaults to `8080`. |
| `TEST_DB_URL` | tests only | See [Contributing](#contributing). |

### 2. Create the database

```bash
createdb go_movie_reservation_api
```

### 3. Run the migrations

Migrations are plain [goose](https://github.com/pressly/goose) files in `sql/schema` and
are also embedded into the binary, so the goose CLI is enough:

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### 4. Run it

```bash
go run .
# Listening on port: 8080
```

### 5. Make your first request

```bash
# create an account
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"secret123"}'

# log in — response contains the user, an access token and a refresh token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"secret123"}'
```

The first account you need to be an admin is whichever one an existing admin promotes.
Set `is_admin = true` for your user directly in the database to bootstrap.

---

## Usage

All routes are registered in a single place, [`handlers/handlers.go`](handlers/handlers.go).
Every route is namespaced under `/api` and the four auth routes under `/api/auth`.

### Public

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/api/auth/signup` | Create an account. Rejects emails that are already registered. |
| `POST` | `/api/auth/login` | Exchange email + password for an access token and a refresh token. |
| `POST` | `/api/auth/refresh` | Exchange a refresh token (Bearer) for a fresh access token. |
| `POST` | `/api/auth/revoke` | Revoke the refresh token supplied in the Bearer header. `204`. |

### Admin

Requires an access token belonging to a user with `is_admin = true`.

| Method | Path | Description |
| --- | --- | --- |
| `PUT` | `/api/users/{userId}` | Grant or revoke admin (`{"is_admin": true}`). An admin cannot change their own status, so the system always keeps at least one. |
| `POST` | `/api/movies` | Create a movie. |
| `POST` | `/api/movies/{movieId}/show-times` | Schedule a movie at a cinema over a date range. |
| `GET` | `/api/admin/movies` | Every movie with its showtimes, for management. |
| `PUT` | `/api/movies/{movieId}` | Update movie details. Rejected while the movie has ongoing showtimes. |
| `DELETE` | `/api/movies/{movieId}` | Delete a movie. Rejected while the movie has ongoing showtimes. |
| `GET` | `/api/show-times/{showTimeId}` | Showtime detail with capacity, bookings and revenue. |
| `POST` | `/api/locations` | Create a location. |
| `GET` | `/api/locations` | List all locations. |
| `GET` | `/api/locations/{locationId}` | One location with its cinemas. |
| `POST` | `/api/cinemas` | Create a cinema with its seat map. |
| `PUT` | `/api/cinemas/{cinemaId}` | Update a cinema, including its seat map. Rejected while the cinema has ongoing showtimes, since that would invalidate seats people already paid for. |

### Authenticated user

Requires any valid access token.

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/api/movies` | Movies currently showing, with their showtimes. Optional filters: `?genre=`, `?time=15:04`, `?date=2006-01-02`. |
| `GET` | `/api/movies/{movieId}` | One movie with its showtimes. |
| `POST` | `/api/seats/reserve` | Hold a seat for 10 minutes. `201` with the reservation. |
| `POST` | `/api/seats/book` | Confirm previously held seats. |
| `GET` | `/api/bookings` | The current user's confirmed bookings. |
| `DELETE` | `/api/bookings/{reservationId}` | Cancel one of the current user's bookings. `204`. |

### Authentication

Login returns both tokens. Both are sent in the `Authorization: Bearer <token>` header —
the access token on ordinary requests, the refresh token on `/api/auth/refresh` and
`/api/auth/revoke`.

```jsonc
{
  "id": "0192...",
  "email": "you@example.com",
  "is_admin": false,
  "token": "eyJhbGciOiJIUzI1NiIs...",        // HS256 JWT, expires in 1 hour
  "refresh_token": "3f9a1c..."               // opaque, stored server-side, expires in 7 days
}
```

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Authorization: Bearer $REFRESH_TOKEN"
# {"token":"eyJhbGciOiJIUzI1NiIs..."}
```

When an access token expires, trade the refresh token for a new one. Refreshing does not
rotate the refresh token, and once a token has been revoked it is rejected on every
subsequent refresh.

### The seat flow

This is the core of the API. A seat is *held* while the user decides, then *booked* once
they commit.

**1. Hold a seat.** The response is the reservation you'll need to book it.

```bash
curl -X POST http://localhost:8080/api/seats/reserve \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"show_time_id":"0192...","seat":"A:4"}'
```

```jsonc
{
  "id": "0192...",              // use this to book the seat
  "show_time_id": "0192...",
  "user_id": "0192...",
  "seat_no": "A:4",
  "status": "available",        // a hold, not yet sold
  "reserved_at": "2026-09-26T12:04:11Z",
  "created_at": "2026-09-26T12:04:11Z",
  "updated_at": "2026-09-26T12:04:11Z"
}
```

The request fails with `422` if the seat doesn't exist in that cinema's seat map, has
already been booked, or is currently held by someone else.

**2. Hold the rest.** Repeat for each seat in the row. The reservation `id` is what you
carry forward, not the seat number.

**3. Book them.** Pass every reservation id at once.

```bash
curl -X POST http://localhost:8080/api/seats/book \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"seat_reservations":["0192...","0192..."]}'
```

Only reservations that belong to the caller and are still `available` are accepted. On
success each row comes back with `status: "booked"` and a null `reserved_at`.

**4. Review or cancel.**

```bash
curl http://localhost:8080/api/bookings -H "Authorization: Bearer $ACCESS_TOKEN"
curl -X DELETE http://localhost:8080/api/bookings/$RESERVATION_ID -H "Authorization: Bearer $ACCESS_TOKEN"
```

**How long a hold lasts:** 10 minutes, counted from `reserved_at`. There is no background
job expiring them — if a hold is stale, the next person to request that seat takes it, and
the original hold is reassigned. A hold is not a booking until step 3 succeeds.

### Seat maps

A cinema's layout is a JSONB document, which keeps irregular auditoriums out of the
relational schema. `total_seats` is denormalised into the document so capacity and revenue
can be computed without counting rows.

```jsonc
{
  "rows": [
    {
      "row": "A",
      "seats": [
        { "number": 1, "type": "regular" },
        { "number": 2, "type": "regular" }
      ]
    },
    { "row": "B", "seats": [ { "number": 1, "type": "vip" } ] }
  ],
  "total_seats": 3,
  "screen": "center"
}
```

Seats are addressed as `"{row}:{number}"`, so `A:4` above is row `A`, seat `4`. The
per-seat `type` is stored but never read — see
[Deliberately out of scope](#deliberately-out-of-scope).

### Creating a cinema

```bash
curl -X POST http://localhost:8080/api/cinemas \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Screen 4",
    "location_id": "0192...",
    "experience_types": ["2D", "3D"],
    "seat_map": { "rows": [...], "total_seats": 3, "screen": "center" }
  }'
```

### Scheduling a showtime

```bash
curl -X POST http://localhost:8080/api/movies/$MOVIE_ID/show-times \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_time": "19:30",
    "price": 25000,
    "description": "opening night",
    "price_currency": "UGX",
    "cinema_id": "0192...",
    "experience_type": "2D",
    "start_date": "2026-09-26",
    "end_date": "2026-10-26"
  }'
```

A showtime is a **date range**, not an instant: a movie playing at a cinema from one date
to another, at one price, in one experience type. That's a simplification — see
[Deliberately out of scope](#deliberately-out-of-scope).

### Error responses

Two shapes. Simple failures:

```json
{ "error": "Show time not found" }
```

Validation failures return `422` with a field-keyed map, so a client can highlight the
offending inputs:

```json
{
  "message": "There were some validation errors",
  "errors": {
    "email": ["The email must be a valid email."],
    "seat_reservations[0]": ["Seat Reservation not found."]
  }
}
```

---

## Architecture

### Layout

```
main.go               entrypoint: config, DB pool, server
handlers/             HTTP layer — routing, middleware, validation, JSON
  handlers.go         every route in the project, in one file
  config.go           ApiConfig: the single dependency-injection point
  middlewareAuth.go   requires a valid access token
  middlewareAdminAuth.go  requires a valid access token and is_admin
  json.go             response and error helpers
  validation.go       validator setup, error message translation
  dbValidatorFuncs.go custom validators that hit the database
  handle*.go          one file per route
internal/
  auth/               JWT creation/validation, password hashing, header parsing
  hashing/            vendored Argon2id primitives
  types/              custom sqlc column types (SeatMap, StringSlice)
  database/           repository layer — entirely sqlc-generated
  testdb/             test database helpers
sql/
  schema/             goose migrations
  queries/            sqlc queries, with -- name: QueryName :one/:many/:exec
  embed.go            embeds schema/*.sql into the binary
```

Dependencies point one way: `handlers` → `internal/*` → `internal/database`. Nothing in
`internal/` knows HTTP exists.

### Dependency Injection

There is no DI framework. `ApiConfig` carries a `database.Querier` **interface** rather
than a concrete `*database.Queries`, so tests can substitute either a transaction-bound
querier or a mock:

```go
type ApiConfig struct {
    DB        database.Querier
    Platform  string
    JwtSecret string
    Validate  *validator.Validate
}
```

A `*pgxpool.Pool` satisfies the same interface, so the production path passes the pool
straight through. See [`handlers/testHelpers_test.go`](handlers/testHelpers_test.go).

### Schema

| Table | Purpose |
| --- | --- |
| `users` | Accounts. Admin is a single `is_admin` boolean — no roles table. |
| `refresh_tokens` | Server-side session store. The token itself is the primary key; `revoked_at` is what makes revocation enforceable. |
| `movies` | Film catalogue. `genre` and `experience_types` are JSON arrays. |
| `locations` | Physical venues, e.g. a mall. |
| `cinemas` | An auditorium within a location, with its own JSONB `seat_map`. |
| `show_times` | A movie playing at a cinema across a date range. |
| `reservations` | One row per seat: the hold, and then the booking. |

All primary keys are `uuidv7()`, generated in SQL rather than the application, which keeps
them time-sortable and therefore index-friendly.

### How the seat model works

`reservations.status` is the whole state machine:

- `available` — a live hold, with `reserved_at` set. It belongs to a user but nobody has
  paid for it.
- `booked` — sold. `reserved_at` is nulled out on confirmation.

```sql
CONSTRAINT unx_show_time_seat_no UNIQUE (show_time_id, seat_no)
```

That constraint is the real guarantee. Because a seat can only have one row per showtime,
two concurrent requests for `A:4` cannot both succeed — one gets a unique violation, no
matter how the application code is written. The handler checks first for a friendlier
error message; the constraint is what makes the check trustworthy.

Expired holds are reclaimed lazily rather than swept. `POST /api/seats/reserve` reads the
existing row, and if `status = 'available'` but `reserved_at` is more than 10 minutes old,
it reassigns the row to the new user instead of failing. There's no cron job and no
per-seat timer.

### Domain invariants

Two rules are enforced in the handlers because the database can't express them:

- A movie with ongoing showtimes cannot be edited or deleted, so details people already
  booked against don't change underneath them.
- A cinema with ongoing showtimes cannot be updated, because changing a seat map would
  invalidate seats that have already been sold.

And one in the schema: `ON DELETE CASCADE` from `show_times`, `locations` and `users` means
deleting a movie, cinema, location or user cleans up the dependent rows for you.

### Where the JSONB is

Nested aggregates are built in SQL with `jsonb_agg` inside `LEFT JOIN LATERAL` rather than
fetched with follow-up queries, which keeps movie + showtime and location + cinema reads
to a single round trip. `json.RawMessage` fields appear in a couple of handlers to stop
sqlc from base64-encoding nested JSON columns on its way through the generated structs.

---

## Retrospective & Future Work

### Things to fix

- [ ] Seat bookings should be handled in a db transaction so that all are rolled back in case of an error
- [ ] Rename the `reservations` table to `bookings`, with `status` values of `pending` and `booked`. The structure can stay as it is, but `bookings` is the more intuitive name — right now routes say `/api/bookings` while the table behind them is called `reservations`, which is needlessly confusing
- [ ] Change every `timestamp` column to `timestamptz` in UTC. This would have prevented the time-zone bug I hit while calculating how long a seat had been reserved; I patched it in `7a47fd6` by converting the `reservations` table, but that migration shouldn't have been necessary
- [ ] Add pagination. No query uses `LIMIT`/`OFFSET`, so `GET /api/movies`, `GET /api/locations` and `GET /api/bookings` all return unbounded result sets

### Deliberately out of scope

- [ ] Payment processing or ticketing provider integration
- [ ] Email or SMS notifications, including booking confirmations and showtime reminders
- [ ] A frontend, or any static file serving — `SetupServer` takes a `filePathRoot` that it currently ignores
- [ ] Seat map editing UI. The layout is uploaded as a single JSONB document rather than drawn in a browser
- [ ] Seat classes and per-seat pricing. The per-seat `type` field is stored but never read, so it can't affect price
- [ ] A background sweeper for expired holds. The lazy reclaim on the next reserve attempt covers the same ground
- [ ] Per-instant showtimes. A showtime is a date range at a single price, so a cinema can't run different prices for a morning and an evening screening of the same film
- [ ] Showtime overlap detection. Nothing currently stops two showtimes from being scheduled in the same cinema at the same time. This is deliberately out of scope because the logic is substantial, but it's worth recording the approach I'd take. When creating a showtime, validate that the cinema is available at the proposed time — this can be done by checking whether any existing booking runs beyond it. Here's a brief proposal:

    1. Add a cleanup time to the cinema table and make it required for each cinema. This field indicates how long it takes the staff to clean up the cinema and have it ready for the next movie. It should also factor in the time it takes the people to exit the cinema.
    2. When looking for existing bookings (which could be a separate table or calculated on the fly), add this time to the movie length to determine when the cinema will next be available for another movie.
    3. The same check can back an API which the UI uses to display the available cinemas when the user selects a time slot.

- [ ] Movie poster uploads. `movies.poster_image_url` exists in the schema but no endpoint sets it
- [ ] Refresh token rotation, and reuse detection
- [ ] Roles and permissions beyond the single `is_admin` boolean
- [ ] Audit logging of admin actions
- [ ] Observability: no structured logging, metrics, tracing or request IDs
- [ ] Rate limiting, CORS, panic recovery, and graceful shutdown
- [ ] Docker, docker-compose, a Makefile, and CI
- [ ] A linter configuration
- [ ] Wider test coverage. Only `handleCreateMovie` has handler tests, and the seat reserve and book flows — the most logic-dense code here — have none
- [ ] Production password hashing parameters. `hashing.ProdParams` is defined and unused; the defaults are still in use

---

## Contributing

Contributions are welcome — issues, bug reports, and pull requests.

### Running the tests

```bash
cp .env.testing.example .env.testing   # required
```

Then set `TEST_DB_URL` in `.env.testing` to a **separate, throwaway** database. The test
helpers run the real migrations and wrap every test in a transaction that is rolled back
when the test finishes, so a test can never leave data behind — but pointing this at a
database you care about is still a bad idea.

```bash
go test ./...
```

Two things to know:

- If `TEST_DB_URL` is empty, tests that need a database skip themselves rather than fail.
  If the `.env.testing` file is *missing*, they fail outright.
- The path to `.env.testing` is relative (`../.env.testing`), so it resolves against the
  package directory of whichever test calls the helper. The only caller today lives in
  `handlers/`, which puts it at the repo root.

### The development loop

The database layer is generated, so SQL changes have two steps:

```bash
# 1. edit sql/queries/*.sql (and add a migration in sql/schema/ if the schema changed)
# 2. regenerate the typed Go
sqlc generate

# 3. rebuild and test
go build ./... && go test ./...
```

Everything in `internal/database/` is generated by sqlc and carries a
`Code generated by sqlc. DO NOT EDIT.` header. **Never edit it by hand** — change the SQL
and regenerate.

### Conventions

- **One route, one file.** `handleSeatBooking.go` contains `handleSeatBooking`, and nothing
  else does. Add new routes the same way.
- **One request struct per handler**, declared unexported at the top of the file, with
  `validate` tags. Tag names are translated to their JSON equivalents, so validation error
  keys match the wire format rather than the Go field names.
- **Migrations need both directions.** Every file in `sql/schema/` should have a
  `-- +goose Up` and a `-- +goose Down` section.
- **Business logic lives in the handler.** There is deliberately no service layer; handlers
  decode, validate, call a query, and respond. Keep it that way unless you're introducing a
  layer for a specific reason.
- **Let the database enforce invariants** where it can, and check in the handler only to
  return a friendlier error.
- **Never commit `.env` or `.env.testing`.** Both are gitignored. `.env.example` and
  `.env.testing.example` are committed and should stay key-complete and value-free.

### Known gaps in the tooling

There's no CI pipeline and no linter config committed yet, so run `go vet ./...` and
`gofmt -l .` locally before opening a pull request. There's also no `LICENSE` file — add
one before publishing this anywhere you care about.
