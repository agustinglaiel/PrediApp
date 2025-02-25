import axios from "axios";

const API_URL = "http://localhost:8054";

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
    const userId = localStorage.getItem("userId");
    if (!userId) {
      throw new Error("User ID not found.");
    }

    // Convertir los valores de p1, p2, y p3 a números
    const prodeDataWithUserId = {
      ...prodeData,
      user_id: parseInt(userId, 10),
      p1: parseInt(prodeData.p1, 10),
      p2: parseInt(prodeData.p2, 10),
      p3: parseInt(prodeData.p3, 10),
    };

    // Imprimir los datos antes de enviarlos
    // console.log("Datos enviados:", prodeDataWithUserId);

    const response = await axios.post(
      `${API_URL}/prodes/session`,
      prodeDataWithUserId,
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
    return response.data;
  } catch (error) {
    console.error("Error creating session prediction:", error);
    console.error("Error details:", error.response || error.message || error);
    throw new Error(
      error.response?.data?.message || "Error creating session prediction."
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
