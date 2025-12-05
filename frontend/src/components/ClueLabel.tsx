import React from "react";
import "./ClueLabel.css";

interface ClueLabelProps {
  label: string;
  clue: string;
}

export const ClueLabel: React.FC<ClueLabelProps> = ({ label, clue }) => {
  return (
    <div className="clue-label">
      <div className="clue-label-text">{label}</div>
      <div className="clue-label-clue">{clue}</div>
    </div>
  );
};
