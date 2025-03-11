// frontendnuevo/src/pages/ProdeSessionResultPage.jsx
import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";

import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import SessionHeader from "../components/pronosticos/SessionHeader";
import DriverResultDisplay from "../components/results/DriverResultDisplay";
import MissingProdeSession from "../components/results/MissingProdeSession"; // Nuevo import

import { getProdeByUserAndSession } from "../api/prodes";
import { getDriverById } from "../api/drivers";
import { getSessionById } from "../api/sessions";

const ProdeSessionResultPage = () => {
  const { session_id } = useParams();
  const navigate = useNavigate();

  const [sessionDetails, setSessionDetails] = useState(null);
  const [prodeData, setProdeData] = useState(null);
  const [drivers, setDrivers] = useState({
    p1: null,
    p2: null,
    p3: null,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showMissingProdeModal, setShowMissingProdeModal] = useState(false); // Estado para el modal

  useEffect(() => {
    const fetchSessionAndProdeData = async () => {
      try {
        setLoading(true);
        setError(null);

        const sessionId = parseInt(session_id, 10);
        const userId = localStorage.getItem("userId");
        if (!userId) {
          throw new Error("Usuario no autenticado. Por favor, inicia sesión.");
        }

        const sessionData = await getSessionById(sessionId);
        const sessionInfo = {
          countryName: sessionData.country_name || "Unknown",
          flagUrl: sessionData.country_name
            ? `/images/flags/${sessionData.country_name.toLowerCase()}.jpg`
            : "/images/flags/default.jpg",
          sessionType: sessionData.session_type || "Qualifying",
          sessionName: sessionData.session_name || "Qualifying",
          dateStart: sessionData.date_start || "2025-01-01T00:00:00Z",
        };

        console.log(
          "ProdeSessionResultPage: sessionInfo desde API:",
          sessionInfo
        );
        setSessionDetails(sessionInfo);

        const prode = await getProdeByUserAndSession(
          parseInt(userId, 10),
          sessionId
        );
        if (!prode) {
          // En lugar de lanzar error, mostramos el modal
          setShowMissingProdeModal(true);
          setLoading(false);
          return;
        }

        if (prode.p4 !== undefined || prode.p5 !== undefined) {
          navigate("/pronosticos");
          return;
        }

        setProdeData(prode);

        const driverPromises = [
          prode.p1 ? getDriverById(prode.p1) : Promise.resolve(null),
          prode.p2 ? getDriverById(prode.p2) : Promise.resolve(null),
          prode.p3 ? getDriverById(prode.p3) : Promise.resolve(null),
        ];

        const [driverP1, driverP2, driverP3] = await Promise.all(
          driverPromises
        );

        setDrivers({
          p1: driverP1 ? driverP1.full_name : null,
          p2: driverP2 ? driverP2.full_name : null,
          p3: driverP3 ? driverP3.full_name : null,
        });
      } catch (err) {
        setError(err.message || "Error al cargar los resultados.");
        console.error("Error en ProdeSessionResultPage:", err);
      } finally {
        if (!showMissingProdeModal) setLoading(false); // Solo desactivamos loading si no mostramos el modal
      }
    };

    fetchSessionAndProdeData();
  }, [session_id, navigate]);

  const handleCloseMissingProdeModal = () => {
    setShowMissingProdeModal(false);
  };

  if (loading || (!showMissingProdeModal && !sessionDetails)) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-50">
        <p className="text-gray-500">Cargando resultados...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-50">
        <p className="text-red-500">{error}</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen bg-gray-50">
      <Header />
      <NavigationBar />
      <main className="flex-grow pt-28 px-4">
        <div className="mt-8 p-2 bg-white rounded-lg shadow-md">
          <SessionHeader
            countryName={sessionDetails?.countryName}
            flagUrl={sessionDetails?.flagUrl}
            sessionName={sessionDetails?.sessionName}
            sessionType={sessionDetails?.sessionType}
            className="mb-4"
          />
          {prodeData && (
            <>
              <div className="flex flex-col gap-4">
                <DriverResultDisplay position="P1" driverName={drivers.p1} />
                <DriverResultDisplay position="P2" driverName={drivers.p2} />
                <DriverResultDisplay position="P3" driverName={drivers.p3} />
              </div>
              {prodeData?.score !== null && prodeData?.score !== undefined && (
                <div className="mt-4 text-center">
                  <p className="text-lg font-semibold text-gray-800">
                    Puntaje obtenido: {prodeData.score} puntos
                  </p>
                </div>
              )}
            </>
          )}
        </div>
      </main>
      <MissingProdeSession
        isOpen={showMissingProdeModal}
        onClose={handleCloseMissingProdeModal}
      />
      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default ProdeSessionResultPage;
