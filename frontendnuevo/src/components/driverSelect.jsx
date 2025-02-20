import React, { useState, useEffect } from "react";
import { getAllDrivers } from "../api/drivers";

// Componente para seleccionar un piloto
const DriverSelect = ({ onSelect }) => {
  const [drivers, setDrivers] = useState([]);
  const [selectedDriver, setSelectedDriver] = useState("");

  useEffect(() => {
    const fetchDrivers = async () => {
      try {
        const driversData = await getAllDrivers();
        setDrivers(driversData);
      } catch (error) {
        console.error("Error fetching drivers:", error);
      }
    };

    fetchDrivers();
  }, []);

  const handleChange = (event) => {
    const selectedDriverId = event.target.value;
    setSelectedDriver(selectedDriverId);
    onSelect(selectedDriverId);
  };

  return (
    <select value={selectedDriver} onChange={handleChange}>
      <option value="">Select a driver</option>
      {drivers.map((driver) => (
        <option key={driver.id} value={driver.id}>
          {driver.full_name}
        </option>
      ))}
    </select>
  );
};

export default DriverSelect;
