// frontendnuevo/src/pages/PronosticosPage.jsx
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

  const groupSessionsByWeekend = (sessions) => {
    const eventsMap = {};
    sessions.forEach((session) => {
      const weekendId = session.weekend_id;
      if (!eventsMap[weekendId]) {
        eventsMap[weekendId] = {
          country: session.country_name || "Unknown",
          circuit: session.circuit_short_name || "Unknown Circuit",
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
        hasPronostico: true,
        prodeSession: null,
        prodeRace: null,
        score: null,
      });
    });
    return Object.values(eventsMap);
  };

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
            if (prode.p4 !== undefined && prode.p5 !== undefined) {
              sess.prodeRace = prode;
              sess.prodeSession = null;
            } else {
              sess.prodeSession = prode;
              sess.prodeRace = null;
            }
            sess.score = prode.score || 0;
          } else {
            sess.score = null;
          }
        } catch (err) {
          console.error(`Error fetching prode for session ${sess.id}:`, err);
          sess.score = null;
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

        const [upcomingRaw, pastRaw] = await Promise.all([
          getUpcomingSessions(),
          getPastSessions(),
        ]);

        console.log("Sesiones pasadas crudas:", pastRaw);

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
    const isPastEvent = new Date(sessionData.date_start) < new Date();
    const isRace =
      sessionData.sessionName === "Race" && sessionData.sessionType === "Race";

    console.log("PronosticosPage: Navegando con sessionData:", sessionData);

    if (isPastEvent && !isRace) {
      navigate(`/pronosticos/result/${sessionData.id}`, { state: sessionData });
    } else if (!isPastEvent) {
      navigate(`/pronosticos/${sessionData.id}`, { state: sessionData });
    }
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
        <UpcomingEvents
          events={upcomingEvents}
          onPronosticoClick={handlePronosticoClick}
        />
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
