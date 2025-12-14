# 🚗 Fleet Tracker System

Real-time vehicle tracking system with geofencing capabilities built with Go, MQTT, PostgreSQL, and RabbitMQ.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![MQTT](https://img.shields.io/badge/MQTT-Mosquitto-3C5280?style=flat&logo=eclipse-mosquitto)](https://mosquitto.org)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.x-FF6600?style=flat&logo=rabbitmq)](https://www.rabbitmq.com)

---

## 📋 Table of Contents

- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [Usage](#-usage)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)
- [Migration Guide](#-migration-guide)
- [Seeding Data](#-seeding-data)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🏗️ Architecture

```
┌─────────────┐
│   Vehicle   │ ──MQTT Publish──▶ ┌──────────────┐
│  (GPS Data) │                   │ MQTT Broker  │
└─────────────┘                   │ (Mosquitto)  │
                                  └──────┬───────┘
                                         │ Subscribe
                                         ▼
                                  ┌──────────────┐
                                  │  Subscriber  │
                                  │   Service    │
                                  └──┬────────┬──┘
                                     │        │
                          Save ◀─────┘        └────▶ Check Geofence
                            │                           │
                            ▼                           │ If Entry
                     ┌─────────────┐                   ▼
                     │ PostgreSQL  │            ┌─────────────┐
                     │  Database   │            │  RabbitMQ   │
                     └──────┬──────┘            │   Queue     │
                            │                   └──────┬──────┘
                            │ Query                    │ Consume
                            ▼                          ▼
                     ┌─────────────┐            ┌─────────────┐
                     │ API Service │            │   Worker    │
                     │   (REST)    │            │   Service   │
                     └─────────────┘            └─────────────┘
```

---

## 🛠️ Tech Stack

### **Backend**
- **Go 1.21+** - Main programming language
- **GORM** - ORM for database operations
- **Gin** - Framework

### **Infrastructure**
- **PostgreSQL 15** - Primary database
- **MQTT (Mosquitto 2)** - Message broker for location data
- **RabbitMQ 3** - Message queue for event processing
- **Docker & Docker Compose** - Containerization

### **Tools**
- **golang-migrate** - Database migration tool (dev only)
- **Paho MQTT** - MQTT client library

---

## 📁 Project Structure

```
tj/
├── .env.example                  # Environment variables
├── docker-compose.yml            # Docker services configuration
├── Makefile                      # Common commands
├── README.md                     # This file
│
├── config/
│   └── config.go                 # Configuration loader
│
├── infra/
|   ├── init-sql/
│        └── 000-init-alternate-query.sql     # Configuration loader
|
├── migrations/
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
│   ├── 000002_rename_ts_unix_to_timestamp.up.sql
│   └── ...                       # Database migrations
│
├── pkg/
│   ├── database/
│   │   ├── db.go                 # GORM connection
│   │   └── migration.go          # Migration helpers
│   ├── model/
│   │   └── models.go             # GORM models
│   ├── mqtt/
│   │   └── init.go               # MQTT client
│   ├── rabbitmq/
│   |   └── rmq.go                # RabbitMQ client
|   └── ...
│
└── services/
    ├── api/                      # REST API service
    │   ├── cmd/
    │   │   └── main.go
    │   ├── internal/
    │   │   ├── controller/
    │   └── Dockerfile
    │
    ├── subscriber/               # MQTT subscriber service
    │   ├── cmd/
    │   │   └── main.go
    │   ├── internal/
    │   │   └── controller/
    │   └── Dockerfile
    │
    ├── worker/                   # Event worker service
    │   ├── cmd/
    │   │   └── main.go
    │   ├── internal/
    │   │   └── controller/
    │   └── Dockerfile
    │
    └── publisher/                # Mock vehicle publisher (for testing)
        ├── cmd/
        │   └── main.go
        ├── internal/
        │   └── controller/
        └── Dockerfile
```

---

## 🚀 Getting Started

### **Prerequisites**

- **Go 1.21+** - [Install Go](https://golang.org/doc/install)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **golang-migrate CLI** (optional) - For running migrations manually

**Windows:**
```bash
choco install golang docker-desktop
```

**macOS:**
```bash
brew install go docker docker-compose
```

**Linux:**
```bash
# Install Go
sudo apt install golang-go

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
```

---

### **Installation**

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/fleet-tracker.git
   cd fleet-tracker
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Install Go dependencies**
   ```bash
   go mod download
   ```

4. **Start all services**
   ```bash
   docker-compose up -d
   ```

5. **Verify services are running**
   ```bash
   docker-compose ps
   ```

   Expected output: (example)
   ```
   NAME                    STATUS              PORTS
   fleet-postgres          Up (healthy)        0.0.0.0:5432->5432/tcp
   fleet-mosquitto         Up                  0.0.0.0:1883->1883/tcp
   fleet-rabbitmq          Up (healthy)        0.0.0.0:5672->5672/tcp, 0.0.0.0:15672->15672/tcp
   fleet-mqtt-ingestor     Up
   fleet-api-service       Up                  0.0.0.0:8080->8080/tcp
   fleet-alert-worker      Up
   fleet-publisher         Up
   ```

---

## 💻 Usage

### **Accessing Services**

| Service | URL | Credentials |
|---------|-----|-------------|
| API Service | http://localhost:8093 | - |
| RabbitMQ Management | http://localhost:15673 | guest / guest|
| PostgreSQL | localhost:5433 | fleetuser / fleetpass |
| MQTT Broker | localhost:1883 | - |

---

## 📚 API Documentation

### **Base URL**
```
http://localhost:8093
```

### **Endpoints**

---

#### **Get Latest Vehicle Location**
```http
GET /vehicles/{vehicle_id}/location
```

**Response:**
```json
{
  "id": 39,
  "vehicle_id": "B1234XYZ",
  "latitude": -6.20867481703723,
  "longitude": 106.84499455217255,
  "timestamp": 1765725257,
  "created_at": "2025-12-14T15:14:17.261121Z"
}
```

---

#### **Get Location History**
```http
GET /vehicles/{vehicle_id}/location/history?start={unix_timestamp}&end={unix_timestamp}&limit={number}
```

**Query Parameters:**
- `start` (optional) - Start timestamp (Unix)
- `end` (optional) - End timestamp (Unix)
- `limit` (optional) - Max records (default: 100, max: 1000)
- `offset` (optional) - Offseting data (default: 0)

**Response:**
```json
[
  {
    "id": 39,
    "vehicle_id": "B1234XYZ",
    "latitude": -6.20867481703723,
    "longitude": 106.84499455217255,
    "timestamp": 1765725257,
    "created_at": "2025-12-14T15:14:17.261121Z"
 }
]
```

**Example:**
```bash
curl "http://localhost:8093/vehicles/B1234XYZ/location/history?limit=10"
```

---

## 🗄️ Database Schema

### **Tables**

#### **vehicle_locations**
Historical GPS location data.

| Column | Type | Description |
|--------|------|-------------|
| id | SERIAL | Primary key |
| vehicle_id | VARCHAR(50) | Foreign key to vehicles |
| latitude | DOUBLE PRECISION | GPS latitude |
| longitude | DOUBLE PRECISION | GPS longitude |
| timestamp | BIGINT | Unix timestamp from device |
| created_at | TIMESTAMP | Record creation time |

**Indexes:**
- `idx_vehicle_locations_vehicle_timestamp` on (vehicle_id, timestamp DESC)
- `idx_timestamp` on (timestamp DESC)

#### **bus_stations**
Circular geofence definitions.

| Column | Type | Description |
|--------|------|-------------|
| id | SERIAL | Primary key |
| name | VARCHAR(100) | Geofence name |
| latitude | DOUBLE PRECISION | Center latitude |
| longitude | DOUBLE PRECISION | Center longitude |
| created_at | TIMESTAMP | Record creation time |

---

## 🔄 Migration Guide

### **Using golang-migrate CLI**

```bash
# Install (Windows)
choco install golang-migrate

# Install (macOS)
brew install golang-migrate

# Install (Linux)
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### **Common Commands**

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq seed_bus_stations
```

### **Data Volumes**

| Environment | Vehicles | Locations/Vehicle | Geofences |
|-------------|----------|-------------------|-----------|
| Development | 10-20 | 100-1000 | 5-10 |
| Staging | 3-5 | 50-100 | 2-5 |
| Production | 0 | 0 | 1-3 (master only) |

---

## 🛠️ Development

### **Run Services Locally**

```bash
# Run subscriber service
go run services/subscriber/cmd/main.go

# Run API service
go run services/api/cmd/main.go

# Run worker service
go run services/worker/cmd/main.go

# Run mock publisher
go run services/publisher/cmd/main.go -vehicle TEST001

# Run migration
cd services/subscriber/cmd

go run migrate/main.go -action up
go run migrate/main.go -action rollback
go run migrate/main.go -action force -version 2
```

### **Development Workflow**

1. **Make changes to code**

2. **Create migration if needed**
   ```bash
   migrate create -ext sql -dir migrations -seq seed_bus_stations
   ```
3. **Rebuild Docker images**
   ```bash
   docker-compose up -d --build
   ```

### **Useful Commands**

```bash
# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f subscriber

# Restart a service
docker-compose restart subscriber

# Stop all services
docker-compose down

# Remove everything (including volumes)
docker-compose down -v
```

---

## 🧪 Testing

**Test API:**
```bash
# Get latest location
curl http://localhost:8093/vehicles/B1234XYZ/location

# Get history
curl "http://localhost:8093/vehicles/B1234XYZ/location/history?limit=10"
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤️ by [anemadhon]