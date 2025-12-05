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
        size="lg"
        className="fw-bold"
        onClick={onClick}
      >
        {label}
      </Button>
      {onDiscard && (
        <Button
          type="button"
          onClick={onDiscard}
          variant="danger"
          className="position-absolute p-0 d-flex align-items-center justify-content-center"
          style={{
            width: "24px",
            height: "24px",
            borderRadius: "50%",
            top: "-8px",
            right: "-8px",
            fontSize: "16px",
            fontWeight: "bold",
            padding: "0 !important",
          }}
        >
          ×
        </Button>
      )}
    </div>
  );
};
