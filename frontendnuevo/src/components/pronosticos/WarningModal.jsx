import React from "react";
import { useNavigate } from "react-router-dom";

const WarningModal = ({ isOpen, onClose }) => {
  const navigate = useNavigate();

  if (!isOpen) return null;

  const handleAccept = () => {
    navigate("/"); // Redirigir al inicio (/)
    onClose(); // Cerrar el modal
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white m-4 p-6 rounded-lg shadow-lg max-w-md w-full">
        <p className="text-center text-gray-800">
          El pronóstico de cada sesión cierra 5 minutos antes del comienzo de la
          misma. Tenlo en cuenta para tus próximos pronósticos!
        </p>
        <div className="mt-4 flex justify-center">
          <button
            onClick={handleAccept}
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
          >
            Aceptar
          </button>
        </div>
      </div>
    </div>
  );
};

export default WarningModal;
