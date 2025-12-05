import React from "react";

interface ClueLabelProps {
  label: string;
  clue: string;
}

export const ClueLabel: React.FC<ClueLabelProps> = ({ label, clue }) => {
  return (
    <div className="d-flex flex-column align-items-center justify-content-center h-100 p-2 bg-light border">
      <div className="fw-bold text-dark">{label}</div>
      <div className="text-secondary small">{clue}</div>
    </div>
  );
};
