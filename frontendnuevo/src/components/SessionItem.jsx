import React from "react";
import DateDisplay from "./DateDisplay";
import AuthModal from "./AuthModal"; // Importamos el nuevo componente (si decides usarlo)

const SessionItem = ({
  date,
  month,
  sessionType,
  startTime,
  endTime,
  hasPronostico,
  onPronosticoClick,
  isModalOpen,
  onCloseModal,
  onContinueToLogin,
}) => {
  return (
    <div className="flex items-center p-3 border-b border-gray-100 last:border-b-0">
      <DateDisplay date={date} month={month} />
      <div className="ml-6 flex-grow">
        <div className="font-semibold">{sessionType}</div>
        <div className="text-sm text-gray-600">
          {startTime}
          {endTime ? ` - ${endTime}` : ""}
        </div>
      </div>
      {hasPronostico !== undefined && (
        <button
          onClick={onPronosticoClick} // Llamamos a la función pasada como prop
          className="bg-orange-300 text-white px-4 py-1 rounded-full text-sm font-medium hover:bg-orange-400"
        >
          Completar pronóstico
        </button>
      )}
      {/* Modal para autenticación (inline) */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg shadow-lg max-w-md w-full mx-4 md:mx-8 lg:mx-16">
            {" "}
            {/* Añadimos márgenes laterales */}
            <p className="text-center text-gray-800 mb-4">
              Para poder realizar esa acción es necesario que inicie sesión.
            </p>
            <div className="flex justify-center space-x-4">
              <button
                onClick={onCloseModal}
                className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 focus:outline-none"
              >
                Cancelar
              </button>
              <button
                onClick={onContinueToLogin}
                className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 focus:outline-none"
              >
                Continuar
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default SessionItem;
