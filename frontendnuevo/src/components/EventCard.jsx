// src/components/EventCard.jsx
import React from "react";
import SessionItem from "./SessionItem";

const EventCard = ({ country, circuit, sessions, flagUrl }) => {
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
          <div className="w-24 h-24 bg-gray-100 rounded-lg"></div>
        </div>
      </div>

      <div className="border-t border-gray-100">
        {sessions.map((session, index) => (
          <SessionItem
            key={index}
            date={session.date}
            month={session.month}
            sessionType={session.type}
            startTime={session.startTime}
            endTime={session.endTime}
            hasPronostico={session.hasPronostico}
          />
        ))}
      </div>
    </div>
  );
};

export default EventCard;
