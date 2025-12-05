import React, { Fragment } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Container, Button, Navbar, Spinner, Alert } from "react-bootstrap";
import { GridButton } from "../components/GridButton";
import { ClueLabel } from "../components/ClueLabel";
import { ActionButton } from "../components/ActionButton";
import { useGameState } from "../hooks/useGameState";
import { postGuess, startGame, leaveRoom } from "../api/gameApi";
import type { Card } from "../api/gameApi";
import "./GameScreen.css";

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
      <Navbar
        bg="dark"
        variant="dark"
        sticky="top"
        className="shadow-sm px-2 px-sm-3"
      >
        <Container fluid className="d-flex flex-wrap align-items-center gap-2">
          <Navbar.Brand href="#" className="fw-bold me-2 me-sm-3">
            CrossClues
          </Navbar.Brand>
          <div className="game-navbar-info text-light flex-grow-1">
            <span>
              <strong>Room:</strong>{" "}
              <code className="room-code">{roomCode}</code>
            </span>
            <span className="d-none d-sm-inline">
              <strong>Player:</strong> {decodeURIComponent(playerName)}
            </span>
            <span>
              <strong>Players:</strong> {players.join(", ")}
            </span>
          </div>
          <div className="game-navbar-buttons">
            <Button
              variant="primary"
              size="sm"
              onClick={handleStartGame}
              className="fw-bold"
            >
              New Game
            </Button>
            <Button
              variant="danger"
              size="sm"
              onClick={handleLeaveRoom}
              className="fw-bold"
            >
              Leave
            </Button>
          </div>
        </Container>
      </Navbar>

      <Container
        fluid
        className="flex-grow-1 d-flex align-items-start justify-content-center py-2 py-sm-3 py-md-4"
      >
        <div className="d-flex flex-column align-items-center w-100">
          {/* Game Over Banner */}
          {gameOver && (
            <Alert variant="success" className="game-alert text-center">
              <Alert.Heading className="fs-5 fs-sm-4">
                🎉 Game Over! 🎉
              </Alert.Heading>
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
            <Alert variant="info" className="game-alert text-center">
              <Alert.Heading className="fs-5 fs-sm-4">
                Waiting for game to start
              </Alert.Heading>
              <p>
                {players.length} player(s) in room. Need at least 2 to start.
              </p>
              {players.length >= 2 && (
                <Button variant="primary" size="sm" onClick={handleStartGame}>
                  Start Game
                </Button>
              )}
            </Alert>
          )}

          <div className="game-grid-container">
            <div
              className="game-grid"
              style={{
                gridTemplateColumns: `1fr repeat(${gridSize}, 1fr)`,
              }}
            >
              {/* Header corner cell */}
              <div className="grid-cell header-cell" />

              {/* Column headers */}
              {columnWords.map((word, idx) => (
                <div key={`col-clue-${idx}`} className="grid-cell header-cell">
                  <ClueLabel label={word} clue={colLabels[idx] || ""} />
                </div>
              ))}

              {/* Grid rows with row headers */}
              {Array.from({ length: gridSize }).map((_, rowIdx) => (
                <Fragment key={`row-${rowIdx}`}>
                  {/* Row header */}
                  <div className="grid-cell header-cell">
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
                        className="grid-cell"
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
          </div>

          {/* Player's cards - only show when game is started and not over */}
          {gameStarted && !gameOver && playerCards.length > 0 && (
            <div className="player-cards-container">
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
