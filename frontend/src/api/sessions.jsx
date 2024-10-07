import axios from "axios";

// Configura la URL base de tu API
const API_BASE_URL = "http://localhost:8060"; // Asegúrate de que coincida con tu backend

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

// Puedes agregar más funciones aquí para otras llamadas a la API relacionadas con sesiones
