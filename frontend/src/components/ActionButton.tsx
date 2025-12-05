import React from "react";
import { Button } from "react-bootstrap";

interface ActionButtonProps {
  label: string;
  onClick?: () => void;
  onDiscard?: () => void;
}

export const ActionButton: React.FC<ActionButtonProps> = ({
  label,
  onClick,
  onDiscard,
}) => {
  return (
    <div className="position-relative d-inline-block">
      <Button
        type="button"
        variant="primary"
        className="fw-bold action-card-btn"
        onClick={onClick}
        style={{
          fontSize: "clamp(1.25rem, 5vw, 1.75rem)",
          padding: "clamp(0.75rem, 3vw, 1.25rem) clamp(1.5rem, 6vw, 2.5rem)",
          minWidth: "clamp(80px, 20vw, 120px)",
          minHeight: "clamp(50px, 12vw, 70px)",
        }}
      >
        {label}
      </Button>
      {onDiscard && (
        <Button
          type="button"
          onClick={onDiscard}
          variant="danger"
          className="position-absolute p-0 d-flex align-items-center justify-content-center action-discard-btn"
          style={{
            width: "clamp(32px, 8vw, 40px)",
            height: "clamp(32px, 8vw, 40px)",
            borderRadius: "50%",
            top: "-8px",
            right: "-8px",
            fontSize: "clamp(20px, 5vw, 28px)",
            fontWeight: "bold",
            zIndex: 20,
          }}
        >
          ×
        </Button>
      )}
    </div>
  );
};
