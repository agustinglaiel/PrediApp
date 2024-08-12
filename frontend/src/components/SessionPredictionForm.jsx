import React, { useState } from "react";
import DriverSelect from "./DriverSelect"; // Importa el componente DriverSelect
import "../styles/RaceWeekendPage.css";

const QualyPredictionForm = ({ sessionType, onSubmit }) => {
  const [formData, setFormData] = useState({
    p1: "",
    p2: "",
    p3: "",
  });

  const handleDriverSelect = (position, driverId) => {
    setFormData({
      ...formData,
      [position]: driverId,
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form className="qualy-prediction-form" onSubmit={handleSubmit}>
      <h3>{sessionType} Prediction</h3>
      <DriverSelect
        label="P1"
        onSelect={(driverId) => handleDriverSelect("p1", driverId)}
      />
      <DriverSelect
        label="P2"
        onSelect={(driverId) => handleDriverSelect("p2", driverId)}
      />
      <DriverSelect
        label="P3"
        onSelect={(driverId) => handleDriverSelect("p3", driverId)}
      />

      <button type="submit">Submit Prediction</button>
    </form>
  );
};

export default QualyPredictionForm;
