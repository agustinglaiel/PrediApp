import React from "react";

export default function YesNoButton({
  label, // Título encima: "Virtual Safety Car" o "Safety Car"
  value, // "yes" o "no"
  onChange, // (newValue) => void
  disabled = false,
}) {
  const handleYes = () => {
    if (!disabled && onChange) {
      onChange("yes");
    }
  };

  const handleNo = () => {
    if (!disabled && onChange) {
      onChange("no");
    }
  };

  return (
    <div className="mb-4 ml-4">
      {/* Etiqueta opcional arriba */}
      {label && (
        <label className="block text-sm font-medium text-black mb-1">
          {label}
        </label>
      )}

      <div className="inline-flex items-center rounded-full border border-gray-300 overflow-hidden">
        {/* Botón "Sí" */}
        <button
          type="button"
          className={`
            px-4 py-2 text-sm focus:outline-none
            ${value === "yes" ? "bg-blue-200" : "bg-white"}
            ${disabled ? "opacity-50 cursor-not-allowed" : ""}
          `}
          onClick={handleYes}
        >
          Sí
        </button>

        {/* Botón "No" */}
        <button
          type="button"
          className={`
            px-4 py-2 text-sm focus:outline-none
            ${value === "no" ? "bg-blue-200" : "bg-white"}
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
