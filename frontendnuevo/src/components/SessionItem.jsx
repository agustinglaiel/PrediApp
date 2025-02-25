import React from "react";
import DateDisplay from "./DateDisplay"; // Importamos el nuevo componente

const SessionItem = ({
  date,
  month,
  sessionType,
  startTime,
  endTime,
  hasPronostico,
}) => {
  console.log(`Rendering SessionItem: date=${date}, month=${month}`); // Depuración
  return (
    <div className="flex items-center p-3 border-b border-gray-100 last:border-b-0">
      <DateDisplay date={date} month={month} /> {/* Usamos DateDisplay aquí */}
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
