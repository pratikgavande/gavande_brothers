# Billing API

Gin-based Go API using MySQL and a simple MVC-style layout.

## API documentation

The OpenAPI 3 specification is in [`docs/openapi.yaml`](docs/openapi.yaml).
Open it in Swagger Editor or import it into an API client to browse and try
the currently implemented endpoints.

## Requirements

- Go 1.22 or newer
- MySQL

## Setup

1. Create a MySQL database, for example:

   ```sql
   CREATE DATABASE billing_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

2. Ensure `.env` exists and set `MYSQL_DSN` and a private, random `JWT_SECRET`
   of at least 32 characters. If creating it for the first time:

   ```powershell
   Copy-Item .env.example .env
   notepad .env
   ```

3. Run the authentication schema in `database/migrations/001_auth_roles.sql`
   against your MySQL database. This creates `roles` and `users`, and seeds
   the protected `superadmin` and `employee` roles. Then start the API:

   ```powershell
   go mod tidy
   go run main.go
   ```

The application loads `.env` from the current working directory. Shell
environment variables take precedence. Do not commit `.env`; it is excluded
by `.gitignore`.

4. Create the first superadmin:

   ```powershell
   go run ./cmd/hashpassword
   ```

   This prompts for an 8-12 character password without displaying it and prints
   a bcrypt hash.
   Replace `REPLACE_WITH_BCRYPT_HASH` in
   `database/queries/create_first_superadmin.sql` with that hash, set a unique
   admin email, and execute the query against the database.
   To add an initial employee manually, repeat the hash step and use
   `database/queries/create_employee.sql` with a different email and hash.

5. Check the API and MySQL connection:

   ```text
   GET http://localhost:8080/api/v1/health
   ```

The health endpoint returns `200` with `{"status":"ok"}` when the API can reach
MySQL. It returns `503` if MySQL is unavailable.

## Authentication and roles

- `POST /api/v1/auth/login` accepts `{"email":"...","password":"..."}` and
  returns a 15-minute HS256 bearer token.
- Send that token as `Authorization: Bearer <token>` to
  `GET /api/v1/auth/me` and the role-management endpoints.
- `GET`, `POST`, `PATCH`, and `DELETE /api/v1/roles` endpoints are restricted
  to the superadmin. System roles cannot be edited or deleted; a role assigned
  to users cannot be deleted.
- Role endpoints are `GET/POST /api/v1/roles`,
  `PATCH /api/v1/roles/:id`, and `DELETE /api/v1/roles/:id`.
- Users and roles are loaded from MySQL on authenticated requests, so changing
  a user's active state or role takes effect without waiting for token expiry.
- Login uses email and bcrypt-hashed passwords. Role permissions and employee
  user-management endpoints are not included yet.

## Billing APIs

The existing `customer`, `bills`, `bill_items`, and `accounting` tables are
queried and written through explicit GORM column mappings. The backend does not
migrate or change the schema. All billing endpoints require the existing JWT:

```http
Authorization: Bearer <access_token>
```

In Postman, first call `POST /api/v1/auth/login`, copy `access_token`, then
select **Authorization → Bearer Token** for the billing requests.

| Method | Path | Query parameters |
|---|---|---|
| GET | `/api/v1/customers` | `page` (default 1), `limit` (default 10, max 100), `search` (name/mobile) |
| GET | `/api/v1/customers/:id` | — |
| GET | `/api/v1/bills` | `page`, `limit`, `customerId`, `date` (`YYYY-MM-DD`), `billNo` |
| GET | `/api/v1/bills/:id` | — |
| GET | `/api/v1/bills/:id/items` | — |
| GET | `/api/v1/bill-items/:id` | — |
| PATCH | `/api/v1/customers/:id` | partial customer fields |
| PATCH | `/api/v1/bills/:id` | permitted bill fields only |
| PATCH | `/api/v1/bill-items/:id` | partial bill-item fields |
| GET | `/api/v1/bills/:id/pdf` | downloads stored bill data as a PDF |
| POST | `/api/v1/bills/:id/send` | emails the generated PDF to the customer's email |
| GET | `/api/v1/accounting` | `page`, `limit`, `date` (`YYYY-MM-DD`), `status` |
| GET | `/api/v1/accounting/:id` | — |

Postman examples (replace the host, IDs, and token as needed):

```http
GET http://localhost:8080/api/v1/customers?page=1&limit=10&search=Gavande
Authorization: Bearer <access_token>

GET http://localhost:8080/api/v1/bills?page=1&limit=10&customerId=1&date=2026-11-06
Authorization: Bearer <access_token>

GET http://localhost:8080/api/v1/bills/1/items
Authorization: Bearer <access_token>

GET http://localhost:8080/api/v1/accounting?page=1&limit=10&date=2026-11-06&status=1
Authorization: Bearer <access_token>
```

List responses include `success`, `message`, `data`, and pagination metadata
(`page`, `limit`, `total`). Single-record responses include `success`,
`message`, and `data`. Invalid parameters return HTTP 400, unauthenticated
requests return HTTP 401, missing records return HTTP 404, and query failures
return HTTP 500. See [`docs/openapi.yaml`](docs/openapi.yaml) for schemas and
all documented responses.

### Partial update examples

Only fields supplied in the JSON body are updated. Customer name and mobile
cannot be blank. Bill updates deliberately exclude bill number, customer ID,
totals, amount received/pending, and item count. Bill-item `date` can be set
to `null` to clear its nullable TEXT column.

```http
PATCH http://localhost:8080/api/v1/customers/1
Authorization: Bearer <access_token>
Content-Type: application/json
```

```json
{
  "mobile": "9876543210",
  "address": "Pune, Maharashtra"
}
```

```http
PATCH http://localhost:8080/api/v1/bills/1
Authorization: Bearer <access_token>
Content-Type: application/json
```

```json
{
  "date": "2026-11-10T09:00:00+05:30",
  "site_name": "Site B",
  "gst_amount": 0
}
```

```http
PATCH http://localhost:8080/api/v1/bill-items/1
Authorization: Bearer <access_token>
Content-Type: application/json
```

```json
{
  "rate": "5500",
  "date": null
}
```

### Download a bill PDF

Use `GET /api/v1/bills/:id/pdf` with the bearer token. A successful response
has `Content-Type: application/pdf` and an attachment filename
`gavande-bros-bill-<bill-number>.pdf`. The PDF reads the stored bill, customer,
and items; GST status/amount and all totals are shown as stored without
recalculation.

```http
GET http://localhost:8080/api/v1/bills/1/pdf
Authorization: Bearer <access_token>
```

In Postman, choose **Send and Download** to save the PDF response.

### Email a bill PDF

The project has no prior email integration. Configure SMTP through the
following `.env` variables before calling `POST /api/v1/bills/:id/send`:

```text
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your-smtp-user
SMTP_PASSWORD=your-smtp-password
SMTP_FROM=billing@example.com
```

The service requires STARTTLS on the configured SMTP connection (or implicit
TLS on port 465). Keep real credentials only in the ignored `.env` file.
Success is returned only after the SMTP server accepts the message; this does
not claim final inbox delivery.

```http
POST http://localhost:8080/api/v1/bills/1/send
Authorization: Bearer <access_token>
```

Invalid input returns HTTP 400, unauthenticated requests HTTP 401, missing
records HTTP 404, invalid customer email HTTP 422, SMTP rejection HTTP 502,
and missing SMTP configuration HTTP 503. Partial updates return HTTP 200.
The existing GET and POST response formats remain unchanged.

### Create customer

`POST /api/v1/customers` requires all database NOT NULL columns. `name` and
`mobile` must be non-empty; the database generates the ID. Empty strings are
accepted for GST number, email, and address.

```http
POST http://localhost:8080/api/v1/customers
Authorization: Bearer <access_token>
Content-Type: application/json
```

```json
{
  "name": "Gavande Bros",
  "mobile": "9876543210",
  "gst_no": "",
  "email": "",
  "address": "Pune, Maharashtra",
  "status": 1
}
```

### Create bill

`POST /api/v1/bills` requires an existing `customer_id` and all database
NOT NULL columns. The API stores supplied amount/GST values as-is; it does not
calculate totals, pending amounts, GST, or bill numbers. The database generates
the bill ID.

```http
POST http://localhost:8080/api/v1/bills
Authorization: Bearer <access_token>
Content-Type: application/json
```

```json
{
  "bill_no": 1001,
  "customer_id": 1,
  "date": "2026-11-10T09:00:00+05:30",
  "gst_status": 1,
  "gst_amount": 0,
  "amount": 15000,
  "amount_received": 0,
  "amount_pending": 15000,
  "engineers_name": "Engineer Name",
  "site_name": "Site A",
  "location": "Pune",
  "item_count": 2,
  "status": 1
}
```

### Add multiple bill items

`POST /api/v1/bills/:id/items` takes a non-empty `items` array. The path bill
must exist. All inserts happen in a single transaction. `date` is nullable
TEXT; rates and amounts are stored as given, with no calculations.

```http
POST http://localhost:8080/api/v1/bills/1/items
Authorization: Bearer <access_token>
Content-Type: application/json
```

```json
{
  "items": [
    {
      "name": "Installation",
      "rate": "5000",
      "unit": "service",
      "amount": 10000,
      "date": "2026-11-10",
      "status": 1
    },
    {
      "name": "Materials",
      "rate": "5000",
      "unit": "lot",
      "amount": 5000,
      "date": null,
      "status": 1
    }
  ]
}
```

Writes return HTTP 201 and the created record(s) in the standard
`success`/`message`/`data` response. Invalid input returns HTTP 400, missing
customer or bill returns HTTP 404, uniqueness conflicts return HTTP 409, and
other database errors return HTTP 500. The existing GET endpoint responses
remain unchanged. See [`docs/openapi.yaml`](docs/openapi.yaml) for schemas.

## Project layout

- `config/`: application configuration
- `controllers/`: HTTP handlers and JSON responses
- `models/`: GORM database connection and, as features are added, data models
- `services/`: business logic to keep controllers focused on HTTP concerns
- `middleware/`: JWT authentication and superadmin authorization
- `database/`: SQL migration and bootstrap query; run these manually
- `cmd/hashpassword/`: hidden-input bcrypt hash generator for the initial admin
- `routes/`: endpoint registration
- `main.go`: application startup and dependency wiring

Expense APIs are not implemented.
