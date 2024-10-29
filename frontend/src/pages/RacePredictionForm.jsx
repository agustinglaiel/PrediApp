import React, { useState } from "react";
import { useParams } from "react-router-dom";
import DriverSelect from "../components/driverSelect";
import { createProdeCarrera } from "../api/prodes";
// import { getAllDrivers } from "../api/drivers";
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

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Validación simple: Asegurar que P1, P2, P3, etc. estén seleccionados
    if (
      !formData.p1 ||
      !formData.p2 ||
      !formData.p3 ||
      !formData.p4 ||
      !formData.p5 ||
      !formData.fastestLap ||
      !formData.dnf
    ) {
      alert("Please fill out all data!");
      return;
    }

    const predictionData = { ...formData, sessionId, userId }; // Incluir userId y sessionId

    try {
      const response = await createProdeCarrera(predictionData);
      console.log("Prediction created successfully:", response);
      onSubmit(predictionData);
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
      formData.dnf
    );
  };

  return (
    <form className="race-prediction-form" onSubmit={handleSubmit}>
      <h3>
        {sessionType} Prediction for Session {sessionId}
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
      <div className="driver-select">
        <label>P4</label>
        <DriverSelect
          onSelect={(driverId) => handleDriverSelect("p4", driverId)}
        />
      </div>
      <div className="driver-select">
        <label>P5</label>
        <DriverSelect
          onSelect={(driverId) => handleDriverSelect("p5", driverId)}
        />
      </div>
      <div className="driver-select">
        <label>Fastest Lap</label>
        <DriverSelect
          onSelect={(driverId) => handleDriverSelect("fastestLap", driverId)}
        />
      </div>

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

      <button type="submit" className={isFormComplete() ? "complete" : ""}>
        Submit Prediction
      </button>
    </form>
  );
};

export default RacePredictionForm;
