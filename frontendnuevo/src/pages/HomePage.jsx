import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom"; // Importamos useNavigate para redirigir
import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import UpcomingEvents from "../components/UpcomingEvents";
import { getUpcomingSessions } from "../api/sessions";

const HomePage = () => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false); // Estado para controlar el modal
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUpcomingSessions = async () => {
      try {
        setLoading(true);
        const data = await getUpcomingSessions();
        const groupedEvents = processSessions(data);
        setEvents(groupedEvents);
      } catch (err) {
        if (err.response && err.response.status === 403) {
          setError(
            "Acceso denegado (403). Verifica los permisos o contacta al soporte."
          );
        } else if (err.response && err.response.status === 401) {
          setError(
            "No autorizado (401). Verifica los permisos o contacta al soporte."
          );
        } else {
          setError(`No se pudieron cargar los eventos: ${err.message}`);
        }
        console.error("Error fetching sessions:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchUpcomingSessions();
  }, []);

  const processSessions = (sessions) => {
    const eventsMap = {};

    sessions.forEach((session) => {
      const key = `${session.country_name}-${session.circuit_short_name}`;
      if (!eventsMap[key]) {
        eventsMap[key] = {
          country: session.country_name,
          circuit: session.circuit_short_name,
          flagUrl: session.country_name
            ? `/images/flags/${session.country_name.toLowerCase()}.jpg`
            : "/images/flags/default.jpg", // Placeholder por defecto si country_name es undefined o vacío
          circuitLayoutUrl: session.country_name
            ? `/images/circuitLayouts/${session.country_name.toLowerCase()}.png`
            : "/images/circuitLayouts/default.png", // Placeholder por defecto para circuitLayout
          sessions: [],
        };
      }

      // Depuración y manejo robusto de date_start
      let day = "1"; // Fallback por defecto si date_start es inválido
      let month = "JAN"; // Fallback por defecto si date_start es inválido
      if (session.date_start && typeof session.date_start === "string") {
        try {
          const [datePart] = session.date_start.split("T"); // Obtener solo la parte de la fecha (e.g., "2025-12-03")
          if (datePart) {
            const [year, monthNum, dayNum] = datePart.split("-");
            day = dayNum; // Día (e.g., "03" → "3")
            // Convertir el número del mes (e.g., "12") a mes en 3 letras (e.g., "DEC")
            const months = [
              "JAN",
              "FEB",
              "MAR",
              "APR",
              "MAY",
              "JUN",
              "JUL",
              "AUG",
              "SEP",
              "OCT",
              "NOV",
              "DEC",
            ];
            month = months[parseInt(monthNum, 10) - 1];
          }
        } catch (error) {
          console.error("Error parsing date_start:", session.date_start, error);
        }
      } else {
        console.warn("date_start is invalid or undefined:", session.date_start);
      }

      const [startTime] = session.date_start
        .split("T")[1]
        .split("-")[0]
        .split(":");
      const [endTime] = session.date_end.split("T")[1].split("-")[0].split(":");

      eventsMap[key].sessions.push({
        date: day, // Día como número (e.g., "3")
        month: month, // Mes en formato "DEC"
        type: session.session_type,
        startTime: `${startTime}:00`,
        endTime: `${endTime}:00`,
        hasPronostico: true,
      });
    });

    return Object.values(eventsMap);
  };

  // Función para manejar el clic en "Completar pronóstico"
  const handlePronosticoClick = () => {
    const token = localStorage.getItem("jwtToken");
    if (token) {
      // Si hay token, redirige a la página de pronósticos (puedes ajustar la ruta)
      navigate("/pronosticos"); // Placeholder, ajusta según necesites
    } else {
      // Si no hay token, muestra el modal
      setIsModalOpen(true);
    }
  };

  // Función para cerrar el modal
  const closeModal = () => {
    setIsModalOpen(false);
  };

  // Función para redirigir al login y cerrar el modal
  const handleContinueToLogin = () => {
    navigate("/login", { replace: true });
    setIsModalOpen(false);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-50">
        <p className="text-gray-600">Cargando eventos...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-50">
        <p className="text-red-600">{error}</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen bg-gray-50">
      <Header />
      <NavigationBar />
      <main className="flex-grow pt-24">
        <UpcomingEvents
          events={events}
          onPronosticoClick={handlePronosticoClick} // Pasamos la función como prop
          isModalOpen={isModalOpen}
          onCloseModal={closeModal}
          onContinueToLogin={handleContinueToLogin} // Pasamos las funciones al modal
        />
      </main>
      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default HomePage;
