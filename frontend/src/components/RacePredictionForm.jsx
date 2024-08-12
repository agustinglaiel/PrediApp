import React, { useState } from "react";
import DriverSelect from "./DriverSelect";
import "../styles/RaceWeekendPage.css";

const RacePredictionForm = ({ sessionType, onSubmit }) => {
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

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form className="race-prediction-form" onSubmit={handleSubmit}>
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

      <button type="submit">Submit Prediction</button>
    </form>
  );
};

export default RacePredictionForm;
