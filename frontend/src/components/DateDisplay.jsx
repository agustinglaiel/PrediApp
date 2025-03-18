// src/components/DateDisplay.jsx
import React from "react";

const DateDisplay = ({ date, month }) => {
  return (
    <div className="bg-gray-200 rounded-lg p-2 text-center w-16">
      <div className="font-bold text-lg">{date}</div>
      <div className="text-xs uppercase text-gray-500">{month}</div>
    </div>
  );
};

export default DateDisplay;
