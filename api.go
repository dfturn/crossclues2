package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// CORS middleware to allow cross-origin requests from frontend
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Helper functions for HTTP responses

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// HTTP Handlers

func handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.RoomCode) == "" {
		writeError(w, http.StatusBadRequest, "Room code is required")
		return
	}

	if strings.TrimSpace(req.PlayerName) == "" {
		writeError(w, http.StatusBadRequest, "Player name is required")
		return
	}

	// Default grid size if not specified
	gridSize := req.GridSize
	if gridSize == 0 {
		gridSize = DefaultGridSize
	}
	if gridSize < MinGridSize || gridSize > MaxGridSize {
		writeError(w, http.StatusBadRequest, "Grid size must be between 3 and 7")
		return
	}

	room, err := CreateRoom(req.RoomCode, gridSize, req.PlayerName)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, CreateRoomResponse{
		RoomCode:   room.RoomCode,
		PlayerName: req.PlayerName,
		Message:    "Room created successfully",
	})
}

func handleJoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract room code from URL: /api/rooms/{roomCode}/join
	path := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "join" {
		writeError(w, http.StatusBadRequest, "Invalid URL")
		return
	}
	roomCode := parts[0]

	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.PlayerName) == "" {
		writeError(w, http.StatusBadRequest, "Player name is required")
		return
	}

	cardsDealt, err := JoinRoom(roomCode, req.PlayerName)
	if err != nil {
		if err.Error() == "room not found" {
			writeError(w, http.StatusNotFound, "Room not found")
		} else {
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, JoinRoomResponse{
		RoomCode:   roomCode,
		PlayerName: req.PlayerName,
		CardsDealt: cardsDealt,
		Message:    "Joined room successfully",
	})
}

func handleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract room code from URL: /api/rooms/{roomCode}/leave
	path := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "leave" {
		writeError(w, http.StatusBadRequest, "Invalid URL")
		return
	}
	roomCode := parts[0]

	var req JoinRoomRequest // Reuse JoinRoomRequest since it only needs playerName
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.PlayerName) == "" {
		writeError(w, http.StatusBadRequest, "Player name is required")
		return
	}

	err := LeaveRoom(roomCode, req.PlayerName)
	if err != nil {
		if err == ErrRoomNotFound {
			writeError(w, http.StatusNotFound, "Room not found")
		} else if err == ErrPlayerNotFound {
			writeError(w, http.StatusNotFound, "Player not found in room")
		} else {
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"roomCode": roomCode,
		"message":  "Left room successfully",
	})
}

func handleStartGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract room code from URL: /api/rooms/{roomCode}/start
	path := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "start" {
		writeError(w, http.StatusBadRequest, "Invalid URL")
		return
	}
	roomCode := parts[0]

	err := StartGame(roomCode)
	if err != nil {
		if err.Error() == "room not found" {
			writeError(w, http.StatusNotFound, "Room not found")
		} else {
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, StartGameResponse{
		RoomCode: roomCode,
		Message:  "Game started",
	})
}

func handleGuess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract room code from URL: /api/rooms/{roomCode}/guess
	path := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "guess" {
		writeError(w, http.StatusBadRequest, "Invalid URL")
		return
	}
	roomCode := parts[0]

	var req GuessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation - detailed validation happens in SubmitGuess
	if req.Row < 0 || req.Row >= MaxGridSize || req.Column < 0 || req.Column >= MaxGridSize {
		writeError(w, http.StatusBadRequest, "Invalid row or column")
		return
	}

	gameOver, err := SubmitGuess(roomCode, req.PlayerName, req.Row, req.Column, req.Correct)
	if err != nil {
		switch err.Error() {
		case "room not found", "player not found in this room":
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, GuessResponse{
		RoomCode: roomCode,
		Message:  "Guess recorded",
		GameOver: gameOver,
	})
}

func handleGetState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract room code from URL: /api/rooms/{roomCode}/state
	path := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "state" {
		writeError(w, http.StatusBadRequest, "Invalid URL")
		return
	}
	roomCode := parts[0]

	playerName := r.URL.Query().Get("playerName")
	if playerName == "" {
		writeError(w, http.StatusBadRequest, "playerName query parameter is required")
		return
	}

	state, err := GetGameState(roomCode, playerName)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

// Router
func handleRooms(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/rooms")

	// POST /api/rooms - Create room
	if path == "" || path == "/" {
		handleCreateRoom(w, r)
		return
	}

	// Remove leading slash
	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) == 2 {
		action := parts[1]
		switch action {
		case "join":
			handleJoinRoom(w, r)
		case "leave":
			handleLeaveRoom(w, r)
		case "start":
			handleStartGame(w, r)
		case "guess":
			handleGuess(w, r)
		case "state":
			handleGetState(w, r)
		default:
			writeError(w, http.StatusNotFound, "Endpoint not found")
		}
		return
	}

	writeError(w, http.StatusNotFound, "Endpoint not found")
}
