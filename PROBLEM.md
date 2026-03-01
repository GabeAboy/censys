# Asset Management Microservice

## Overview
Our Web App team is building a new dashboard to help users categorize their "Assets" (IP addresses and Services). We need a robust, scalable backend service to manage these assets and their associated risk metadata.

## Task Summary
Build a RESTful Microservice in Go that allows a UI to manage a collection of Assets.

## Core Features

### Asset Registration
* As a UI, I want to POST a new asset (IP Address, Hostname, and a list of Open Ports) to the system.

### Automated Risk Scoring
* As a Security Analyst, when an asset is added, the system should automatically assign a Risk_Level (Low, Medium, High) based on the ports:
  * High: Ports 22 (SSH), 3389 (RDP), or 21 (FTP).
  * Medium: Port 443 (HTTPS) with a specific flag (e.g., "expired cert").
  * Low: All other ports.

### Tagging API
* As a UI, I need to PUT or PATCH custom tags (e.g., "Staging", "Critical") onto an existing asset.

### Search & Filter
* As a UI, I want to GET all assets filtered by a specific Tag or Risk_Level.

## Your Service requirements

### Persistence
* Store assets in a database of your choice explaining the reasoning behind such a choice.

### Contract
* Provide an OpenAPI/Swagger specification or a Postman Collection.

### Bonus points
* Load new random assets into the DB on a schedule to simulate new data ingestion

## Deliverables

### Source Code
* A link to a private GitHub repository or a ZIP file.

### Setup Guide
* Clear instructions on how to build and run the service.

### Must run with
* `docker-compose up`

### Include database migrations as needed

### Test Guide
* Automated test framework, or simple clear manual steps

### Short Summary
* A few bullet points (5-10) on what you would do next if you had another 40 hours to make it production ready