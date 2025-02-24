// src/components/SessionItem.jsx
import React from "react";

const SessionItem = ({
  date,
  month,
  sessionType,
  startTime,
  endTime,
  hasPronostico,
}) => {
  return (
    <div className="flex items-center p-3 border-b border-gray-100 last:border-b-0">
      <div className="bg-gray-200 rounded-lg p-2 text-center w-16">
        <div className="font-bold text-lg">{date}</div>
        <div className="text-xs uppercase text-gray-500">{month}</div>
      </div>

      <div className="ml-6 flex-grow">
        <div className="font-semibold">{sessionType}</div>
        <div className="text-sm text-gray-600">
          {startTime}
          {endTime ? ` - ${endTime}` : ""}
        </div>
      </div>

      {hasPronostico !== undefined && (
        <button className="bg-orange-300 text-white px-4 py-1 rounded-full text-sm font-medium">
          Completar pronóstico
        </button>
      )}
    </div>
  );
};

export default SessionItem;
