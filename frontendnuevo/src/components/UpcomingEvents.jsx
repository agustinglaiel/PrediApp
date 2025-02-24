import React from "react";
import EventCard from "./EventCard";

const UpcomingEvents = ({ events }) => {
  return (
    <div className="px-4 mt-12">
      {" "}
      {/* Cambiamos mt-6 a mt-12 para más separación */}
      <h2 className="text-2xl font-bold mb-4">Próximos eventos</h2>
      {events.map((event, index) => (
        <EventCard
          key={index}
          country={event.country}
          circuit={event.circuit}
          sessions={event.sessions}
          flagUrl={event.flagUrl}
        />
      ))}
    </div>
  );
};

export default UpcomingEvents;
