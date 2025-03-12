// frontendnuevo/src/pages/ProdeSessionResultPage.jsx
import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";

import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import SessionHeader from "../components/pronosticos/SessionHeader";
import DriverResultDisplay from "../components/results/DriverResultDisplay";
import MissingProdeSession from "../components/results/MissingProdeSession";

import { getProdeByUserAndSession } from "../api/prodes";
import { getDriverById } from "../api/drivers";
import { getSessionById } from "../api/sessions";
import { getTopNDriversInSession } from "../api/results";

const ProdeSessionResultPage = () => {
  const { session_id } = useParams();
  const navigate = useNavigate();

  const [sessionDetails, setSessionDetails] = useState(null);
  const [prodeData, setProdeData] = useState(null);
  const [userDrivers, setUserDrivers] = useState({
    p1: null,
    p2: null,
    p3: null,
  });
  const [realDrivers, setRealDrivers] = useState({
    p1: null,
    p2: null,
    p3: null,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showMissingProdeModal, setShowMissingProdeModal] = useState(false);

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

        // Obtener detalles de la sesión
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
        setSessionDetails(sessionInfo);

        // Obtener pronóstico del usuario
        const prode = await getProdeByUserAndSession(
          parseInt(userId, 10),
          sessionId
        );
        if (!prode) {
          setShowMissingProdeModal(true); // Mostrar modal si no hay pronóstico
          setProdeData(null); // Asegurar que no haya datos residuales
        } else {
          if (prode.p4 !== undefined || prode.p5 !== undefined) {
            navigate("/pronosticos"); // Redirigir si es carrera
            return;
          }
          setProdeData(prode);

          const userDriverPromises = [
            prode.p1 ? getDriverById(prode.p1) : Promise.resolve(null),
            prode.p2 ? getDriverById(prode.p2) : Promise.resolve(null),
            prode.p3 ? getDriverById(prode.p3) : Promise.resolve(null),
          ];
          const [driverP1, driverP2, driverP3] = await Promise.all(
            userDriverPromises
          );
          setUserDrivers({
            p1: driverP1 ? driverP1.full_name : null,
            p2: driverP2 ? driverP2.full_name : null,
            p3: driverP3 ? driverP3.full_name : null,
          });
        }

        // Obtener los top 3 drivers reales
        const topDrivers = await getTopNDriversInSession(sessionId, 3);
        const realDriverPromises = topDrivers.map((driver) =>
          getDriverById(driver.driver_id)
        );
        const realDriverData = await Promise.all(realDriverPromises);

        setRealDrivers({
          p1: realDriverData[0]?.full_name || "No disponible",
          p2: realDriverData[1]?.full_name || "No disponible",
          p3: realDriverData[2]?.full_name || "No disponible",
        });
      } catch (err) {
        setError(err.message || "Error al cargar los resultados.");
        console.error("Error en ProdeSessionResultPage:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchSessionAndProdeData();
  }, [session_id, navigate]);

  const handleCloseMissingProdeModal = () => {
    setShowMissingProdeModal(false);
    navigate("/pronosticos"); // Redirigir al cerrar el modal
  };

  if (loading) {
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
        {/* Pronóstico del usuario */}
        <SessionHeader
          countryName={sessionDetails?.countryName}
          flagUrl={sessionDetails?.flagUrl}
          sessionName={sessionDetails?.sessionName}
          sessionType={sessionDetails?.sessionType}
          className="mt-4"
        />
        <div className="mt-4 p-2 bg-white rounded-lg shadow-md">
          {prodeData ? (
            <>
              <h3 className="text-lg font-semibold text-gray-800 mt-2 mb-4 ml-4">
                Tu Pronóstico
              </h3>
              <div className="flex flex-col gap-4">
                <DriverResultDisplay
                  position="P1"
                  driverName={userDrivers.p1}
                />
                <DriverResultDisplay
                  position="P2"
                  driverName={userDrivers.p2}
                />
                <DriverResultDisplay
                  position="P3"
                  driverName={userDrivers.p3}
                />
              </div>
              {prodeData?.score !== null && prodeData?.score !== undefined && (
                <div className="mt-4 text-center">
                  <p className="text-lg font-semibold text-gray-800">
                    Puntaje obtenido: {prodeData.score} puntos
                  </p>
                </div>
              )}
            </>
          ) : (
            !showMissingProdeModal && (
              <p className="text-center text-gray-600">
                No hay pronóstico disponible para esta sesión.
              </p>
            )
          )}
        </div>

        {/* Resultados reales */}
        <div className="mt-4 p-2 bg-white rounded-lg shadow-md">
          <h3 className="text-lg font-semibold text-gray-800 mt-2 mb-4 ml-4">
            Resultados Reales
          </h3>
          <div className="flex flex-col gap-4">
            <DriverResultDisplay position="P1" driverName={realDrivers.p1} />
            <DriverResultDisplay position="P2" driverName={realDrivers.p2} />
            <DriverResultDisplay position="P3" driverName={realDrivers.p3} />
          </div>
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
