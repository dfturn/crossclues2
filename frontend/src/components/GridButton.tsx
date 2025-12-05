import React from "react";

interface GridButtonProps {
  label?: string;
  guessed?: boolean;
  discarded?: boolean;
}

export const GridButton: React.FC<GridButtonProps> = ({
  label,
  guessed = false,
  discarded = false,
}) => {
  return (
    <div
      className={`w-100 h-100 fw-bold d-flex align-items-center justify-content-center rounded border-2 ${
        guessed
          ? "bg-primary border-primary text-white"
          : discarded
          ? "bg-danger border-danger text-white"
          : "bg-light border-secondary text-secondary"
      }`}
      style={{ minHeight: "60px" }}
    >
      {discarded ? "✗" : guessed ? label : null}
    </div>
  );
};
