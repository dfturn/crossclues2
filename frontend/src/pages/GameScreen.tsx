import React, { Fragment } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Container, Button, Navbar, Spinner, Alert } from "react-bootstrap";
import { GridButton } from "../components/GridButton";
import { ClueLabel } from "../components/ClueLabel";
import { ActionButton } from "../components/ActionButton";
import { useGameState } from "../hooks/useGameState";
import { postGuess, startGame, leaveRoom } from "../api/gameApi";
import type { Card } from "../api/gameApi";

// Helper to create card label from row/column indices (e.g., "A3" for row 0, column 2)
const getCardLabel = (card: Card): string => {
  const rowLabel = String.fromCharCode(65 + card.row); // A, B, C, ...
  const colLabel = String(card.column + 1); // 1, 2, 3, ...
  return `${rowLabel}${colLabel}`;
};

export const GameScreen: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const roomCode = searchParams.get("room") || "";
  const playerName = searchParams.get("player") || "";
  const { gameState, error, loading, refetch } = useGameState(
    roomCode,
    playerName
  );

  const handleLeaveRoom = async () => {
    await leaveRoom({ roomCode, playerName });
    navigate("/");
  };

  const handleStartGame = async () => {
    const res = await startGame(roomCode);
    if (!res.success) {
      alert(res.message || "Failed to start game");
    } else {
      await refetch();
    }
  };

  // Handlers for ActionButtons
  const handleGuess = async (card: Card) => {
    if (!gameState) return;
    try {
      await postGuess({
        roomCode,
        playerName,
        row: card.row,
        column: card.column,
        correct: true,
      });
      await refetch();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to submit guess");
    }
  };

  const handleDiscard = async (card: Card) => {
    if (!gameState) return;
    console.log("Discardig card:", card);
    try {
      await postGuess({
        roomCode,
        playerName,
        row: card.row,
        column: card.column,
        correct: false,
      });
      await refetch();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to discard");
    }
  };

  // Loading state
  if (loading && !gameState) {
    return (
      <div className="min-vh-100 d-flex align-items-center justify-content-center bg-light">
        <Spinner animation="border" role="status">
          <span className="visually-hidden">Loading...</span>
        </Spinner>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="min-vh-100 d-flex align-items-center justify-content-center bg-light">
        <Alert variant="danger">
          <Alert.Heading>Error</Alert.Heading>
          <p>{error}</p>
          <Button variant="outline-danger" onClick={() => navigate("/")}>
            Back to Home
          </Button>
        </Alert>
      </div>
    );
  }

  // Use gameState data
  const gridSize = gameState?.gridSize || 5;
  const rowWords = gameState?.rowWords || [];
  const columnWords = gameState?.columnWords || [];
  const grid = gameState?.grid || [];
  const playerCards = gameState?.playerCards || [];
  const gameStarted = gameState?.gameStarted || false;
  const gameOver = gameState?.gameOver || false;
  const players = gameState?.players || [];

  // Generate row/column labels (A, B, C... and 1, 2, 3...)
  const rowLabels = Array.from({ length: gridSize }, (_, i) =>
    String.fromCharCode(65 + i)
  );
  const colLabels = Array.from({ length: gridSize }, (_, i) => String(i + 1));

  return (
    <div className="min-vh-100 d-flex flex-column bg-light">
      <Navbar bg="dark" variant="dark" sticky="top" className="shadow-sm">
        <Container>
          <Navbar.Brand href="#" className="fw-bold">
            CrossClues
          </Navbar.Brand>
          <Navbar.Text className="me-auto text-light">
            <span className="me-4">
              <strong>Room:</strong>{" "}
              <code className="bg-light text-dark px-2 py-1 rounded user-select-all">
                {roomCode}
              </code>
            </span>
            <span className="me-4">
              <strong>Player:</strong> {decodeURIComponent(playerName)}
            </span>
            <span>
              <strong>Players:</strong> {players.join(", ")}
            </span>
          </Navbar.Text>
          <Button
            variant="primary"
            size="sm"
            onClick={handleStartGame}
            className="fw-bold me-2"
          >
            New Game
          </Button>
          <Button
            variant="danger"
            size="sm"
            onClick={handleLeaveRoom}
            className="fw-bold"
          >
            Leave Room
          </Button>
        </Container>
      </Navbar>

      <Container
        fluid
        className="flex-grow-1 d-flex align-items-center justify-content-center py-5"
      >
        <div className="d-flex flex-column align-items-center w-100">
          {/* Game Over Banner */}
          {gameOver && (
            <Alert variant="success" className="mb-4 text-center">
              <Alert.Heading>🎉 Game Over! 🎉</Alert.Heading>
              <p className="mb-0">
                You guessed{" "}
                <strong>
                  {grid.flat().filter((cell) => cell?.guessedCorrectly).length}
                </strong>{" "}
                out of <strong>{gridSize * gridSize}</strong> tiles correctly!
              </p>
            </Alert>
          )}

          {/* Waiting for game to start */}
          {!gameStarted && !gameOver && (
            <Alert variant="info" className="mb-4 text-center">
              <Alert.Heading>Waiting for game to start</Alert.Heading>
              <p>
                {players.length} player(s) in room. Need at least 2 to start.
              </p>
              {players.length >= 2 && (
                <Button variant="primary" onClick={handleStartGame}>
                  Start Game
                </Button>
              )}
            </Alert>
          )}

          <div
            style={{
              display: "grid",
              gridTemplateColumns: `1fr repeat(${gridSize}, 1fr)`,
              gap: "0.5rem",
              padding: "2rem",
              backgroundColor: "white",
              borderRadius: "0.5rem",
              boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)",
            }}
          >
            {/* Header corner cell */}
            <div
              style={{
                aspectRatio: "1",
                backgroundColor: "#e9ecef",
                border: "1px solid #dee2e6",
              }}
            />

            {/* Column headers */}
            {columnWords.map((word, idx) => (
              <div
                key={`col-clue-${idx}`}
                style={{
                  aspectRatio: "1",
                  backgroundColor: "#e9ecef",
                  border: "1px solid #dee2e6",
                }}
              >
                <ClueLabel label={word} clue={colLabels[idx] || ""} />
              </div>
            ))}

            {/* Grid rows with row headers */}
            {Array.from({ length: gridSize }).map((_, rowIdx) => (
              <Fragment key={`row-${rowIdx}`}>
                {/* Row header */}
                <div
                  style={{
                    aspectRatio: "1",
                    backgroundColor: "#e9ecef",
                    border: "1px solid #dee2e6",
                  }}
                >
                  <ClueLabel
                    label={rowWords[rowIdx] || ""}
                    clue={rowLabels[rowIdx] || ""}
                  />
                </div>

                {/* Grid cells for this row */}
                {Array.from({ length: gridSize }).map((_, colIdx) => {
                  const cell = grid[rowIdx]?.[colIdx];
                  const cellLabel = `${rowLabels[rowIdx]}${colLabels[colIdx]}`;
                  return (
                    <div
                      key={`cell-${rowIdx}-${colIdx}`}
                      style={{ aspectRatio: "1" }}
                    >
                      <GridButton
                        label={cellLabel}
                        guessed={cell?.guessedCorrectly}
                        discarded={cell?.discardedByMe}
                      />
                    </div>
                  );
                })}
              </Fragment>
            ))}
          </div>

          {/* Player's cards - only show when game is started and not over */}
          {gameStarted && !gameOver && playerCards.length > 0 && (
            <div className="d-flex gap-3 mt-4">
              {playerCards.map((card, idx) => (
                <ActionButton
                  key={`card-${idx}`}
                  label={getCardLabel(card)}
                  onClick={() => handleGuess(card)}
                  onDiscard={() => handleDiscard(card)}
                />
              ))}
            </div>
          )}
        </div>
      </Container>
    </div>
  );
};
