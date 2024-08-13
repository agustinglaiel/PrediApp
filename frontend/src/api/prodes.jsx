import axios from "axios";

const API_URL = "http://localhost:8070";

// Crear un prode de carrera
export const createProdeCarrera = async (prodeData) => {
  try {
    const response = await axios.post(`${API_URL}/prodes/carrera`, prodeData);
    return response.data;
  } catch (error) {
    throw new Error(
      error.response.data.message || "Error creating race prediction."
    );
  }
};

// Crear un prode de sesión
export const createProdeSession = async (prodeData) => {
  try {
    const response = await axios.post(`${API_URL}/prodes/session`, prodeData);
    return response.data;
  } catch (error) {
    throw new Error(
      error.response.data.message || "Error creating session prediction."
    );
  }
};

// Obtener un prode de carrera por ID
export const getProdeCarreraByID = async (id) => {
  try {
    const response = await axios.get(`${API_URL}/prode-carrera/${id}`);
    return response.data;
  } catch (error) {
    throw new Error(
      error.response.data.message || "Error fetching race prediction."
    );
  }
};

// Obtener un prode de sesión por ID
export const getProdeSessionByID = async (id) => {
  try {
    const response = await axios.get(`${API_URL}/prode-session/${id}`);
    return response.data;
  } catch (error) {
    throw new Error(
      error.response.data.message || "Error fetching session prediction."
    );
  }
};

// Actualizar un prode de carrera
export const updateProdeCarrera = async (id, prodeData) => {
  try {
    const response = await axios.put(
      `${API_URL}/prode-carrera/${id}`,
      prodeData
    );
    return response.data;
  } catch (error) {
    throw new Error(
      error.response.data.message || "Error updating race prediction."
    );
  }
};

// Actualizar un prode de sesión
export const updateProdeSession = async (id, prodeData) => {
  try {
    const response = await axios.put(
      `${API_URL}/prode-session/${id}`,
      prodeData
    );
    return response.data;
  } catch (error) {
    throw new Error(
      error.response.data.message || "Error updating session prediction."
    );
  }
};

// Eliminar un prode por ID
export const deleteProdeByID = async (id, userID) => {
  try {
    await axios.delete(`${API_URL}/prode/${id}`, {
      params: { userID },
    });
  } catch (error) {
    throw new Error(
      error.response.data.message || "Error deleting prediction."
    );
  }
};
