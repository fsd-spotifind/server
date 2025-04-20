# Spotifind Backend

A Go-based backend service for the Spotifind application, providing API endpoints and database functionality

GitHub: [link](https://github.com/fsd-spotifind/server

## Main Functionalities

### User's Spotify Stats

The backend's Spotify integration computes user's stats through the `Recently Played` endpoint, based on user's top tracks, artists and albums. The stats are computed over the following periods - weekly, monthly, semi-annual, annual with a cron job scheduled to run at fixed times.

### Song of the Day (SOTD)

The backend's Spotify integration gets user's recently played tracks and uses them as recommendations for the Song of the Day. User can then choose one of the tracks as his or her Song of the Day, with short note and mood description attached to it.

### Befriend

The backend's `friend` endpoint allows users to befriend each other by clicking on the `add friend` button on the user's profile page. The backend will store the friendship in the database as well as the status of the friendship (pending, active, removed).

## Project Structure

```
├── cmd/            # Application entrypoints
│   └── api/        # Main API server
├── internal/       # Private application code
└──  prisma/        # Database schema and migrations
```

### Internal

```
├── database/     # Database operations
├── logger/       # For logging
├── models/       # Database models
├── server/       # Contains the routes, handlers, and cron jobs
├── spotify/      # Spotify API operations
└── utils/        # Common utility functions
```

## Prerequisites

- Go 1.23.2 or higher
- Docker and Docker Compose
- Make

## Getting Started

1. Clone the repository

```bash
git clone https://github.com/fsd-spotifind/server.git
```

2. Copy the environment file

```bash
cp .env.example .env
```

3. Update the `.env` file with the right configurations

## Local Development

### Building and Running

Build the application:

```bash
make build
```

Run the application:

```bash
make run
```

Run the application with live reload:

```bash
make watch
```

### Database

Start the database container:

```bash
make docker-run
```

Stop the database container:

```bash
make docker-down
```

## Links

- [Postman collection for this project](https://interstellar-sunset-265479.postman.co/workspace/Team-Workspace~6b3fa817-fb80-4793-8aa6-211438648900/collection/23441224-71cd4c35-c61d-4dfd-8894-23ee85622817?action=share&creator=23441224)

## References

- [Spotify API](https://developer.spotify.com/documentation/web-api)
- [Prisma Go Client](https://github.com/prisma/prisma-client-go)
