package main

import (
	"testing"
)

func TestCreateRoom(t *testing.T) {
	ClearRooms()

	// Test creating a new room with grid size and player
	room, err := CreateRoom("TEST123", 5, "Alice")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if room.RoomCode != "TEST123" {
		t.Errorf("Expected room code TEST123, got %s", room.RoomCode)
	}
	if room.GridSize != 5 {
		t.Errorf("Expected grid size 5, got %d", room.GridSize)
	}
	if room.GameStarted {
		t.Error("Expected GameStarted to be false")
	}
	if len(room.Players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(room.Players))
	}
	if room.Players[0] != "Alice" {
		t.Errorf("Expected first player to be Alice, got %s", room.Players[0])
	}
	// First player should have 2 cards (since there's only 1 player)
	if len(room.PlayerHands["Alice"]) != 2 {
		t.Errorf("Expected Alice to have 2 cards, got %d", len(room.PlayerHands["Alice"]))
	}
	// Verify row and column words are set
	for i := 0; i < room.GridSize; i++ {
		if room.RowWords[i] == "" {
			t.Errorf("Row word %d is empty", i)
		}
		if room.ColumnWords[i] == "" {
			t.Errorf("Column word %d is empty", i)
		}
	}
	// Verify grid is initialized
	if len(room.Grid) != room.GridSize {
		t.Errorf("Expected grid to have %d rows, got %d", room.GridSize, len(room.Grid))
	}

	// Test creating duplicate room
	_, err = CreateRoom("TEST123", 5, "Bob")
	if err != ErrRoomExists {
		t.Errorf("Expected ErrRoomExists, got %v", err)
	}
}

func TestCreateRoomDifferentGridSizes(t *testing.T) {
	ClearRooms()

	// Test with grid size 3
	room3, err := CreateRoom("GRID3", 3, "Player1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if room3.GridSize != 3 {
		t.Errorf("Expected grid size 3, got %d", room3.GridSize)
	}
	if len(room3.RowWords) != 3 {
		t.Errorf("Expected 3 row words, got %d", len(room3.RowWords))
	}
	if len(room3.Grid) != 3 || len(room3.Grid[0]) != 3 {
		t.Error("Expected 3x3 grid")
	}
	if len(room3.CardDeck)+len(room3.PlayerHands["Player1"]) != 9 {
		t.Errorf("Expected 9 total cards (3x3), got %d", len(room3.CardDeck)+len(room3.PlayerHands["Player1"]))
	}

	// Test with grid size 7
	room7, err := CreateRoom("GRID7", 7, "Player2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if room7.GridSize != 7 {
		t.Errorf("Expected grid size 7, got %d", room7.GridSize)
	}
	if len(room7.Grid) != 7 || len(room7.Grid[0]) != 7 {
		t.Error("Expected 7x7 grid")
	}
}

func TestJoinRoom(t *testing.T) {
	ClearRooms()

	// Create a room first
	CreateRoom("JOINTEST", 5, "Alice")

	// Test joining the room - Bob should get 2 cards (2 players total)
	cardsDealt, err := JoinRoom("JOINTEST", "Bob")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cardsDealt != 2 {
		t.Errorf("Expected 2 cards dealt, got %d", cardsDealt)
	}

	// Test joining with same name
	_, err = JoinRoom("JOINTEST", "Alice")
	if err != ErrPlayerExists {
		t.Errorf("Expected ErrPlayerExists, got %v", err)
	}

	// Test joining non-existent room
	_, err = JoinRoom("NONEXISTENT", "Charlie")
	if err != ErrRoomNotFound {
		t.Errorf("Expected ErrRoomNotFound, got %v", err)
	}

	// Test joining with third player - should get 2 cards (3 players = 2 cards each)
	cardsDealt, err = JoinRoom("JOINTEST", "Charlie")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cardsDealt != 2 {
		t.Errorf("Expected 2 cards dealt for 3rd player, got %d", cardsDealt)
	}

	// Verify players
	room, _ := getRoom("JOINTEST")
	if len(room.Players) != 3 {
		t.Errorf("Expected 3 players, got %d", len(room.Players))
	}
}

func TestJoinRoomCardLimit(t *testing.T) {
	ClearRooms()

	// Create room and add 3 more players
	CreateRoom("CARDLIMIT", 5, "P1")
	JoinRoom("CARDLIMIT", "P2")
	JoinRoom("CARDLIMIT", "P3")

	// 4th player should get only 1 card
	cardsDealt, err := JoinRoom("CARDLIMIT", "P4")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cardsDealt != 1 {
		t.Errorf("Expected 1 card dealt for 4th player, got %d", cardsDealt)
	}

	// 5th player should also get only 1 card
	cardsDealt, err = JoinRoom("CARDLIMIT", "P5")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cardsDealt != 1 {
		t.Errorf("Expected 1 card dealt for 5th player, got %d", cardsDealt)
	}
}

func TestStartGame(t *testing.T) {
	ClearRooms()

	// Create room with 2 players
	CreateRoom("STARTTEST", 5, "Alice")
	JoinRoom("STARTTEST", "Bob")

	// Test starting game
	err := StartGame("STARTTEST")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	room, _ := getRoom("STARTTEST")
	if !room.GameStarted {
		t.Error("Expected GameStarted to be true")
	}

	// Test starting game again (should work - restarts the game)
	err = StartGame("STARTTEST")
	if err != nil {
		t.Errorf("Expected no error on restart, got %v", err)
	}

	// Verify game is reset with new cards
	room, _ = getRoom("STARTTEST")
	if !room.GameStarted {
		t.Error("Expected GameStarted to be true after restart")
	}
	if room.GameOver {
		t.Error("Expected GameOver to be false after restart")
	}
}

func TestStartGameNotEnoughPlayers(t *testing.T) {
	ClearRooms()

	// Create room with only 1 player (creator)
	CreateRoom("ONEPLAYERTEST", 5, "Alice")

	err := StartGame("ONEPLAYERTEST")
	if err != ErrNotEnoughPlayers {
		t.Errorf("Expected ErrNotEnoughPlayers, got %v", err)
	}
}

func TestStartGameRoomNotFound(t *testing.T) {
	ClearRooms()

	err := StartGame("NONEXISTENT")
	if err != ErrRoomNotFound {
		t.Errorf("Expected ErrRoomNotFound, got %v", err)
	}
}

func TestSubmitGuess(t *testing.T) {
	ClearRooms()

	// Setup game
	CreateRoom("GUESSTEST", 5, "Alice")
	JoinRoom("GUESSTEST", "Bob")
	StartGame("GUESSTEST")

	room, _ := getRoom("GUESSTEST")

	// Get Alice's first card
	aliceCard := room.PlayerHands["Alice"][0]

	// Submit correct guess
	gameOver, err := SubmitGuess("GUESSTEST", "Alice", aliceCard.Row, aliceCard.Column, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if gameOver {
		t.Error("Expected game not to be over yet")
	}

	// Verify cell is marked as guessed correctly
	room, _ = getRoom("GUESSTEST")
	if !room.Grid[aliceCard.Row][aliceCard.Column].GuessedCorrectly {
		t.Error("Expected cell to be guessed correctly")
	}

	// Verify Alice drew a new card (should still have 2 cards)
	if len(room.PlayerHands["Alice"]) != 2 {
		t.Errorf("Expected Alice to have 2 cards, got %d", len(room.PlayerHands["Alice"]))
	}
}

func TestSubmitGuessDiscard(t *testing.T) {
	ClearRooms()

	// Setup game
	CreateRoom("DISCARDTEST", 5, "Alice")
	JoinRoom("DISCARDTEST", "Bob")
	StartGame("DISCARDTEST")

	room, _ := getRoom("DISCARDTEST")
	aliceCard := room.PlayerHands["Alice"][0]

	// Submit incorrect guess (discard)
	_, err := SubmitGuess("DISCARDTEST", "Alice", aliceCard.Row, aliceCard.Column, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify cell is marked as discarded by Alice
	room, _ = getRoom("DISCARDTEST")
	if room.Grid[aliceCard.Row][aliceCard.Column].DiscardedBy != "Alice" {
		t.Errorf("Expected cell to be discarded by Alice, got %s", room.Grid[aliceCard.Row][aliceCard.Column].DiscardedBy)
	}
}

func TestSubmitGuessErrors(t *testing.T) {
	ClearRooms()

	// Test room not found
	_, err := SubmitGuess("NONEXISTENT", "Alice", 0, 0, true)
	if err != ErrRoomNotFound {
		t.Errorf("Expected ErrRoomNotFound, got %v", err)
	}

	// Create room but don't start game
	CreateRoom("NOTSTARTEDTEST", 5, "Alice")
	JoinRoom("NOTSTARTEDTEST", "Bob")

	_, err = SubmitGuess("NOTSTARTEDTEST", "Alice", 0, 0, true)
	if err != ErrGameNotStarted {
		t.Errorf("Expected ErrGameNotStarted, got %v", err)
	}

	// Start game and test player not found
	StartGame("NOTSTARTEDTEST")
	_, err = SubmitGuess("NOTSTARTEDTEST", "Charlie", 0, 0, true)
	if err != ErrPlayerNotFound {
		t.Errorf("Expected ErrPlayerNotFound, got %v", err)
	}

	// Test player doesn't have card
	_, err = SubmitGuess("NOTSTARTEDTEST", "Alice", 4, 4, true)
	if err != ErrNoCard {
		t.Errorf("Expected ErrNoCard, got %v", err)
	}
}

func TestGetGameState(t *testing.T) {
	ClearRooms()

	// Setup game with grid size 4
	CreateRoom("STATETEST", 4, "Alice")
	JoinRoom("STATETEST", "Bob")
	StartGame("STATETEST")

	// Get Alice's state
	state, err := GetGameState("STATETEST", "Alice")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if state.RoomCode != "STATETEST" {
		t.Errorf("Expected room code STATETEST, got %s", state.RoomCode)
	}
	if state.GridSize != 4 {
		t.Errorf("Expected grid size 4, got %d", state.GridSize)
	}
	if !state.GameStarted {
		t.Error("Expected GameStarted to be true")
	}
	if state.GameOver {
		t.Error("Expected GameOver to be false")
	}
	if len(state.RowWords) != 4 {
		t.Errorf("Expected 4 row words, got %d", len(state.RowWords))
	}
	if len(state.ColumnWords) != 4 {
		t.Errorf("Expected 4 column words, got %d", len(state.ColumnWords))
	}
	if len(state.PlayerCards) != 2 {
		t.Errorf("Expected 2 player cards, got %d", len(state.PlayerCards))
	}
	if len(state.Grid) != 4 {
		t.Errorf("Expected 4 rows in grid, got %d", len(state.Grid))
	}
	if len(state.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(state.Players))
	}
}

func TestGetGameStateErrors(t *testing.T) {
	ClearRooms()

	// Test room not found
	_, err := GetGameState("NONEXISTENT", "Alice")
	if err != ErrRoomNotFound {
		t.Errorf("Expected ErrRoomNotFound, got %v", err)
	}

	// Create room and test player not found
	CreateRoom("STATETEST2", 5, "Alice")

	_, err = GetGameState("STATETEST2", "Bob")
	if err != ErrPlayerNotFound {
		t.Errorf("Expected ErrPlayerNotFound, got %v", err)
	}
}

func TestDiscardedByMeVisibility(t *testing.T) {
	ClearRooms()

	// Setup game
	CreateRoom("VISTEST", 5, "Alice")
	JoinRoom("VISTEST", "Bob")
	StartGame("VISTEST")

	room, _ := getRoom("VISTEST")
	aliceCard := room.PlayerHands["Alice"][0]

	// Alice discards a card
	SubmitGuess("VISTEST", "Alice", aliceCard.Row, aliceCard.Column, false)

	// Get Alice's state - should see discardedByMe = true
	aliceState, _ := GetGameState("VISTEST", "Alice")
	if !aliceState.Grid[aliceCard.Row][aliceCard.Column].DiscardedByMe {
		t.Error("Expected Alice to see discardedByMe = true")
	}

	// Get Bob's state - should see discardedByMe = false
	bobState, _ := GetGameState("VISTEST", "Bob")
	if bobState.Grid[aliceCard.Row][aliceCard.Column].DiscardedByMe {
		t.Error("Expected Bob to see discardedByMe = false")
	}
}

func TestGameOver(t *testing.T) {
	ClearRooms()

	// Setup game with small grid for easier testing
	CreateRoom("GAMEOVERTEST", 3, "Alice")
	JoinRoom("GAMEOVERTEST", "Bob")
	StartGame("GAMEOVERTEST")

	room, _ := getRoom("GAMEOVERTEST")

	// Manually set all cells except one to guessed correctly
	room.mu.Lock()
	for row := 0; row < room.GridSize; row++ {
		for col := 0; col < room.GridSize; col++ {
			if row == 0 && col == 0 {
				continue // Leave this one
			}
			room.Grid[row][col].GuessedCorrectly = true
		}
	}
	// Give Alice a card for [0][0]
	room.PlayerHands["Alice"] = []Card{{Row: 0, Column: 0}}
	room.mu.Unlock()

	// Submit the final guess
	gameOver, err := SubmitGuess("GAMEOVERTEST", "Alice", 0, 0, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !gameOver {
		t.Error("Expected game to be over")
	}

	// Verify game is over
	room, _ = getRoom("GAMEOVERTEST")
	if !room.GameOver {
		t.Error("Expected room.GameOver to be true")
	}

	// Try to submit another guess - should fail
	room.mu.Lock()
	room.PlayerHands["Alice"] = []Card{{Row: 1, Column: 1}}
	room.mu.Unlock()

	_, err = SubmitGuess("GAMEOVERTEST", "Alice", 1, 1, true)
	if err != ErrGameOver {
		t.Errorf("Expected ErrGameOver, got %v", err)
	}
}

func TestRoomMethods(t *testing.T) {
	room := &Room{
		RoomCode:    "METHODTEST",
		GridSize:    5,
		Players:     []string{"Alice", "Bob"},
		PlayerHands: make(map[string][]Card),
		CardDeck:    []Card{{Row: 0, Column: 0}, {Row: 1, Column: 1}, {Row: 2, Column: 2}},
		Grid:        make([][]Cell, 5),
	}
	for i := 0; i < 5; i++ {
		room.Grid[i] = make([]Cell, 5)
	}
	room.PlayerHands["Alice"] = []Card{}
	room.PlayerHands["Bob"] = []Card{}

	// Test GetCardsPerPlayer
	if room.GetCardsPerPlayer() != 2 {
		t.Errorf("Expected 2 cards per player for 2 players, got %d", room.GetCardsPerPlayer())
	}

	// Test DrawCard
	room.DrawCard("Alice")
	if len(room.PlayerHands["Alice"]) != 1 {
		t.Errorf("Expected Alice to have 1 card, got %d", len(room.PlayerHands["Alice"]))
	}
	if len(room.CardDeck) != 2 {
		t.Errorf("Expected 2 cards in deck, got %d", len(room.CardDeck))
	}

	// Test HasPlayer
	if !room.HasPlayer("Alice") {
		t.Error("Expected HasPlayer to return true for Alice")
	}
	if room.HasPlayer("Charlie") {
		t.Error("Expected HasPlayer to return false for Charlie")
	}

	// Test HasCard
	if !room.HasCard("Alice", 0, 0) {
		t.Error("Expected HasCard to return true for Alice's card")
	}
	if room.HasCard("Alice", 1, 1) {
		t.Error("Expected HasCard to return false for card Alice doesn't have")
	}

	// Test RemoveCardFromHand
	removed := room.RemoveCardFromHand("Alice", 0, 0)
	if !removed {
		t.Error("Expected RemoveCardFromHand to return true")
	}
	if len(room.PlayerHands["Alice"]) != 0 {
		t.Errorf("Expected Alice to have 0 cards, got %d", len(room.PlayerHands["Alice"]))
	}

	// Test RemoveCardFromHand for non-existent card
	removed = room.RemoveCardFromHand("Alice", 9, 9)
	if removed {
		t.Error("Expected RemoveCardFromHand to return false")
	}

	// Test CheckGameOver (should return false - not all cells are filled)
	if room.CheckGameOver() {
		t.Error("Expected CheckGameOver to return false")
	}

	// Mark all cells as guessed or discarded
	for row := 0; row < room.GridSize; row++ {
		for col := 0; col < room.GridSize; col++ {
			room.Grid[row][col].GuessedCorrectly = true
		}
	}
	if !room.CheckGameOver() {
		t.Error("Expected CheckGameOver to return true")
	}
}

func TestGetCardsPerPlayer(t *testing.T) {
	// Test with different player counts
	tests := []struct {
		playerCount int
		expected    int
	}{
		{1, 2},
		{2, 2},
		{3, 2},
		{4, 1},
		{5, 1},
		{6, 1},
	}

	for _, tc := range tests {
		room := &Room{Players: make([]string, tc.playerCount)}
		got := room.GetCardsPerPlayer()
		if got != tc.expected {
			t.Errorf("GetCardsPerPlayer() with %d players = %d, want %d", tc.playerCount, got, tc.expected)
		}
	}
}

func TestDrawCardLimit(t *testing.T) {
	room := &Room{
		RoomCode:    "DRAWLIMIT",
		GridSize:    5,
		Players:     []string{"Alice"},
		PlayerHands: make(map[string][]Card),
		CardDeck:    []Card{{Row: 0, Column: 0}, {Row: 1, Column: 1}, {Row: 2, Column: 2}},
	}
	room.PlayerHands["Alice"] = []Card{{Row: 3, Column: 3}, {Row: 4, Column: 4}}

	// Alice already has 2 cards, should not draw more (1 player = 2 cards max)
	drawn := room.DrawCard("Alice")
	if drawn {
		t.Error("Expected DrawCard to return false when at card limit")
	}
	if len(room.PlayerHands["Alice"]) != 2 {
		t.Errorf("Expected Alice to still have 2 cards, got %d", len(room.PlayerHands["Alice"]))
	}
	if len(room.CardDeck) != 3 {
		t.Errorf("Expected deck to still have 3 cards, got %d", len(room.CardDeck))
	}
}

func TestDrawCardEmptyDeck(t *testing.T) {
	room := &Room{
		RoomCode:    "EMPTYDECK",
		GridSize:    5,
		Players:     []string{"Alice"},
		PlayerHands: make(map[string][]Card),
		CardDeck:    []Card{},
	}
	room.PlayerHands["Alice"] = []Card{}

	// Deck is empty, should not crash and return false
	drawn := room.DrawCard("Alice")
	if drawn {
		t.Error("Expected DrawCard to return false when deck is empty")
	}
	if len(room.PlayerHands["Alice"]) != 0 {
		t.Errorf("Expected Alice to have 0 cards, got %d", len(room.PlayerHands["Alice"]))
	}
}

func TestDrawCardWith4Players(t *testing.T) {
	room := &Room{
		RoomCode:    "FOURPLAYERS",
		GridSize:    5,
		Players:     []string{"P1", "P2", "P3", "P4"},
		PlayerHands: make(map[string][]Card),
		CardDeck:    []Card{{Row: 0, Column: 0}, {Row: 1, Column: 1}, {Row: 2, Column: 2}},
	}
	room.PlayerHands["P1"] = []Card{{Row: 3, Column: 3}}

	// P1 has 1 card and with 4 players, max is 1, so should not draw more
	drawn := room.DrawCard("P1")
	if drawn {
		t.Error("Expected DrawCard to return false when P1 is at limit with 4 players")
	}
	if len(room.PlayerHands["P1"]) != 1 {
		t.Errorf("Expected P1 to still have 1 card, got %d", len(room.PlayerHands["P1"]))
	}
}
