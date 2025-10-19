# Quotes API


A simple RESTful API for managing quotes built with Go and Gin. This service allows you to retrieve, add, and manage inspirational quotes.

## Features

- Get all quotes
- Retrieve a random quote
- Fetch a specific quote by ID
- Add new quotes
- Persistent storage using JSON file

## Installation

1. Clone the repository:
    ```
    git clone https://github.com/yourusername/quote-api.git
    cd quote-api
    ```

2. Install dependencies:
    ```
    go mod tidy
    ```

3. Run the application:
    ```
    go run main.go
    ```

The server will start on port 8080.

## API Endpoints

### GET /quotes
Returns all available quotes.

**Response:**
```json
[
  {
     "id": "1",
     "quote": "The only way to do great work is to love what you do."
  },
  {
     "id": "2",
     "quote": "Life is what happens when you're busy making other plans."
  }
]
```

### GET /quote
Returns a random quote.

**Response:**
```json
{
  "id": "1",
  "quote": "The only way to do great work is to love what you do."
}
```

### GET /quote/:id
Returns a specific quote by ID.

**Response:**
```json
{
  "id": "2",
  "quote": "Life is what happens when you're busy making other plans."
}
```

### POST /quotes
Adds a new quote.

**Request:**
```json
{
  "quote": "Be yourself; everyone else is already taken."
}
```

**Response:**
```json
{
  "id": "3",
  "quote": "Be yourself; everyone else is already taken."
}
```

## Data Persistence

Quotes are stored in a `quotes.json` file and loaded when the application starts. When you add new quotes, they're automatically saved to this file.

## Example Usage

Add a new quote:
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"quote": "The only way to do great work is to love what you do."}'
```

Get all quotes:
```bash
curl http://localhost:8080/quotes
```

Get a random quote:
```bash
curl http://localhost:8080/quote
```

Get a quote by ID:
```bash
curl http://localhost:8080/quote/1
```

## Project Structure
```
quote-api/
├── main.go      # Application entry point and API implementation
├── quotes.json  # Persistent storage for quotes
├── go.mod       # Go module definition
├── go.sum       # Go module checksum
└── README.md    # Project documentation
```

## Dependencies
# Quotes API

A simple RESTful API for managing inspirational quotes, built with Go and the Gin web framework. This API allows you to retrieve, add, and search quotes with persistent storage in a JSON file.

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
    "quote": "The only way to do great work is to love what you do."
  },
  {
    "id": "2",
    "quote": "Life is what happens when you're busy making other plans."
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
  "quote": "Life is what happens when you're busy making other plans."
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
  "quote": "The only way to do great work is to love what you do."
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
  -d '{"quote": "The only way to do great work is to love what you do."}'
```

Response:
```json
{
  "id": "1",
  "quote": "The only way to do great work is to love what you do."
}
```

## Data Storage

Quotes are stored in a quotes.json file in the same directory as the application. If the file doesn't exist, it will be created automatically with an empty array when the application starts.

## Project Structure

```
quote-api/
├── main.go      # Application entry point and API implementation
├── quotes.json  # Persistent storage for quotes
├── go.mod       # Go module definition
├── go.sum       # Go module checksum
└── README.md    # Project documentation
```

## Dependencies

- [Gin Web Framework](https://github.com/gin-gonic/gin) 
- HTTP web framework - HTTP web framework