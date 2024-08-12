import React, { useState, useEffect } from "react";
import RacePredictionForm from "./RacePredictionForm";
import SessionPredictionForm from "./SessionPredictionForm";
import "../styles/RaceWeekendPage.css"; // Importa el archivo CSS correspondiente

const RaceWeekendPage = () => {
  const [weekendType, setWeekendType] = useState(""); // Esto se puede setear dinámicamente según el tipo de fin de semana
  const [submitted, setSubmitted] = useState(false);

  useEffect(() => {
    // Aquí puedes determinar el tipo de fin de semana (común o sprint) dinámicamente
    // Ejemplo: setWeekendType('sprint') o setWeekendType('regular');
    setWeekendType("regular"); // O 'sprint'
  }, []);

  const handleSubmit = (data) => {
    // Lógica para manejar la sumisión de los datos del formulario
    console.log("Form submitted:", data);
    setSubmitted(true);
  };

  return (
    <div className="race-weekend-page">
      <h2>Race Weekend Predictions</h2>

      {submitted ? (
        <div className="confirmation">
          <p>Your predictions have been submitted!</p>
          <button onClick={() => setSubmitted(false)}>
            Modify Predictions
          </button>
        </div>
      ) : (
        <div className="forms-container">
          {weekendType === "regular" && (
            <>
              <SessionPredictionForm
                sessionType="Qualifying"
                onSubmit={handleSubmit}
              />
              <RacePredictionForm sessionType="Race" onSubmit={handleSubmit} />
            </>
          )}

          {weekendType === "sprint" && (
            <>
              <SessionPredictionForm
                sessionType="Qualifying"
                onSubmit={handleSubmit}
              />
              <SessionPredictionForm
                sessionType="Sprint Qualifying"
                onSubmit={handleSubmit}
              />
              <SessionPredictionForm
                sessionType="Sprint Race"
                onSubmit={handleSubmit}
              />
              <RacePredictionForm sessionType="Race" onSubmit={handleSubmit} />
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default RaceWeekendPage;
