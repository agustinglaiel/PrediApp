// src/components/SessionItem.jsx
import React from "react";
import DateDisplay from "./DateDisplay";
import AuthModal from "./AuthModal";

const SessionItem = ({
  sessionId,
  date,
  month,
  sessionName,
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
  isPastEvent = false,
  score,
}) => {
  // Verificamos si existe prode
  const hasProde =
    (sessionType !== "Race" && prodeSession) ||
    (sessionType === "Race" && prodeRace);

  return (
    <div className="flex items-center p-3 border-b border-gray-100 last:border-b-0">
      <DateDisplay date={date} month={month} />

      <div className="ml-6 flex-grow">
        <div className="font-semibold">{sessionName}</div>
        <div className="text-sm text-gray-600">
          {startTime}
          {endTime ? ` - ${endTime}` : ""}
        </div>
      </div>

      {isPastEvent ? (
        // Si el evento es pasado, mostramos "Evento finalizado" en negrita
        <span className="font-bold text-black text-sm mr-4">
          {score !== null && score !== undefined ? score : "0"} Puntos
        </span>
      ) : (
        // Evento futuro => mostramos botón
        hasPronostico !== undefined && (
          <button
            onClick={onPronosticoClick}
            className={`
              px-4 py-1
              rounded-full
              text-sm font-medium
              transition-colors duration-200
              whitespace-nowrap
              ${
                hasProde
                  ? "bg-white text-green-500 border border-green-500 hover:bg-green-50"
                  : "bg-orange-300 text-white hover:bg-orange-400"
              }
            `}
          >
            {hasProde ? "Actualizar pronóstico" : "Completar pronóstico"}
          </button>
        )
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
