// ProdeSessionPage.jsx
import React, { useState, useEffect } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";

import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import SessionHeader from "../components/pronosticos/SessionHeader";
import Top3FormHeader from "../components/pronosticos/Top3FormHeader";
import DriverSelect from "../components/pronosticos/DriverSelect";
import SubmitButton from "../components/pronosticos/SubmitButton";
import WarningModal from "../components/pronosticos/WarningModal";

import { getAllDrivers } from "../api/drivers";
import { createProdeSession } from "../api/prodes"; // <--- Usamos esta función

// Simulamos isRaceSession para verificar si la sesión es Race (no deberíamos usar ProdeSessionPage si es Race)
const isRaceSession = (sessionName, sessionType) => {
  return sessionName === "Race" && sessionType === "Race";
};

const ProdeSessionPage = () => {
  const { session_id } = useParams();
  const navigate = useNavigate();
  const { state } = useLocation(); // Si venimos de HomePage con navigate("/pronosticos/...", { state: ... })

  // Lista de pilotos
  const [allDrivers, setAllDrivers] = useState([]);
  const [loadingDrivers, setLoadingDrivers] = useState(true);
  const [driversError, setDriversError] = useState(null);

  // Detalles de la sesión
  const [sessionDetails, setSessionDetails] = useState(() => {
    if (state) {
      return {
        countryName: state.countryName,
        flagUrl: state.flagUrl,
        sessionType: state.sessionTypeng,
        sessionName: state.sessionNameng,
        dateStart: state.dateStart,
      };
    }
  });

  // Estado para P1, P2, P3
  const [formData, setFormData] = useState({
    P1: null,
    P2: null,
    P3: null,
  });

  // Ver si está completo
  const isFormComplete = formData.P1 && formData.P2 && formData.P3;

  // Modal warning
  const [showWarningModal, setShowWarningModal] = useState(false);

  // Cargar pilotos
  useEffect(() => {
    async function fetchDrivers() {
      try {
        setLoadingDrivers(true);
        const response = await getAllDrivers();
        setAllDrivers(response);
      } catch (err) {
        setDriversError(`Error cargando pilotos: ${err.message}`);
      } finally {
        setLoadingDrivers(false);
      }
    }
    fetchDrivers();

    // Si quieres cargar la sesión del backend en caso de no tener state:
    // getSessionById(session_id).then(resp => setSessionDetails(resp))
  }, [session_id]);

  // Mostrar warning si faltan < 5 min
  useEffect(() => {
    const now = new Date();
    const sessionStart = new Date(sessionDetails.dateStart);
    const fiveMinutesInMs = 5 * 60 * 1000;
    const diff = sessionStart - now;
    if (diff <= fiveMinutesInMs && diff > 0) {
      setShowWarningModal(true);
    }
  }, [sessionDetails.dateStart]);

  // Manejar cambio de pilotos
  const handleDriverChange = (position, value) => {
    setFormData((prev) => ({ ...prev, [position]: value }));
  };

  // Evitar elegir mismo piloto
  const driversForP1 = allDrivers.filter(
    (d) => d.id !== formData.P2 && d.id !== formData.P3
  );
  const driversForP2 = allDrivers.filter(
    (d) => d.id !== formData.P1 && d.id !== formData.P3
  );
  const driversForP3 = allDrivers.filter(
    (d) => d.id !== formData.P1 && d.id !== formData.P2
  );

  // Enviar pronóstico (sesión no Race)
  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const payload = {
        session_id,
        p1: formData.P1,
        p2: formData.P2,
        p3: formData.P3,
      };

      // Llamamos a createProdeSession
      const response = await createProdeSession(payload);
      console.log("ProdeSession response:", response);

      // Podrías mostrar un mensaje o navegar
      navigate("/");
    } catch (err) {
      console.error("Error createProdeSession:", err.message);
      // Manejar error en la UI
    }
  };

  const handleCloseModal = () => setShowWarningModal(false);

  if (loadingDrivers) {
    return <div>Cargando pilotos...</div>;
  }
  if (driversError) {
    return <div>{driversError}</div>;
  }

  // Si la sesión es Race, no deberíamos mostrar este form
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

        {/* Solo mostramos este form si NO es Race */}
        {!isRace && (
          <div className="mt-4 p-4 bg-white rounded-lg shadow-md">
            <Top3FormHeader sessionType={sessionDetails.sessionType} />

            <form
              onSubmit={handleSubmit}
              disabled={showWarningModal}
              className="flex flex-col gap-4"
            >
              {/* P1, P2, P3 */}
              <DriverSelect
                position="P1"
                value={formData.P1}
                onChange={(val) => handleDriverChange("P1", val)}
                disabled={showWarningModal}
                drivers={driversForP1}
              />
              <DriverSelect
                position="P2"
                value={formData.P2}
                onChange={(val) => handleDriverChange("P2", val)}
                disabled={showWarningModal}
                drivers={driversForP2}
              />
              <DriverSelect
                position="P3"
                value={formData.P3}
                onChange={(val) => handleDriverChange("P3", val)}
                disabled={showWarningModal}
                drivers={driversForP3}
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

        <WarningModal isOpen={showWarningModal} onClose={handleCloseModal} />
      </main>

      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default ProdeSessionPage;
