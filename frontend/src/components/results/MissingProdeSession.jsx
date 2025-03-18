// frontendnuevo/src/components/results/MissingProdeSession.jsx
import React from "react";
import { useNavigate } from "react-router-dom";

const MissingProdeSession = ({ isOpen, onClose }) => {
  const navigate = useNavigate();

  if (!isOpen) return null;

  const handleAccept = () => {
    onClose();
    navigate("/pronosticos"); // Redirige al inicio (página de pronósticos)
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
      <div className="bg-white rounded-lg shadow-lg p-6 w-80">
        <h2 className="text-lg font-bold text-gray-800 mb-4">Sin pronóstico</h2>
        <p className="text-gray-600 mb-6">
          Usted no tiene un pronóstico para esta sesión.
        </p>
        <div className="flex justify-end">
          <button
            onClick={handleAccept}
            className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg"
          >
            Aceptar
          </button>
        </div>
      </div>
    </div>
  );
};

export default MissingProdeSession;
