Recipe Swapper API

This is a backend service written in Go that allows users to upload recipes and find healthier ingredient swaps to reduce calorie intake. The service is built with a clean, layered architecture suitable for a professional development environment.
Features

    Recipe Upload: Submit new recipes via a RESTful API endpoint.

    Ingredient Parsing: Automatically parses ingredient strings into quantity, unit, and name.

    Healthier Swaps: Integrates with the Spoonacular API to find calorie information and suggest healthier alternatives for common ingredients.

    Database Storage: Persists all recipes and their ingredients in a MySQL database.

    Structured Logging: Provides clear, JSON-formatted logs for monitoring.

    Concurrency: Uses Go routines and channels to process ingredient swaps efficiently.

Project Architecture

The project follows a clean, layered architecture (often called Ports and Adapters) to strictly separate concerns, making the application modular, testable, and maintainable:

    cmd/: The application's entry point (main.go). It handles startup and orchestration.

    internal/: Contains all the core application logic.

        api/: Handles HTTP routing and translates requests into service calls (the "Port").

        config/: Manages configuration loading from environment variables (.env).

        database/: The GORM adapter that connects to and interacts with MySQL.

        models/: Defines the core data structures (structs).

        services/: Contains the business logic (the "Adapter"), including parsing and API calls.

    go.mod: Manages project dependencies.

Setup and Installation

Follow these steps to get the application running locally.
Prerequisites

    Go (version 1.21 or later)

    MySQL Server

    A Spoonacular API Key (available via free registration)

1. Configure Environment Variables

Create a file named .env in the root directory of the project and add your secrets:

# .env
SPOONACULAR_API_KEY="YOUR_SPOONACULAR_API_KEY_HERE"
# Optional: Database connection string is set to a default value if DSN is not provided.
# DB_DSN="recipe_user:password@tcp(127.0.0.1:3306)/recipe_db?charset=utf8mb4&parseTime=True&loc=Local"

2. Set Up the Database

Run the provided SQL script to create the necessary database and user credentials:

mysql -u root -p < setup.sql

3. Install Dependencies & Verify Integrity

Run go mod tidy to download all required libraries and synchronize your module files.

go mod tidy

4. Run Tests

Demonstrate code quality by running the included unit tests.

go test ./...

5. Run the Application

The server will start and be accessible at http://localhost:8080.

go run ./cmd/main.go

API Usage (Client Examples)

Use curl commands to interact with the API and test the business logic.
1. Upload a New Recipe (POST)

This endpoint uploads a recipe, saves it to the database, and returns the original recipe along with healthier swap suggestions calculated by the service.

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

This endpoint retrieves all recipes that have been saved in the database.

Endpoint: GET /recipes

curl -X GET http://localhost:8080/recipes

