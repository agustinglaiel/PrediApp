// frontendnuevo/src/pages/ProdeRaceResultPage.jsx
import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";

import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import SessionHeader from "../components/pronosticos/SessionHeader";
import DriverResultDisplay from "../components/results/DriverResultDisplay"; // Cambiado a este componente
import YesNoButton from "../components/pronosticos/YesNoButton";
import MissingProdeSession from "../components/results/MissingProdeSession";

import { getProdeByUserAndSession } from "../api/prodes";
import { getDriverById } from "../api/drivers";
import { getSessionById } from "../api/sessions";
import { getTopNDriversInSession } from "../api/results";

const ProdeRaceResultPage = () => {
  const { session_id } = useParams();
  const navigate = useNavigate();

  const [sessionDetails, setSessionDetails] = useState(null);
  const [prodeData, setProdeData] = useState(null);
  const [userPrediction, setUserPrediction] = useState({
    p1: null,
    p2: null,
    p3: null,
    p4: null,
    p5: null,
    vsc: null,
    sc: null,
    dnf: null,
  });
  const [realResults, setRealResults] = useState({
    p1: null,
    p2: null,
    p3: null,
    p4: null,
    p5: null,
    vsc: null,
    sc: null,
    dnf: null,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showMissingProdeModal, setShowMissingProdeModal] = useState(false);

  useEffect(() => {
    const fetchRaceData = async () => {
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
          sessionType: sessionData.session_type || "Race",
          sessionName: sessionData.session_name || "Race",
          dateStart: sessionData.date_start || "2025-01-01T00:00:00Z",
        };
        setSessionDetails(sessionInfo);

        // Verificar que sea una carrera
        if (
          sessionData.session_type !== "Race" ||
          sessionData.session_name !== "Race"
        ) {
          navigate("/pronosticos");
          return;
        }

        // Obtener pronóstico del usuario
        const prode = await getProdeByUserAndSession(
          parseInt(userId, 10),
          sessionId
        );
        if (!prode) {
          setShowMissingProdeModal(true);
          setProdeData(null);
        } else {
          setProdeData(prode);
          const driverPromises = [
            prode.p1 ? getDriverById(prode.p1) : Promise.resolve(null),
            prode.p2 ? getDriverById(prode.p2) : Promise.resolve(null),
            prode.p3 ? getDriverById(prode.p3) : Promise.resolve(null),
            prode.p4 ? getDriverById(prode.p4) : Promise.resolve(null),
            prode.p5 ? getDriverById(prode.p5) : Promise.resolve(null),
          ];
          const [p1, p2, p3, p4, p5] = await Promise.all(driverPromises);
          setUserPrediction({
            p1: p1?.full_name || "No seleccionado",
            p2: p2?.full_name || "No seleccionado",
            p3: p3?.full_name || "No seleccionado",
            p4: p4?.full_name || "No seleccionado",
            p5: p5?.full_name || "No seleccionado",
            vsc: prode.vsc || false,
            sc: prode.sc || false,
            dnf: prode.dnf || 0,
          });
        }

        // Obtener resultados reales
        const topDrivers = await getTopNDriversInSession(sessionId, 5);
        if (topDrivers.length === 0) {
          setRealResults({
            p1: "No disponible",
            p2: "No disponible",
            p3: "No disponible",
            p4: "No disponible",
            p5: "No disponible",
            vsc: sessionData.vsc || false,
            sc: sessionData.sc || false,
            dnf: sessionData.dnf || 0,
          });
        } else {
          const realDriverPromises = topDrivers.map((driver) =>
            getDriverById(driver.driver_id)
          );
          const realDriverData = await Promise.all(realDriverPromises);
          setRealResults({
            p1: realDriverData[0]?.full_name || "No disponible",
            p2: realDriverData[1]?.full_name || "No disponible",
            p3: realDriverData[2]?.full_name || "No disponible",
            p4: realDriverData[3]?.full_name || "No disponible",
            p5: realDriverData[4]?.full_name || "No disponible",
            vsc: sessionData.vsc || false,
            sc: sessionData.sc || false,
            dnf: sessionData.dnf || 0,
          });
        }
      } catch (err) {
        setError(err.message || "Error al cargar los resultados.");
        console.error("Error en ProdeRaceResultPage:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchRaceData();
  }, [session_id, navigate]);

  const handleCloseMissingProdeModal = () => {
    setShowMissingProdeModal(false);
    navigate("/pronosticos");
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
              <h3 className="text-lg font-semibold text-gray-800 mt-2 mb-2 ml-4">
                Tu Pronóstico
              </h3>
              <div className="flex flex-col gap-4">
                <DriverResultDisplay
                  position="P1"
                  driverName={userPrediction.p1}
                />
                <DriverResultDisplay
                  position="P2"
                  driverName={userPrediction.p2}
                />
                <DriverResultDisplay
                  position="P3"
                  driverName={userPrediction.p3}
                />
                <DriverResultDisplay
                  position="P4"
                  driverName={userPrediction.p4}
                />
                <DriverResultDisplay
                  position="P5"
                  driverName={userPrediction.p5}
                />
                <div className="flex flex-row gap-14 ml-4 mb-1">
                  <YesNoButton
                    label="Virtual Safety Car"
                    value={userPrediction.vsc}
                    onChange={() => {}}
                    disabled={true}
                  />
                  <YesNoButton
                    label="Safety Car"
                    value={userPrediction.sc}
                    onChange={() => {}}
                    disabled={true}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-black mt-2 mb-1 ml-4">
                    DNF
                  </label>
                  <input
                    type="number"
                    value={userPrediction.dnf}
                    disabled={true}
                    className="border border-gray-300 p-2 rounded w-24 ml-4 bg-gray-100"
                  />
                </div>
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
                No hay pronóstico disponible para esta carrera.
              </p>
            )
          )}
        </div>

        {/* Resultados reales */}
        <div className="mt-4 p-2 bg-white rounded-lg shadow-md">
          <h3 className="text-lg font-semibold text-gray-800 mt-2 mb-4 ml-4">
            Resultados
          </h3>
          <div className="flex flex-col gap-4">
            <DriverResultDisplay position="P1" driverName={realResults.p1} />
            <DriverResultDisplay position="P2" driverName={realResults.p2} />
            <DriverResultDisplay position="P3" driverName={realResults.p3} />
            <DriverResultDisplay position="P4" driverName={realResults.p4} />
            <DriverResultDisplay position="P5" driverName={realResults.p5} />
            <div className="flex flex-row gap-14 ml-4 mb-1">
              <YesNoButton
                label="Virtual Safety Car"
                value={realResults.vsc}
                onChange={() => {}}
                disabled={true}
              />
              <YesNoButton
                label="Safety Car"
                value={realResults.sc}
                onChange={() => {}}
                disabled={true}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-black mt-2 mb-1 ml-4">
                DNF
              </label>
              <input
                type="number"
                value={realResults.dnf}
                disabled={true}
                className="border border-gray-300 p-2 rounded w-24 ml-4 bg-gray-100"
              />
            </div>
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

export default ProdeRaceResultPage;
