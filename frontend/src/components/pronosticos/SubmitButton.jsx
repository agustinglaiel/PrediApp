import React from "react";

const SubmitButton = ({
  isDisabled,
  onClick,
  label = "Enviar pronóstico",
  className,
}) => {
  return (
    <button
      onClick={onClick}
      disabled={isDisabled}
      className={`px-4 py-2 rounded-full text-sm font-medium transition-colors duration-200 ${
        isDisabled
          ? "bg-gray-300 text-gray-500 cursor-not-allowed"
          : "bg-indigo-600 text-white hover:bg-indigo-700"
      } ${className || ""}`}
    >
      {label}
    </button>
  );
};

export default SubmitButton;
