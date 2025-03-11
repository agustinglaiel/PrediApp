// frontendnuevo/src/pages/ProdeSessionResultPage.jsx
import React, { useState, useEffect } from "react";
import { useParams, useLocation, useNavigate } from "react-router-dom";

import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import SessionHeader from "../components/pronosticos/SessionHeader";
import SessionResultForm from "../components/results/SessionResultForm";

import { getProdeByUserAndSession } from "../api/prodes";
import { getDriverById } from "../api/drivers";

const ProdeSessionResultPage = () => {
  const { session_id } = useParams();
  const { state } = useLocation();
  const navigate = useNavigate();

  const [sessionDetails, setSessionDetails] = useState(() => {
    if (state) {
      return {
        countryName: state.countryName || "Unknown",
        flagUrl: state.flagUrl || "/images/flags/default.jpg",
        sessionType: state.sessionType || "Qualifying",
        sessionName: state.sessionName || "Qualifying",
        dateStart: state.dateStart || "2025-01-01T00:00:00Z",
      };
    }
    return {
      countryName: "Unknown",
      flagUrl: "/images/flags/default.jpg",
      sessionType: "Qualifying",
      sessionName: "Qualifying",
      dateStart: "2025-01-01T00:00:00Z",
    };
  });

  const [prodeData, setProdeData] = useState(null);
  const [drivers, setDrivers] = useState({
    p1: null,
    p2: null,
    p3: null,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchProdeAndDrivers = async () => {
      try {
        setLoading(true);
        setError(null);

        const userId = localStorage.getItem("userId");
        if (!userId) {
          throw new Error("Usuario no autenticado. Por favor, inicia sesión.");
        }

        // Obtener el pronóstico del usuario para esta sesión
        const prode = await getProdeByUserAndSession(
          parseInt(userId, 10),
          parseInt(session_id, 10)
        );

        if (!prode) {
          throw new Error("No se encontró un pronóstico para esta sesión.");
        }

        // Verificar si es un pronóstico de carrera (p4 y p5 existen)
        if (prode.p4 !== undefined || prode.p5 !== undefined) {
          // Si es un pronóstico de carrera, podrías redirigir a otra página en el futuro
          navigate("/pronosticos");
          return;
        }

        setProdeData(prode);

        // Obtener los datos de los pilotos
        const driverPromises = [
          prode.p1 ? getDriverById(prode.p1) : Promise.resolve(null),
          prode.p2 ? getDriverById(prode.p2) : Promise.resolve(null),
          prode.p3 ? getDriverById(prode.p3) : Promise.resolve(null),
        ];

        const [driverP1, driverP2, driverP3] = await Promise.all(
          driverPromises
        );

        setDrivers({
          p1: driverP1,
          p2: driverP2,
          p3: driverP3,
        });
      } catch (err) {
        setError(err.message || "Error al cargar los resultados.");
        console.error("Error en ProdeSessionResultPage:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchProdeAndDrivers();
  }, [session_id, navigate]);

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
        <SessionHeader
          countryName={sessionDetails.countryName}
          flagUrl={sessionDetails.flagUrl}
          sessionName={sessionDetails.sessionName}
          sessionType={sessionDetails.sessionType}
          className="mt-6"
        />
        <SessionResultForm
          sessionType={sessionDetails.sessionType}
          drivers={drivers}
        />
        {prodeData?.score !== null && prodeData?.score !== undefined && (
          <div className="mt-4 text-center">
            <p className="text-lg font-semibold text-gray-800">
              Puntaje obtenido: {prodeData.score} puntos
            </p>
          </div>
        )}
      </main>
      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default ProdeSessionResultPage;
