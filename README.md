# Friends-Management

## Overview
This project implements a simple Friend Management API service using Go. It simulates core social networking features like managing friendships, subscriptions, blocks, and update notifications. The API is designed around a set of user stories and aims to be extensible and production-ready.

---

## User stories
- As a user, I need an API to create a friend connection between two email
addresses.
- As a user, I need an API to retrieve the friends list for an email address.
- As a user, I need an API to retrieve the common friends list between two
email addresses.
- As a user, I need an API to subscribe to updates from an email address.
- As a user, I need an API to block updates from an email address.
- As a user, I need an API to retrieve all email addresses that can receive
updates from an email address.
---

## Technical choices

- Go v1.24+
- Gin – HTTP framework
- GORM – ORM for PostgreSQL
- Docker + Railway

---

## How to run

### Option 1: Docker Compose 

1. Clone the repo:
```bash
git clone https://github.com/tienkhoa03/Friends-Management.git
```
2. Run Docker Compose

```bash
docker compose up --build
```
This command:
- Builds the Go API server
- Runs the service on localhost:8080
- Automatically initialize the database with default schema/config

### Option 2: Manual

1. Clone the repo:
```bash
git clone https://github.com/tienkhoa03/Friends-Management.git
```
2. Set up dependencies
```bash
cd be
go mod tidy
```
3. Run the APIs:
```bash
cd be/cmd/server
go run main.go
```

- API runs on: http://localhost:8080
- You must manually ensure the database is running and accessible (configured via .env)

### Option 3: Dockerfile (Single-container build)
1. Build the Docker image:
```bash
docker build -t friends-api .
```
2. Run the container:
```bash
docker run -p 8080:8080 --env-file .env friends-api
```
- API runs on: http://localhost:8080
- You must manually ensure the database is running and accessible (configured via .env)

---

## API Endpoints

### **Users**

| Method | Endpoint           | Description    |
| ------ | ------------------ | -------------  |
| GET    | /api/users         | Get all users  |
| GET    | /api/users/{id}    | Get user by ID |
| POST   | /api/users         | Create new user|
| PUT    | /api/users/{id}    | Update user    |
| DELETE | /api/users/{id}    | Delete user    |

### **Friendship**

| Method | Endpoint                       | Description                               |
| ------ | ------------------------------ | ----------------------------------------- |
| POST   | /api/friendship                | Create new friendship                     |
| GET    | /api/friendship/friends        | Retrieve friends list for an email address|
| GET    | /api/friendship/common-friends | Retrieve common friends list between      |

### **Subscription**

| Method | Endpoint           | Description             |
| ------ | ------------------ | ----------------------  |
| POST   | /api/subscription  | Create new subscription |

### **Block**

| Method | Endpoint           | Description                    |
| ------ | ------------------ | -----------------------------  |
| POST   | /api/block         | Create new block relationship  |

### **Notification**

| Method | Endpoint                | Description            |
| ------ | ----------------------- | -------------          |
| POST    | /api/update-recipients | Get update recipients  |

---

## Project Structure

```
be
├── .github/workflows/   # CI/CD workflows (e.g., GitHub Actions)
├── .vscode/             # VSCode editor settings and workspace configuration
├── api/                 # API layer: routing, handlers, middleware
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # HTTP middlewares (e.g., logging, auth)
│   └── router/          # API route definitions
├── cmd/server/          # Main application entry point (e.g., main.go)
├── config/              # Application configurations and environment loading
├── constant/            # Constant values, enums, status codes
├── internal/            # Internal application logic (domain-driven design)
│   ├── domain/
│   │   ├── dto/         # Data Transfer Objects (used between layers)
│   │   └── entity/      # Core business entities / database models
│   ├── repository/      # Repository interfaces and their implementations
│   └── service/         # Business logic and use cases
├── pkg/                 # Reusable helper packages (e.g., JWT, hashing, utils)
├── .env.template        # Template for environment variables
├── .gitignore           # Git ignore rules
├── docker-compose.yml   # Docker multi-service configuration
├── Dockerfile           # Docker build instructions for the app
├── go.mod               # Go module file (dependencies and module name)
└── go.sum               # Go checksum file (dependency verification)

```

---

## Environment Variables

| Key          | Description        |
| ------------ | ------------------ |
| PORT         | Port server        |
| DB\_USER     | Database username  |
| DB\_PASSWORD | Database password  |
| DB\_NAME     | Database name      |
| DB\_HOST     | Database host      |


Create `.env` file base on `.env.template`.

---

## Testing
- Unit tests cover all core services and edge cases
- Run tests with:
```bash
cd be
go test ./...
```

---

## Deployment

Deploy using **Railway**:

1. Connect your GitHub repository.
2. Set up environment variables according to `.env.template`.
3. Railway will automatically build and deploy the container.

---

## Production

The API has been deployed and is publicly available at:

[https://friends-management-production.up.railway.app/swagger/index.html#/](https://friends-management-production.up.railway.app/swagger/index.html#/)

---
## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---
