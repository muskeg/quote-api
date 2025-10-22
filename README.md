# Quotes API

A simple RESTful API for managing "inspirational" quotes, built with Go and the Gin web framework. This API allows you to retrieve, add, and search quotes with persistent storage in a JSON file.

## Features

- Get all available quotes
- Retrieve a random quote
- Fetch a specific quote by ID
- Add new quotes (with automatic ID generation)
- Persistent storage using JSON files

## Prerequisites

- Go 1.16 or higher
- Git

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/quote-api.git
cd quote-api
```

2. Install dependencies:

```bash
go mod download
```

3. Build the application:

```bash
go build
```

## Usage

### Running the Server

Run the server with:

```bash
./quote-api
```

The server will start on port 8080 by default.

### API Endpoints

#### Get All Quotes
```
GET /quotes
```

Example:
```bash
curl http://localhost:8080/quotes
```

Response:
```json
[
  {
    "id": "1",
    "quote": "Oh great, another incident. Must be working as intended."
  },
  {
    "id": "2",
    "quote": "Our uptime is mostly emotional."
  }
]
```

#### Get a Random Quote
```
GET /quote
```

Example:
```bash
curl http://localhost:8080/quote
```

Response:
```json
{
  "id": "2",
  "quote": "Our uptime is mostly emotional."
}
```

#### Get a Specific Quote by ID
```
GET /quote/:id
```

Example:
```bash
curl http://localhost:8080/quote/1
```

Response:
```json
{
  "id": "1",
  "quote": "Oh great, another incident. Must be working as intended."
}
```

#### Add a New Quote
```
POST /quotes
```

Example:
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"quote": "Oh great, another incident. Must be working as intended."}'
```

Response:
```json
{
  "id": "1",
  "quote": "Oh great, another incident. Must be working as intended."
}
```

## Data Storage

Quotes are stored in a quotes.json file in the data directory. If the file doesn't exist, it will be created automatically with an example array when the application starts.


## Project Structure

```
quote-api/
├── .gitignore
├── Dockerfile              # Docker container build instructions
├── README.md               # Project documentation
├── config-default.yaml     # Default application config (copied to data/ on first run)
├── docker-compose.yml      # Docker Compose service definition
├── entry.sh                # Entrypoint script for Docker
├── go.mod                  # Go module definition and dependencies
├── go.sum                  # Go module checksums
├── main.go                 # Application entry point and API implementation
├── quotes-default.json     # Default quotes (copied to data/ on first run)

```

## Dependencies

- [Gin Web Framework](https://github.com/gin-gonic/gin) 
- HTTP web framework - HTTP web framework

## Docker Compose

Run and test on Docker Compose

```
version: "3.8"

services:
  quote-api:
    image: ghcr.io/muskeg/quote-api:main
    container_name: quote-api
    ports:
      - "8080:8080"
    environment:
      GIN_MODE: "release" # Production mode, one of "debug", "release", or "test"
      PORT: "8080" # Port inside the container
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```