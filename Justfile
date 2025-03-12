run:
    go run cmd/movie-service/main.go

migrate-diff:
    atlas migrate diff --env gorm

migrate-local:
    atlas migrate apply --dir "file://migrations" --url "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"