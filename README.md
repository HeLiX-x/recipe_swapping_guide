# Recipe Swapper API

This is a backend service written in Go that allows users to upload recipes and find healthier ingredient swaps to reduce calorie intake. The service is built with a clean, layered architecture suitable for a professional development environment.

## Features

* **Recipe Upload**: Submit new recipes via a RESTful API endpoint.
* **Ingredient Parsing**: Automatically parses ingredient strings into quantity, unit, and name.
* **Healthier Swaps**: Integrates with the Spoonacular API to suggest healthier alternatives.
* **Database Storage**: Persists recipes and ingredients in a MySQL database using GORM.
* **Structured Logging**: Provides clear, JSON-formatted logs for monitoring.
* **Concurrency**: Uses goroutines to process ingredient swaps efficiently.
* **Performance Profiling**: Exposes Go's `pprof` endpoints for real-time performance analysis.

## Project Architecture

The project follows a clean, layered architecture to separate concerns:

* `cmd/`: The application's entry point (`main.go`).
* `internal/`: Contains all the core application logic.
    * `api/`: Handles HTTP routing and requests.
    * `config/`: Manages configuration loading from environment variables.
    * `database/`: The GORM adapter for MySQL.
    * `models/`: Defines the core data structures.
    * `services/`: Contains the business logic.
* `go.mod`: Manages project dependencies.

## Setup and Installation

Follow these steps to get the application running locally.

### Prerequisites

* Go (version 1.21 or later)
* MySQL Server
* A Spoonacular API Key

### 1. Configure Environment Variables

This project uses a `.env` file for configuration. A template is provided in `.env.example`.

First, copy the example file:
```bash
cp .env.example .env

2. Set Up the Database

Run the provided SQL script to create the database and user:
mysql -u root -p < setup.sql

3. Install Dependencies
go mod tidy

4. Run Tests
go test ./...

5. Run the Application

The API server will start on http://localhost:8080, and the pprof server will start on http://localhost:6060.
go run ./cmd/main.go

API Usage

1. Upload a New Recipe (POST)

Endpoint: POST /upload
curl -X POST http://localhost:8080/upload \
-H "Content-Type: application/json" \
-d '{
  "title": "Creamy Chicken Pasta",
  "ingredients": [
    {"name": "1 cup whole milk"},
    {"name": "2 tbsp butter"},
    {"name": "1/2 cup sour cream"}
  ],
  "instructions": "Mix ingredients and bake for 30 minutes."
}'

2. Get All Recipes (GET)

Endpoint: GET /recipes
curl -X GET http://localhost:8080/recipes

3. Performance Profiling (pprof)

The application exposes pprof endpoints for performance analysis on a separate port.

    Index: http://localhost:6060/debug/pprof/

    CPU Profile: go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

    Memory Profile: go tool pprof http://localhost:6060/debug/pprof/heap

---

### **Final Step: The `.env.example` File**

Now, just create a **new file** in your project directory named `.env.example` and put the following text inside it:

Environment variables for the Recipe Swapper API

Copy this file to .env and fill in your details

SPOONACULAR_API_KEY="YOUR_SPOONACULAR_API_KEY_HERE"
DB_DSN="recipe_user:password@tcp(127.0.0.1:3306)/recipe_db?charset=utf8mb4&parseTime=True&loc=Local"

Once you've saved both files, commit them to Git, and your project will be perfectly documented and ready to show off.
