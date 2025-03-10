// src/pages/PronosticosPage.jsx
import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import UpcomingEvents from "../components/UpcomingEvents";
import PastEvents from "../components/PastEvents";

import { getUpcomingSessions, getPastSessions } from "../api/sessions";
import { getProdeByUserAndSession } from "../api/prodes";

const PronosticosPage = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [upcomingEvents, setUpcomingEvents] = useState([]);
  const [pastEvents, setPastEvents] = useState([]);

  const navigate = useNavigate();

  // Agrupa las sesiones por weekend_id
  const groupSessionsByWeekend = (sessions) => {
    const eventsMap = {};

    sessions.forEach((session) => {
      const weekendId = session.weekend_id;

      if (!eventsMap[weekendId]) {
        eventsMap[weekendId] = {
          country: session.country_name,
          circuit: session.circuit_short_name,
          flagUrl: session.country_name
            ? `/images/flags/${session.country_name.toLowerCase()}.jpg`
            : "/images/flags/default.jpg",
          circuitLayoutUrl: session.country_name
            ? `/images/circuitLayouts/${session.country_name.toLowerCase()}.png`
            : "/images/circuitLayouts/default.png",
          sessions: [],
        };
      }

      const dateStartObj = new Date(session.date_start);
      const dateEndObj = new Date(session.date_end);

      // Para eliminar AM/PM, añadimos hour12: false
      const startTime = dateStartObj.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
      });
      const endTime = dateEndObj.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
      });

      eventsMap[weekendId].sessions.push({
        id: session.id,
        date: dateStartObj.getDate().toString(),
        month: dateStartObj
          .toLocaleString("en", { month: "short" })
          .toUpperCase(),
        sessionName: session.session_name,
        sessionType: session.session_type,
        startTime,
        endTime,
        date_start: session.date_start,
        hasPronostico: true, // seguirás ajustando según tu lógica
        prodeSession: null,
        prodeRace: null,
      });
    });

    return Object.values(eventsMap);
  };

  // Rellena cada sesión con los datos de Prode si el usuario está logueado
  const fillProdeData = async (eventsArray) => {
    const userId = localStorage.getItem("userId");
    if (!userId) return eventsArray;

    for (const event of eventsArray) {
      const prodePromises = event.sessions.map(async (sess) => {
        try {
          const prode = await getProdeByUserAndSession(
            parseInt(userId, 10),
            sess.id
          );
          if (prode) {
            // Si hay p4 y p5, asumimos que es carrera
            if (prode.p4 !== undefined && prode.p5 !== undefined) {
              sess.prodeRace = prode;
              sess.prodeSession = null;
            } else {
              sess.prodeSession = prode;
              sess.prodeRace = null;
            }
          }
        } catch {
          // Podrías ignorar errores 404
        }
      });
      await Promise.all(prodePromises);
    }

    return eventsArray;
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);

        // Llamamos a ambos endpoints
        const [upcomingRaw, pastRaw] = await Promise.all([
          getUpcomingSessions(),
          getPastSessions(),
        ]);

        const upcomingGrouped = groupSessionsByWeekend(upcomingRaw || []);
        const pastGrouped = groupSessionsByWeekend(pastRaw || []);

        const upcomingWithProde = await fillProdeData(upcomingGrouped);
        const pastWithProde = await fillProdeData(pastGrouped);

        setUpcomingEvents(upcomingWithProde);
        setPastEvents(pastWithProde);
      } catch (err) {
        setError(`Error al cargar datos: ${err.message}`);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handlePronosticoClick = (sessionData) => {
    navigate(`/pronosticos/${sessionData.id}`, { state: sessionData });
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-50">
        <p className="text-gray-500">Cargando sesiones...</p>
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

      <main className="flex-grow pt-24">
        {/* Próximos eventos */}
        <UpcomingEvents
          events={upcomingEvents}
          onPronosticoClick={handlePronosticoClick}
        />

        {/* Eventos anteriores */}
        <PastEvents
          events={pastEvents}
          onPronosticoClick={handlePronosticoClick}
        />
      </main>

      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default PronosticosPage;
