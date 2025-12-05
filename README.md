# Crossclues

A cooperative word-association board game where players give clues to help their team guess the correct grid cells.

## Game Overview

Crossclues is a multiplayer game where:

- Players join a room and are dealt cards representing grid cells (e.g., "A3" for row A, column 3)
- Each row and column has a word clue displayed on the grid
- Players give one-word clues to help teammates guess which cell they have
- The goal is to correctly guess as many cells as possible as a team

### Game Rules

- **2-3 players**: Each player holds 2 cards
- **4+ players**: Each player holds 1 card
- Players can mark a guess as correct (✓) or discard it (✗)
- The game ends when all cells have been guessed or discarded

## Tech Stack

- **Backend**: Go (standard library HTTP server)
- **Frontend**: React + TypeScript + Vite
- **UI Framework**: React Bootstrap

## Project Structure

```
├── main.go          # HTTP server entry point
├── api.go           # HTTP handlers and routing
├── game.go          # Game logic and room management
├── schema.go        # Data models and types
├── game_test.go     # Unit tests
├── frontend/        # React frontend application
│   ├── src/
│   │   ├── api/         # API client functions
│   │   ├── components/  # Reusable UI components
│   │   ├── hooks/       # Custom React hooks
│   │   └── pages/       # Page components
│   └── ...
└── README.md
```

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- npm

### Backend

```bash
# Build the server
go build

# Run the server (listens on port 8080)
./crossclues2

# Run tests
go test -v
```

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start development server (port 5173)
npm run dev

# Build for production
npm run build
```

### Running Both

1. Start the backend server:

   ```bash
   ./crossclues2
   ```

2. In another terminal, start the frontend:

   ```bash
   cd frontend && npm run dev
   ```

3. Open http://localhost:5173 in your browser

## API Endpoints

| Method | Endpoint                               | Description            |
| ------ | -------------------------------------- | ---------------------- |
| POST   | `/api/rooms`                           | Create a new room      |
| POST   | `/api/rooms/{code}/join`               | Join an existing room  |
| POST   | `/api/rooms/{code}/leave`              | Leave a room           |
| POST   | `/api/rooms/{code}/start`              | Start/restart the game |
| POST   | `/api/rooms/{code}/guess`              | Submit a guess         |
| GET    | `/api/rooms/{code}/state?playerName=X` | Get game state         |

## Docker Deployment

Build and run the containerized application:

```bash
# Build the Docker image
docker build -t crossclues .

# Run the container
docker run -p 8080:8080 crossclues
```

Then open http://localhost:8080 in your browser.

## Cloud Deployment (Google Cloud)

This project is deployed automatically via a Google Cloud Build trigger on pushes to the `main` branch. The build produces a container image and deploys the service to Cloud Run in `us-west1`.

- Deployment trigger: pushes to `main` branch (Cloud Build)
- Cloud Run service: `crossclues2` (region: `us-west1`)

You can view the deployed service and metrics in the Google Cloud Console:

https://console.cloud.google.com/run/detail/us-west1/crossclues2/observability/metrics?project=crossclues

If you need to deploy manually (for debugging or one-off deploys), you can build and push the image yourself and deploy to Cloud Run using gcloud. Example:

```bash
# Build and push (use a project/registry appropriate for your org)
docker build -t gcr.io/<PROJECT_ID>/crossclues2:latest .
docker push gcr.io/<PROJECT_ID>/crossclues2:latest

# Deploy to Cloud Run (replace <PROJECT_ID> and --region as needed)
gcloud run deploy crossclues2 \
   --image gcr.io/<PROJECT_ID>/crossclues2:latest \
   --region us-west1 \
   --platform managed \
   --allow-unauthenticated
```

Notes:

- The automated Cloud Build trigger and the Cloud Run service are managed in your Google Cloud project; if you need the trigger configuration or IAM changes, update them in the Cloud Console or via Infrastructure-as-Code in your repo.
- If Cloud Build fails due to cached dependencies (for example npm EEXIST errors), try clearing the Cloud Build cache or running the build with `--no-cache`, or switch to a non-Alpine Node base image in the Dockerfile as described in the development notes above.

## Configuration

- **Grid Size**: Configurable from 3x3 to 5x5 when creating a room
- **Backend Port**: 8080 (hardcoded in main.go)
- **Frontend Port**: 5173 (Vite default, development only)
