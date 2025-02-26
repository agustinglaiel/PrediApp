import React from "react";
import DateDisplay from "./DateDisplay";
import AuthModal from "./AuthModal"; // Importamos el nuevo componente (si decides usarlo)

const SessionItem = ({
  sessionId, // Pasado desde HomePage.jsx
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
  prodeSession,
  prodeRace,
}) => {
  // Determinar si hay un pronóstico (ya sea de sesión o carrera)
  const hasProde =
    (sessionType !== "Race" && prodeSession) ||
    (sessionType === "Race" && prodeRace);

  return (
    <div className="flex items-center p-3 border-b border-gray-100 last:border-b-0">
      <DateDisplay date={date} month={month} />
      <div className="ml-6 flex-grow">
        <div className="font-semibold">{sessionType}</div>
        <div className="text-sm text-gray-600">
          {startTime}
          {endTime ? ` - ${endTime}` : ""}
        </div>
        {prodeRace && (
          <div className="text-xs text-gray-500">
            Pronóstico Carrera: {prodeRace.P1}, {prodeRace.P2}, {prodeRace.P3}
          </div>
        )}
      </div>
      {hasPronostico !== undefined && (
        <button
          onClick={onPronosticoClick}
          className={`px-4 py-1 rounded-full text-sm font-medium transition-colors duration-200 ${
            hasProde
              ? "bg-white text-green-500 border border-green-500 hover:bg-green-50"
              : "bg-orange-300 text-white hover:bg-orange-400"
          }`}
        >
          {hasProde ? "Actualizar pronóstico" : "Completar pronóstico"}
        </button>
      )}
      <AuthModal
        isOpen={isModalOpen}
        onClose={onCloseModal}
        onContinueToLogin={onContinueToLogin}
      />
    </div>
  );
};

export default SessionItem;
