import React from "react";

const SessionHeader = ({
  countryName,
  flagUrl,
  sessionName,
  sessionType,
  className,
}) => {
  // Determinar el texto a mostrar según sessionName y sessionType
  let headerText;
  if (sessionName === sessionType) {
    headerText = `${countryName} Main ${sessionType}`;
  } else {
    headerText = `${countryName} ${sessionName}`;
  }

  return (
    <div className={`flex items-center p-4 rounded-lg ${className || ""}`}>
      <img
        src={flagUrl}
        alt={`${countryName} flag`}
        className="w-12 h-8 mr-4"
      />
      <div>
        <div className="font-semibold">{headerText}</div>
      </div>
    </div>
  );
};

export default SessionHeader;
