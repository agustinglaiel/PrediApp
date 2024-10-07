import React, { createContext, useState, useEffect } from "react";
import { getUpcomingSessions } from "../api/sessions"; // Asegúrate de crear este archivo más adelante

// Crear el contexto
export const SessionsContext = createContext();

// Crear el proveedor
export const SessionsProvider = ({ children }) => {
  const [sessions, setSessions] = useState([]); // Inicializamos como un array vacío
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null); // Añadimos manejo de errores

  // Función para obtener las próximas sesiones desde el backend
  const fetchUpcomingSessions = async () => {
    try {
      setLoading(true);
      setError(null); // Reseteamos el estado de error
      const data = await getUpcomingSessions();

      // Verificar que data no sea null o undefined
      if (data && Array.isArray(data)) {
        setSessions(data);
      } else {
        setSessions([]); // Si no es un array, inicializar con array vacío
      }
    } catch (error) {
      console.error("Error fetching upcoming sessions:", error);
      setError("Error fetching upcoming sessions");
    } finally {
      setLoading(false); // Se asegura que loading siempre sea false al finalizar
    }
  };

  // useEffect para obtener las sesiones cuando el componente se monta
  useEffect(() => {
    fetchUpcomingSessions();
  }, []);

  return (
    <SessionsContext.Provider value={{ sessions, setSessions, loading, error }}>
      {children}
    </SessionsContext.Provider>
  );
};
