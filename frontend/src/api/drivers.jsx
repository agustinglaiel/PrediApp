import axios from "axios";

// Configura la instancia de Axios con la URL base del backend
const api = axios.create({
  baseURL: "http://localhost:8051", // Cambia el puerto si tu backend está en otro puerto
});

// Función para obtener todos los pilotos desde el backend
export const getAllDrivers = async () => {
  try {
    const response = await api.get("/drivers");
    return response.data; // Retorna los datos de los pilotos
  } catch (error) {
    console.error("Error fetching drivers:", error);
    throw error; // Lanza el error para que pueda manejarse en el frontend si es necesario
  }
};
