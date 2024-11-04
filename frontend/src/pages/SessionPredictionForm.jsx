import React, { useState } from "react";
import DriverSelect from "../components/driverSelect"; // Asegúrate de la ruta correcta
import "../styles/RaceWeekendPage.css";
import { createProdeSession } from "../api/prodes"; // Asegúrate de importar la función correcta

const SessionPredictionForm = () => {
  const sessionId = 3; // Valor definido para sessionId
  const [formData, setFormData] = useState({
    p1: null,
    p2: null,
    p3: null,
  });

  const handleDriverSelect = (position, driverId) => {
    setFormData((prevData) => ({
      ...prevData,
      [position]: driverId,
    }));
  };

  const isFormComplete = () => {
    return formData.p1 && formData.p2 && formData.p3;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isFormComplete()) {
      try {
        // Crear el payload con los datos de la predicción y el sessionId
        const prodeData = {
          ...formData,
          session_id: sessionId,
        };

        // Llamar a la función para crear el prode de sesión
        const response = await createProdeSession(prodeData);
        alert("Prediction submitted successfully!");
        console.log(response); // Puedes manejar la respuesta como desees
      } catch (error) {
        console.error("Error submitting prediction:", error);
        alert(
          error.response?.data?.message ||
            "An error occurred while submitting your prediction."
        );
      }
    } else {
      alert("Please fill out all fields before submitting.");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h3>Submit Your Prediction</h3>
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
