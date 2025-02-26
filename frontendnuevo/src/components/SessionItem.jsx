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

  // Depuración: Imprimir si hay un pronóstico para session_id = 6
  if (sessionId === 6) {
    console.log(
      `Session 6 - ProdeSession: ${
        prodeSession ? JSON.stringify(prodeSession) : "null"
      }, ProdeRace: ${prodeRace ? JSON.stringify(prodeRace) : "null"}`
    );
  }

  return (
    <div className="flex items-center p-3 border-b border-gray-100 last:border-b-0">
      <DateDisplay date={date} month={month} />
      <div className="ml-6 flex-grow">
        <div className="font-semibold">{sessionType}</div>
        <div className="text-sm text-gray-600">
          {startTime}
          {endTime ? ` - ${endTime}` : ""}
        </div>
        {prodeSession && (
          <div className="text-xs text-gray-500">
            Pronóstico: {prodeSession.P1}, {prodeSession.P2}, {prodeSession.P3}
          </div>
        )}
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
              ? "bg-green-500 text-white hover:bg-green-600"
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
