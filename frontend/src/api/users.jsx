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

// Función de registro que guarda el token en localStorage
export const signUp = async (userData) => {
  try {
    const response = await axios.post(`${API_URL}/api/signup`, userData);
    const { token, refresh_token } = response.data;

    // Almacenar ambos tokens y establecer el token de acceso en las solicitudes
    if (token && refresh_token) {
      localStorage.setItem("jwtToken", token);
      localStorage.setItem("refreshToken", refresh_token);
      setAuthToken(token);
    }

    return response.data;
  } catch (error) {
    console.error("Signup error:", error.response?.data || error.message);
    throw new Error(error.response?.data?.message || "Error signing up.");
  }
};

// Función de inicio de sesión que guarda el token en localStorage
export const login = async (userData) => {
  try {
    const response = await axios.post(`${API_URL}/api/login`, userData);
    const data = response.data;

    // Verificamos si data existe y es un objeto
    if (!data || typeof data !== "object") {
      throw new Error("Respuesta del servidor inválida.");
    }

    // Extraemos token, refresh_token e id de la respuesta
    const { token, refresh_token, id: userId } = data;

    // Verificamos si los campos necesarios existen
    if (!token || !refresh_token || !userId) {
      throw new Error(
        "Faltan datos necesarios en la respuesta del servidor. Verifica que el backend devuelva 'token', 'refresh_token' e 'id'."
      );
    }

    // Almacenar ambos tokens y el userId
    localStorage.setItem("jwtToken", token);
    localStorage.setItem("refreshToken", refresh_token);
    localStorage.setItem("userId", userId);
    setAuthToken(token);

    return data;
  } catch (error) {
    console.error("Login error:", {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status,
    });
    throw new Error(
      error.response?.data?.message ||
        "Error logging in. Verifica tus credenciales o contacta al soporte."
    );
  }
};

export const refreshToken = async () => {
  try {
    const refreshToken = localStorage.getItem("refreshToken");
    if (!refreshToken) {
      throw new Error("No refresh token available. Please log in again.");
    }

    const response = await axios.post(`${API_URL}/api/refresh`, {
      refresh_token: refreshToken,
    });
    const { token } = response.data;

    if (token) {
      localStorage.setItem("jwtToken", token);
      setAuthToken(token);
      return token;
    }

    throw new Error("No token received from refresh endpoint.");
  } catch (error) {
    console.error(
      "Refresh token error:",
      error.response?.data || error.message
    );
    throw new Error(error.response?.data?.message || "Error refreshing token.");
  }
};

export const logout = async () => {
  try {
    const refreshToken = localStorage.getItem("refreshToken");
    if (refreshToken) {
      await axios.post(`${API_URL}/api/signout`, {
        refresh_token: refreshToken,
      });
    }

    // Limpiar localStorage y encabezados
    localStorage.removeItem("jwtToken");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("userId");
    setAuthToken(null);
  } catch (error) {
    console.error("Logout error:", error.response?.data || error.message);
    throw new Error(error.response?.data?.message || "Error logging out.");
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
    console.error("Get user error:", error.response?.data || error.message);
    throw new Error(error.response?.data?.message || "Error fetching user.");
  }
};

// Actualizar usuario por ID (requiere token en el encabezado)
export const updateUserById = async (id, userData) => {
  try {
    const response = await axios.put(`${API_URL}/users/${id}`, userData);
    return response.data;
  } catch (error) {
    console.error("Update user error:", error.response?.data || error.message);
    throw new Error(error.response?.data?.message || "Error updating user.");
  }
};

// Eliminar usuario por ID (requiere token en el encabezado)
export const deleteUserById = async (id) => {
  try {
    await axios.delete(`${API_URL}/users/${id}`);
  } catch (error) {
    console.error("Delete user error:", error.response?.data || error.message);
    throw new Error(error.response?.data?.message || "Error deleting user.");
  }
};

// Obtener todos los usuarios (requiere token en el encabezado)
export const getUsers = async () => {
  try {
    const response = await axios.get(`${API_URL}/users/`);
    return response.data;
  } catch (error) {
    console.error("Get users error:", error.response?.data || error.message);
    throw new Error(error.response?.data?.message || "Error fetching users.");
  }
};

axios.interceptors.response.use(
  (response) => response, // Pasar respuestas exitosas sin cambios
  async (error) => {
    const originalRequest = error.config;
    if (
      error.response?.status === 401 &&
      !originalRequest._retry && // Evitar bucles infinitos
      originalRequest.url !== `${API_URL}/api/refresh` && // No reintentar en el refresh
      originalRequest.url !== `${API_URL}/api/login` && // No reintentar en login
      originalRequest.url !== `${API_URL}/api/signup` // No reintentar en signup
    ) {
      originalRequest._retry = true;
      try {
        const newToken = await refreshToken();
        originalRequest.headers["Authorization"] = `Bearer ${newToken}`;
        return axios(originalRequest); // Reintentar la solicitud original con el nuevo token
      } catch (refreshError) {
        console.error("Failed to refresh token:", refreshError);
        await logout(); // Forzar logout si falla la renovación
        window.location.href = "/login"; // Redirigir al login
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(error); // Rechazar otros errores
  }
);

// Verificar si hay un token almacenado al cargar la aplicación
const token = localStorage.getItem("jwtToken");
if (token) {
  setAuthToken(token);
  console.log("JWT Token establecido desde almacenamiento local");
} else {
  console.warn("No se encontró un JWT Token almacenado");
}
