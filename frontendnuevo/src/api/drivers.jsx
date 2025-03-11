import axios from "axios";

// Configura la instancia de Axios con la URL base del backend
const api = axios.create({
  baseURL: "http://localhost:8080", // Cambia el puerto si tu backend está en otro puerto
});

// Función para obtener todos los pilotos desde el backend
export const getAllDrivers = async () => {
  try {
    const response = await api.get("/drivers", {
      headers: { Authorization: `Bearer ${localStorage.getItem("jwtToken")}` },
    });
    console.log("Pilotos recibidos desde el backend:", response.data); // Log para depuración
    return response.data; // Retorna los datos de los pilotos
  } catch (error) {
    console.error("Error fetching drivers - Detalle:", {
      message: error.message,
      code: error.code,
      config: error.config,
      request: error.request,
      response: error.response,
    });
    throw new Error(
      error.response?.data?.message ||
        "Error fetching drivers from the backend."
    ); // Lanza un error personalizado
  }
};

export const getDriverById = async (driverId) => {
  try {
    const response = await api.get(`/drivers/${driverId}`, {
      headers: { Authorization: `Bearer ${localStorage.getItem("jwtToken")}` },
    });
    console.log(`Piloto con ID ${driverId} recibido:`, response.data);
    return response.data; // Retorna el objeto ResponseDriverDTO
  } catch (error) {
    console.error(`Error fetching driver with ID ${driverId}:`, {
      message: error.message,
      response: error.response,
    });
    throw new Error(
      error.response?.data?.message || `Error fetching driver ${driverId}.`
    );
  }
};