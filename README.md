# migrator ![Build](https://github.com/lukaszbudnik/migrator/workflows/Build/badge.svg) ![Docker](https://github.com/lukaszbudnik/migrator/workflows/Docker%20Image%20CI/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/lukaszbudnik/migrator)](https://goreportcard.com/report/github.com/lukaszbudnik/migrator) [![codecov](https://codecov.io/gh/lukaszbudnik/migrator/branch/main/graph/badge.svg)](https://codecov.io/gh/lukaszbudnik/migrator)

Super fast and lightweight DB migration tool written in Go. migrator outperforms other market-leading DB migration frameworks by a few orders of magnitude when comparing both execution time and memory consumption.

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

The official docker image is available on:
- Docker Hub: [lukasz/migrator](https://hub.docker.com/r/lukasz/migrator)
- GitHub Container Registry: [ghcr.io/lukaszbudnik/migrator](https://github.com/lukaszbudnik/migrator/pkgs/container/migrator)

```bash
docker pull lukasz/migrator:latest
```

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

### 3. Play around with migrator

Set the port and create your first migration:

```bash
MIGRATOR_PORT=8181
COMMIT_SHA="your-version-here"

# Create new version
curl -s -d @- http://localhost:$MIGRATOR_PORT/v2/service <<EOF | jq
{
  "query": "mutation CreateVersion(\$input: VersionInput!) {
    createVersion(input: \$input) {
      version { id, name }
    }
  }",
  "variables": {
    "input": {
      "versionName": "$COMMIT_SHA"
    }
  }
}
EOF

# Fetch migrator versions
curl -s -d @- http://localhost:$MIGRATOR_PORT/v2/service <<EOF | jq -r ".data.versions"
{
  "query": "
  query Versions {
    versions {
        id,
        name,
        created,
      }
  }",
  "operationName": "Versions"
}
EOF

# Fetch tenants
curl -s -d @- http://localhost:$MIGRATOR_PORT/v2/service <<EOF | jq -r ".data.tenants"
{
  "query": "
  query Tenants {
    tenants {
        name
      }
  }",
  "operationName": "Tenants"
}
EOF

# Create new tenant
TENANT_NAME="newcustomer$RANDOM"
VERSION_NAME="create-tenant-$TENANT_NAME"
curl -s -d @- http://localhost:$MIGRATOR_PORT/v2/service <<EOF | jq -r '.data.createTenant'
{
  "query": "
  mutation CreateTenant(\$input: TenantInput!) {
    createTenant(input: \$input) {
      version {
        id,
        name,
      }
      summary {
        startedAt
        duration
        tenants
        migrationsGrandTotal
        scriptsGrandTotal
      }
    }
  }",
  "operationName": "CreateTenant",
  "variables": {
    "input": {
      "versionName": "$VERSION_NAME",
      "tenantName": "$TENANT_NAME"
    }
  }
}
EOF
```

> **💡 Tip**: For a complete GraphQL schema and production deployment guides, see the [📡 API](#-api) and [📚 Tutorials](#-tutorials) sections below.

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

## 📄 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## 🙏 Credits

This project is based on [https://github.com/lukaszbudnik/migrator](https://github.com/lukaszbudnik/migrator)