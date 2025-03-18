// ProdeRacePage.jsx
import React, { useState, useEffect } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";

import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import SessionHeader from "../components/pronosticos/SessionHeader";
import Top5FormHeader from "../components/pronosticos/Top5FormHeader";
import DriverSelect from "../components/pronosticos/DriverSelect";
import YesNoButton from "../components/pronosticos/YesNoButton";
import SubmitButton from "../components/pronosticos/SubmitButton";
import WarningModal from "../components/pronosticos/WarningModal";

import { getAllDrivers } from "../api/drivers";
import { createProdeCarrera } from "../api/prodes";

const isRaceSession = (sessionName, sessionType) => {
  return sessionName === "Race" && sessionType === "Race";
};

const ProdeRacePage = () => {
  const { session_id } = useParams();
  const { state } = useLocation();
  const navigate = useNavigate();

  const [allDrivers, setAllDrivers] = useState([]);
  const [loadingDrivers, setLoadingDrivers] = useState(true);
  const [driversError, setDriversError] = useState(null);

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
      return {
        countryName: "Hungary",
        flagUrl: "/images/flags/hungary.jpg",
        sessionType: "Race",
        sessionName: "Race",
        dateStart: "2025-12-02T04:00:00-03:00",
      };
    }
  });

  const [formData, setFormData] = useState({
    P1: null,
    P2: null,
    P3: null,
    P4: null,
    P5: null,
    vsc: false,
    sc: false,
    dnf: 0,
  });

  const isFormComplete =
    formData.P1 &&
    formData.P2 &&
    formData.P3 &&
    formData.P4 &&
    formData.P5 &&
    formData.dnf >= 0;

  const [showWarningModal, setShowWarningModal] = useState(false);

  // Cargar pilotos
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

  // Mostrar modal si faltan <5 min
  useEffect(() => {
    const now = new Date();
    const sessionStart = new Date(sessionDetails.dateStart);
    const fiveMinutes = 5 * 60 * 1000;
    const diff = sessionStart - now;
    if (diff <= fiveMinutes && diff > 0) {
      setShowWarningModal(true);
    }
  }, [sessionDetails.dateStart]);

  // Manejar selects
  const handleDriverChange = (position, value) => {
    setFormData((prev) => ({ ...prev, [position]: value }));
  };
  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  // Filtrar
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

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log("ProdeRacePage submit:", formData, session_id);

    try {
      const payload = {
        session_id,
        p1: formData.P1,
        p2: formData.P2,
        p3: formData.P3,
        p4: formData.P4,
        p5: formData.P5,
        vsc: formData.vsc,
        sc: formData.sc,
        dnf: formData.dnf,
      };

      const response = await createProdeCarrera(payload);
      console.log("ProdeCarrera response:", response);
      navigate("/");
    } catch (error) {
      console.error("Error en createProdeCarrera:", error.message);
    }
  };

  const handleCloseModal = () => setShowWarningModal(false);

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
    <div className="flex flex-col min-h-screen">
      <Header />
      <NavigationBar />

      {/* main con flex-grow para empujar el footer */}
      <main className="flex-grow pt-28 px-4">
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
            <form
              onSubmit={handleSubmit}
              disabled={showWarningModal}
              className="flex flex-col gap-4 mt-4"
            >
              {/* P1..P5 */}
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

              <div className="flex flex-row gap-14 ml-4 mb-4">
                <YesNoButton
                  label="Virtual Safety Car"
                  value={formData.vsc}
                  onChange={(newVal) => handleChange("vsc", newVal)}
                  disabled={showWarningModal}
                />
                <YesNoButton
                  label="Safety Car"
                  value={formData.sc}
                  onChange={(newVal) => handleChange("sc", newVal)}
                  disabled={showWarningModal}
                />
              </div>

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
        <WarningModal isOpen={showWarningModal} onClose={handleCloseModal} />
      </main>

      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default ProdeRacePage;
