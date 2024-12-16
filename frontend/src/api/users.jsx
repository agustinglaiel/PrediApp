import axios from "axios";

const API_URL = "http://localhost:8080";

// Función para establecer el token JWT en el encabezado de autorización
const setAuthToken = (token) => {
  if (token) {
    axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
  } else {
    delete axios.defaults.headers.common["Authorization"];
  }
};

// Función para establecer el RefreshToken en el encabezado de autorización
const setRefreshToken = (refreshToken) => {
  if (refreshToken) {
    axios.defaults.headers.common["Refresh-Token"] = refreshToken;
  } else {
    delete axios.defaults.headers.common["Refresh-Token"];
  }
};

// Función de registro que guarda el token en localStorage
export const signUp = async (userData) => {
  try {
    const response = await axios.post(`${API_URL}/users/signup`, userData);
    const { token } = response.data;
    console.log(token);

    // Almacenar el token y establecerlo en las solicitudes
    if (token) {
      localStorage.setItem("jwtToken", token);
      setAuthToken(token);
    }

    return response.data;
  } catch (error) {
    throw new Error(error.response.data.message || "Error signing up.");
  }
};

// Función de inicio de sesión que guarda el token en localStorage
export const login = async (userData) => {
  try {
    const response = await axios.post(`${API_URL}/users/login`, userData);
    const { token, id: userId } = response.data;

    // Almacenar el token y establecerlo en las solicitudes
    if (token) {
      localStorage.setItem("jwtToken", token);
      localStorage.setItem("refresToken", refreshToken);
      localStorage.setItem("userId", userId);

      setAuthToken(token);
      setRefreshToken(refreshToken);

      console.log("Token:", token);
      console.log("Refresh Token:", refreshToken);
      console.log("User ID:", userId); // Imprimir el userId para verificar
    }

    return response.data;
  } catch (error) {
    throw new Error(error.response.data.message || "Error logging in.");
  }
};

// Obtener usuario por ID (requiere token en el encabezado)
export const getUserById = async (id) => {
  try {
    const token = localStorage.getItem("jwtToken");
    const response = await axios.get(`${API_URL}/users/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    throw new Error(error.response.data.message || "Error fetching user.");
  }
};

// Actualizar usuario por ID (requiere token en el encabezado)
export const updateUserById = async (id, userData) => {
  try {
    const response = await axios.put(`${API_URL}/users/${id}`, userData);
    return response.data;
  } catch (error) {
    throw new Error(error.response.data.message || "Error updating user.");
  }
};

// Eliminar usuario por ID (requiere token en el encabezado)
export const deleteUserById = async (id) => {
  try {
    await axios.delete(`${API_URL}/users/${id}`);
  } catch (error) {
    throw new Error(error.response.data.message || "Error deleting user.");
  }
};

// Obtener todos los usuarios (requiere token en el encabezado)
export const getUsers = async () => {
  try {
    const response = await axios.get(`${API_URL}/users`);
    return response.data;
  } catch (error) {
    throw new Error(error.response.data.message || "Error fetching users.");
  }
};

// Verificar si hay un token almacenado al cargar la aplicación
const token = localStorage.getItem("jwtToken");
if (token) {
  setAuthToken(token);
  console.log("JWT Token establecido desde almacenamiento local");
} else {
  console.warn("No se encontró un JWT Token almacenado");
}

const refreshToken = localStorage.getItem("refreshToken");
if (refreshToken) {
  setRefreshToken(refreshToken);
  console.log("Refresh Token establecido desde almacenamiento local");
} else {
  console.warn("No se encontró un Refresh Token almacenado");
}
