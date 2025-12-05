package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// In-memory storage
var (
	rooms   = make(map[string]*Room)
	roomsMu sync.RWMutex
)

// Word list for the game
var wordList = []string{
	"DISEASE",
	"OLD",
	"SNAIL",
	"VACATION",
	"BICYCLE",
	"AIR",
	"GREEN",
	"SANDWICH",
	"VASE",
	"TEACHER",
	"COSTUME",
	"LIGHT",
	"WOLF",
	"YELLOW",
	"SUGAR",
	"RAT",
	"PEN",
	"HELICOPTER",
	"WOOD",
	"BOOK",
	"BUS",
	"FIRE",
	"BLUE",
	"RED",
	"PALM TREE",
	"RUBBER BAND",
	"YOUNG",
	"SMALL",
	"BOX",
	"COUCH",
	"OCTOPUS",
	"SHOVEL",
	"WAR",
	"TREASURE",
	"GOAT",
	"ORANGE",
	"PLANET",
	"CANE",
	"WHEAT",
	"WHITE",
	"ANGER",
	"SNAKE",
	"COLD",
	"MEAN",
	"STRAWBERRY",
	"DESERT",
	"CHILD",
	"SPACESHIP",
	"ARROW",
	"GREY",
	"TRAILER",
	"AIRPLANE",
	"HOUSE",
	"LENTILS",
	"COOK",
	"SOLDIER",
	"BUTTER",
	"AVOCADO",
	"JUNGLE",
	"FOUNTAIN",
	"DUCK",
	"CAMEL",
	"CAT",
	"SPRING",
	"ARROW",
	"SECURITY",
	"BOAT",
	"DIAMOND",
	"HERO",
	"PAPER",
	"RED",
	"EARTH",
	"JOY",
	"BABY BOTTLE",
	"DAY",
	"HEAVY",
	"BLACK",
	"SALAD",
	"EYE",
	"WATER",
	"HEAD",
	"CHEESE",
	"CASTLE",
	"AIRPORT",
	"MOUNTAIN",
	"BEACH",
	"CAMERA",
	"CAULIFLOWER",
	"TRAIN",
	"TOMATO",
	"WINTER",
	"SCHOOL",
	"SNOW",
	"PEAR",
	"ANKLE",
	"MOUTH",
	"SUMMER",
	"SURPRISE",
	"PRESIDENT",
	"MAP",
	"SOAP",
	"PAINTING",
	"EAR",
	"HORSE",
	"PEAK",
	"DOCTOR",
	"OCEAN",
	"JAIL",
	"TRAVEL",
	"VETERINARIAN",
	"SMART",
	"HISTORY",
	"DISGUST",
	"SUITCASE",
	"QUEEN",
	"FAST",
	"DRAGON",
	"DINOSAUR",
	"UNICORN",
	"PIANO",
	"SLOW",
	"FIREFIGHTER",
	"GUITAR",
	"TOY",
	"ARM",
	"DETECTIVE",
	"ENEMY",
	"KNIGHT",
	"FOOT",
	"HELMET",
	"BROWN",
	"SPIDER",
	"CHICKEN",
	"GLASSES",
	"NURSE",
	"PIRATE",
	"FISH",
	"PIGEON",
	"HONEY",
	"PIG",
	"STONE",
	"CAPE",
	"HAPPINESS",
	"CAKE",
	"MEAL",
	"WEIRD",
	"ICE",
	"DOG",
	"FRENCH FRIES",
	"SADNESS",
	"SHIRT",
	"TALL",
	"MONKEY",
	"FEAR",
	"BALL",
	"MUSHROOM",
	"MOON",
	"MARS",
	"RING",
	"PLATE",
	"SHARK",
	"LAPTOP",
	"CIRCUS",
	"CHOCOLATE",
	"ROBOT",
	"DESSERT",
	"FRIENDS",
	"KING",
	"COW",
	"BAG",
	"LIGHT BULB",
	"WARDROBE",
	"BEAR",
	"MUSTACHE",
	"NIGHT",
	"NATION",
	"MAN",
	"WOMAN",
	"WIND",
	"NICE",
	"HOLE",
	"PARACHUTE",
	"LEG",
	"BANANA",
	"CHEST",
	"HAT",
	"UGLY",
	"TIME",
	"HOT",
	"PRETTY",
	"PINK",
	"RADISH",
	"ZOO",
	"MOTORCYCLE",
	"ROAD",
	"GROUP",
	"LETTER",
	"AUTUMN",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Error definitions
var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrRoomExists       = errors.New("room already exists")
	ErrPlayerExists     = errors.New("player name already taken in this room")
	ErrPlayerNotFound   = errors.New("player not found in this room")
	ErrGameNotStarted   = errors.New("game has not started")
	ErrGameOver         = errors.New("game is already over")
	ErrNotEnoughPlayers = errors.New("need at least 2 players to start")
	ErrNoCard           = errors.New("player does not have a card for this cell")
)

// Helper functions

func getRoom(roomCode string) (*Room, bool) {
	roomsMu.RLock()
	defer roomsMu.RUnlock()
	room, exists := rooms[roomCode]
	return room, exists
}

func shuffleWords() []string {
	shuffled := make([]string, len(wordList))
	copy(shuffled, wordList)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func createCardDeck(gridSize int) []Card {
	cards := make([]Card, 0, gridSize*gridSize)
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			cards = append(cards, Card{Row: row, Column: col})
		}
	}
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

// Room methods

// GetCardsPerPlayer returns how many cards each player should have based on player count
func (r *Room) GetCardsPerPlayer() int {
	if len(r.Players) <= 3 {
		return 2
	}
	return 1
}

func (r *Room) DrawCard(playerName string) bool {
	maxCards := r.GetCardsPerPlayer()
	if len(r.CardDeck) > 0 && len(r.PlayerHands[playerName]) < maxCards {
		card := r.CardDeck[0]
		r.CardDeck = r.CardDeck[1:]
		r.PlayerHands[playerName] = append(r.PlayerHands[playerName], card)
		return true
	}
	return false
}

func (r *Room) RemoveCardFromHand(playerName string, row, col int) bool {
	hand := r.PlayerHands[playerName]
	for i, card := range hand {
		if card.Row == row && card.Column == col {
			r.PlayerHands[playerName] = append(hand[:i], hand[i+1:]...)
			return true
		}
	}
	return false
}

func (r *Room) CheckGameOver() bool {
	for row := 0; row < r.GridSize; row++ {
		for col := 0; col < r.GridSize; col++ {
			cell := r.Grid[row][col]
			if !cell.GuessedCorrectly && cell.DiscardedBy == "" {
				return false
			}
		}
	}
	return true
}

func (r *Room) HasPlayer(playerName string) bool {
	for _, player := range r.Players {
		if player == playerName {
			return true
		}
	}
	return false
}

func (r *Room) HasCard(playerName string, row, col int) bool {
	for _, card := range r.PlayerHands[playerName] {
		if card.Row == row && card.Column == col {
			return true
		}
	}
	return false
}

// Game logic functions

// CreateRoom creates a new room with the given code, grid size, and first player
func CreateRoom(roomCode string, gridSize int, playerName string) (*Room, error) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	if _, exists := rooms[roomCode]; exists {
		return nil, ErrRoomExists
	}

	// Initialize row and column words
	words := shuffleWords()
	rowWords := make([]string, gridSize)
	columnWords := make([]string, gridSize)
	for i := 0; i < gridSize; i++ {
		rowWords[i] = words[i]
		columnWords[i] = words[i+gridSize]
	}

	// Initialize grid
	grid := make([][]Cell, gridSize)
	for row := 0; row < gridSize; row++ {
		grid[row] = make([]Cell, gridSize)
		for col := 0; col < gridSize; col++ {
			grid[row][col] = Cell{
				GuessedCorrectly: false,
				DiscardedBy:      "",
			}
		}
	}

	// Create and shuffle deck
	cardDeck := createCardDeck(gridSize)

	room := &Room{
		RoomCode:    roomCode,
		GridSize:    gridSize,
		Players:     []string{playerName},
		GameStarted: false,
		GameOver:    false,
		RowWords:    rowWords,
		ColumnWords: columnWords,
		Grid:        grid,
		CardDeck:    cardDeck,
		PlayerHands: make(map[string][]Card),
	}

	// Deal cards to first player (2 cards since only 1 player)
	room.PlayerHands[playerName] = []Card{}
	room.DrawCard(playerName)
	room.DrawCard(playerName)

	rooms[roomCode] = room
	return room, nil
}

// JoinRoom adds a player to an existing room and deals them cards
func JoinRoom(roomCode, playerName string) (cardsDealt int, err error) {
	room, exists := getRoom(roomCode)
	if !exists {
		return 0, ErrRoomNotFound
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if room.HasPlayer(playerName) {
		return 0, ErrPlayerExists
	}

	room.Players = append(room.Players, playerName)
	room.PlayerHands[playerName] = []Card{}

	// Deal cards based on player count (including new player)
	// 2 cards if 3 or fewer players, 1 card if 4 or more
	cardsToAdd := room.GetCardsPerPlayer()
	for i := 0; i < cardsToAdd; i++ {
		if room.DrawCard(playerName) {
			cardsDealt++
		}
	}

	return cardsDealt, nil
}

// LeaveRoom removes a player from a room
func LeaveRoom(roomCode, playerName string) error {
	room, exists := getRoom(roomCode)
	if !exists {
		return ErrRoomNotFound
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if !room.HasPlayer(playerName) {
		return ErrPlayerNotFound
	}

	// Remove player from players list
	for i, player := range room.Players {
		if player == playerName {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			break
		}
	}

	// Return player's cards to the deck
	if cards, exists := room.PlayerHands[playerName]; exists {
		room.CardDeck = append(room.CardDeck, cards...)
		delete(room.PlayerHands, playerName)
	}

	return nil
}

// StartGame starts or restarts the game in a room
func StartGame(roomCode string) error {
	room, exists := getRoom(roomCode)
	if !exists {
		return ErrRoomNotFound
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if len(room.Players) < 2 {
		return ErrNotEnoughPlayers
	}

	// Reset the game state for a new game
	// Shuffle new words
	words := shuffleWords()
	for i := 0; i < room.GridSize; i++ {
		room.RowWords[i] = words[i]
		room.ColumnWords[i] = words[i+room.GridSize]
	}

	// Reset grid
	for row := 0; row < room.GridSize; row++ {
		for col := 0; col < room.GridSize; col++ {
			room.Grid[row][col] = Cell{
				GuessedCorrectly: false,
				DiscardedBy:      "",
			}
		}
	}

	// Create new shuffled deck
	room.CardDeck = createCardDeck(room.GridSize)

	// Deal new cards to all players
	for _, player := range room.Players {
		room.PlayerHands[player] = []Card{}
		cardsPerPlayer := room.GetCardsPerPlayer()
		for i := 0; i < cardsPerPlayer; i++ {
			room.DrawCard(player)
		}
	}

	room.GameStarted = true
	room.GameOver = false

	return nil
}

// SubmitGuess processes a guess from a player
func SubmitGuess(roomCode, playerName string, row, col int, correct bool) (gameOver bool, err error) {
	room, exists := getRoom(roomCode)
	if !exists {
		return false, ErrRoomNotFound
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if !room.GameStarted {
		return false, ErrGameNotStarted
	}

	if room.GameOver {
		return false, ErrGameOver
	}

	if !room.HasPlayer(playerName) {
		return false, ErrPlayerNotFound
	}

	if !room.HasCard(playerName, row, col) {
		return false, ErrNoCard
	}

	// Remove card from hand
	room.RemoveCardFromHand(playerName, row, col)

	// Update grid
	if correct {
		room.Grid[row][col].GuessedCorrectly = true
	} else {
		room.Grid[row][col].DiscardedBy = playerName
	}

	// Draw a new card
	room.DrawCard(playerName)

	// Check if game is over
	room.GameOver = room.CheckGameOver()

	return room.GameOver, nil
}

// GetGameState returns the game state from a specific player's perspective
func GetGameState(roomCode, playerName string) (*GameStateResponse, error) {
	room, exists := getRoom(roomCode)
	if !exists {
		return nil, ErrRoomNotFound
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	if !room.HasPlayer(playerName) {
		return nil, ErrPlayerNotFound
	}

	// Build grid response (player-specific view)
	gridResponse := make([][]CellResponse, room.GridSize)
	for row := 0; row < room.GridSize; row++ {
		gridResponse[row] = make([]CellResponse, room.GridSize)
		for col := 0; col < room.GridSize; col++ {
			cell := room.Grid[row][col]
			gridResponse[row][col] = CellResponse{
				GuessedCorrectly: cell.GuessedCorrectly,
				DiscardedByMe:    cell.DiscardedBy == playerName,
			}
		}
	}

	// Get player's cards
	playerCards := room.PlayerHands[playerName]
	if playerCards == nil {
		playerCards = []Card{}
	}

	return &GameStateResponse{
		RoomCode:    room.RoomCode,
		GridSize:    room.GridSize,
		GameStarted: room.GameStarted,
		GameOver:    room.GameOver,
		RowWords:    room.RowWords,
		ColumnWords: room.ColumnWords,
		PlayerCards: playerCards,
		Grid:        gridResponse,
		Players:     room.Players,
	}, nil
}

// ClearRooms clears all rooms (useful for testing)
func ClearRooms() {
	roomsMu.Lock()
	defer roomsMu.Unlock()
	rooms = make(map[string]*Room)
}
