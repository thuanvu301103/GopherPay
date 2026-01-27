# Auth Service

## API Contract Document

- Update document:
```Bash
swag init -g cmd/api/main.go
```

- Endpoint: `http://localhost:3000/docs/index.html`

## Database Migration

### Using Atlas
- Change the `DB_AUTO_MIGRATE` in `.env` to `true`
- Create migration file: 
```bash
atlas migrate diff <file_name> --env gorm
```
- Run migration:
```bash
atlas migrate apply --env gorm --url "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

### Using AutoMigration
- Change the `DB_AUTO_MIGRATE` in `.env` to `false`
- Run the server

### Run Service
```Bash
go run cmd/api/main.go
```