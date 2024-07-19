# Movie-booking-Golang
Movie_booking_system_backend

This is a backend application built using Go with the Gin framework and MongoDB Atlas as the database. It provides APIs for managing movies, showtimes, and booking seats.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Usage](#usage)
5. [Endpoints](#endpoints)
6. [License](#license)

## Prerequisites

Before running the application, ensure you have the following installed:

- Go programming language (version 1.15 or higher)
- MongoDB Atlas account and a configured cluster
- Git (optional, for cloning the repository)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/movie-booking-system.git

2. Navigate to the project directory

   ```bash
   cd movie-booking-system

3. Install Dependencies

   ```bash
   go mod tidy

## Configuration

1. Mongodb atlas setup

- Create a MongoDB Atlas cluster.
- Whitelist your IP address in MongoDB Atlas to allow access.
- Obtain your MongoDB connection string.
- Environment Variables:

2. Create a .env file in the root directory of your project.

- Add the MongoDB connection string and any other sensitive information:

  ```bash
  MONGO_URI=mongodb+srv://<username>:<password>@<clustername>.mongodb.net/<dbname>?retryWrites=true&w=majority

Replace username, password, clustername, and <dbname> with your actual MongoDB Atlas credentials and cluster details.

## Usage

1. Run the application

  ```bash
 go run main.go


2. The server will start at http://localhost:8080 by default.
  ```

## Endpoints

- GET /movies: Retrieve all movies.

- GET /movies/:id: Retrieve a specific movie by ID.

- POST /movies: Add a new movie.

- PUT /movies/:id: Update a movie by ID.

- DELETE /movies/:id: Delete a movie by ID.

- POST /movies/:id/book: Book a seat for a specific movie.

Replace :id with the actual ID of the movie.


### License

This project is licensed under the MIT License - see the LICENSE file for details.


### Explanation:

- **Markdown Formatting**: Uses Markdown syntax to structure the document with headers, lists, code blocks, and links.
- **Instructions**: Each section (Prerequisites, Installation, Configuration, Usage, Endpoints, Testing, Deployment, License) corresponds to a specific aspect of setting up, running, and managing your application.
- **Environment Variables**: Instructions for setting up environment variables, including sensitive information like the MongoDB connection string.
- **Endpoints**: Describes available API endpoints with HTTP methods and URL paths.
- **License**: Mentions the licensing information for the project.





 

