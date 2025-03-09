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
import { createProdeSession } from "../api/prodes";

const isRaceSession = (sessionName, sessionType) => {
  return sessionName === "Race" && sessionType === "Race";
};

const ProdeSessionPage = () => {
  const { session_id } = useParams();
  const navigate = useNavigate();
  const { state } = useLocation();

  const [allDrivers, setAllDrivers] = useState([]);
  const [loadingDrivers, setLoadingDrivers] = useState(true);
  const [driversError, setDriversError] = useState(null);

  const [sessionDetails, setSessionDetails] = useState(() => {
    if (state) {
      return {
        countryName: state.countryName || "Hungary",
        flagUrl: state.flagUrl || "/images/flags/hungary.jpg",
        sessionType: state.sessionType || "Qualifying",
        sessionName: state.sessionName || "Qualifying",
        dateStart: state.dateStart || "2025-12-02T04:00:00-03:00",
      };
    } else {
      return {
        countryName: "Hungary",
        flagUrl: "/images/flags/hungary.jpg",
        sessionType: "Qualifying",
        sessionName: "Qualifying",
        dateStart: "2025-12-02T04:00:00-03:00",
      };
    }
  });

  const [formData, setFormData] = useState({
    P1: null,
    P2: null,
    P3: null,
  });

  const isFormComplete = formData.P1 && formData.P2 && formData.P3;
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
  }, [session_id]);

  // Warning si faltan <5 min
  useEffect(() => {
    const now = new Date();
    const sessionStart = new Date(sessionDetails.dateStart);
    const fiveMin = 5 * 60 * 1000;
    if (sessionStart - now <= fiveMin && sessionStart - now > 0) {
      setShowWarningModal(true);
    }
  }, [sessionDetails.dateStart]);

  const handleDriverChange = (position, value) => {
    setFormData((prev) => ({ ...prev, [position]: value }));
  };

  const driversForP1 = allDrivers.filter(
    (d) => d.id !== formData.P2 && d.id !== formData.P3
  );
  const driversForP2 = allDrivers.filter(
    (d) => d.id !== formData.P1 && d.id !== formData.P3
  );
  const driversForP3 = allDrivers.filter(
    (d) => d.id !== formData.P1 && d.id !== formData.P2
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        session_id,
        p1: formData.P1,
        p2: formData.P2,
        p3: formData.P3,
      };
      const response = await createProdeSession(payload);
      console.log("ProdeSession response:", response);
      navigate("/");
    } catch (err) {
      console.error("Error createProdeSession:", err.message);
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
      <main className="flex-grow pt-28 px-4">
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
