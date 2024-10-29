import React, { useState } from "react";
import DriverSelect from "../components/driverSelect"; // Asegúrate de la ruta correcta
import "../styles/RaceWeekendPage.css";

const SessionPredictionForm = ({ sessionType, onSubmit }) => {
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
    if (isFormComplete()) {
      onSubmit(formData);
    } else {
      alert("Please fill out all fields before submitting.");
    }
  };

  // Validación para verificar si todos los campos están completos
  const isFormComplete = () => {
    return formData.p1 && formData.p2 && formData.p3;
  };

  return (
    <form className="race-prediction-form" onSubmit={handleSubmit}>
      <h3>
        {sessionType} Prediction for Session {sessionType}
      </h3>
      <div className="driver-select">
        <label>P1</label>
        <DriverSelect
          onSelect={(driverId) => handleDriverSelect("p1", driverId)}
        />
      </div>
      <div className="driver-select">
        <label>P2</label>
        <DriverSelect
          onSelect={(driverId) => handleDriverSelect("p2", driverId)}
        />
      </div>
      <div className="driver-select">
        <label>P3</label>
        <DriverSelect
          onSelect={(driverId) => handleDriverSelect("p3", driverId)}
        />
      </div>

      {/* Botón de envío habilitado/deshabilitado según el estado de isFormComplete */}
      <button
        type="submit"
        className={isFormComplete() ? "complete" : ""}
        disabled={!isFormComplete()}
      >
        Submit Prediction
      </button>
    </form>
  );
};

export default SessionPredictionForm;
