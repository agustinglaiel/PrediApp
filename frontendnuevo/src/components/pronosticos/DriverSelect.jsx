// DriverSelect.jsx
import React from "react";

const DriverSelect = ({
  position,
  value,
  onChange,
  disabled,
  drivers = [],
}) => {
  return (
    <div className="mb-4">
      <label className="block text-sm font-medium text-gray-700">
        {position}
      </label>
      <select
        value={value || ""}
        onChange={(e) => onChange(parseInt(e.target.value) || null)}
        disabled={disabled}
        className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
      >
        <option value="">Selecciona un piloto</option>
        {drivers.map((driver) => (
          <option key={driver.id} value={driver.id}>
            {driver.full_name}
          </option>
        ))}
      </select>
    </div>
  );
};

export default DriverSelect;
