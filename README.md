# migrator ![Build](https://github.com/lukaszbudnik/migrator/workflows/Build/badge.svg) ![Docker](https://github.com/lukaszbudnik/migrator/workflows/Docker%20Image%20CI/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/lukaszbudnik/migrator)](https://goreportcard.com/report/github.com/lukaszbudnik/migrator) [![codecov](https://codecov.io/gh/lukaszbudnik/migrator/branch/main/graph/badge.svg)](https://codecov.io/gh/lukaszbudnik/migrator)

Super fast and lightweight DB migration tool written in Go. migrator outperforms other market-leading DB migration frameworks by a few orders of magnitude when comparing both execution time and memory consumption.

## Table of Contents

- [Key Features](#-key-features)
- [Supported Databases](#-supported-databases)
- [Installation](#-installation)
  - [Docker Installation](#docker-installation)
  - [Manual Installation](#manual-installation)
  - [Dashboard Access](#dashboard-access)
- [Overview](#-overview)
- [Quick Start Guide](#-quick-start-guide)
- [Project Structure](#-project-structure)
- [Configuration](#-configuration)
  - [Configuration File](#configuration-file)
  - [Dashboard Configuration](#dashboard-configuration)
- [License](#-license)
- [Credits](#-credits)

## ✨ Key Features

- **🚀 Ultra Performance**: Orders of magnitude faster than other migration tools
- **☁️ Multi-Cloud Storage**: Read migrations from local disk, AWS S3, or Azure Blob Storage
- **🏢 Multi-Tenant Ready**: Built-in support for multi-schema, multi-tenant SaaS applications
- **📡 GraphQL API**: Modern HTTP GraphQL service with comprehensive query capabilities
- **📊 Observability**: Built-in Prometheus metrics and health checks
- **🐳 Container Native**: Ultra-lightweight 30MB Docker image, perfect for microservices
- **🔧 Legacy Migration**: Sync existing migrations from other frameworks seamlessly
- **🎯 CI/CD Ready**: Easy integration into continuous deployment pipelines

## 🗄️ Supported Databases

### Relational Databases
- **PostgreSQL** 9.6+ (and flavours: Amazon RDS/Aurora, Google CloudSQL)
- **MySQL** 5.7+ (and flavours: MariaDB, TiDB, Percona, Amazon RDS/Aurora, Google CloudSQL)
- **Microsoft SQL Server** 2008+
- **SQLite** (file-based database)

### NoSQL Databases
- **MongoDB** 4.0+

### File Formats
- **Excel** (.xlsx, .xls)
- **CSV** (Comma-Separated Values)
- **JSON** (JavaScript Object Notation)
- **SQL** (Structured Query Language dumps)

## 📦 Installation

### Docker Installation

The official docker image is available on:
- Docker Hub: [lukasz/migrator](https://hub.docker.com/r/lukasz/migrator)
- GitHub Container Registry: [ghcr.io/lukaszbudnik/migrator](https://github.com/lukaszbudnik/migrator/pkgs/container/migrator)

```bash
docker pull lukasz/migrator:latest
```

### Manual Installation

To run migrator manually, you'll need to have Go installed on your system:

1. **Install Go 1.24+**
   Download and install Go from the [official website](https://golang.org/dl/).

2. **Clone the repository**
   ```bash
   git clone https://github.com/xahmd/migrator-enhanced
   cd migrator-enhanced
   ```

3. **Run the application**
   You can run migrator directly using Go:
   ```bash
   go run migrator.go
   ```
   
   Or build and run the binary:
   ```bash
   go build -o migrator migrator.go
   ./migrator
   ```

4. **Configuration**
   By default, migrator looks for a `migrator.yaml` configuration file in the current directory. You can specify a different configuration file using the `-configFile` flag:
   ```bash
   go run migrator.go -configFile=/path/to/your/config.yaml
   ```

### Dashboard Access

The migrator includes a web-based dashboard for easier migration management. Once the application is running, you can access the dashboard at:
- http://localhost:8080/static/ (default port)

The dashboard allows you to:
- Upload migration files in various formats (SQL, Excel, CSV, JSON)
- Select source and target formats
- Trigger migration processes
- View migration results and download reports
## 🎯 Overview

migrator manages and versions all DB changes for you and completely eliminates manual and error-prone administrative tasks. migrator versions can be used for auditing and compliance purposes. migrator not only supports single schemas, but also comes with multi-schema support out of the box, making it an ideal DB migrations solution for multi-tenant SaaS products.

migrator runs as a HTTP GraphQL service and can be easily integrated into existing continuous integration and continuous deployment pipelines. migrator can also sync existing migrations from legacy frameworks making the technology switch even more straightforward.

migrator supports the following multi-tenant databases:

- PostgreSQL and all its flavours
- MySQL and all its flavours
- Microsoft SQL Server
- MongoDB

migrator supports reading DB migrations from:

- local folder (any Docker/Kubernetes deployments)
- AWS S3
- Azure Blob Containers

The official docker image is available on:

- docker hub at: [lukasz/migrator](https://hub.docker.com/r/lukasz/migrator)
- alternative mirror at: [ghcr.io/lukaszbudnik/migrator](https://github.com/lukaszbudnik/migrator/pkgs/container/migrator)

The image is ultra lightweight and has a size of 30MB. Ideal for micro-services deployments!

## 🚀 Quick Start Guide

You can apply your first migrations with migrator in literally a few seconds. There is a ready-to-use docker-compose file which sets up migrator and test databases.

### 1. Get the migrator project

Get the source code:

```bash
git clone https://github.com/xahmd/migrator-enhanced
cd migrator-enhanced
```

### 2. Start migrator and test DB containers

You can either use Docker or run manually:

#### Option A: Using Docker (Recommended for testing)

Start migrator and setup test DB containers using docker-compose:

```bash
docker-compose -f ./test/docker-compose.yaml up
```

docker-compose will start and configure the following services:

1. `migrator` - service using latest official migrator image, listening on port `8181`
2. `migrator-dev` - service built from local branch, listening on port `8282`
3. `postgres` - PostgreSQL service, listening on port `5432`
4. `mysql` - MySQL service, listening on port `3306`
5. `mssql` - MS SQL Server, listening on port `1433`
6. `mongodb` - MongoDB service, listening on port `27017`

#### Option B: Manual Installation

Alternatively, you can run migrator manually:

1. Make sure you have Go 1.24+ installed
2. Configure your database connection in `migrator.yaml`
3. Run the application:
   ```bash
   go run migrator.go
   ```

The application will start on port 8080 by default.

## 📁 Project Structure
```
migrator/
├── converter/          # File conversion logic
├── server/             # Web server implementation
├── static/             # Frontend files (HTML, CSS, JS)
├── uploads/            # Directory for uploaded files
└── migrator.go         # Main application entry point
```

## 🔧 Configuration

The application runs on port 8080 by default. To change the port, modify the configuration in `migrator.go` or use the `port` setting in `migrator.yaml`.

### Configuration File

The `migrator.yaml` file contains the main configuration settings:

```yaml
baseLocation: test/migrations  # Base directory for migration files
driver: postgres               # Database driver (postgres, mysql, mssql, mongodb)
dataSource: "host=localhost user=postgres password=yourpassword dbname=migrator_test port=5432 sslmode=disable"  # Database connection string
singleMigrations:
  - ref                       # Single migration directories
  - config
tenantMigrations:
  - tenants                   # Tenant migration directories
port: 8080                   # HTTP server port
```

### Dashboard Configuration

The web dashboard is served from the `/static/` endpoint and includes:
- File upload functionality for various formats (SQL, Excel, CSV, JSON)
- Format selection for source and target migrations
- Progress tracking for migration operations
- Result display with download options for reports and migrated files

## 📄 License
This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## 🙏 Credits

This project is based on [https://github.com/lukaszbudnik/migrator](https://github.com/lukaszbudnik/migrator)