// ProdeSessionPage.jsx
import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getAllDrivers } from "../api/drivers"; // <--- API que ya tienes
import SessionHeader from "../components/pronosticos/SessionHeader";
import Top3FormHeader from "../components/pronosticos/Top3FormHeader";
import DriverSelect from "../components/pronosticos/DriverSelect";
import SubmitButton from "../components/pronosticos/SubmitButton";
import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import WarningModal from "../components/pronosticos/WarningModal";

// Simulamos isRaceSession desde el frontend para consistencia con el backend
const isRaceSession = (sessionName, sessionType) => {
  return sessionName === "Race" && sessionType === "Race";
};

const ProdeSessionPage = () => {
  const { session_id } = useParams();
  const navigate = useNavigate();

  // Almacenar la lista completa de pilotos
  const [allDrivers, setAllDrivers] = useState([]);
  const [loadingDrivers, setLoadingDrivers] = useState(true);
  const [driversError, setDriversError] = useState(null);

  // Estado para la sesión (dummy)
  const [sessionDetails, setSessionDetails] = useState({
    countryName: "Hungary",
    flagUrl: "/images/flags/hungary.jpg",
    sessionType: "Qualifying",
    sessionName: "Qualifying",
    dateStart: "2025-12-02T04:00:00-03:00",
  });

  // Estado para lo que el usuario selecciona en P1/P2/P3
  const [formData, setFormData] = useState({
    P1: null,
    P2: null,
    P3: null,
  });

  // Manejar si el formulario está completo (los 3 pilotos)
  const isFormComplete = formData.P1 && formData.P2 && formData.P3;

  // Modal warning
  const [showWarningModal, setShowWarningModal] = useState(false);

  // 1. Al montar el componente, obtenemos la lista de pilotos una sola vez
  useEffect(() => {
    async function fetchDrivers() {
      try {
        setLoadingDrivers(true);
        const response = await getAllDrivers();
        setAllDrivers(response);
      } catch (error) {
        setDriversError(`Error cargando pilotos: ${error.message}`);
      } finally {
        setLoadingDrivers(false);
      }
    }
    fetchDrivers();

    // También aquí podría ir la lógica para getSessionById(session_id), etc.
  }, [session_id]);

  // 2. Chequear el timeRemaining (dummy) para mostrar modal
  useEffect(() => {
    const now = new Date();
    const sessionStart = new Date(sessionDetails.dateStart);
    const fiveMinutesInMs = 5 * 60 * 1000;
    const timeDifference = sessionStart - now;
    if (timeDifference <= fiveMinutesInMs && timeDifference > 0) {
      setShowWarningModal(true);
    }
  }, [sessionDetails.dateStart]);

  // 3. Cuando el usuario cambia P1/P2/P3
  const handleDriverChange = (position, value) => {
    setFormData((prev) => ({
      ...prev,
      [position]: value,
    }));
  };

  // 4. Evitar que un mismo piloto se seleccione varias veces
  //    => Filtramos la lista "allDrivers" para cada dropdown.
  const driversForP1 = allDrivers.filter(
    (d) => d.id !== formData.P2 && d.id !== formData.P3
  );
  const driversForP2 = allDrivers.filter(
    (d) => d.id !== formData.P1 && d.id !== formData.P3
  );
  const driversForP3 = allDrivers.filter(
    (d) => d.id !== formData.P1 && d.id !== formData.P2
  );

  // 5. Manejar envío
  const handleSubmit = (e) => {
    e.preventDefault();
    console.log("Enviando pronóstico con:", formData);
    // Llamada al backend...
    navigate("/");
  };

  // Cerrar modal
  const handleCloseModal = () => {
    setShowWarningModal(false);
  };

  // 6. Render
  if (loadingDrivers) {
    return <div>Cargando pilotos...</div>;
  }
  if (driversError) {
    return <div>{driversError}</div>;
  }

  const isRace = isRaceSession(
    sessionDetails.sessionName,
    sessionDetails.sessionType
  );

  return (
    <div>
      <Header />
      <NavigationBar />
      <main className="pt-28 px-4">
        <SessionHeader
          countryName={sessionDetails.countryName}
          flagUrl={sessionDetails.flagUrl}
          sessionName={sessionDetails.sessionName}
          sessionType={sessionDetails.sessionType}
          className="mt-6"
        />

        {!isRace && (
          <div className="mt-4 p-4 bg-white rounded-lg shadow-md">
            <Top3FormHeader sessionType={sessionDetails.sessionType} />

            <form
              onSubmit={handleSubmit}
              disabled={showWarningModal}
              className="flex flex-col gap-4"
            >
              <DriverSelect
                position="P1"
                value={formData.P1}
                onChange={(value) => handleDriverChange("P1", value)}
                disabled={showWarningModal}
                drivers={driversForP1} // <--- Filtrado
              />
              <DriverSelect
                position="P2"
                value={formData.P2}
                onChange={(value) => handleDriverChange("P2", value)}
                disabled={showWarningModal}
                drivers={driversForP2} // <--- Filtrado
              />
              <DriverSelect
                position="P3"
                value={formData.P3}
                onChange={(value) => handleDriverChange("P3", value)}
                disabled={showWarningModal}
                drivers={driversForP3} // <--- Filtrado
              />

              <SubmitButton
                isDisabled={!isFormComplete || showWarningModal}
                onClick={handleSubmit}
                label="Enviar pronóstico"
                className="mt-4"
              />
            </form>
          </div>
        )}

        <button onClick={() => navigate("/")} className="mt-4">
          Volver
        </button>

        <WarningModal isOpen={showWarningModal} onClose={handleCloseModal} />
      </main>
    </div>
  );
};

export default ProdeSessionPage;
