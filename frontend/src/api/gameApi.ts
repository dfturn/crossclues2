// Game API - communicates with Go backend

// Use relative URL for production (served from same origin) or localhost for development
const API_BASE = import.meta.env.DEV ? "http://localhost:8080/api" : "/api";

// Types matching Go backend schema

export interface Card {
  row: number;
  column: number;
}

export interface CellResponse {
  guessedCorrectly: boolean;
  discardedByMe: boolean;
}

export interface GameState {
  roomCode: string;
  gridSize: number;
  gameStarted: boolean;
  gameOver: boolean;
  rowWords: string[];
  columnWords: string[];
  playerCards: Card[];
  grid: CellResponse[][];
  players: string[];
}

export interface CreateRoomResponse {
  roomCode: string;
  playerName: string;
  message: string;
}

export interface JoinRoomResponse {
  roomCode: string;
  playerName: string;
  cardsDealt: number;
  message: string;
}

export interface GuessResponse {
  roomCode: string;
  message: string;
  gameOver: boolean;
}

export interface ErrorResponse {
  error: string;
}

// --- Fetch game state from backend ---
export async function fetchGameState(
  roomCode: string,
  playerName: string
): Promise<GameState> {
  const response = await fetch(
    `${API_BASE}/rooms/${roomCode}/state?playerName=${encodeURIComponent(
      playerName
    )}`
  );
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json();
    throw new Error(errorData.error || "Failed to fetch game state");
  }
  return response.json();
}

// --- Submit a guess (correct or discard) ---
export async function postGuess(payload: {
  roomCode: string;
  playerName: string;
  row: number;
  column: number;
  correct: boolean;
}): Promise<GuessResponse> {
  const response = await fetch(`${API_BASE}/rooms/${payload.roomCode}/guess`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      playerName: payload.playerName,
      row: payload.row,
      column: payload.column,
      correct: payload.correct,
    }),
  });
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json();
    throw new Error(errorData.error || "Failed to submit guess");
  }
  return response.json();
}

// --- Create room API ---
export async function createRoom(payload: {
  roomCode: string;
  playerName: string;
  gridSize?: number;
}): Promise<{ success: boolean; message: string }> {
  const response = await fetch(`${API_BASE}/rooms`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      roomCode: payload.roomCode,
      playerName: payload.playerName,
      gridSize: payload.gridSize || 5,
    }),
  });
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json();
    return { success: false, message: errorData.error };
  }
  const data: CreateRoomResponse = await response.json();
  return { success: true, message: data.message };
}

// --- Join room API ---
export async function joinRoom(payload: {
  roomCode: string;
  playerName: string;
}): Promise<{ success: boolean; message: string; cardsDealt?: number }> {
  const response = await fetch(`${API_BASE}/rooms/${payload.roomCode}/join`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ playerName: payload.playerName }),
  });
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json();
    return { success: false, message: errorData.error };
  }
  const data: JoinRoomResponse = await response.json();
  return { success: true, message: data.message, cardsDealt: data.cardsDealt };
}

// --- Leave room API ---
export async function leaveRoom(payload: {
  roomCode: string;
  playerName: string;
}): Promise<{ success: boolean; message: string }> {
  const response = await fetch(`${API_BASE}/rooms/${payload.roomCode}/leave`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ playerName: payload.playerName }),
  });
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json();
    return { success: false, message: errorData.error };
  }
  const data = await response.json();
  return { success: true, message: data.message };
}

// --- Start game API ---
export async function startGame(
  roomCode: string
): Promise<{ success: boolean; message: string }> {
  const response = await fetch(`${API_BASE}/rooms/${roomCode}/start`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  });
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json();
    return { success: false, message: errorData.error };
  }
  const data = await response.json();
  return { success: true, message: data.message };
}
