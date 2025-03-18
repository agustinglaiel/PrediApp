// frontendnuevo/src/components/results/DriverResultDisplay.jsx
import React from "react";

const DriverResultDisplay = ({ position, driverName }) => {
  return (
    <div className="mb-2 ml-4">
      {/* Reducido mb-4 a mb-2 */}
      <label className="block text-sm font-medium text-black">{position}</label>
      <div className="mt-1 block w-full py-1 px-3 text-gray-700 bg-gray-100 border border-gray-300 rounded-md">
        {/* Reducido py-2 a py-1 */}
        {driverName || "No seleccionado"}
      </div>
    </div>
  );
};

export default DriverResultDisplay;
