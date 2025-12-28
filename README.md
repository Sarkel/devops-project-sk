# DevOps Project SK

A comprehensive DevOps demonstration project integrating modern software development practices, containerization, monitoring, and infrastructure automation.

## 🚀 Technologies

- **Backend:** Go (Golang)
- **Frontend:** React + Vite
- **Database:** PostgreSQL
- **Migrations:** Goose
- **SQL Code Gen:** SQLC
- **Communication:** MQTT (Mosquitto)
- **Containerization:** Docker & Docker Compose
- **Infrastructure as Code (IaC):** Terraform (Azure)
- **CI/CD:** GitHub Actions
- **Monitoring & Observability:** Prometheus, Grafana, Loki, Promtail

## 📁 Project Structure

- `app/` - Backend source code (Go)
  - `cmd/api` - API service
  - `cmd/crawler` - Data crawling service
  - `cmd/reader` - MQTT data processing service
- `web/` - Frontend application (React)
- `common/` - Common Go libraries (logger, db, config, mqtt)
- `docker/` - Docker configurations, Dockerfiles, and service configs (nginx, prometheus, etc.)
- `infra/` - Terraform infrastructure definitions (Azure)
- `migrations/` - SQL database migrations
- `scripts/` - Helper shell scripts
- `seeder/` - Test data generation tool
- `docker-compose.yml` - Local environment definition

## 🛠️ Local Setup

### Prerequisites
- Docker & Docker Compose
- Go (optional, for development)
- Node.js & npm (optional, for development)
- Make

### Quick Start
1. Copy the environment variables template:
   ```bash
   make env-setup
   ```
2. (Optional) Adjust values in the `.env` file.

3. Start the entire technology stack:
   ```bash
   docker compose up -d
   ```

The application will be available at:
- Frontend: `http://localhost:${PROXY_PORT}` (default 80)
- Grafana: `http://localhost:${PROXY_PORT}/monitoring/`

## 📊 Monitoring

The project includes a built-in monitoring stack:
- **Prometheus:** Metrics collection.
- **Grafana:** Data visualization (metrics and logs).
- **Loki:** Log aggregation.
- **Promtail:** Shipping logs to Loki.

## 🏗️ Infrastructure and CI/CD

### Terraform
The `infra/` directory contains Terraform configuration for deploying infrastructure to Azure, including:
- Resource Group
- Virtual Network (VNet) and subnets
- Virtual Machine (VM)
- Azure Container Registry (ACR)
- Key Vault for secrets management

### GitHub Actions
- **CI Pipeline (`ci.yml`):** Runs on every Pull Request and push to the `master` branch. Performs backend and frontend tests, verifies SQLC generation, and builds Docker images.
- **CD Pipeline (`cd.yml`):** Runs after creating a version tag (e.g., `v1.0.0`). Automatically deploys infrastructure using Terraform and publishes images to the registry.

## 📜 Makefile - Useful Commands

The project uses a `Makefile` to automate repetitive tasks:

| Command | Description |
|-----------|------|
| `make build-app` | Builds all Go services |
| `make build-web` | Builds the React application |
| `make test-app` | Runs backend tests with coverage report |
| `make test-web` | Runs frontend tests |
| `make migrate-up` | Executes database migrations |
| `make sqlc` | Generates Go code from SQL queries |
| `make env-setup` | Creates a `.env` file from default settings |
| `make clean` | Removes build artifacts |

## ⚖️ License

Project distributed under the license contained in the `LICENSE` file.
