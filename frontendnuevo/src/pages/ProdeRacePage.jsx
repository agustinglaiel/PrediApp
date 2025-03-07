import React, { useState, useEffect } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";
import { getAllDrivers } from "../api/drivers";

// IMPORTAMOS UNA VERSIÓN NUEVA DE HEADERS
import SessionHeader from "../components/pronosticos/SessionHeader";
import Top5FormHeader from "../components/pronosticos/Top5FormHeader";

import DriverSelect from "../components/pronosticos/DriverSelect";
import SubmitButton from "../components/pronosticos/SubmitButton";
import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import WarningModal from "../components/pronosticos/WarningModal";

const isRaceSession = (sessionName, sessionType) => {
  return sessionName === "Race" && sessionType === "Race";
};

const ProdeRacePage = () => {
  // 1) Leemos el id de la URL y el state (por si venimos de navigate)
  const { session_id } = useParams();
  const { state } = useLocation();
  const navigate = useNavigate();

  // 2) Listado de pilotos
  const [allDrivers, setAllDrivers] = useState([]);
  const [loadingDrivers, setLoadingDrivers] = useState(true);
  const [driversError, setDriversError] = useState(null);

  // 3) Detalles de la sesión (dummy + fallback a state)
  const [sessionDetails, setSessionDetails] = useState(() => {
    if (state) {
      return {
        countryName: state.countryName || "Hungary",
        flagUrl: state.flagUrl || "/images/flags/hungary.jpg",
        sessionType: state.sessionType || "Race",
        sessionName: state.sessionName || "Race",
        dateStart: state.dateStart || "2025-12-02T04:00:00-03:00",
      };
    } else {
      // Fallback dummy si no hay state
      return {
        countryName: "Hungary",
        flagUrl: "/images/flags/hungary.jpg",
        sessionType: "Race",
        sessionName: "Race",
        dateStart: "2025-12-02T04:00:00-03:00",
      };
    }
  });

  // 4) Campos del formulario de carrera:
  const [formData, setFormData] = useState({
    P1: null,
    P2: null,
    P3: null,
    P4: null,
    P5: null,
    vsc: "", // "yes" o "no"
    sc: "", // "yes" o "no"
    dnf: 0, // número de 0 a 20
  });

  // 5) Ver si el formulario está completo
  const isFormComplete =
    formData.P1 &&
    formData.P2 &&
    formData.P3 &&
    formData.P4 &&
    formData.P5 &&
    formData.vsc &&
    formData.sc &&
    formData.dnf >= 0;

  // 6) Modal de advertencia
  const [showWarningModal, setShowWarningModal] = useState(false);

  // 7) Cargar los pilotos al montar
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
  }, [session_id]);

  // 8) Mostrar warning si faltan <5min para la sesión
  useEffect(() => {
    const now = new Date();
    const sessionStart = new Date(sessionDetails.dateStart);
    const fiveMinutesInMs = 5 * 60 * 1000;
    const diff = sessionStart - now;
    if (diff <= fiveMinutesInMs && diff > 0) {
      setShowWarningModal(true);
    }
  }, [sessionDetails.dateStart]);

  // 9) Manejar cambios en driver selects (P1..P5)
  const handleDriverChange = (position, value) => {
    setFormData((prev) => ({
      ...prev,
      [position]: value,
    }));
  };

  // 10) Manejar cambios en VSC, SC y DNF
  const handleChange = (field, value) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  // 11) Filtrado para evitar elegir el mismo piloto en P1..P5
  const driversForP1 = allDrivers.filter(
    (d) =>
      d.id !== formData.P2 &&
      d.id !== formData.P3 &&
      d.id !== formData.P4 &&
      d.id !== formData.P5
  );
  const driversForP2 = allDrivers.filter(
    (d) =>
      d.id !== formData.P1 &&
      d.id !== formData.P3 &&
      d.id !== formData.P4 &&
      d.id !== formData.P5
  );
  const driversForP3 = allDrivers.filter(
    (d) =>
      d.id !== formData.P1 &&
      d.id !== formData.P2 &&
      d.id !== formData.P4 &&
      d.id !== formData.P5
  );
  const driversForP4 = allDrivers.filter(
    (d) =>
      d.id !== formData.P1 &&
      d.id !== formData.P2 &&
      d.id !== formData.P3 &&
      d.id !== formData.P5
  );
  const driversForP5 = allDrivers.filter(
    (d) =>
      d.id !== formData.P1 &&
      d.id !== formData.P2 &&
      d.id !== formData.P3 &&
      d.id !== formData.P4
  );

  // 12) Manejar submit
  const handleSubmit = (e) => {
    e.preventDefault();
    console.log("ProdeRacePage submit:", formData, session_id);
    // Llamada al backend...
    navigate("/"); // Volver a Home
  };

  const handleCloseModal = () => {
    setShowWarningModal(false);
  };

  // Render de Loading / error
  if (loadingDrivers) {
    return <div>Cargando pilotos...</div>;
  }
  if (driversError) {
    return <div>{driversError}</div>;
  }

  // 13) Ver si es “Race”
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

        {isRace && (
          <div className="mt-4 p-4 bg-white rounded-lg shadow-md">
            <Top5FormHeader sessionType={sessionDetails.sessionType} />

            {/* FORM con layout en columna */}
            <form
              onSubmit={handleSubmit}
              disabled={showWarningModal}
              className="flex flex-col gap-4 mt-4"
            >
              {/* P1..P5 (labels en DriverSelect.jsx) */}
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
              <DriverSelect
                position="P4"
                value={formData.P4}
                onChange={(val) => handleDriverChange("P4", val)}
                disabled={showWarningModal}
                drivers={driversForP4}
              />
              <DriverSelect
                position="P5"
                value={formData.P5}
                onChange={(val) => handleDriverChange("P5", val)}
                disabled={showWarningModal}
                drivers={driversForP5}
              />

              {/* Virtual Safety Car (Sí/No) */}
              <div>
                <label className="block text-sm font-medium text-black mb-1 ml-4">
                  Virtual Safety Car
                </label>
                <select
                  className="border border-gray-300 p-2 rounded ml-4"
                  value={formData.vsc}
                  onChange={(e) => handleChange("vsc", e.target.value)}
                  disabled={showWarningModal}
                >
                  <option value="">Selecciona</option>
                  <option value="yes">Sí</option>
                  <option value="no">No</option>
                </select>
              </div>

              {/* Safety Car (Sí/No) */}
              <div>
                <label className="block text-sm font-medium text-black mb-1 ml-4">
                  Safety Car
                </label>
                <select
                  className="border border-gray-300 p-2 rounded ml-4"
                  value={formData.sc}
                  onChange={(e) => handleChange("sc", e.target.value)}
                  disabled={showWarningModal}
                >
                  <option value="">Selecciona</option>
                  <option value="yes">Sí</option>
                  <option value="no">No</option>
                </select>
              </div>

              {/* DNF (0..20) */}
              <div>
                <label className="block text-sm font-medium text-black mb-1 ml-4">
                  DNF
                </label>
                <input
                  type="number"
                  min="0"
                  max="20"
                  className="border border-gray-300 p-2 rounded w-24 ml-4"
                  value={formData.dnf}
                  onChange={(e) =>
                    handleChange("dnf", parseInt(e.target.value, 10) || 0)
                  }
                  disabled={showWarningModal}
                />
              </div>

              <SubmitButton
                isDisabled={!isFormComplete || showWarningModal}
                onClick={handleSubmit}
                label="Enviar pronóstico"
                className="mt-4"
              />
            </form>
          </div>
        )}

        {/* <button onClick={() => navigate("/")} className="mt-4">
          Volver
        </button> */}

        <WarningModal isOpen={showWarningModal} onClose={handleCloseModal} />
      </main>
      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default ProdeRacePage;
