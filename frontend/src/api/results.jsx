import axios from "axios";

const API_URL = "http://localhost:8080";

// New function to get results ordered by position for a session
export const getResultsOrderedByPosition = async (sessionID) => {
  try {
    const response = await axios.get(
      `${API_URL}/results/session/${sessionID}`,
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
    console.log(`Results for session ${sessionID}:`, response.data);
    return response.data; // Returns array of ResponseResultDTO
  } catch (error) {
    console.error(`Error fetching ordered results for session ${sessionID}:`, {
      message: error.message,
      response: error.response?.data,
    });
    throw new Error(
      error.response?.data?.message ||
        "Error fetching ordered results for the session."
    );
  }
};

// New function to get top N drivers in a session
// frontendnuevo/src/api/results.jsx
export const getTopNDriversInSession = async (sessionID, n) => {
  try {
    // Validate inputs
    if (!sessionID || !n) {
      throw new Error("Session ID and number of drivers are required");
    }
    if (n <= 0) {
      throw new Error("Number of drivers must be greater than 0");
    }

    const response = await axios.get(
      `${API_URL}/results/session/${sessionID}/top/${n}`,
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
    console.log(`Top ${n} drivers for session ${sessionID}:`, response.data);
    return response.data; // Returns array of TopDriverDTO
  } catch (error) {
    if (error.response?.status === 404) {
      console.log(
        `No results found for session ${sessionID}, returning empty array`
      );
      return []; // Devolver array vacío en lugar de lanzar error
    }
    console.error(`Error fetching top ${n} drivers for session ${sessionID}:`, {
      message: error.message,
      response: error.response?.data,
    });
    throw new Error(
      error.response?.data?.message ||
        `Error fetching top ${n} drivers for the session.`
    );
  }
};
