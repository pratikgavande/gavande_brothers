# Billing API

Gin-based Go API using MySQL and a simple MVC-style layout.

## Requirements

- Go 1.22 or newer
- MySQL

## Setup

1. Create a MySQL database, for example:

   ```sql
   CREATE DATABASE billing_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

2. Copy `.env.example` to `.env` and update `MYSQL_DSN` with your local
   MySQL credentials. In PowerShell:

   ```powershell
   Copy-Item .env.example .env
   notepad .env
   ```

3. Download dependencies and start the API:

   ```powershell
   go mod tidy
   go run main.go
   ```

The application loads `.env` from the current working directory. Variables
already set in the shell take precedence over values in `.env`. Do not commit
`.env`; it is excluded by `.gitignore`.

4. Check the API and MySQL connection:

   ```text
   GET http://localhost:8080/api/v1/health
   ```

The health endpoint returns `200` with `{"status":"ok"}` when the API can reach
MySQL. It returns `503` if MySQL is unavailable.

## Project layout

- `config/`: application configuration
- `controllers/`: HTTP handlers and JSON responses
- `models/`: GORM database connection and, as features are added, data models
- `services/`: business logic to keep controllers focused on HTTP concerns
- `routes/`: endpoint registration
- `main.go`: application startup and dependency wiring

Business features such as authentication, users, roles, billing, and expenses
have not been implemented yet.
