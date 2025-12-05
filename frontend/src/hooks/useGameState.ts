import { useState, useEffect, useCallback } from "react";
import { fetchGameState } from "../api/gameApi";
import type { GameState } from "../api/gameApi";

export function useGameState(roomCode: string, playerName: string) {
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  // Refetch function for manual refresh
  const refetch = useCallback(async () => {
    if (!roomCode || !playerName) return;
    try {
      const state = await fetchGameState(roomCode, playerName);
      setGameState(state);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to load game state"
      );
    } finally {
      setLoading(false);
    }
  }, [roomCode, playerName]);

  useEffect(() => {
    if (!roomCode || !playerName) return;

    refetch();

    const interval = setInterval(() => refetch(), 1000);

    // Cleanup interval on unmount
    return () => clearInterval(interval);
  }, [roomCode, playerName, refetch]);

  return { gameState, error, loading, refetch };
}
