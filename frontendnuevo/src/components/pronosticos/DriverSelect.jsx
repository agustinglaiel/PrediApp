import React, { useState, useEffect } from "react";
import { getAllDrivers } from "../../api/drivers";

const DriverSelect = ({ position, value, onChange }) => {
  const [drivers, setDrivers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchDrivers = async () => {
      try {
        setLoading(true);
        const driverList = await getAllDrivers();
        setDrivers(driverList);
      } catch (err) {
        setError(`Error cargando pilotos: ${err.message}`);
      } finally {
        setLoading(false);
      }
    };

    fetchDrivers();
  }, []);

  if (loading) return <div>Cargando pilotos...</div>;
  if (error) return <div>{error}</div>;

  return (
    <div className="mb-4">
      <label className="block text-sm font-medium text-gray-700">
        {position}
      </label>
      <select
        value={value || ""}
        onChange={(e) => onChange(parseInt(e.target.value) || null)}
        className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
      >
        <option value="">Selecciona un piloto</option>
        {drivers.map((driver) => (
          <option key={driver.ID} value={driver.ID}>
            {driver.FullName} ({driver.TeamName})
          </option>
        ))}
      </select>
    </div>
  );
};

export default DriverSelect;
