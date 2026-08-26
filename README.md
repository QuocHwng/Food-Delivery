# Food Delivery — Restaurant Order System

A microservices-based food delivery system built with Go (Gin), ReactJS, PostgreSQL, and RabbitMQ.

## Tech Stack

| Layer | Technology |
|---|---|
| Backend (Gateway + Services) | Go + Gin + sqlx |
| Message Broker | RabbitMQ |
| Real-time | WebSocket (gorilla/websocket) |
| Database | PostgreSQL |
| Payment | VNPay Sandbox |
| Frontend | ReactJS + Vite + TypeScript |
| Container | Docker + Docker Compose |
| Local Cloud | Floci |

## Services

| Service | Port | Description |
|---|---|---|
| API Gateway | 8080 | JWT auth, reverse proxy, rate limiting |
| User Service | 8081 | Auth, profile, addresses |
| Order Service | 8082 | Orders, state machine, coupons |
| Payment Service | 8083 | VNPay sandbox integration |
| Restaurant Service | 8084 | Menu, dashboard, statistics |
| Notification Service | 8085 | WebSocket real-time push |

## Quick Start

```bash
# Clone repo
git clone https://github.com/QuocHung52/food-delivery.git
cd food-delivery

# Start infrastructure (PostgreSQL + RabbitMQ)
make infra-up

# Run database migrations
make migrate-up

# Start all services
make up
```

## Project Structure

```
├── backend/          # All Go microservices
│   ├── gateway/      # API Gateway
│   ├── services/     # Individual microservices
│   └── shared/       # Shared Go packages
├── frontend/         # ReactJS + Vite
├── infrastructure/   # Migrations, Docker, Floci
├── docs/             # Architecture & API docs
├── scripts/          # Dev helper scripts
└── docker-compose.yml
```

## Development

See [docs/architecture.md](docs/architecture.md) for full architecture details.
