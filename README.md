# GateForge

[![Go](https://img.shields.io/badge/Go-%3E%3D1.22-blue?logo=go)](https://go.dev)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.29+--blue?logo=kubernetes)](https://kubernetes.io)
[![Docker](https://img.shields.io/badge/Docker-Multi--stage-blue?logo=docker)](https://www.docker.com)
[![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-Tracing_%26_Metrics-orange?logo=opentelemetry)](https://opentelemetry.io)
[![License](https://img.shields.io/badge/License-Apache%202.0-green?logo=apache)](https://opensource.org/licenses/Apache-2.0)

> Production-ready microservice template with built-in auth, full observability stack, and Kubernetes deployment manifests.

---

## Architecture Overview

```
┌──────────────┐     ┌──────────────┐     ┌──────────────────────────────────┐
│   Web SPA     │────▶│   Ingress    │────▶│        Go Service (gRPC/REST)     │
│  (React TS)   │     │  (Nginx/     │     │                                  │
│  Registration │     │   Istio)     │     │  ┌─────────┐  ┌──────────────┐  │
│  & Login      │     └──────────────┘     │  │ gRPC    │  │  Auth / JWT   │  │
└──────────────┘                           │  │ Gateway │  │  Handler      │  │
                                           │  └────┬────┘  └──────┬───────┘  │
                                           │       │              │          │
                                           │  ┌────▼────┐  ┌──────▼───────┐  │
                                           │  │Metrics  │  │  PostgreSQL  │  │
                                           │  │/metrics │  │              │  │
                                           │  └─────────┘  └──────────────┘  │
                                           │                                  │
                                           │  ┌──────────────────────────┐    │
                                           │  │  OpenTelemetry Collector │    │
                                           │  └──┬─────────┬────────┬───┘    │
                                           │     │         │        │       │
                                           │     ▼         ▼        ▼       │
                                           │  Prometheus  Jaeger   Loki    │
                                           │  (Metrics)  (Traces) (Logs)   │
                                           └──────────────────────────────────┘
```

**Data flow:** Browser → Ingress → gRPC/REST Gateway → Business Logic → PostgreSQL  
**Telemetry flow:** Services → OTel SDK → Collector → Prometheus / Jaeger / Loki → Grafana

---

## Key Features

- **Clean Architecture** — Strict separation of `cmd/`, `internal/`, `pkg/` layers; dependency inversion enforced by interfaces
- **gRPC-First Internal Contracts** — Protocol Buffers defined services for inter-service communication; REST Gateway for external clients
- **JWT Authentication** — Stateless auth service with access/refresh token rotation, bcrypt password hashing, and role-based access control
- **Multi-Stage Docker Builds** — Distroless / minimal-alpine Go images (~15 MB base), hardened Nginx builds for frontend
- **Production Kubernetes Manifests** — Deployments, Services, Ingress, ConfigMaps, Secrets, HPA, PDBs — ready for `kubectl apply -f k8s/`
- **Observability Out-of-the-Box** — Structured JSON logging, Prometheus `/metrics`, distributed tracing via OpenTelemetry/Jaeger
- **Infrastructure-as-Code Ready** — `docker-compose.yaml` for local dev; K8s manifests for production

---

## Project Layout

```
gateforge/
├── cmd/                          # Application entrypoints
│   ├── server/                   # Main entrypoint for the Go service
│   └── migrate/                  # Database migration runner
├── internal/                     # Business logic (private to the module)
│   ├── auth/                     # JWT auth, token management, password hashing
│   ├── handler/                  # HTTP/gRPC handlers and middleware
│   ├── service/                  # Core business logic layer
│   ├── repository/               # PostgreSQL data access layer
│   ├── model/                    # Domain models and DTOs
│   └── telemetry/                # OpenTelemetry instrumentation, middleware
├── proto/                        # Protobuf definitions (.proto files)
│   └── v1/                       # Service contracts v1
├── web/                          # React + TypeScript SPA (frontend)
│   ├── src/
│   │   ├── components/           # Reusable UI components
│   │   ├── pages/                # Registration, Login, Dashboard
│   │   ├── services/             # API client (fetch/gRPC-web)
│   │   └── app/                  # App shell, routing, providers
│   ├── public/
│   ├── package.json
│   └── tsconfig.json
├── k8s/                          # Kubernetes manifests
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml               # Base64-encoded; use Sealed Secrets / SOPS in prod
│   ├── postgres-deployment.yaml
│   ├── postgres-service.yaml
│   ├── app-deployment.yaml
│   ├── app-service.yaml
│   ├── ingress.yaml
│   ├── prometheus-deployment.yaml
│   ├── prometheus-service.yaml
│   ├── grafana-deployment.yaml
│   ├── grafana-service.yaml
│   ├── jaeger-deployment.yaml
│   ├── jaeger-service.yaml
│   └── hpa.yaml
├── docker-compose.yaml           # Local dev stack (Postgres, Jaeger, Prometheus)
├── Dockerfile.go                 # Multi-stage Go builder → distroless runtime
├── Dockerfile.web                # Multi-stage React build → Nginx serve
├── Makefile                      # Common targets: build, test, docker, deploy, migrate
├── go.mod / go.sum               # Go module definitions
├── .env.example                  # Environment variable template
├── LICENSE                       # Apache 2.0
└── README.md
```

---

## Observability Stack

| Signal    | Tool                  | Endpoint / Method                          |
|-----------|-----------------------|--------------------------------------------|
| **Metrics** | Prometheus          | `GET /metrics` (Go service), scraped via ServiceMonitor |
| **Traces**  | Jaeger / OTel       | gRPC/HTTP spans exported to OTel Collector   |
| **Logs**    | Structured JSON     | JSON-formatted logs shipped to Loki/Grafana  |
| **Dashboards** | Grafana         | Pre-configured dashboards for metrics & traces |

**Instrumentation:** All services use the OpenTelemetry Go SDK (`go.opentelemetry.io/otel`). Middleware captures request duration, status codes, and metadata. Custom business events are instrumented via the `internal/telemetry` package.

---

## Quick Start / Local Development

### Prerequisites

- Go ≥ 1.22
- Node.js ≥ 18 / npm
- Docker & Docker Compose
- `protoc` (for protobuf regeneration)

### 1. Start Infrastructure

```bash
docker-compose up -d postgres jaeger prometheus grafana
```

This spins up PostgreSQL, Jaeger, Prometheus, and Grafana with pre-configured scraping and dashboard provisioning.

### 2. Run the Backend

```bash
cp .env.example .env
# Edit .env with local database credentials

make backend-build
make backend-run
# Or: go run ./cmd/server
```

The service starts on `:8080` (REST) and `:9090` (gRPC). Metrics available at `:8080/metrics`. Traces exported to Jaeger at `:14268`.

### 3. Run the Frontend

```bash
cd web
npm install
npm run dev
```

Frontend available at `http://localhost:3000`. Registration and login forms are fully functional against the backend API.

### 4. Verify Observability

- **Grafana:** `http://localhost:3000` (admin/admin)
- **Jaeger:** `http://localhost:16686`
- **Prometheus:** `http://localhost:9090`

---

## Deployment (Kubernetes)

### 1. Build Docker Images

```bash
make docker-build
# Produces: gateforge/gateforge:latest
#           gateforge/gateforge-web:latest
```

Push to your registry:

```bash
docker push registry.example.com/gateforge/gateforge:latest
docker push registry.example.com/gateforge/gateforge-web:latest
```

### 2. Update Manifests

Edit `k8s/secret.yaml` with base64-encoded credentials. Update image tags in `k8s/app-deployment.yaml` and `k8s/postgres-deployment.yaml` if using custom registries.

### 3. Deploy to Cluster

```bash
kubectl apply -f k8s/
```

Verify:

```bash
kubectl get pods -n gateforge
kubectl get svc -n gateforge
kubectl get ingress -n gateforge
```

### 4. Run Migrations

```bash
kubectl apply -f k8s/migrate-job.yaml
# Or locally: make migrate
```

---

## Environment Variables

| Variable                  | Default                  | Description                                      |
|---------------------------|--------------------------|--------------------------------------------------|
| `APP_ENV`                 | `development`            | `development` / `production`                     |
| `APP_PORT`                | `8080`                   | REST API port                                    |
| `GRPC_PORT`               | `9090`                   | gRPC server port                                 |
| `DATABASE_URL`            | `postgres://localhost:5432/gateforge` | PostgreSQL connection string          |
| `DATABASE_SSLMODE`        | `disable`                | `disable` / `require` / `verify-full`            |
| `JWT_SECRET`              | *(required)*             | Signing key for access tokens                    |
| `JWT_REFRESH_SECRET`      | *(required)*             | Signing key for refresh tokens                   |
| `JWT_ACCESS_TTL`          | `15m`                    | Access token TTL                                 |
| `JWT_REFRESH_TTL`         | `7d`                     | Refresh token TTL                                |
| `BCRYPT_COST`             | `12`                     | bcrypt hashing cost factor                       |
| `OTEL_EXPORTER_ENDPOINT`  | `http://localhost:4317`  | OTel Collector gRPC endpoint                     |
| `PROMETHEUS_ADDR`         | `:8080/metrics`          | Prometheus metrics bind address                  |
| `LOG_LEVEL`               | `info`                   | `debug` / `info` / `warn` / `error`              |
| `LOG_FORMAT`              | `json`                   | `json` / `text`                                  |
| `INGRESS_HOST`            | `localhost`              | Ingress host for Kubernetes deployment           |

---

## Roadmap / Hackathon Ready

This template is architected to scale from hackathon MVP to production-grade platform:

- **[ ]** Add additional microservices (e.g., `user-service`, `notification-service`) with shared protobuf contracts
- **[ ]** Integrate Redis for session caching and rate limiting
- **[ ]** Add CI/CD pipeline (GitHub Actions → Docker Buildx → K8s rollout)
- **[ ]** Implement gRPC load balancing with client-side service discovery
- **[ ]** Add e2e tests (Playwright for frontend, Go testing for backend)
- **[ ]** Pod security policies / OPA Gatekeeper policies for production hardening
- **[ ]** ArgoCD / Flux CD for GitOps continuous deployment
- **[ ]** Multi-cluster deployment with federation
- **[ ]** Audit logging and compliance trail
- **[ ]** Feature flag integration (LaunchDarkly / Unleash)

---

## License

Distributed under the Apache License 2.0. See `LICENSE` for full terms.

```
Copyright 2026 GateForge Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
