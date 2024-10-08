import React, { useContext } from "react";
import { useNavigate } from "react-router-dom";
import { SessionsContext } from "../contexts/SessionsContext";
import "../styles/SessionsList.css";

// Función para agrupar sesiones por circuito
const groupSessionsByLocation = (sessions) => {
  const grouped = {};

  sessions.forEach((session) => {
    const key = `${session.location}-${session.country_name}`;

    if (!grouped[key]) {
      grouped[key] = {
        location: session.location,
        country: session.country_name,
        sessions: [],
      };
    }

    grouped[key].sessions.push(session);
  });

  return Object.values(grouped);
};

const Home = () => {
  const { sessions, loading, error } = useContext(SessionsContext);
  const navigate = useNavigate(); // Para redireccionar

  const validCombinations = {
    "Sprint Qualifying": "Qualifying",
    Sprint: "Race",
    Qualifying: "Qualifying",
    Race: "Race",
  };

  // Filtrar las sesiones en el frontend
  const filteredSessions = sessions.filter(
    (session) =>
      validCombinations[session.session_name] === session.session_type
  );

  if (loading) {
    return <div>Loading sessions...</div>;
  }

  if (error) {
    return <div>Error loading sessions: {error}</div>;
  }

  const groupedSessions = groupSessionsByLocation(filteredSessions);

  // Función para manejar la redirección según el tipo de sesión
  const handlePredictionClick = (session) => {
    if (session.session_name === "Race" && session.session_type === "Race") {
      navigate(`/prode/race/${session.id}`);
    } else {
      navigate(`/prode/session/${session.id}`);
    }
  };

  return (
    <div className="page-container">
      <h2 className="page-title">Upcoming Events</h2>
      <div className="sessions-container">
        {groupedSessions.map((group, index) => (
          <div key={index} className="location-group">
            <div className="location-header">
              <h3>{group.location}</h3>
              <p>{group.country}</p>
            </div>
            <div className="session-island">
              <ul>
                {group.sessions.map((session) => (
                  <li key={session.id} className="session-item">
                    <div className="session-details">
                      <span className="session-date">
                        {new Date(session.date_start).toLocaleDateString(
                          "en-GB",
                          {
                            day: "numeric",
                            month: "short",
                          }
                        )}
                      </span>
                      <span className="session-name">
                        {session.session_name}
                      </span>
                      <span className="session-time">
                        {new Date(session.date_start).toLocaleTimeString([], {
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </span>
                    </div>
                    <button
                      className="prediction-button"
                      onClick={() => handlePredictionClick(session)}
                    >
                      Completar pronóstico
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Home;
