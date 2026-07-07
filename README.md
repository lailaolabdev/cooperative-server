# Cooperative API

Go Gin + MongoDB API for public cooperative map data, admin authentication, and cooperative CRUD.

## Run locally

From the repository root:

```bash
docker compose up --build
cd cooperative-website
pnpm dev
```

Admin accounts are stored in MongoDB with bcrypt password hashes. Create the first account interactively:

```bash
go run ./cmd/create-admin
```

To run Gin without Docker, start MongoDB and run:

```bash
go run ./cmd/api
```

The application automatically loads `cooperative-service/.env`. Real environment variables override values from that file.

Seed the 58 cooperative records (safe to run repeatedly):

```bash
go run ./cmd/seed-cooperatives --file ../cooperatives_58.json
```

Seed SCU records from the LASCU workbook and remove deprecated saving/service types:

```bash
go run ./cmd/seed-scu --file "../LASCU members and SCU in Laos 2025.xlsx"
```

## Routes

- `POST /api/v1/auth/admin/login`
- `GET /api/v1/cooperatives?type=agriculture&province=VC`
- `GET /api/v1/cooperatives/:id`
- `POST /api/v1/admin/cooperatives` (Bearer token)
- `PUT /api/v1/admin/cooperatives/:id` (Bearer token)
- `DELETE /api/v1/admin/cooperatives/:id` (Bearer token)
