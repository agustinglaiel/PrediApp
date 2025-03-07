// EventCard.jsx
import React from "react";
import SessionItem from "./SessionItem";

const EventCard = ({
  country,
  circuit,
  sessions,
  flagUrl,
  circuitLayoutUrl,
  onPronosticoClick,
  isModalOpen,
  onCloseModal,
  onContinueToLogin,
}) => {
  // Función local para combinar datos del evento y la sesión
  const handlePronosticoClickLocal = (session) => {
    if (!onPronosticoClick) return;

    // Empaquetamos la info
    const sessionData = {
      id: session.id,
      sessionName: session.sessionName || session.session_name,
      sessionType: session.sessionType || session.type,
      dateStart: session.date_start || null,
      countryName: country,
      flagUrl,
      circuitName: circuit,
    };

    console.log(
      "EventCard: Llamando onPronosticoClick con sessionData:",
      sessionData
    );
    onPronosticoClick(sessionData);
  };

  return (
    <div className="bg-white rounded-xl shadow-sm mb-6 overflow-hidden">
      <div className="p-4 flex items-center">
        <div className="bg-gray-100 w-16 h-12 rounded-lg flex items-center justify-center overflow-hidden">
          {flagUrl && (
            <img
              src={flagUrl}
              alt={`${country} flag`}
              className="w-full h-auto"
            />
          )}
        </div>
        <div className="ml-4">
          <h3 className="font-bold text-xl">{country}</h3>
          <p className="text-gray-600">{circuit}</p>
        </div>
        <div className="ml-auto">
          {circuitLayoutUrl && (
            <img
              src={circuitLayoutUrl}
              alt={`${circuit} layout`}
              className="w-24 h-24 rounded-lg object-contain"
            />
          )}
        </div>
      </div>

      <div className="border-t border-gray-100">
        {sessions.map((session, index) => (
          <SessionItem
            key={index}
            sessionId={session.id} // Pasamos sessionId
            date={session.date}
            month={session.month}
            sessionType={session.type}
            startTime={session.startTime}
            endTime={session.endTime}
            hasPronostico={session.hasPronostico}
            isModalOpen={isModalOpen}
            onCloseModal={onCloseModal}
            onContinueToLogin={onContinueToLogin}
            prodeSession={session.prodeSession}
            prodeRace={session.prodeRace}
            // Importante: cuando se hace clic, armamos la data final y llamamos la función
            onPronosticoClick={() =>
              handlePronosticoClickLocal({
                ...session,
                // si tu backend la guardó como session_name en lugar de sessionName
                sessionName: session.sessionName || session.session_name,
                sessionType: session.sessionType || session.type,
              })
            }
          />
        ))}
      </div>
    </div>
  );
};

export default EventCard;
