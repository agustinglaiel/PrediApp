import axios from "axios";

const API_BASE_URL = "http://localhost:8080"; // Asegúrate de que coincida con tu backend

// Obtener el token desde localStorage
const getAuthToken = () => localStorage.getItem("jwtToken");

// Función para obtener las próximas sesiones
export const getUpcomingSessions = async () => {
  try {
    const token = getAuthToken();
    if (!token) {
      throw new Error(
        "No se encontró un token de autenticación. Por favor, inicia sesión."
      );
    }

    const response = await axios.get(`${API_BASE_URL}/sessions/upcoming`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
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

// Función para obtener la sesión por ID
export const getSessionById = async (sessionId) => {
  try {
    const token = getAuthToken();
    if (!token) {
      throw new Error(
        "No se encontró un token de autenticación. Por favor, inicia sesión."
      );
    }

    const response = await axios.get(`${API_BASE_URL}/sessions/${sessionId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data; // Retornamos los datos de la sesión
  } catch (error) {
    throw new Error("Error fetching session by ID: " + error.message);
  }
};
