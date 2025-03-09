import axios from "axios";

const API_BASE_URL = "http://localhost:8080"; // Asegúrate de que coincida con tu backend

// Función para obtener las próximas sesiones (ahora pública, sin token)
export const getUpcomingSessions = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/sessions/upcoming`);
    return response.data;
  } catch (error) {
    console.error("Error fetching upcoming sessions:", error);
    if (error.response) {
      throw new Error(
        `Error en la API: ${error.response.data.message || error.message}`
      );
    }
    throw error;
  }
};

export const getPastSessions = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/sessions/lasts`);
    return response.data;
  } catch (error) {
    console.error("Error fetching past sessions:", error);
    throw error;
  }
};

// Función para obtener la sesión por ID (puedes mantenerla como privada si es necesario, pero la haremos pública también por simplicidad)
export const getSessionById = async (sessionId) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/sessions/${sessionId}`);
    return response.data; // Retornamos los datos de la sesión
  } catch (error) {
    throw new Error("Error fetching session by ID: " + error.message);
  }
};
