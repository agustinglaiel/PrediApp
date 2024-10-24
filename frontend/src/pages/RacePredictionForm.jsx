import React, { useState } from "react";
import { useParams } from "react-router-dom";
import DriverSelect from "./DriverSelect";
import { createProdeCarrera } from "../api/prodes";
import "../styles/RaceWeekendPage.css";

const RacePredictionForm = ({ sessionType, onSubmit }) => {
  const { sessionId } = useParams(); // Obtenemos sessionId desde la URL
  const [userId] = useState(1); // Temporal, luego lo reemplazarás con el userId real

  const [formData, setFormData] = useState({
    p1: "",
    p2: "",
    p3: "",
    p4: "",
    p5: "",
    fastestLap: "",
    vsc: false,
    sc: false,
    dnf: 0,
  });

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData({
      ...formData,
      [name]: type === "checkbox" ? checked : value,
    });
  };

  const handleDriverSelect = (position, driverId) => {
    setFormData({
      ...formData,
      [position]: driverId,
    });
  };

  // Marcar esta función como async para poder usar await dentro de ella
  const handleSubmit = async (e) => {
    e.preventDefault();

    // Validación simple: Asegurar que P1, P2 y P3 estén seleccionados
    if (
      !formData.p1 ||
      !formData.p2 ||
      !formData.p3 ||
      !formData.p4 ||
      !formData.p5 ||
      !formData.fastestLap ||
      !formData.sc ||
      !formData.vsc ||
      !formData.dnf
    ) {
      alert("Please fill out all data!");
      return;
    }

    const predictionData = { ...formData, sessionId, userId }; // Incluir userId y sessionId

    try {
      // Llamar a la función de creación de prode con los datos del formulario
      const response = await createProdeCarrera(predictionData);
      console.log("Prediction created successfully:", response);
      onSubmit(predictionData); // Puedes actualizar la UI o redirigir al usuario aquí
    } catch (error) {
      console.error("Error creating prediction:", error.message);
    }
  };

  const isFormComplete = () => {
    return (
      formData.p1 &&
      formData.p2 &&
      formData.p3 &&
      formData.p4 &&
      formData.p5 &&
      formData.fastestLap &&
      formData.sc &&
      formData.vsc &&
      formData.dnf
    );
  };

  return (
    <form className="race-prediction-form" onSubmit={handleSubmit}>
      <h3>
        {sessionType} Prediction for Session {sessionId}
      </h3>
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
      <DriverSelect
        label="P4"
        onSelect={(driverId) => handleDriverSelect("p4", driverId)}
      />
      <DriverSelect
        label="P5"
        onSelect={(driverId) => handleDriverSelect("p5", driverId)}
      />
      <DriverSelect
        label="Fastest Lap"
        onSelect={(driverId) => handleDriverSelect("fastestLap", driverId)}
      />

      <div className="checkbox-group">
        <label>
          <input
            type="checkbox"
            name="vsc"
            checked={formData.vsc}
            onChange={handleChange}
          />
          Virtual Safety Car
        </label>
        <label>
          <input
            type="checkbox"
            name="sc"
            checked={formData.sc}
            onChange={handleChange}
          />
          Safety Car
        </label>
      </div>

      <div className="dnf-group">
        <label>DNF (Did Not Finish):</label>
        <input
          type="number"
          name="dnf"
          value={formData.dnf}
          min="0"
          max="20"
          onChange={handleChange}
        />
      </div>

      {/* Cambiar el color del botón según si el formulario está completo o no */}
      <button type="submit" className={isFormComplete() ? "complete" : ""}>
        Submit Prediction
      </button>
    </form>
  );
};

export default RacePredictionForm;
