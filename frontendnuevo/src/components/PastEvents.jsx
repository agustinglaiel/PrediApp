// src/components/PastEvents.jsx
import React from "react";
import EventCard from "./EventCard";

const PastEvents = ({
  events,
  onPronosticoClick,
  isModalOpen,
  onCloseModal,
  onContinueToLogin,
}) => {
  return (
    <div className="px-4 mt-12">
      <h2 className="text-2xl font-bold mb-4">Eventos anteriores</h2>
      {events.length === 0 ? (
        <p className="text-gray-500 text-sm italic">
          No hay eventos pasados registrados.
        </p>
      ) : (
        events.map((event, index) => (
          <EventCard
            key={index}
            country={event.country}
            circuit={event.circuit}
            sessions={event.sessions}
            flagUrl={event.flagUrl}
            circuitLayoutUrl={event.circuitLayoutUrl}
            isPastEvent={true}
            onPronosticoClick={onPronosticoClick}
            isModalOpen={isModalOpen}
            onCloseModal={onCloseModal}
            onContinueToLogin={onContinueToLogin}
          />
        ))
      )}
    </div>
  );
};

export default PastEvents;
