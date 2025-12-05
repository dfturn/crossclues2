# Build stage for frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Build stage for backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN go build -o crossclues2

# Production stage
FROM alpine:latest
WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/crossclues2 .

# Copy frontend build to serve as static files
COPY --from=frontend-builder /app/frontend/dist ./static

EXPOSE 8080

CMD ["./crossclues2"]
