import axios from "axios";

const API_URL = "http://localhost:8080";

// Crear un prode de carrera
export const createProdeCarrera = async (prodeData) => {
  try {
    const userId = localStorage.getItem("userId");
    if (!userId) {
      throw new Error("User ID not found in localStorage");
    }

    const payload = {
      user_id: parseInt(userId, 10),
      session_id: parseInt(prodeData.session_id, 10),
      p1: parseInt(prodeData.p1, 10),
      p2: parseInt(prodeData.p2, 10),
      p3: parseInt(prodeData.p3, 10),
      p4: parseInt(prodeData.p4, 10),
      p5: parseInt(prodeData.p5, 10),
      vsc: !!prodeData.vsc, // forzar bool
      sc: !!prodeData.sc, // forzar bool
      dnf: parseInt(prodeData.dnf, 10),
      // fastest_lap: parseInt(prodeData.fastest_lap, 10) || 0
    };

    const response = await axios.post(`${API_URL}/prodes/carrera`, payload, {
      headers: { "Content-Type": "application/json" },
    });
    return response.data;
  } catch (error) {
    console.error("Error creating/updating race prediction:", error);
    const errorMsg =
      error.response?.data?.message || "Error creating race prediction.";
    throw new Error(errorMsg);
  }
};

// Crear un prode de sesión
export const createProdeSession = async (prodeData) => {
  try {
    const userId = localStorage.getItem("userId");
    if (!userId) {
      throw new Error("User ID not found in localStorage");
    }

    const payload = {
      user_id: parseInt(userId, 10),
      session_id: parseInt(prodeData.session_id, 10),
      p1: parseInt(prodeData.p1, 10),
      p2: parseInt(prodeData.p2, 10),
      p3: parseInt(prodeData.p3, 10),
    };

    const response = await axios.post(`${API_URL}/prodes/session`, payload, {
      headers: { "Content-Type": "application/json" },
    });
    return response.data;
  } catch (error) {
    console.error("Error creating/updating session prediction:", error);
    const errorMsg =
      error.response?.data?.message || "Error creating session prediction.";
    throw new Error(errorMsg);
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

// Obtener un prode de carrera por user_id y session_id
export const getRaceProdeByUserAndSession = async (userId, sessionId) => {
  try {
    const token = localStorage.getItem("jwtToken");
    if (!token) {
      throw new Error("Authentication token not found. Please log in.");
    }

    const response = await axios.get(
      `${API_URL}/prodes/carrera/user/${userId}/session/${sessionId}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        // Aceptar todos los códigos 200 para manejar manualmente
        validateStatus: (status) => status === 200,
      }
    );

    // Verificar si la respuesta está vacía, null, o un objeto vacío
    if (
      !response.data ||
      Object.keys(response.data).length === 0 ||
      (typeof response.data === "object" &&
        Object.values(response.data).every((v) => v === null || v === ""))
    ) {
      return null; // Devolver null silenciosamente si no hay datos
    }
    console.log(
      `Successfully fetched race prode for user ${userId}, session ${sessionId}:`,
      response.data
    );
    return response.data;
  } catch (error) {
    // Silenciar errores 404/400 si los hubiera (aunque el backend ahora devuelve 200)
    if (
      error.response &&
      (error.response.status === 404 || error.response.status === 400)
    ) {
      return null; // Devolver null silenciosamente para 404/400 (por si acaso)
    }
    console.log(
      `Unexpected error fetching race prode for user ${userId}, session ${sessionId}:`,
      error.message
    );
    throw new Error(error.message || "Error fetching race prediction.");
  }
};

// // Obtener un prode de sesión por user_id y session_id
// export const getSessionProdeByUserAndSession = async (userId, sessionId) => {
//   try {
//     const token = localStorage.getItem("jwtToken");
//     if (!token) {
//       throw new Error("Authentication token not found. Please log in.");
//     }

//     const response = await axios.get(
//       `${API_URL}/prodes/session/user/${userId}/session/${sessionId}`,
//       {
//         headers: {
//           Authorization: `Bearer ${token}`,
//           "Content-Type": "application/json",
//         },
//         // Aceptar todos los códigos 200 para manejar manualmente
//         validateStatus: (status) => status === 200,
//       }
//     );

//     // Verificar si la respuesta está vacía, null, o un objeto vacío
//     if (
//       !response.data ||
//       Object.keys(response.data).length === 0 ||
//       (typeof response.data === "object" &&
//         Object.values(response.data).every((v) => v === null || v === ""))
//     ) {
//       return null; // Devolver null silenciosamente si no hay datos
//     }
//     console.log(
//       `Successfully fetched session prode for user ${userId}, session ${sessionId}:`,
//       response.data
//     );
//     return response.data;
//   } catch (error) {
//     // Silenciar errores 404/400 si los hubiera (aunque el backend ahora devuelve 200)
//     if (
//       error.response &&
//       (error.response.status === 404 || error.response.status === 400)
//     ) {
//       return null; // Devolver null silenciosamente para 404/400 (por si acaso)
//     }
//     console.log(
//       `Unexpected error fetching session prode for user ${userId}, session ${sessionId}:`,
//       error.message
//     );
//     throw new Error(error.message || "Error fetching session prediction.");
//   }
// };

// Obtener un prode de sesión por user_id y session_id

export const getProdeByUserAndSession = async (userId, sessionId) => {
  try {
    const token = localStorage.getItem("jwtToken");
    if (!token) {
      throw new Error("Authentication token not found. Please log in.");
    }

    const response = await axios.get(
      `${API_URL}/prodes/user/${userId}/session/${sessionId}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        // Aceptar 200 para manejar manualmente
        validateStatus: (status) => status === 200,
      }
    );

    // Verificar si la respuesta es un array vacío o null
    if (
      !response.data ||
      (Array.isArray(response.data) && response.data.length === 0)
    ) {
      return null; // Devolver null silenciosamente si no hay datos
    }
    // Asumimos que el array contiene un solo objeto (el pronóstico)
    // Determinar si es un pronóstico de carrera o sesión basado en los campos
    const prode = response.data[0];

    return prode; // Devolver el primer (y único) elemento del array
  } catch (error) {
    console.log(
      `Unexpected error fetching prode for user ${userId}, session ${sessionId}:`,
      error.message
    );
    throw new Error(error.message || "Error fetching prode prediction.");
  }
};
