# SpotifyClone

Microservices music streaming platform built with Go and Angular.

## Tech Stack

**Backend:** Go, Gin Framework, JWT, OpenTelemetry
**Frontend:** Angular 21, TypeScript, Bootstrap
**Databases:** MongoDB, Redis, Cassandra, Neo4j
**Infrastructure:** Docker, Jaeger

## Architecture

| Service | Port | Database | Description |
|---------|------|----------|-------------|
| API Gateway | 8080 | - | Routing, auth middleware |
| Users Service | 8001 | MongoDB | Authentication, profiles |
| Content Service | 8002 | MongoDB | Songs, albums, artists, genres |
| Ratings Service | 8003 | Redis | Song ratings |
| Subscriptions Service | 8004 | Redis | Artist/genre follows |
| Notifications Service | 8005 | Cassandra | User notifications |
| Recommendation Service | 8006 | Neo4j | Song recommendations |
| Frontend | 4200 | - | Angular SPA |

## Features

- JWT authentication with OTP and magic link support
- Music catalog management (CRUD)
- Rating system with Redis caching
- Artist/genre subscriptions
- Real-time notifications
- Graph-based recommendations
- Distributed tracing (Jaeger)
- Swagger API documentation

## Quick Start

```bash
docker-compose up --build
```

**With TLS:**
```bash
./generate-certs.sh
docker-compose -f docker-compose.yml -f docker-compose.tls.yml up --build
```

## URLs

- Frontend: http://localhost:4200
- API: http://localhost:8080/api/v1
- Swagger: http://localhost:8080/swagger/index.html
- Jaeger: http://localhost:16686
