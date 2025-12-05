import React from "react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Container, Form, Button, Row, Col, Card } from "react-bootstrap";
import { createRoom, joinRoom } from "../api/gameApi";

const MIN_GRID_SIZE = 3;
const MAX_GRID_SIZE = 5;
const DEFAULT_GRID_SIZE = 5;

export const RoomCreation: React.FC = () => {
  const [joinRoomCode, setJoinRoomCode] = useState("");
  const [joinPlayerName, setJoinPlayerName] = useState("");
  const [createPlayerName, setCreatePlayerName] = useState("");
  const [gridSize, setGridSize] = useState(DEFAULT_GRID_SIZE);
  const navigate = useNavigate();

  // Keep the document title as set in index.html ("Crossclues").

  const handleCreateRoom = async (e: React.FormEvent) => {
    e.preventDefault();
    if (createPlayerName.trim()) {
      // Generate a random room code
      const newRoomCode = Math.random()
        .toString(36)
        .substring(2, 8)
        .toUpperCase();
      // Call API to create the room on the server
      const res = await createRoom({
        roomCode: newRoomCode,
        playerName: createPlayerName,
        gridSize,
      });
      if (res.success) {
        navigate(
          `/game?room=${newRoomCode}&player=${encodeURIComponent(
            createPlayerName
          )}`
        );
      } else {
        alert(res.message || "Failed to create room");
      }
    }
  };

  const handleJoinRoom = async (e: React.FormEvent) => {
    e.preventDefault();
    if (joinRoomCode.trim() && joinPlayerName.trim()) {
      const res = await joinRoom({
        roomCode: joinRoomCode,
        playerName: joinPlayerName,
      });
      if (res.success) {
        navigate(
          `/game?room=${joinRoomCode}&player=${encodeURIComponent(
            joinPlayerName
          )}`
        );
      } else {
        alert(res.message || "Failed to join room");
      }
    }
  };

  return (
    <div className="min-vh-100 d-flex align-items-center justify-content-center bg-gradient">
      <Container>
        {/* Visible page header with high-contrast background so it's always readable */}
        <Card className="mb-5 shadow-sm">
          <Card.Body className="p-4 bg-light text-dark rounded">
            <div className="text-center">
              <h1 className="display-4 fw-bold mb-2">Crossclues</h1>
              <p className="lead mb-0">
                Crossclues is a cooperative word-association board game. Give
                one-word clues to help your team guess the correct grid cells
                before the deck runs out.
              </p>
            </div>
          </Card.Body>
        </Card>

        <Row className="justify-content-center g-4">
          {/* Join Room Card */}
          <Col md={5}>
            <Card className="shadow-lg border-0 h-100">
              <Card.Body className="p-4">
                <h4 className="fw-bold mb-4 text-center">Join a Room</h4>
                <Form onSubmit={handleJoinRoom}>
                  <Form.Group className="mb-3">
                    <Form.Label className="fw-bold">Your Name</Form.Label>
                    <Form.Control
                      type="text"
                      value={joinPlayerName}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setJoinPlayerName(e.target.value)
                      }
                      placeholder="Enter your name"
                      size="lg"
                    />
                  </Form.Group>
                  <Form.Group className="mb-4">
                    <Form.Label className="fw-bold">Room Code</Form.Label>
                    <Form.Control
                      type="text"
                      value={joinRoomCode}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setJoinRoomCode(e.target.value.toUpperCase())
                      }
                      placeholder="Enter room code"
                      maxLength={6}
                      size="lg"
                      className="text-uppercase"
                    />
                  </Form.Group>
                  <Button
                    type="submit"
                    variant="success"
                    size="lg"
                    className="w-100 fw-bold"
                  >
                    Join Room
                  </Button>
                </Form>
              </Card.Body>
            </Card>
          </Col>

          {/* Create Room Card */}
          <Col md={5}>
            <Card className="shadow-lg border-0 h-100">
              <Card.Body className="p-4">
                <h4 className="fw-bold mb-4 text-center">Create a Room</h4>
                <Form onSubmit={handleCreateRoom}>
                  <Form.Group className="mb-3">
                    <Form.Label className="fw-bold">Your Name</Form.Label>
                    <Form.Control
                      type="text"
                      value={createPlayerName}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setCreatePlayerName(e.target.value)
                      }
                      placeholder="Enter your name"
                      size="lg"
                    />
                  </Form.Group>
                  <Form.Group className="mb-4">
                    <Form.Label className="fw-bold">
                      Grid Size: {gridSize}x{gridSize}
                    </Form.Label>
                    <Form.Range
                      min={MIN_GRID_SIZE}
                      max={MAX_GRID_SIZE}
                      value={gridSize}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setGridSize(Number(e.target.value))
                      }
                    />
                    <div className="d-flex justify-content-between text-muted small">
                      <span>
                        {MIN_GRID_SIZE}x{MIN_GRID_SIZE}
                      </span>
                      <span>
                        {MAX_GRID_SIZE}x{MAX_GRID_SIZE}
                      </span>
                    </div>
                  </Form.Group>
                  <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    className="w-100 fw-bold"
                  >
                    Create Room
                  </Button>
                </Form>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </div>
  );
};
