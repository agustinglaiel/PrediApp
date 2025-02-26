import React from "react";

const WarningModal = ({ isOpen, onClose }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg shadow-lg max-w-md w-full mx-4 md:mx-8 lg:mx-16">
        <p className="text-center text-gray-800 mb-6">
          El pronóstico de cada sesión cierra 5 minutos antes del comienzo de la
          misma. ¡Tenlo en cuenta para tus próximos pronósticos!
        </p>
        <div className="flex justify-center">
          <button
            onClick={onClose}
            className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600 focus:outline-none"
          >
            Volver al inicio
          </button>
        </div>
      </div>
    </div>
  );
};

export default WarningModal;
