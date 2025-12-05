package main

import "sync"

// Game configuration defaults
const DefaultGridSize = 5
const MinGridSize = 3
const MaxGridSize = 5

// Data Models

type Card struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

type Cell struct {
	GuessedCorrectly bool   `json:"guessedCorrectly"`
	DiscardedBy      string `json:"-"` // Internal: who discarded this cell
}

type CellResponse struct {
	GuessedCorrectly bool `json:"guessedCorrectly"`
	DiscardedByMe    bool `json:"discardedByMe"`
}

type Room struct {
	RoomCode    string            `json:"roomCode"`
	GridSize    int               `json:"gridSize"`
	Players     []string          `json:"players"`
	GameStarted bool              `json:"gameStarted"`
	GameOver    bool              `json:"gameOver"`
	RowWords    []string          `json:"rowWords"`
	ColumnWords []string          `json:"columnWords"`
	Grid        [][]Cell          `json:"-"`
	CardDeck    []Card            `json:"-"`
	PlayerHands map[string][]Card `json:"-"`
	mu          sync.RWMutex
}

// Request/Response types

type CreateRoomRequest struct {
	RoomCode   string `json:"roomCode"`
	GridSize   int    `json:"gridSize"`
	PlayerName string `json:"playerName"`
}

type CreateRoomResponse struct {
	RoomCode   string `json:"roomCode"`
	PlayerName string `json:"playerName"`
	Message    string `json:"message"`
}

type JoinRoomRequest struct {
	PlayerName string `json:"playerName"`
}

type JoinRoomResponse struct {
	RoomCode   string `json:"roomCode"`
	PlayerName string `json:"playerName"`
	CardsDealt int    `json:"cardsDealt"`
	Message    string `json:"message"`
}

type StartGameResponse struct {
	RoomCode string `json:"roomCode"`
	Message  string `json:"message"`
}

type GuessRequest struct {
	PlayerName string `json:"playerName"`
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	Correct    bool   `json:"correct"`
}

type GuessResponse struct {
	RoomCode string `json:"roomCode"`
	Message  string `json:"message"`
	GameOver bool   `json:"gameOver"`
}

type GameStateResponse struct {
	RoomCode       string           `json:"roomCode"`
	GridSize       int              `json:"gridSize"`
	GameStarted    bool             `json:"gameStarted"`
	GameOver       bool             `json:"gameOver"`
	CorrectGuesses int              `json:"correctGuesses"`
	TotalCells     int              `json:"totalCells"`
	RowWords       []string         `json:"rowWords"`
	ColumnWords    []string         `json:"columnWords"`
	PlayerCards    []Card           `json:"playerCards"`
	Grid           [][]CellResponse `json:"grid"`
	Players        []string         `json:"players"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
