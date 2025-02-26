import React, { useState, useEffect } from "react";
import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import UpcomingEvents from "../components/UpcomingEvents";
import { getUpcomingSessions } from "../api/sessions";
import { getProdeByUserAndSession } from "../api/prodes"; // Usamos el nuevo endpoint

// Simulamos isRaceSession desde el frontend para consistencia con el backend
const isRaceSession = (sessionName, sessionType) => {
  return sessionName === "Race" && sessionType === "Race";
};

const HomePage = () => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchUpcomingSessionsAndProdes = async () => {
      try {
        setLoading(true);
        setEvents([]); // Limpiar el estado antes de fetch para evitar datos cacheados
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

    fetchUpcomingSessionsAndProdes();
  }, []);

  const processSessions = (sessions) => {
    const eventsMap = {};
    let prodeCount = 0;

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

      let day = "1";
      let month = "JAN";
      if (session.date_start && typeof session.date_start === "string") {
        try {
          const [datePart] = session.date_start.split("T");
          if (datePart) {
            const [year, monthNum, dayNum] = datePart.split("-");
            day = dayNum;
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
            month = months[parseInt(monthNum, 10) - 1] || "JAN";
          }
        } catch (error) {
          console.error("Error parsing date_start:", session.date_start, error);
        }
      }

      const [startTime] = session.date_start
        .split("T")[1]
        .split("-")[0]
        .split(":");
      const [endTime] = session.date_end.split("T")[1].split("-")[0].split(":");

      eventsMap[weekendId].sessions.push({
        id: session.id,
        date: day,
        month: month,
        type: session.session_type,
        startTime: `${startTime}:00`,
        endTime: `${endTime}:00`,
        hasPronostico: true,
        prodeSession: null,
        prodeRace: null,
        sessionName: session.session_name, // Añadimos sessionName para usar en isRaceSession
        sessionType: session.session_type, // Añadimos sessionType para usar en isRaceSession
      });
    });

    // Depuración: Contar cuántos pronósticos se fetch y manejar silenciosamente
    const userId = localStorage.getItem("userId");
    if (userId) {
      const updatedEventsMap = { ...eventsMap };
      const processedSessions = new Set(); // Para evitar duplicados

      Object.values(eventsMap).forEach((event) => {
        const prodePromises = event.sessions.map((session) => {
          const sessionKey = `${userId}-${session.id}-${session.type}`;
          if (processedSessions.has(sessionKey)) {
            return Promise.resolve({
              prode: null,
              error: null,
              isRace: isRaceSession(session.sessionName, session.sessionType),
            }); // Evitar solicitudes duplicadas sin logs
          }
          processedSessions.add(sessionKey);

          return getProdeByUserAndSession(parseInt(userId, 10), session.id)
            .then((prode) => {
              if (prode) {
                prodeCount++;
                // Determinar si es un prode de carrera o sesión basado en los campos
                if (prode.p4 && prode.p5 && prode.fastest_lap !== undefined) {
                  return { prode: prode, error: null, isRace: true };
                } else {
                  return { prode: prode, error: null, isRace: false };
                }
              }
              return {
                prode: null,
                error: null,
                isRace: isRaceSession(session.sessionName, session.sessionType),
              };
            })
            .catch((error) => {
              // Silenciar 404 y devolver null sin logs (aunque ahora el backend debería devolver 200)
              if (error.response && error.response.status === 404) {
                return {
                  prode: null,
                  error: null,
                  isRace: isRaceSession(
                    session.sessionName,
                    session.sessionType
                  ),
                };
              }
              console.error(
                `Unexpected error fetching prode for session ${session.id}:`,
                error
              );
              return {
                prode: null,
                error: null,
                isRace: isRaceSession(session.sessionName, session.sessionType),
              };
            });
        });

        Promise.all(prodePromises)
          .then((results) => {
            results.forEach((result, index) => {
              const session = event.sessions[index];
              if (result.isRace) {
                session.prodeRace = result.prode === null ? null : result.prode;
                session.prodeSession = null;
              } else {
                session.prodeSession =
                  result.prode === null ? null : result.prode;
                session.prodeRace = null;
              }
            });
            setEvents(
              Object.values(updatedEventsMap).sort((a, b) => {
                const dateA = new Date(
                  a.sessions[0].date_start || "2025-01-01"
                );
                const dateB = new Date(
                  b.sessions[0].date_start || "2025-01-01"
                );
                return dateA - dateB;
              })
            );
          })
          .catch((error) =>
            console.error("Error fetching prodes in batch:", error)
          );
      });
    } else {
      return Object.values(eventsMap).sort((a, b) => {
        const dateA = new Date(a.sessions[0].date_start || "2025-01-01");
        const dateB = new Date(b.sessions[0].date_start || "2025-01-01");
        return dateA - dateB;
      });
    }

    return Object.values(eventsMap).sort((a, b) => {
      const dateA = new Date(a.sessions[0].date_start || "2025-01-01");
      const dateB = new Date(b.sessions[0].date_start || "2025-01-01");
      return dateA - dateB;
    });
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
