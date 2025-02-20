import axios from "axios";

// Configura la URL base de tu API
const API_BASE_URL = "http://localhost:8056"; // Asegúrate de que coincida con tu backend

// Función para obtener las próximas sesiones
export const getUpcomingSessions = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/sessions/upcoming`);
    return response.data;
  } catch (error) {
    console.error("Error fetching upcoming sessions:", error);
    throw error;
  }
};

// Función para obtener la sesión por ID
export const getSessionById = async (sessionId) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/sessions/${sessionId}`);
    return response.data; // Retornamos los datos de la sesión
  } catch (error) {
    throw new Error("Error fetching session by ID");
  }
};
