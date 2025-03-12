# ITV Movie Service API

A RESTful web service built with Go for managing movies, genres, languages, and authentication. This API provides endpoints for user authentication, movie management, and other related operations.

## 📜 API Documentation

Detailed API documentation is available in [Postman Docs](https://documenter.getpostman.com/view/27903487/2sAYk8wPTw).

## 🚀 Tech Stack

The project utilizes the following technologies:

- **Go** – The primary language used for development
- **Gin** – Web framework for routing
- **Uber-FX** – Dependency injection framework
- **GORM** – ORM for database interactions
- **AtlasGO** – Database versioned migration management
- **Docker** – For containerization and deployment
- **JWT** – Security and authentication

## 📌 Project Structure

```
├── cmd
│   └── movie-service
│       └── main.go          # Entry point of the application
├── config/                  # Configuration files
│   ├── config.yml
│   ├── local.yml
│   └── release.yml
├── internal
│   ├── api
│   │   ├── handlers/        # API request handlers
│   │   ├── middlewares/     # Middleware for authentication and permissions
│   │   ├── routes/          # API route definitions
│   │   ├── services/        # Business logic layer
│   ├── config/              # Configuration utilities
│   ├── models/              # Database models
│   ├── pkg
│   │   ├── jwt/             # JWT token utilities
│   │   ├── utils/           # Helper utilities and logging
│   └── storage
│       └── database
│           ├── postgres.go  # PostgreSQL connection setup
│           ├── repository.go # Generic repository interface
│           └── repositories/ # Repository implementations
├── migrations/              # Database migration files
├── Dockerfile               # Docker configuration
├── docker-compose.yml       # Docker Compose configuration
├── go.mod                   # Go module file
├── go.sum                   # Go dependencies
```

## 🔥 Why AtlasGO Instead of GORM AutoMigrate?

AtlasGO is used for database migrations instead of GORM's `AutoMigrate` due to its ability to:

- Provide structured, version-controlled migrations.
- Prevent accidental schema changes.
- Enable rollback capabilities for failed migrations.
- Maintain database consistency across different environments.

## 📦 Repository Pattern

This project follows the **Repository Pattern** to:

- Decouple database logic from business logic.
- Enhance testability by allowing mock implementations.
- Improve maintainability by abstracting data access.

## 📡 API Endpoints

### 🛠️ Authentication

- **POST** `/api/v1/auth/register-admin` – Register an admin user
- **POST** `/api/v1/auth/register` – Register a user
- **POST** `/api/v1/auth/login` – Login and obtain JWT token
- **POST** `/api/v1/auth/refresh` – Refresh JWT token
- **POST** `/api/v1/auth/logout` – Logout user
- **GET** `/api/v1/auth/admin/users` – Fetch all users (Admin only)
- **PUT** `/api/v1/auth/admin/status` – Update user status
- **DELETE** `/api/v1/auth/admin/users/{id}` – Delete a user

### 🎬 Movies

- **POST** `/api/v1/movies` – Add a new movie
- **GET** `/api/v1/movies` – Get all movies (search, pagination supported)
- **GET** `/api/v1/movies/{id}` – Get a specific movie
- **PUT** `/api/v1/movies/{id}` – Update movie details
- **DELETE** `/api/v1/movies/{id}` – Delete a movie

### 🎭 Genres

- **POST** `/api/v1/genres` – Create a new genre
- **GET** `/api/v1/genres` – Get all genres
- **GET** `/api/v1/genres/{id}` – Get a specific genre
- **PUT** `/api/v1/genres/{id}` – Update genre details
- **DELETE** `/api/v1/genres/{id}` – Delete a genre

### 🌎 Languages

- **POST** `/api/v1/languages` – Create a new language
- **GET** `/api/v1/languages` – Get all languages
- **GET** `/api/v1/languages/{id}` – Get a specific language
- **PUT** `/api/v1/languages/{id}` – Update language details
- **DELETE** `/api/v1/languages/{id}` – Delete a language

### 🌍 Countries

- **POST** `/api/v1/countries` – Create a new country
- **GET** `/api/v1/countries` – Get all countries
- **GET** `/api/v1/countries/{id}` – Get a specific country
- **PUT** `/api/v1/countries/{id}` – Update country details
- **DELETE** `/api/v1/countries/{id}` – Delete a country

## 🏗️ Deployment (Docker)

To build and run the project using Docker:

```sh
docker build -t itv-movie-service .
docker run -p 8080:8080 itv-movie-service
```

or better use docker compose

```sh
docker compose up -d
```

## 🛡️ Security

- JWT authentication is used for securing endpoints.
- Role-based access control for admin and users.