import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import UpcomingEvents from "../components/UpcomingEvents";
import { getUpcomingSessions } from "../api/sessions";

const HomePage = () => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUpcomingSessions = async () => {
      try {
        setLoading(true);
        const token = localStorage.getItem("jwtToken");
        if (!token) {
          navigate("/login", { replace: true });
          return;
        }

        const data = await getUpcomingSessions();
        console.log("Upcoming sessions:", data);
        const groupedEvents = processSessions(data);
        setEvents(groupedEvents);
      } catch (err) {
        if (err.response && err.response.status === 403) {
          setError(
            "Acceso denegado (403). Verifica tu rol o token. Redirigiendo al login..."
          );
          console.error("Detalles del error 403:", {
            status: err.response.status,
            data: err.response.data,
            headers: err.response.headers,
          });
          setTimeout(() => navigate("/login", { replace: true }), 2000);
        } else if (err.response && err.response.status === 401) {
          setError("No autorizado (401). Redirigiendo al login...");
          console.error("Detalles del error 401:", err.response);
          setTimeout(() => navigate("/login", { replace: true }), 2000);
        } else {
          setError(`No se pudieron cargar los eventos: ${err.message}`);
        }
        console.error("Error fetching sessions:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchUpcomingSessions();
  }, [navigate]);

  const processSessions = (sessions) => {
    const eventsMap = {};
    console.log("Processing sessions:", sessions);

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
            : "/images/circuitLayout/default.jpg", // Placeholder por defecto para circuitLayout
          sessions: [],
        };
        console.log(
          `Generated flagUrl for ${session.country_name}:`,
          eventsMap[key].flagUrl
        );
        console.log(
          `Generated circuitLayoutUrl for ${session.country_name}:`,
          eventsMap[key].circuitLayoutUrl
        );
      }

      const [date, time] = session.date_start.split("T")[0].split("-");
      const [startTime] = session.date_start
        .split("T")[1]
        .split("-")[0]
        .split(":");
      const [endTime] = session.date_end.split("T")[1].split("-")[0].split(":");

      eventsMap[key].sessions.push({
        date: date.split("-")[2],
        month: date.split("-")[1]?.toUpperCase().padStart(3, "0"),
        type: session.session_type,
        startTime: `${startTime}:00`,
        endTime: `${endTime}:00`,
        hasPronostico: true,
      });
    });

    return Object.values(eventsMap);
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
        <UpcomingEvents events={events} />
      </main>
      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default HomePage;
