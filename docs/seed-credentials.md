# Dev Seed Credentials

> For development and testing only. Do not use in production.

## System Users

| Email | Password | Role | How seeded |
|---|---|---|---|
| admin@lodge.dev | Admin@123 | admin | Seeded at startup by `main.go` (`SeedSuperAdmin`) |
| manager@lodge.dev | Manager@123 | manager | Migration `000014_seed_dev_users` |
| receptionist@lodge.dev | Reception@123 | receptionist | Migration `000014_seed_dev_users` |
| cleaner@lodge.dev | Cleaner@123 | cleaner | Migration `000014_seed_dev_users` |

## Running Migrations

```bash
# Apply all migrations
make migrate-up

# Roll back all
make migrate-down

# Roll back one step
migrate -path ./migrations -database "postgres://postgres:<password>@localhost:5432/lodge-management-system?sslmode=disable" down 1
```

## Database

- **DB Name:** `lodge-management-system`
- **Connection:** configured via `.env` (see `.env.example`)

> **Note:** The `admin@lodge.dev` user is NOT created by a migration — it is seeded automatically every time the server starts via `SeedSuperAdmin` in `main.go`. Run `make migrate-up` first, then `make run`.
