# Censys - Asset Management System

A modern asset management system for tracking network assets, open ports, and risk levels. Built with Go, PostgreSQL, and React.

## Testing Instructions

1. Start the service: `make up`
2. Open a browser and navigate to `http://localhost:3000` and `http://localhost:8080/swagger/index.html#/`
3. Click the "Add Asset" button and enter a new asset information
4. Click the "Search" button and enter the asset information you just entered
5. Verify that the new asset is displayed in the search results
6. Click the "Delete" button next to the new asset
7. Click the "Search" button again and verify that the new asset is no longer displayed in the search results

## Architecture

- **Backend**: Go 1.26 with Gin web framework
- **Database**: PostgreSQL 15
- **Frontend**: React with Tailwind CSS
- **Containerization**: Docker & Docker Compose

## Prerequisites

- [Make]
- [Docker]
- [Go 1.26+]
- [Node.js 18+]

## Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd censys
   ```

2. **Start all services**
   ```bash
   make up
   ```
   
   Or without Make:
   ```bash
   docker compose up -d
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - API: http://localhost:8080

4. **View logs**
   ```bash
   make logs          # All services
   ```

5. **Stop all services**
   ```bash
   make down
   ```

# Install dependencies
go mod download

# Run the server

go run main.go

```
The API will start on http://localhost:8080

#### 3. Frontend Setup

```bash
cd web

# Install dependencies
npm install

# Start development server
npm start
```

The web app will start on http://localhost:3000

## 🗄️ Database

### Schema

The database automatically initializes with the following tables:

- **`assets`** - Network assets (IP, hostname, risk level)
- **`ports`** - Open ports for each asset
- **`tags`** - Tags for categorizing assets
- **`asset_tags`** - Many-to-many relationship between assets and tags
