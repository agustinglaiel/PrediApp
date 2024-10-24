import React from "react";

const DriverSelect = ({ drivers = [], label, onSelect }) => {
  const handleChange = (e) => {
    const driverId = e.target.value;
    onSelect(driverId);
  };

  return (
    <div className="driver-select">
      <label>{label}</label>
      <select onChange={handleChange}>
        <option value="">Select a driver</option>
        {drivers.map((driver) => (
          <option key={driver.id} value={driver.id}>
            {driver.name}
          </option>
        ))}
      </select>
    </div>
  );
};

export default DriverSelect;
