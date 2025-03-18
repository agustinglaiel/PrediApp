// YesNoButton.jsx
import React from "react";

export default function YesNoButton({
  label,
  value, // boolean: true o false
  onChange, // (bool) => void
  disabled = false,
}) {
  const handleYes = () => {
    if (!disabled && onChange) {
      onChange(true); // Al hacer clic en “Sí”, seteamos true
    }
  };

  const handleNo = () => {
    if (!disabled && onChange) {
      onChange(false); // Al hacer clic en “No”, seteamos false
    }
  };

  // Si value === true => resaltar el botón “Sí”; si value === false => resaltar “No”.
  return (
    <div className="mb-4 ml-4">
      {label && (
        <label className="block text-sm font-medium text-black mb-1">
          {label}
        </label>
      )}
      <div className="inline-flex items-center rounded-full border border-gray-300 overflow-hidden">
        <button
          type="button"
          className={`
            px-4 py-2 text-sm focus:outline-none
            ${value === true ? "bg-blue-200" : "bg-white"}
            ${disabled ? "opacity-50 cursor-not-allowed" : ""}
          `}
          onClick={handleYes}
        >
          Sí
        </button>

        <button
          type="button"
          className={`
            px-4 py-2 text-sm focus:outline-none
            ${value === false ? "bg-blue-200" : "bg-white"}
            ${disabled ? "opacity-50 cursor-not-allowed" : ""}
          `}
          onClick={handleNo}
        >
          No
        </button>
      </div>
    </div>
  );
}
